package proxy

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/kevinhury/membrane/config"
	"github.com/kevinhury/membrane/config/actions"
)

// Hook struct
type Hook struct {
	Config *config.Configuration
}

// PreHook func
func (h Hook) PreHook(r *http.Request, w http.ResponseWriter, plugin config.Plugin) error {
	act, ok := plugin.Action.(actions.Proxy)
	if !ok {
		return errors.New("Unsupported action")
	}

	target := h.Config.Service(act.OutboundEndpoint).URL

	url, err := url.Parse(target)
	if err != nil {
		return err
	}

	r.URL.Host = url.Host
	r.URL.Scheme = url.Scheme
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	r.Host = url.Host

	return nil
}
