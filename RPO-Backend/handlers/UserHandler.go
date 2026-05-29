package handlers

import (
	"encoding/json"
	"go-back/service"
	"net/http"
	"strconv"
)

type UserHandler struct {
	userService *service.UserService
}

type UserRequest struct {
	ID           int    `json:"id"`
	Login        string `json:"login"`
	Name         string `json:"name"`
	PasswordHash string `json:"password_hash"`
	IsAdmin      int    `json:"is_admin"`
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GetAllUsers
// @Summary List users
// @Description Returns all user records (includes password hash fields as stored).
// @Tags users
// @Produce json
// @Success 200 {array} repository.User
// @Failure 405 {string} string "Method not allowed"
// @Failure 500 {string} string "Internal server error"
// @Router /users/all [get]
func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	users, err := h.userService.GetAllUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

// GetUserByID
// @Summary Get user by ID
// @Tags users
// @Produce json
// @Param id query int true "User ID"
// @Success 200 {object} repository.User
// @Failure 400 {string} string "Invalid id or not found"
// @Failure 405 {string} string "Method not allowed"
// @Router /users/get [get]
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idRaw := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idRaw)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetUserByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// CreateUser
// @Summary Create user
// @Description Admin only. Expects pre-hashed password in password_hash.
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body UserRequest true "login, name, password_hash, is_admin"
// @Success 201 {object} map[string]string
// @Failure 400 {string} string "Validation error"
// @Failure 401 {string} string "Missing or invalid token"
// @Failure 403 {string} string "Not an admin"
// @Failure 405 {string} string "Method not allowed"
// @Router /users/create [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tokenString, err := extractBearerToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var req UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.userService.CreateUser(tokenString, req.Login, req.Name, req.PasswordHash, req.IsAdmin)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created"})
}

// UpdateUser
// @Summary Update user
// @Description Admin only.
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body UserRequest true "id, name, password_hash, is_admin"
// @Success 200 {object} map[string]string
// @Failure 400 {string} string "Validation error"
// @Failure 401 {string} string "Missing or invalid token"
// @Failure 403 {string} string "Not an admin"
// @Failure 405 {string} string "Method not allowed"
// @Router /users/update [put]
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tokenString, err := extractBearerToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var req UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.userService.UpdateUser(tokenString, req.ID, req.Name, req.PasswordHash, req.IsAdmin)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User updated"})
}

// DeleteUser
// @Summary Delete user
// @Description Admin only.
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param id query int true "User ID"
// @Success 200 {object} map[string]string
// @Failure 400 {string} string "Invalid id"
// @Failure 401 {string} string "Missing or invalid token"
// @Failure 403 {string} string "Not an admin"
// @Failure 405 {string} string "Method not allowed"
// @Router /users/delete [delete]
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tokenString, err := extractBearerToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	idRaw := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idRaw)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	err = h.userService.DeleteUser(tokenString, id)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User deleted"})
}
