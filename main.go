package main

import (
	"fmt"
	"net/http"

	"github.com/kevinhury/membrane/proxy/server"

	"github.com/kevinhury/membrane/httpserver"
)

func rootHandlerFunc(w http.ResponseWriter, r *http.Request) {
	pr := server.NewWithConfigFile("config.yaml", false)
	err := pr.Serve(w, r)
	if err != nil {
		fmt.Printf("Received error from root handler: %s\n", err)
		w.Write([]byte(err.Error()))
	}
}

func main() {
	fmt.Println("Starting server..")
	httpserver.StartServer(rootHandlerFunc, ":8082")
}
