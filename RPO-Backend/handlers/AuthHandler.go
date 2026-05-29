package handlers

import (
	"encoding/json"
	"fmt"
	"go-back/service"
	"net/http"
)

type AuthHandler struct {
	authService *service.AuthService
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token   string `json:"token"`
	Message string `json:"message"`
	UserID  int    `json:"user_id"`
}

type ValidateTokenRequest struct {
	Token string `json:"token"`
}

type ValidateTokenResponse struct {
	Valid   bool   `json:"valid"`
	Message string `json:"message"`
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// GetToken
// @Summary Issue JWT
// @Description Authenticates a user with login and password and returns a JWT access token.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Credentials"
// @Success 200 {object} LoginResponse "JWT and user id"
// @Failure 400 {string} string "Invalid request body or missing fields"
// @Failure 401 {string} string "Invalid login or password"
// @Failure 405 {string} string "Method not allowed"
// @Router /auth/login [post]
func (h *AuthHandler) GetToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	// fmt.Printf("Received login request: %+v\n", req)

	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Login == "" || req.Password == "" {
		http.Error(w, "Login and password are required", http.StatusBadRequest)
		return
	}

	user, err := h.authService.AuthenticateUser(req.Login, req.Password)
	if err != nil {
		http.Error(w, "Invalid login or password", http.StatusUnauthorized)
		return
	}

	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(LoginResponse{
	// 	Message: "Authorization successful",
	// 	UserID:  user.ID,
	// })

	token, err := h.authService.GenerateJWT(user)

	response := LoginResponse{
		Token:   token,
		Message: "Authorization successful",
		UserID:  user.ID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// ValidateToken
// @Summary Validate JWT
// @Description Verifies the Bearer token sent in the Authorization header.
// @Tags auth
// @Produce json
// @Param Authorization header string true "Bearer JWT" default(Bearer )
// @Success 200 {object} ValidateTokenResponse "Token is valid"
// @Failure 400 {string} string "Missing or malformed Authorization header"
// @Failure 401 {object} ValidateTokenResponse "Invalid or expired token"
// @Failure 405 {string} string "Method not allowed"
// @Router /auth/validate [post]
func (h *AuthHandler) ValidateToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tokenString, err := extractBearerToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = h.authService.ValidateJWT(tokenString)
	if err != nil {
		fmt.Printf("Token validation error: %v\n", err)
		response := ValidateTokenResponse{
			Valid:   false,
			Message: "Invalid or expired token",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := ValidateTokenResponse{
		Valid:    true,
		Message:  "Token is valid",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
