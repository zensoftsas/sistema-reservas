package middleware

import (
	"net/http"
	"strings"
)

// CORSMiddleware handles Cross-Origin Resource Sharing (CORS) for the API
// allowedOrigins can be a comma-separated list of origins or "*" for all origins
func CORSMiddleware(allowedOrigins string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// Check if origin is allowed
			if allowedOrigins == "*" {
				w.Header().Set("Access-Control-Allow-Origin", "*")
			} else if origin != "" {
				// Parse allowed origins (comma-separated)
				origins := strings.Split(allowedOrigins, ",")
				for _, allowedOrigin := range origins {
					allowedOrigin = strings.TrimSpace(allowedOrigin)
					if origin == allowedOrigin {
						w.Header().Set("Access-Control-Allow-Origin", origin)
						break
					}
				}
			}

			// Set other CORS headers
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept, Origin")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Max-Age", "3600")
			w.Header().Set("Access-Control-Expose-Headers", "Content-Length")

			// Handle preflight OPTIONS request
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			// Call the next handler
			next.ServeHTTP(w, r)
		})
	}
}
