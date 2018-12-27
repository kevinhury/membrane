package middlewares

import "net/http"

// Middleware type
type Middleware func(http.HandlerFunc) http.HandlerFunc

// Chain middlewares in order
func Chain(handler http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	if len(middlewares) == 0 {
		return handler
	}

	wrapped := handler
	for i := len(middlewares) - 1; i >= 0; i-- {
		wrapped = middlewares[i](wrapped)
	}

	return wrapped
}

// ChainSlice middleware accepts slice as list of middlewares
func ChainSlice(handler http.HandlerFunc, middlewares []Middleware) http.HandlerFunc {
	return Chain(handler, middlewares...)
}
