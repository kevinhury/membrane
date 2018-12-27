package httpserver

import (
	"net/http"

	"github.com/kevinhury/membrane/httpserver/logger"
	"github.com/kevinhury/membrane/httpserver/middlewares"
)

func rootHandlerFunc(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}

// StartServer func
func StartServer() {
	m := http.NewServeMux()

	root := middlewares.Chain(rootHandlerFunc, logger.Logger)
	m.HandleFunc("/", root)

	s := &http.Server{
		Addr:    ":8000",
		Handler: m,
	}

	s.ListenAndServe()
}
