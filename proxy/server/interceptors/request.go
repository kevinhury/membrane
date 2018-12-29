package interceptors

import (
	"io"
	"net/http"

	"github.com/kevinhury/membrane/config"
)

// RequestModifier func
func RequestModifier(r *http.Request, pipelines []config.Pipeline) (io.ReadCloser, int64) {
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
