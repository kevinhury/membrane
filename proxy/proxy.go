package proxy

import "net/http"

// RequestInterceptor struct
type RequestInterceptor struct{}

// ResponseInterceptor struct
type ResponseInterceptor struct{}

// Proxy interface
type Proxy interface {
	Serve(w http.ResponseWriter, r *http.Request) error
}
