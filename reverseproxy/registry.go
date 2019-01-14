package reverseproxy

import (
	"errors"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/kevinhury/membrane/config/actions"

	"github.com/kevinhury/membrane/config"
	"github.com/kevinhury/membrane/reverseproxy/hooks"
)

// Registry struct
type Registry struct {
	config    *config.Configuration
	inHosts   map[string]*httputil.ReverseProxy
	outHosts  map[string]*httputil.ReverseProxy
	prehooks  map[string]hooks.PreHook
	posthooks map[string]hooks.PostHook
}

// NewWithConfig func
func NewWithConfig(content []byte) *Registry {
	conf, err := config.NewWithData(content)
	if err != nil {
		return nil
	}

	inHosts := make([]string, len(conf.ConfigMap.InboundEndpoints))
	for i := 0; i < len(inHosts); i++ {
		inHosts[i] = conf.ConfigMap.InboundEndpoints[0].Host
	}
	outHosts := make([]string, len(conf.ConfigMap.OutboundEndpoints))
	for i := 0; i < len(outHosts); i++ {
		outHosts[i] = conf.ConfigMap.OutboundEndpoints[0].URL
	}

	return &Registry{
		config:    conf,
		inHosts:   initHosts(inHosts),
		outHosts:  initHosts(outHosts),
		prehooks:  hooks.Prehooks(conf),
		posthooks: hooks.Posthooks(conf),
	}
}

func initHosts(names []string) map[string]*httputil.ReverseProxy {
	hosts := make(map[string]*httputil.ReverseProxy, len(names))

	for i := 0; i < len(names); i++ {
		rawurl := names[i]
		url, err := url.Parse(rawurl)
		if err != nil {
			continue
		}
		hosts[rawurl] = httputil.NewSingleHostReverseProxy(url)
	}

	return hosts
}

// ConfigMap func
func (reg *Registry) ConfigMap() *config.Configuration {
	return reg.config
}

// Serve func
func (reg *Registry) Serve(w http.ResponseWriter, r *http.Request) error {
	log.Printf("ReverseProxy::Serve r.Host(%s) r.URL.Path(%s) r.Method(%s)", r.Host, r.URL.Path, r.Method)
	pipelines := reg.config.Pipelines(r.Host, r.URL.Path, r.Method)
	url, err := reg.parseTarget(r, pipelines)
	if err != nil {
		return err
	}

	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.ModifyResponse = func(resp *http.Response) error {
		reg.runPostHooks(resp, pipelines)
		return nil
	}
	reg.runPreHooks(w, r, pipelines)

	proxy.ServeHTTP(w, r)

	return nil
}

func (reg *Registry) runPreHooks(w http.ResponseWriter, r *http.Request, pipelines []config.Pipeline) {
	for i := 0; i < len(pipelines); i++ {
		pipeline := pipelines[i]
		for j := 0; j < len(pipeline.Plugins); j++ {
			plugin := pipeline.Plugins[j]
			if hook, ok := reg.prehooks[plugin.Name]; ok {
				hook.PreHook(r, w, plugin)
			}
		}
	}
}

func (reg *Registry) runPostHooks(resp *http.Response, pipelines []config.Pipeline) {
	for i := 0; i < len(pipelines); i++ {
		pipeline := pipelines[i]
		for j := 0; j < len(pipeline.Plugins); j++ {
			plugin := pipeline.Plugins[j]
			if hook, ok := reg.posthooks[plugin.Name]; ok {
				hook.PostHook(resp, plugin)
			}
		}
	}
}

// SetConfig func
func (reg *Registry) SetConfig(content []byte) error {
	conf, err := config.NewWithData(content)
	if err != nil {
		return err
	}

	// TODO: Need to wait for all pending requests and then change
	reg.config = conf

	return nil
}

func (reg *Registry) parseTarget(req *http.Request, pipelines []config.Pipeline) (*url.URL, error) {
	var target string

	if len(pipelines) == 0 {
		return nil, errors.New("Unsupported URL")
	}

	for idx := range pipelines {
		p := pipelines[idx]
		plugs := p.PluginsMatchingName("proxy")
		for idx := range plugs {
			plugin := plugs[idx]
			if action, ok := plugin.Action.(actions.Proxy); ok {
				target = reg.config.Service(action.OutboundEndpoint).URL
				break
			} else {
				log.Printf("Could not parse plugin %+v\n", plugin)
			}
		}
	}

	if target == "" {
		return req.URL, nil
	}

	url, err := url.Parse(target)
	if err != nil {
		return nil, err
	}

	log.Printf("Proxying to url %s\n", url)

	return url, nil
}
