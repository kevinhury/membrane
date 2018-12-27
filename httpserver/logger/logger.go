package logger

import (
	"log"
	"net/http"
	"time"
)

// Middleware logger
func Middleware(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		h.ServeHTTP(w, r)
		endTime := time.Since(startTime)

		log.Printf("%s %s %v", r.URL, r.Method, endTime)
	})
}
