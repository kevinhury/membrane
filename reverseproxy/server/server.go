package server

import (
	"errors"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/kevinhury/membrane/reverseproxy/hooks/restransform"

	"github.com/kevinhury/membrane/reverseproxy"
	"github.com/kevinhury/membrane/reverseproxy/hooks/jwt"
	"github.com/kevinhury/membrane/reverseproxy/hooks/proxy"
	"github.com/kevinhury/membrane/reverseproxy/hooks/reqtransform"

	"github.com/kevinhury/membrane/config/actions"

	"github.com/kevinhury/membrane/config"
	"github.com/kevinhury/membrane/reverseproxy/hooks"
)

// ReverseProxy struct
type ReverseProxy struct {
	config    *config.Configuration
	inHosts   map[string]*httputil.ReverseProxy
	outHosts  map[string]*httputil.ReverseProxy
	prehooks  map[string]hooks.PreHook
	posthooks map[string]hooks.PostHook
}

// NewWithConfig func
func NewWithConfig(content []byte) reverseproxy.Registry {
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

	return &ReverseProxy{
		config:    conf,
		inHosts:   initHosts(inHosts),
		outHosts:  initHosts(outHosts),
		prehooks:  initPrehooks(conf),
		posthooks: initPosthooks(conf),
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

func initPrehooks(conf *config.Configuration) map[string]hooks.PreHook {
	return map[string]hooks.PreHook{
		"proxy":             proxy.Hook{Config: conf},
		"request-transform": reqtransform.Hook{Config: conf},
		"jwt":               jwt.Hook{Config: conf},
	}
}

func initPosthooks(conf *config.Configuration) map[string]hooks.PostHook {
	return map[string]hooks.PostHook{
		"response-transform": restransform.Hook{Config: conf},
	}
}

// ConfigMap func
func (rp *ReverseProxy) ConfigMap() *config.Configuration {
	return rp.config
}

// Serve func
func (rp *ReverseProxy) Serve(w http.ResponseWriter, r *http.Request) error {
	log.Printf("ReverseProxy::Serve r.Host(%s) r.URL.Path(%s) r.Method(%s)", r.Host, r.URL.Path, r.Method)
	pipelines := rp.config.Pipelines(r.Host, r.URL.Path, r.Method)
	url, err := rp.parseTarget(r, pipelines)
	if err != nil {
		return err
	}

	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.ModifyResponse = func(resp *http.Response) error {
		rp.runPostHooks(resp, pipelines)
		return nil
	}
	rp.runPreHooks(w, r, pipelines)

	proxy.ServeHTTP(w, r)

	return nil
}

func (rp *ReverseProxy) runPreHooks(w http.ResponseWriter, r *http.Request, pipelines []config.Pipeline) {
	for i := 0; i < len(pipelines); i++ {
		pipeline := pipelines[i]
		for j := 0; j < len(pipeline.Plugins); j++ {
			plugin := pipeline.Plugins[j]
			if hook, ok := rp.prehooks[plugin.Name]; ok {
				hook.PreHook(r, w, plugin)
			}
		}
	}
}

func (rp *ReverseProxy) runPostHooks(resp *http.Response, pipelines []config.Pipeline) {
	for i := 0; i < len(pipelines); i++ {
		pipeline := pipelines[i]
		for j := 0; j < len(pipeline.Plugins); j++ {
			plugin := pipeline.Plugins[j]
			if hook, ok := rp.posthooks[plugin.Name]; ok {
				hook.PostHook(resp, plugin)
			}
		}
	}
}

// SetConfig func
func (rp *ReverseProxy) SetConfig(content []byte) error {
	conf, err := config.NewWithData(content)
	if err != nil {
		return err
	}

	// TODO: Need to wait for all pending requests and then change
	rp.config = conf

	return nil
}

func (rp *ReverseProxy) parseTarget(req *http.Request, pipelines []config.Pipeline) (*url.URL, error) {
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
				target = rp.config.Service(action.OutboundEndpoint).URL
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
