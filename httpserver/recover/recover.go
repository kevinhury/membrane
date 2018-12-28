package recover

import "net/http"

// Middleware Recover
func Middleware(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Paniced"))
			}
		}()
		h.ServeHTTP(w, r)
	})
}
