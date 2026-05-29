package service

import (
	"fmt"
	"go-back/repository"
)

type KeyService struct {
	keyRepo     *repository.KeyRepository
	authService *AuthService
}

func NewKeyService(keyRepo *repository.KeyRepository, authService *AuthService) *KeyService {
	return &KeyService{
		keyRepo:     keyRepo,
		authService: authService,
	}
}

func (s *KeyService) GetAllKeys() ([]repository.Key, error) {
	keys, err := s.keyRepo.GetAllKeys()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve all keys: %w", err)
	}
	return keys, nil
}

func (s *KeyService) GetKeyByID(id int) (*repository.Key, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid key ID: must be greater than 0")
	}

	key, err := s.keyRepo.GetKeyByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve key with ID %d: %w", id, err)
	}
	return key, nil
}

func (s *KeyService) CreateKey(tokenString, value string) error {
	if err := ensureAdmin(s.authService, tokenString); err != nil {
		return err
	}
	if value == "" {
		return fmt.Errorf("key value cannot be empty")
	}

	if err := s.keyRepo.CreateKey(value); err != nil {
		return fmt.Errorf("failed to create key: %w", err)
	}
	return nil
}

func (s *KeyService) UpdateKey(tokenString string, id int, value string) error {
	if err := ensureAdmin(s.authService, tokenString); err != nil {
		return err
	}
	if id <= 0 {
		return fmt.Errorf("invalid key ID: must be greater than 0")
	}
	if value == "" {
		return fmt.Errorf("key value cannot be empty")
	}

	if err := s.keyRepo.UpdateKey(id, value); err != nil {
		return fmt.Errorf("failed to update key with ID %d: %w", id, err)
	}
	return nil
}

func (s *KeyService) DeleteKey(tokenString string, id int) error {
	if err := ensureAdmin(s.authService, tokenString); err != nil {
		return err
	}
	if id <= 0 {
		return fmt.Errorf("invalid key ID: must be greater than 0")
	}

	if err := s.keyRepo.DeleteKeyByID(id); err != nil {
		return fmt.Errorf("failed to delete key with ID %d: %w", id, err)
	}
	return nil
}

