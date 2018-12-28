package httpserver

import (
	"log"
	"net/http"

	"github.com/kevinhury/membrane/httpserver/logger"
	"github.com/kevinhury/membrane/httpserver/middlewares"
	"github.com/kevinhury/membrane/httpserver/recover"
)

// StartServer func
func StartServer(h http.HandlerFunc, addr string) {
	m := http.NewServeMux()

	root := middlewares.Chain(h, recover.Middleware, logger.Middleware)
	m.HandleFunc("/", root)

	s := &http.Server{
		Addr:    addr,
		Handler: m,
	}

	log.Fatal(s.ListenAndServe())
}
