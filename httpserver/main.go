package httpserver

import (
	"net/http"

	"github.com/kevinhury/membrane/httpserver/logger"
	"github.com/kevinhury/membrane/httpserver/middlewares"
	"github.com/kevinhury/membrane/httpserver/recover"
)

func rootHandlerFunc(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}

// StartServer func
func StartServer() {
	m := http.NewServeMux()

	root := middlewares.Chain(rootHandlerFunc, recover.Middleware, logger.Middleware)
	m.HandleFunc("/", root)

	s := &http.Server{
		Addr:    ":8000",
		Handler: m,
	}

	s.ListenAndServe()
}
