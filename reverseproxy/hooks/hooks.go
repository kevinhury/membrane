package hooks

import (
	"net/http"

	"github.com/kevinhury/membrane/config"
)

// PreRequest type
type PreRequest func(*http.Request, http.ResponseWriter, config.Plugin) error

// PostRequest type
type PostRequest func(resp *http.Response, plugin config.Plugin) error

// PreHook interface
type PreHook interface {
	PreHook(*http.Request, http.ResponseWriter, config.Plugin) error
}

// PostHook interface
type PostHook interface {
	PostHook(resp *http.Response, plugin config.Plugin) error
}
