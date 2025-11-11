package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"version-1-0/internal/delivery/http/middleware"
	"version-1-0/internal/usecase/user"
)

// UserHandler handles HTTP requests related to user operations
type UserHandler struct {
	createUserUC *user.CreateUserUseCase
	getUserUC    *user.GetUserUseCase
	listUsersUC  *user.ListUsersUseCase
	updateUserUC *user.UpdateUserUseCase
	deleteUserUC *user.DeleteUserUseCase
}

// NewUserHandler creates a new instance of UserHandler
func NewUserHandler(
	createUserUC *user.CreateUserUseCase,
	getUserUC *user.GetUserUseCase,
	listUsersUC *user.ListUsersUseCase,
	updateUserUC *user.UpdateUserUseCase,
	deleteUserUC *user.DeleteUserUseCase,
) *UserHandler {
	return &UserHandler{
		createUserUC: createUserUC,
		getUserUC:    getUserUC,
		listUsersUC:  listUsersUC,
		updateUserUC: updateUserUC,
		deleteUserUC: deleteUserUC,
	}
}

// CreateUser godoc
// @Summary      Crear usuario
// @Description  Crear un nuevo usuario (admin, doctor o patient)
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        user  body      dto.CreateUserRequest  true  "Datos del usuario"
// @Success      201  {object}  dto.UserResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Router       /api/users [post]
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

// GetProfile godoc
// @Summary      Obtener perfil del usuario autenticado
// @Description  Retorna información del usuario actual basado en el token JWT
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  dto.UserResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Router       /api/users/me [get]
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

// Update handles the HTTP request for updating a user
// Method: PUT
// Requires: JWT token (admin can update anyone, users can update themselves)
// URL parameter: id in path
// Response: 200 OK with updated user data
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Verify HTTP method is PUT
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user ID from URL path (format: /api/users/{id})
	path := r.URL.Path
	userID := strings.TrimPrefix(path, "/api/users/")
	if userID == "" || userID == path {
		http.Error(w, "User ID is required in path", http.StatusBadRequest)
		return
	}

	// Get authenticated user info from context
	authenticatedUserID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	authenticatedUserRole, ok := r.Context().Value(middleware.RoleKey).(string)
	if !ok {
		http.Error(w, "Role not found in context", http.StatusUnauthorized)
		return
	}

	// Decode request body
	var req user.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Execute use case
	ctx := context.Background()
	response, err := h.updateUserUC.Execute(ctx, userID, authenticatedUserID, authenticatedUserRole, req)
	if err != nil {
		if err.Error() == "insufficient permissions to update this user" {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
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

// Delete handles the HTTP request for deleting (deactivating) a user
// Method: DELETE
// Requires: JWT token with admin role
// URL parameter: id in path
// Response: 204 No Content on success
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Verify HTTP method is DELETE
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user ID from query parameter
	userID := r.URL.Query().Get("id")
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Get authenticated user info from context
	authenticatedUserID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	authenticatedUserRole, ok := r.Context().Value(middleware.RoleKey).(string)
	if !ok {
		http.Error(w, "Role not found in context", http.StatusUnauthorized)
		return
	}

	// Execute use case
	ctx := context.Background()
	err := h.deleteUserUC.Execute(ctx, userID, authenticatedUserID, authenticatedUserRole)
	if err != nil {
		if err.Error() == "only administrators can delete users" {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		if err.Error() == "cannot delete your own account" {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		if err.Error() == "user not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if err.Error() == "user is already inactive" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response (204 No Content)
	w.WriteHeader(http.StatusNoContent)
}
