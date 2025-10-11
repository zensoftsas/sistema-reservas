package middleware

import (
	"log"
	"net/http"
	"runtime/debug"
)

// RecoveryMiddleware recovers from panics and returns a 500 Internal Server Error
// This prevents the entire server from crashing when a handler panics
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Defer a function to recover from panic
		defer func() {
			if err := recover(); err != nil {
				// Log the panic error
				log.Printf("PANIC RECOVERED: %v", err)

				// Log the full stack trace for debugging
				log.Printf("Stack trace:\n%s", debug.Stack())

				// Return 500 Internal Server Error to the client
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		// Execute the next handler
		next.ServeHTTP(w, r)
	})
}
