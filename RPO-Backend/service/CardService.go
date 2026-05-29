package service

import (
	"fmt"
	"go-back/repository"
)

type CardService struct {
	cardRepo    *repository.CardRepository
	authService *AuthService
}

func NewCardService(cardRepo *repository.CardRepository, authService *AuthService) *CardService {
	return &CardService{
		cardRepo:    cardRepo,
		authService: authService,
	}
}

func (s *CardService) GetAllCards() ([]repository.Card, error) {
	cards, err := s.cardRepo.GetAllCards()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve all cards: %w", err)
	}
	return cards, nil
}

func (s *CardService) GetCardByID(id int) (*repository.Card, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid card ID: must be greater than 0")
	}

	card, err := s.cardRepo.GetCardByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve card with ID %d: %w", id, err)
	}
	return card, nil
}

func (s *CardService) CreateCard(tokenString, cardNumber string, balance float64, isBlocked int, ownerName string, keyID *int) error {
	if err := ensureAdmin(s.authService, tokenString); err != nil {
		return err
	}
	if cardNumber == "" {
		return fmt.Errorf("card number cannot be empty")
	}
	if ownerName == "" {
		return fmt.Errorf("owner name cannot be empty")
	}
	if isBlocked != 0 && isBlocked != 1 {
		return fmt.Errorf("is_blocked must be 0 or 1")
	}

	if err := s.cardRepo.CreateCard(cardNumber, balance, isBlocked, ownerName, keyID); err != nil {
		return fmt.Errorf("failed to create card: %w", err)
	}
	return nil
}

func (s *CardService) UpdateCard(tokenString string, id int, balance float64, isBlocked int, ownerName string, keyID *int) error {
	if err := ensureAdmin(s.authService, tokenString); err != nil {
		return err
	}
	if id <= 0 {
		return fmt.Errorf("invalid card ID: must be greater than 0")
	}
	if ownerName == "" {
		return fmt.Errorf("owner name cannot be empty")
	}
	if isBlocked != 0 && isBlocked != 1 {
		return fmt.Errorf("is_blocked must be 0 or 1")
	}

	if err := s.cardRepo.UpdateCard(id, balance, isBlocked, ownerName, keyID); err != nil {
		return fmt.Errorf("failed to update card with ID %d: %w", id, err)
	}
	return nil
}

func (s *CardService) UpdateCardBalance(tokenString string, id int, balance float64) error {
	if err := ensureAdmin(s.authService, tokenString); err != nil {
		return err
	}
	if id <= 0 {
		return fmt.Errorf("invalid card ID: must be greater than 0")
	}

	if err := s.cardRepo.UpdateCardBalance(id, balance); err != nil {
		return fmt.Errorf("failed to update card balance for ID %d: %w", id, err)
	}
	return nil
}

func (s *CardService) DeleteCard(tokenString string, id int) error {
	if err := ensureAdmin(s.authService, tokenString); err != nil {
		return err
	}
	if id <= 0 {
		return fmt.Errorf("invalid card ID: must be greater than 0")
	}

	if err := s.cardRepo.DeleteCardByID(id); err != nil {
		return fmt.Errorf("failed to delete card with ID %d: %w", id, err)
	}
	return nil
}

