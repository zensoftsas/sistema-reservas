package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// JWT secret key - TEMPORARY: should be moved to configuration in production
// Must match the secret used in login.go
const jwtSecret = "tu-super-secret-key-cambiar-en-produccion"

// ContextKey is a custom type for context keys to avoid collisions
type ContextKey string

// UserIDKey is the context key for storing the authenticated user ID
const UserIDKey ContextKey = "user_id"

// RoleKey is the context key for storing the user role
const RoleKey ContextKey = "user_role"

// AuthMiddleware validates JWT tokens and adds user information to the request context
// Requires a valid Bearer token in the Authorization header
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Read Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		// Verify Bearer prefix
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
			return
		}

		// Extract token (remove "Bearer " prefix)
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse and validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Verify signing method is HMAC
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		// Extract user_id from claims
		userID, ok := claims["user_id"].(string)
		if !ok {
			http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
			return
		}

		// Extract role from claims
		userRole, ok := claims["role"].(string)
		if !ok {
			http.Error(w, "Invalid role in token", http.StatusUnauthorized)
			return
		}

		// Add user_id and role to context and execute next handler
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		ctx = context.WithValue(ctx, RoleKey, userRole)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireRole returns a middleware that checks if the authenticated user has the required role
// Must be used AFTER AuthMiddleware as it depends on the user being authenticated
func RequireRole(requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get user role from context (set by previous middleware or modify AuthMiddleware)
			role, ok := r.Context().Value(RoleKey).(string)
			if !ok {
				http.Error(w, "Role not found in context", http.StatusForbidden)
				return
			}

			// Check if user has the required role
			if role != requiredRole {
				http.Error(w, "Insufficient permissions", http.StatusForbidden)
				return
			}

			// User has required role, proceed
			next.ServeHTTP(w, r)
		})
	}
}
