package api

import (
	"crypto/subtle"
	"net/http"
	"strings"
)

func RequireBearerToken(expected string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if expected == "" {
			http.Error(w, "API not configured", http.StatusServiceUnavailable)
			return
		}

		const prefix = "Bearer "
		header := r.Header.Get("Authorization")
		if !strings.HasPrefix(header, prefix) {
			// Bearer not set
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		provided := header[len(prefix):]
		if subtle.ConstantTimeCompare([]byte(provided), []byte(expected)) != 1 {
			// Bearer Incorrect
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
