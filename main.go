package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/kevinhury/membrane/proxy"
	"github.com/kevinhury/membrane/proxy/server"

	"github.com/kevinhury/membrane/httpserver"
)

func handler(pr proxy.Proxy) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := pr.Serve(w, r)
		if err != nil && err.Error() == "Unsupported URL" {
			log.Printf("%s %s %s %s\n", err.Error(), r.Host, r.URL.Path, r.Method)
			w.WriteHeader(http.StatusNotFound)
		} else if err != nil {
			log.Printf("[erro] Received error from root handler: %s\n", err)
			w.Write([]byte(err.Error()))
		}
	}
}

func main() {
	fmt.Println("Starting server..")
	pr := server.NewWithConfigFile("config.yaml", false)
	if pr == nil {
		panic("Could not set up proxy")
	}
	httpserver.StartServer(handler(pr), ":3000")
}
