package handlers

import (
	"encoding/json"
	"go-back/service"
	"net/http"
	"strconv"
)

type TransactionHandler struct {
	txService *service.TransactionService
}

type TransactionRequest struct {
	ID         int     `json:"id"`
	Amount     float64 `json:"amount"`
	CardID     int     `json:"card_id"`
	TerminalID int     `json:"terminal_id"`
}

type AuthorizationRequest struct {
	CardNumber           string  `json:"card_number"`
	TerminalSerialNumber string  `json:"terminal_serial_number"`
	Amount               float64 `json:"amount"`
	Operation            string  `json:"operation"`
}

func NewTransactionHandler(txService *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{txService: txService}
}

// GetAllTransactions
// @Summary List transactions
// @Tags transactions
// @Produce json
// @Success 200 {array} repository.Transaction
// @Failure 405 {string} string "Method not allowed"
// @Failure 500 {string} string "Internal server error"
// @Router /transactions/all [get]
func (h *TransactionHandler) GetAllTransactions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	txs, err := h.txService.GetAllTransactions()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(txs)
}

// GetTransactionByID
// @Summary Get transaction by ID
// @Tags transactions
// @Produce json
// @Param id query int true "Transaction ID"
// @Success 200 {object} repository.Transaction
// @Failure 400 {string} string "Invalid id or not found"
// @Failure 405 {string} string "Method not allowed"
// @Router /transactions/get [get]
func (h *TransactionHandler) GetTransactionByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idRaw := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idRaw)
	if err != nil {
		http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
		return
	}

	tx, err := h.txService.GetTransactionByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tx)
}

// CreateTransaction
// @Summary Create transaction
// @Description Admin only. Records amount for a card at a terminal.
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body TransactionRequest true "amount, card_id, terminal_id"
// @Success 201 {object} map[string]string
// @Failure 400 {string} string "Validation error"
// @Failure 401 {string} string "Missing or invalid token"
// @Failure 403 {string} string "Not an admin"
// @Failure 405 {string} string "Method not allowed"
// @Router /transactions/create [post]
func (h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tokenString, err := extractBearerToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var req TransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.txService.CreateTransaction(tokenString, req.Amount, req.CardID, req.TerminalID)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Transaction created"})
}

// DeleteTransaction
// @Summary Delete transaction
// @Description Admin only.
// @Tags transactions
// @Produce json
// @Security BearerAuth
// @Param id query int true "Transaction ID"
// @Success 200 {object} map[string]string
// @Failure 400 {string} string "Invalid id"
// @Failure 401 {string} string "Missing or invalid token"
// @Failure 403 {string} string "Not an admin"
// @Failure 405 {string} string "Method not allowed"
// @Router /transactions/delete [delete]
func (h *TransactionHandler) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
		return
	}

	err = h.txService.DeleteTransaction(tokenString, id)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Transaction deleted"})
}

// AuthorizeTransaction
// @Summary Process card transaction
// @Description Processes a withdraw or deposit operation for a card at a terminal.
// @Tags transactions
// @Accept json
// @Produce json
// @Param request body AuthorizationRequest true "card_number, terminal_serial_number, amount, operation (withdraw or deposit)"
// @Success 200 {object} service.AuthorizationResponse
// @Failure 400 {string} string "Invalid request"
// @Failure 405 {string} string "Method not allowed"
// @Router /transactions/authorize [post]
func (h *TransactionHandler) AuthorizeTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AuthorizationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result := h.txService.AuthorizeTransaction(req.CardNumber, req.TerminalSerialNumber, req.Amount, req.Operation)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
