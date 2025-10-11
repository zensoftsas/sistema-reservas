package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"version-1-0/internal/delivery/http/middleware"
	"version-1-0/internal/usecase/user"
)

// UserHandler handles HTTP requests related to user operations
type UserHandler struct {
	createUserUC *user.CreateUserUseCase
	getUserUC    *user.GetUserUseCase
	listUsersUC  *user.ListUsersUseCase
}

// NewUserHandler creates a new instance of UserHandler
func NewUserHandler(createUserUC *user.CreateUserUseCase, getUserUC *user.GetUserUseCase, listUsersUC *user.ListUsersUseCase) *UserHandler {
	return &UserHandler{
		createUserUC: createUserUC,
		getUserUC:    getUserUC,
		listUsersUC:  listUsersUC,
	}
}

// Create handles the HTTP request for creating a new user
// Method: POST
// Request body: JSON with user creation data
// Response: 201 Created with user data (without password) or error
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Verify HTTP method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode JSON request body
	var req user.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Execute use case
	ctx := context.Background()
	response, err := h.createUserUC.Execute(ctx, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GetByID handles the HTTP request for retrieving a user by ID
// Method: GET
// URL parameter: id in query string (?id=xxx)
// Response: 200 OK with user data (without password) or error
func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// Verify HTTP method is GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user ID from query parameter
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Execute use case
	ctx := context.Background()
	response, err := h.getUserUC.Execute(ctx, id)
	if err != nil {
		// Check if error is "user not found"
		if err.Error() == "user not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetMe handles the HTTP request for getting the authenticated user's profile
// Method: GET
// Requires: JWT token in Authorization header
// Response: 200 OK with user data or 401 Unauthorized
func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	// Verify HTTP method is GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	// Execute use case
	ctx := context.Background()
	response, err := h.getUserUC.Execute(ctx, userID)
	if err != nil {
		if err.Error() == "user not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// List handles the HTTP request for listing all users with pagination
// Method: GET
// Requires: JWT token with admin role
// Query params: limit (optional, default 20), offset (optional, default 0)
// Response: 200 OK with paginated users list
func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	// Verify HTTP method is GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse query parameters for pagination
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	var limit, offset int

	// Parse limit (default 20)
	if limitStr != "" {
		fmt.Sscanf(limitStr, "%d", &limit)
	}

	// Parse offset (default 0)
	if offsetStr != "" {
		fmt.Sscanf(offsetStr, "%d", &offset)
	}

	// Create request
	req := user.ListUsersRequest{
		Limit:  limit,
		Offset: offset,
	}

	// Execute use case
	ctx := context.Background()
	response, err := h.listUsersUC.Execute(ctx, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
