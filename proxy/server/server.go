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
	target := "http://localhost:3112"
	url, err := url.Parse(target)
	if err != nil {
		return err
	}

	log.Printf("Proxying to url %s\n", url)
	proxy := httputil.NewSingleHostReverseProxy(url)

	r.URL.Host = url.Host
	r.URL.Scheme = url.Scheme
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	r.Host = url.Host

	proxy.ServeHTTP(w, r)

	return nil
}

// NewWithConfigFile func
func NewWithConfigFile(fileName string) proxy.Proxy {
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil
	}

	conf, err := config.Parse(content)
	if err != nil {
		return nil
	}

	return &ReverseProxy{config: conf}
}
