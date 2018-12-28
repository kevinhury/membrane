package server

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/kevinhury/membrane/proxy"
	"github.com/kevinhury/membrane/proxy/config"
)

// ReverseProxy struct
type ReverseProxy struct {
	config *config.Configuration
}

// Serve func
func (rp *ReverseProxy) Serve(w http.ResponseWriter, r *http.Request) error {
	url, err := rp.parseTarget(r)
	if err != nil {
		return err
	}

	proxy := httputil.NewSingleHostReverseProxy(url)

	r.URL.Host = url.Host
	r.URL.Scheme = url.Scheme
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	r.Host = url.Host

	proxy.ServeHTTP(w, r)

	return nil
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

func (rp *ReverseProxy) parseTarget(req *http.Request) (*url.URL, error) {
	var target string

	pipelines := rp.config.Pipelines(req.Host, req.URL.Path)
	if len(pipelines) == 0 {
		return req.URL, nil
	}

	for _, p := range pipelines {
		for _, policy := range p.Policies {
			if policy.Name != "proxy" {
				continue
			}
			name := policy.Action.ServiceEndpoint
			target = rp.config.Service(name).URL
			break
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

func (rp *ReverseProxy) matchingEndpoint() bool {
	return false
}
