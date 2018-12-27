package recover

import "net/http"

// Middleware Recover
func Middleware(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Respond internal server error
			}
		}()
		h.ServeHTTP(w, r)
	})
}
