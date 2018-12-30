package server

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/kevinhury/membrane/config/actions"

	"github.com/kevinhury/membrane/config"
	"github.com/kevinhury/membrane/proxy"
	"github.com/kevinhury/membrane/proxy/server/interceptors"
)

// ReverseProxy struct
type ReverseProxy struct {
	config *config.Configuration
}

// NewWithConfigFile func
func NewWithConfigFile(fileName string, watch bool) proxy.Proxy {
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil
	}

	conf, err := config.NewWithData(content)
	if err != nil {
		return nil
	}

	return &ReverseProxy{config: conf}
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
	proxy.ModifyResponse = interceptors.ResponseModifier(pipelines)
	interceptors.RequestModifier(r, pipelines)

	r.URL.Host = url.Host
	r.URL.Scheme = url.Scheme
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	r.Host = url.Host

	proxy.ServeHTTP(w, r)

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
