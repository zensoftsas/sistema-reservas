package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"version-1-0/internal/usecase/auth"
)

// AuthHandler handles HTTP requests related to authentication operations
type AuthHandler struct {
	loginUC *auth.LoginUseCase
}

// NewAuthHandler creates a new instance of AuthHandler
func NewAuthHandler(loginUC *auth.LoginUseCase) *AuthHandler {
	return &AuthHandler{
		loginUC: loginUC,
	}
}

// Login godoc
// @Summary      Login de usuario
// @Description  Autenticar usuario y obtener token JWT
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        credentials  body      auth.LoginRequest  true  "Email y password"
// @Success      200  {object}  auth.LoginResponse
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Router       /api/auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// Verify HTTP method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode JSON request body
	var req auth.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Execute use case
	ctx := context.Background()
	response, err := h.loginUC.Execute(ctx, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
