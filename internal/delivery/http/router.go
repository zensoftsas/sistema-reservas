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
	listUsersHandler := http.HandlerFunc(userHandler.List)
	listUsersWithAuth := middleware.AuthMiddleware(jwtSecret)(listUsersHandler)
	listUsersWithRole := middleware.RequireRole("admin")(listUsersWithAuth)
	mux.Handle("/api/users/list", listUsersWithRole)

	// Update user - requires authentication (admin can update anyone, users can update themselves)
	updateUserHandler := http.HandlerFunc(userHandler.Update)
	updateUserWithAuth := middleware.AuthMiddleware(jwtSecret)(updateUserHandler)
	mux.Handle("/api/users/", updateUserWithAuth)

	// Delete user endpoint will use query param: /api/users/delete?id=xxx
	deleteHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/users/delete" {
			userHandler.Delete(w, r)
		} else {
			http.Error(w, "Not found", http.StatusNotFound)
		}
	})
	deleteWithRole := middleware.RequireRole("admin")(deleteHandler)
	deleteWithAuth := middleware.AuthMiddleware(jwtSecret)(deleteWithRole)
	mux.Handle("/api/users/delete", deleteWithAuth)

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
