package server

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"

	"github.com/kevinhury/membrane/config"
	"github.com/kevinhury/membrane/proxy"
)

// ReverseProxy struct
type ReverseProxy struct {
	config *config.Configuration
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
	proxy.Transport = rp.getTransport(r, pipelines)

	r.URL.Host = url.Host
	r.URL.Scheme = url.Scheme
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	r.Host = url.Host
	r.Body, r.ContentLength = rp.getRequestBody(r, pipelines)

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

func (rp *ReverseProxy) parseTarget(req *http.Request, pipelines []config.Pipeline) (*url.URL, error) {
	var target string

	if len(pipelines) == 0 {
		return nil, errors.New("Unsupported URL")
	}

	for _, p := range pipelines {
		for _, plugin := range p.Plugins {
			if plugin.Name != "proxy" {
				continue
			}
			// Check if type assertions are checked in if statement
			if name, ok := plugin.Action["outboundEndpoint"].(string); ok {
				target = rp.config.Service(name).URL
				break
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

func (rp *ReverseProxy) getRequestBody(r *http.Request, pipelines []config.Pipeline) (io.ReadCloser, int64) {
	var target string

	for _, p := range pipelines {
		for _, plugin := range p.Plugins {
			if plugin.Name != "request-transform" {
				continue
			}
			target = ""
			break
		}
	}
	if target == "" {
		return r.Body, r.ContentLength
	}
	// TODO set the body here
	// b := ioutil.ReadAll(r.Body);
	// r.Body = ioutil.NopCloser(bytes.NewBuffer(b))
	return r.Body, r.ContentLength
}

func (rp *ReverseProxy) getTransport(r *http.Request, pipelines []config.Pipeline) http.RoundTripper {
	var target *config.Plugin

	for _, p := range pipelines {
		for _, plugin := range p.Plugins {
			if plugin.Name != "proxy" {
				continue
			}
			target = &plugin
			break
		}
	}
	if target == nil {
		return http.DefaultTransport
	}

	return &customTransport{
		plugin: target,
	}
}

type customTransport struct {
	plugin *config.Plugin
	http.RoundTripper
}

func (t *customTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	resp, err := t.RoundTripper.RoundTrip(r)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}
	b = bytes.Replace(b, []byte("server"), []byte("schmerver"), -1)
	body := ioutil.NopCloser(bytes.NewReader(b))
	resp.Body = body
	resp.ContentLength = int64(len(b))
	resp.Header.Set("Content-Length", strconv.Itoa(len(b)))
	return resp, nil
}

// func (ct *customTransport) RoundTrip(r *http.Request) (*http.Response, error) {
// 	response, err := http.DefaultTransport.RoundTrip(r)
// 	if err != nil {
// 		return nil, err
// 	}

// 	body, err := httputil.DumpResponse(response, true)
// 	if err != nil {
// 		return nil, err
// 	}

// 	log.Printf("Intercepted Response %s", string(body))
// 	return response, nil
// }
