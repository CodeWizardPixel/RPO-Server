package handlers

import (
	"encoding/json"
	"go-back/service"
	"net/http"
	"strconv"
)

type KeyHandler struct {
	keyService *service.KeyService
}

type KeyRequest struct {
	ID    int    `json:"id"`
	Value string `json:"value"`
}

func NewKeyHandler(keyService *service.KeyService) *KeyHandler {
	return &KeyHandler{keyService: keyService}
}

// GetAllKeys
// @Summary List keys
// @Tags keys
// @Produce json
// @Success 200 {array} repository.Key
// @Failure 405 {string} string "Method not allowed"
// @Failure 500 {string} string "Internal server error"
// @Router /keys/all [get]
func (h *KeyHandler) GetAllKeys(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	keys, err := h.keyService.GetAllKeys()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(keys)
}

// GetKeyByID
// @Summary Get key by ID
// @Tags keys
// @Produce json
// @Param id query int true "Key ID"
// @Success 200 {object} repository.Key
// @Failure 400 {string} string "Invalid id or not found"
// @Failure 405 {string} string "Method not allowed"
// @Router /keys/get [get]
func (h *KeyHandler) GetKeyByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idRaw := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idRaw)
	if err != nil {
		http.Error(w, "Invalid key ID", http.StatusBadRequest)
		return
	}

	key, err := h.keyService.GetKeyByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(key)
}

// CreateKey
// @Summary Create key
// @Description Admin only.
// @Tags keys
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body KeyRequest true "value"
// @Success 201 {object} map[string]string
// @Failure 400 {string} string "Validation error"
// @Failure 401 {string} string "Missing or invalid token"
// @Failure 403 {string} string "Not an admin"
// @Failure 405 {string} string "Method not allowed"
// @Router /keys/create [post]
func (h *KeyHandler) CreateKey(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tokenString, err := extractBearerToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var req KeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.keyService.CreateKey(tokenString, req.Value)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Key created"})
}

// UpdateKey
// @Summary Update key
// @Description Admin only.
// @Tags keys
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body KeyRequest true "id, value"
// @Success 200 {object} map[string]string
// @Failure 400 {string} string "Validation error"
// @Failure 401 {string} string "Missing or invalid token"
// @Failure 403 {string} string "Not an admin"
// @Failure 405 {string} string "Method not allowed"
// @Router /keys/update [put]
func (h *KeyHandler) UpdateKey(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tokenString, err := extractBearerToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var req KeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.keyService.UpdateKey(tokenString, req.ID, req.Value)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Key updated"})
}

// DeleteKey
// @Summary Delete key
// @Description Admin only.
// @Tags keys
// @Produce json
// @Security BearerAuth
// @Param id query int true "Key ID"
// @Success 200 {object} map[string]string
// @Failure 400 {string} string "Invalid id"
// @Failure 401 {string} string "Missing or invalid token"
// @Failure 403 {string} string "Not an admin"
// @Failure 405 {string} string "Method not allowed"
// @Router /keys/delete [delete]
func (h *KeyHandler) DeleteKey(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "Invalid key ID", http.StatusBadRequest)
		return
	}

	err = h.keyService.DeleteKey(tokenString, id)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Key deleted"})
}
