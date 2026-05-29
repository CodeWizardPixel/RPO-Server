package service

import (
	"fmt"
	"go-back/repository"
)

type TransactionService struct {
	txRepo      *repository.TransactionRepository
	cardRepo    *repository.CardRepository
	authService *AuthService
}

type AuthorizationResponse struct {
	Authorized bool    `json:"authorized"`
	Message    string  `json:"message"`
	CardNumber string  `json:"card_number,omitempty"`
	Operation  string  `json:"operation,omitempty"`
	Balance    float64 `json:"balance,omitempty"`
}

func NewTransactionService(txRepo *repository.TransactionRepository, cardRepo *repository.CardRepository, authService *AuthService) *TransactionService {
	return &TransactionService{
		txRepo:      txRepo,
		cardRepo:    cardRepo,
		authService: authService,
	}
}

func (s *TransactionService) GetAllTransactions() ([]repository.Transaction, error) {
	txs, err := s.txRepo.GetAllTransactions()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve all transactions: %w", err)
	}
	return txs, nil
}

func (s *TransactionService) GetTransactionByID(id int) (*repository.Transaction, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid transaction ID: must be greater than 0")
	}

	tx, err := s.txRepo.GetTransactionByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve transaction with ID %d: %w", id, err)
	}
	return tx, nil
}

func (s *TransactionService) CreateTransaction(tokenString string, amount float64, cardID, terminalID int) error {
	if err := ensureAdmin(s.authService, tokenString); err != nil {
		return err
	}
	if cardID <= 0 {
		return fmt.Errorf("invalid card ID: must be greater than 0")
	}
	if terminalID <= 0 {
		return fmt.Errorf("invalid terminal ID: must be greater than 0")
	}

	if err := s.txRepo.CreateTransaction(amount, cardID, terminalID, "withdraw"); err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}
	return nil
}

func (s *TransactionService) DeleteTransaction(tokenString string, id int) error {
	if err := ensureAdmin(s.authService, tokenString); err != nil {
		return err
	}
	if id <= 0 {
		return fmt.Errorf("invalid transaction ID: must be greater than 0")
	}

	if err := s.txRepo.DeleteTransactionByID(id); err != nil {
		return fmt.Errorf("failed to delete transaction with ID %d: %w", id, err)
	}
	return nil
}

func (s *TransactionService) AuthorizeTransaction(cardNumber, terminalSerialNumber string, amount float64, operation string) *AuthorizationResponse {
	if cardNumber == "" {
		return &AuthorizationResponse{
			Authorized: false,
			Message:    "Card number cannot be empty",
		}
	}

	if terminalSerialNumber == "" {
		return &AuthorizationResponse{
			Authorized: false,
			Message:    "Terminal serial number cannot be empty",
			CardNumber: cardNumber,
			Operation:  operation,
		}
	}

	if amount <= 0 {
		return &AuthorizationResponse{
			Authorized: false,
			Message:    "Transaction amount must be greater than 0",
			CardNumber: cardNumber,
			Operation:  operation,
		}
	}

	if operation != "withdraw" && operation != "deposit" {
		return &AuthorizationResponse{
			Authorized: false,
			Message:    "Operation must be withdraw or deposit",
			CardNumber: cardNumber,
			Operation:  operation,
		}
	}

	card, err := s.txRepo.ProcessCardOperation(cardNumber, terminalSerialNumber, amount, operation)
	if err != nil {
		return &AuthorizationResponse{
			Authorized: false,
			Message:    err.Error(),
			CardNumber: cardNumber,
			Operation:  operation,
		}
	}

	return &AuthorizationResponse{
		Authorized: true,
		Message:    "Transaction processed successfully",
		CardNumber: cardNumber,
		Operation:  operation,
		Balance:    card.Balance,
	}
}
