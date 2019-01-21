package hooks

import (
	"net/http"

	"github.com/kevinhury/membrane/config"
	"github.com/kevinhury/membrane/reverseproxy/hooks/cors"
	"github.com/kevinhury/membrane/reverseproxy/hooks/jwt"
	jwtExtract "github.com/kevinhury/membrane/reverseproxy/hooks/jwt-extract"
	"github.com/kevinhury/membrane/reverseproxy/hooks/proxy"
	"github.com/kevinhury/membrane/reverseproxy/hooks/reqtransform"
	"github.com/kevinhury/membrane/reverseproxy/hooks/restransform"
)

// PreHook interface
type PreHook interface {
	PreHook(*http.Request, http.ResponseWriter, config.Plugin) error
}

// PostHook interface
type PostHook interface {
	PostHook(resp *http.Response, plugin config.Plugin) error
}

// Prehooks func
func Prehooks(conf *config.Configuration) map[string]PreHook {
	return map[string]PreHook{
		"proxy":             proxy.Hook{Config: conf},
		"request-transform": reqtransform.Hook{Config: conf},
		"jwt":               jwt.Hook{Config: conf},
		"jwt-extract":       jwtExtract.Hook{Config: conf},
		"cors":              cors.Hook{Config: conf},
	}
}

// Posthooks func
func Posthooks(conf *config.Configuration) map[string]PostHook {
	return map[string]PostHook{
		"response-transform": restransform.Hook{Config: conf},
	}
}
