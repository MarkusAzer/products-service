package middleware

import (
	"net/http"
)

// ValidateHeaderType validate that header is application/json
func ValidateHeaderType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "" {
			//TODO: check this package github.com/golang/gddo/tree/master/httputil/header >>> ParseValueAndParams
			value := r.Header.Get("Content-Type")
			if value != "application/json" {
				msg := "Content-Type header is not application/json"
				http.Error(w, msg, http.StatusUnsupportedMediaType)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
