package reverseproxy

import (
	"net/http"

	"github.com/kevinhury/membrane/config"
)

// Registry interface
type Registry interface {
	Serve(w http.ResponseWriter, r *http.Request) error
	SetConfig([]byte) error
	ConfigMap() *config.Configuration
}
