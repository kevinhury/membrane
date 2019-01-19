package reverseproxy

import (
	"errors"
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
	proxies   map[string]*httputil.ReverseProxy
	prehooks  map[string]hooks.PreHook
	posthooks map[string]hooks.PostHook
}

// NewWithConfig func
func NewWithConfig(content []byte) *Registry {
	conf, err := config.NewWithData(content)
	if err != nil {
		return nil
	}

	outHosts := make([]string, len(conf.ConfigMap.OutboundEndpoints))
	for i := 0; i < len(outHosts); i++ {
		outHosts[i] = conf.ConfigMap.OutboundEndpoints[0].URL
	}

	return &Registry{
		config:    conf,
		proxies:   initProxies(outHosts),
		prehooks:  hooks.Prehooks(conf),
		posthooks: hooks.Posthooks(conf),
	}
}

func initProxies(names []string) map[string]*httputil.ReverseProxy {
	proxies := make(map[string]*httputil.ReverseProxy, len(names))

	for i := 0; i < len(names); i++ {
		rawurl := names[i]
		url, err := url.Parse(rawurl)
		if err != nil {
			continue
		}
		proxies[rawurl] = httputil.NewSingleHostReverseProxy(url)
	}

	return proxies
}

// ConfigMap func
func (reg *Registry) ConfigMap() *config.Configuration {
	return reg.config
}

// Serve func
func (reg *Registry) Serve(w http.ResponseWriter, r *http.Request) error {
	pipelines := reg.config.Pipelines(r.Host, r.URL.Path, r.Method)
	target, err := reg.matchOutboundEndpoint(r, pipelines)
	if err != nil {
		return err
	}
	proxy := reg.proxies[target.URL]
	proxy.ModifyResponse = func(resp *http.Response) error {
		reg.runPostHooks(resp, pipelines)
		return nil
	}
	err = reg.runPreHooks(w, r, pipelines)
	if err != nil {
		return nil
	}

	proxy.ServeHTTP(w, r)

	return nil
}

func (reg *Registry) runPreHooks(w http.ResponseWriter, r *http.Request, pipelines []config.Pipeline) error {
	for i := 0; i < len(pipelines); i++ {
		pipeline := pipelines[i]
		for j := 0; j < len(pipeline.Plugins); j++ {
			plugin := pipeline.Plugins[j]
			if hook, ok := reg.prehooks[plugin.Name]; ok {
				err := hook.PreHook(r, w, plugin)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (reg *Registry) runPostHooks(resp *http.Response, pipelines []config.Pipeline) {
	for i := 0; i < len(pipelines); i++ {
		pipeline := pipelines[i]
		for j := 0; j < len(pipeline.Plugins); j++ {
			plugin := pipeline.Plugins[j]
			if hook, ok := reg.posthooks[plugin.Name]; ok {
				err := hook.PostHook(resp, plugin)
				if err != nil {
					return
				}
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

	reg.config = conf

	return nil
}

func (reg *Registry) matchOutboundEndpoint(req *http.Request, pipelines []config.Pipeline) (*config.OutboundEndpoint, error) {
	if len(pipelines) == 0 {
		return nil, errors.New("Unsupported URL")
	}

	for idx := range pipelines {
		p := pipelines[idx]
		plugs := p.PluginsMatchingName("proxy")
		for idx := range plugs {
			plugin := plugs[idx]
			if action, ok := plugin.Action.(actions.Proxy); ok {
				return reg.config.Service(action.OutboundEndpoint), nil
			}
		}
	}

	return nil, errors.New("Non matching outbound endpoint")
}
