package http

import (
	"net/http"

	"version-1-0/internal/delivery/http/handler"
	"version-1-0/internal/delivery/http/middleware"
)

// SetupRouter configures and returns the HTTP router with all application routes
func SetupRouter(userHandler *handler.UserHandler, authHandler *handler.AuthHandler, jwtSecret string) http.Handler {
	// Create a new HTTP multiplexer
	mux := http.NewServeMux()

	// Register user routes
	mux.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			userHandler.Create(w, r)
		} else if r.Method == http.MethodGet {
			userHandler.GetByID(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Register authentication routes
	mux.HandleFunc("/api/auth/login", authHandler.Login)

	// Register protected user routes
	protectedUserRoutes := http.HandlerFunc(userHandler.GetMe)
	protectedUserRoutesWithAuth := middleware.AuthMiddleware(jwtSecret)(protectedUserRoutes)
	mux.Handle("/api/users/me", protectedUserRoutesWithAuth)

	// Register admin-only routes
	// List users - requires admin role
	// Apply middlewares in correct order: Auth first, then Role check
	listUsersHandler := http.HandlerFunc(userHandler.List)
	listUsersWithRole := middleware.RequireRole("admin")(listUsersHandler)
	listUsersWithAuth := middleware.AuthMiddleware(jwtSecret)(listUsersWithRole)
	mux.Handle("/api/users/list", listUsersWithAuth)

	// Register health check route
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Sistema de Reservas - API Running"))
	})

	// Apply middlewares in order: Recovery -> Logging -> Handlers
	// Recovery middleware wraps everything to catch panics
	withRecovery := middleware.RecoveryMiddleware(mux)

	// Logging middleware wraps the recovery middleware to log all requests
	withLogging := middleware.LoggingMiddleware(withRecovery)

	return withLogging
}
