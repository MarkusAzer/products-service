package middleware

import (
	"net/http"
)

// SetResHeaderType Set response header type as application/json
func SetResHeaderType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		//w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
