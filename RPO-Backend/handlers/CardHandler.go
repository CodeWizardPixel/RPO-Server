package handlers

import (
	"encoding/json"
	"go-back/service"
	"net/http"
	"strconv"
)

type CardHandler struct {
	cardService *service.CardService
}

type CardRequest struct {
	ID         int     `json:"id"`
	CardNumber string  `json:"card_number"`
	Balance    float64 `json:"balance"`
	IsBlocked  int     `json:"is_blocked"`
	OwnerName  string  `json:"owner_name"`
	KeyID      *int    `json:"key_id"`
}

type CardBalanceRequest struct {
	ID      int     `json:"id"`
	Balance float64 `json:"balance"`
}

func NewCardHandler(cardService *service.CardService) *CardHandler {
	return &CardHandler{cardService: cardService}
}

// GetAllCards
// @Summary List cards
// @Tags cards
// @Produce json
// @Success 200 {array} repository.Card
// @Failure 405 {string} string "Method not allowed"
// @Failure 500 {string} string "Internal server error"
// @Router /cards/all [get]
func (h *CardHandler) GetAllCards(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cards, err := h.cardService.GetAllCards()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cards)
}

// GetCardByID
// @Summary Get card by ID
// @Tags cards
// @Produce json
// @Param id query int true "Card ID"
// @Success 200 {object} repository.Card
// @Failure 400 {string} string "Invalid id or not found"
// @Failure 405 {string} string "Method not allowed"
// @Router /cards/get [get]
func (h *CardHandler) GetCardByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idRaw := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idRaw)
	if err != nil {
		http.Error(w, "Invalid card ID", http.StatusBadRequest)
		return
	}

	card, err := h.cardService.GetCardByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(card)
}

// CreateCard
// @Summary Create card
// @Description Admin only.
// @Tags cards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CardRequest true "card_number, balance, is_blocked, owner_name, optional key_id"
// @Success 201 {object} map[string]string
// @Failure 400 {string} string "Validation error"
// @Failure 401 {string} string "Missing or invalid token"
// @Failure 403 {string} string "Not an admin"
// @Failure 405 {string} string "Method not allowed"
// @Router /cards/create [post]
func (h *CardHandler) CreateCard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tokenString, err := extractBearerToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var req CardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.cardService.CreateCard(tokenString, req.CardNumber, req.Balance, req.IsBlocked, req.OwnerName, req.KeyID)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Card created"})
}

// UpdateCard
// @Summary Update card
// @Description Admin only. Updates balance, block flag, owner, and key link.
// @Tags cards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CardRequest true "id and fields to update"
// @Success 200 {object} map[string]string
// @Failure 400 {string} string "Validation error"
// @Failure 401 {string} string "Missing or invalid token"
// @Failure 403 {string} string "Not an admin"
// @Failure 405 {string} string "Method not allowed"
// @Router /cards/update [put]
func (h *CardHandler) UpdateCard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tokenString, err := extractBearerToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var req CardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.cardService.UpdateCard(tokenString, req.ID, req.Balance, req.IsBlocked, req.OwnerName, req.KeyID)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Card updated"})
}

// UpdateCardBalance
// @Summary Update card balance only
// @Description Admin only. Sets balance to the given value.
// @Tags cards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CardBalanceRequest true "id and balance"
// @Success 200 {object} map[string]string
// @Failure 400 {string} string "Validation error"
// @Failure 401 {string} string "Missing or invalid token"
// @Failure 403 {string} string "Not an admin"
// @Failure 405 {string} string "Method not allowed"
// @Router /cards/balance [put]
func (h *CardHandler) UpdateCardBalance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tokenString, err := extractBearerToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var req CardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.cardService.UpdateCardBalance(tokenString, req.ID, req.Balance)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Card balance updated"})
}

// DeleteCard
// @Summary Delete card
// @Description Admin only.
// @Tags cards
// @Produce json
// @Security BearerAuth
// @Param id query int true "Card ID"
// @Success 200 {object} map[string]string
// @Failure 400 {string} string "Invalid id"
// @Failure 401 {string} string "Missing or invalid token"
// @Failure 403 {string} string "Not an admin"
// @Failure 405 {string} string "Method not allowed"
// @Router /cards/delete [delete]
func (h *CardHandler) DeleteCard(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "Invalid card ID", http.StatusBadRequest)
		return
	}

	err = h.cardService.DeleteCard(tokenString, id)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Card deleted"})
}
