package cors

import (
	"net/http"

	"github.com/kevinhury/membrane/config"
	"github.com/kevinhury/membrane/config/actions"
)

// Hook struct
type Hook struct {
	Config *config.Configuration
}

// PreHook func
func (h Hook) PreHook(r *http.Request, w http.ResponseWriter, plugin config.Plugin) error {
	action := plugin.Action.(actions.Cors)

	origin := "*"
	methods := "HEAD,GET,POST,PUT,PATCH,DELETE"

	if action.Origin != "" {
		origin = action.Origin
	}
	if action.Methods != "" {
		methods = action.Methods
	}

	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Access-Control-Allow-Methods", methods)

	if action.Headers != "" {
		w.Header().Set("Access-Control-Allow-Headers", action.Headers)
	}

	return nil
}
