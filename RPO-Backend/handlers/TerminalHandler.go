package handlers

import (
	"encoding/json"
	"go-back/service"
	"net/http"
	"strconv"
)

type TerminalHandler struct {
	terminalService *service.TerminalService
}

type TerminalRequest struct {
	ID           int    `json:"id"`
	SerialNumber string `json:"serial_number"`
	Address      string `json:"address"`
	Name         string `json:"name"`
}

func NewTerminalHandler(terminalService *service.TerminalService) *TerminalHandler {
	return &TerminalHandler{
		terminalService: terminalService,
	}
}

// GetAllTerminals
// @Summary List terminals
// @Description Returns all payment terminals.
// @Tags terminals
// @Produce json
// @Success 200 {array} repository.Terminal
// @Failure 405 {string} string "Method not allowed"
// @Failure 500 {string} string "Internal server error"
// @Router /terminals/all [get]
func (h *TerminalHandler) GetAllTerminals(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	terminals, err := h.terminalService.GetAllTerminals()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(terminals)
}

// GetTerminalByID
// @Summary Get terminal by ID
// @Description Returns a single terminal by numeric id.
// @Tags terminals
// @Produce json
// @Param id query int true "Terminal ID"
// @Success 200 {object} repository.Terminal
// @Failure 400 {string} string "Invalid id or not found"
// @Failure 405 {string} string "Method not allowed"
// @Router /terminals/get [get]
func (h *TerminalHandler) GetTerminalByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idRaw := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idRaw)
	if err != nil {
		http.Error(w, "Invalid terminal ID", http.StatusBadRequest)
		return
	}

	terminal, err := h.terminalService.GetTerminalByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(terminal)
}

// CreateTerminal
// @Summary Create terminal
// @Description Creates a terminal. Admin JWT required (is_admin=1).
// @Tags terminals
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body TerminalRequest true "Terminal fields"
// @Success 201 {object} map[string]string
// @Failure 400 {string} string "Validation error"
// @Failure 401 {string} string "Missing or invalid token"
// @Failure 403 {string} string "Not an admin"
// @Failure 405 {string} string "Method not allowed"
// @Router /terminals/create [post]
func (h *TerminalHandler) CreateTerminal(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tokenString, err := extractBearerToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var req TerminalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.terminalService.CreateTerminal(tokenString, req.SerialNumber, req.Address, req.Name)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Terminal created"})
}

// UpdateTerminal
// @Summary Update terminal
// @Description Updates a terminal by id. Admin JWT required.
// @Tags terminals
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body TerminalRequest true "Must include id and fields"
// @Success 200 {object} map[string]string
// @Failure 400 {string} string "Validation error"
// @Failure 401 {string} string "Missing or invalid token"
// @Failure 403 {string} string "Not an admin"
// @Failure 405 {string} string "Method not allowed"
// @Router /terminals/update [put]
func (h *TerminalHandler) UpdateTerminal(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tokenString, err := extractBearerToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var req TerminalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.terminalService.UpdateTerminal(tokenString, req.ID, req.SerialNumber, req.Address, req.Name)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Terminal updated"})
}

// DeleteTerminal
// @Summary Delete terminal
// @Description Deletes a terminal by id. Admin JWT required.
// @Tags terminals
// @Produce json
// @Security BearerAuth
// @Param id query int true "Terminal ID"
// @Success 200 {object} map[string]string
// @Failure 400 {string} string "Invalid id"
// @Failure 401 {string} string "Missing or invalid token"
// @Failure 403 {string} string "Not an admin"
// @Failure 405 {string} string "Method not allowed"
// @Router /terminals/delete [delete]
func (h *TerminalHandler) DeleteTerminal(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "Invalid terminal ID", http.StatusBadRequest)
		return
	}

	err = h.terminalService.DeleteTerminal(tokenString, id)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Terminal deleted"})
}