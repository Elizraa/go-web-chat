package middlewares

import (
	"net/http"
)

// Middleware function to respond with correct content-type headers
func MiddlewareJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
