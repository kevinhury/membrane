package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/fsnotify/fsnotify"
	"github.com/kevinhury/membrane/reverseproxy"

	"github.com/kevinhury/membrane/httpserver"
)

func handler(reg *reverseproxy.Registry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := reg.Serve(w, r)
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

	fileName := "config.yaml"
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic("Could not read configfile")
	}

	reg := reverseproxy.NewWithConfig(content)
	if reg == nil {
		panic("Could not set up registry")
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("ERROR", err)
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			// watch for events
			case event := <-watcher.Events:
				if event.Op == fsnotify.Write {
					content, err := ioutil.ReadFile(event.Name)
					if err != nil {
						fmt.Println("ERROR Reading watched file", event.Name)
					}
					reg.SetConfig(content)
				}

				// watch for errors
			case err := <-watcher.Errors:
				fmt.Println("ERROR", err)
			}
		}
	}()

	if err := watcher.Add(fileName); err != nil {
		panic(err)
	}

	httpserver.StartServer(handler(reg), ":3000")
}
