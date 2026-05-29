package service

import (
	"fmt"
	"go-back/repository"
)

type TerminalService struct {
	terminalRepo *repository.TerminalRepository
	authService  *AuthService
}

func NewTerminalService(terminalRepo *repository.TerminalRepository, authService *AuthService) *TerminalService {
	return &TerminalService{
		terminalRepo: terminalRepo,
		authService:  authService,
	}
}

func (s *TerminalService) GetAllTerminals() ([]repository.Terminal, error) {
	terminals, err := s.terminalRepo.GetAllTerminals()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve all terminals: %w", err)
	}

	return terminals, nil
}

func (s *TerminalService) GetTerminalByID(id int) (*repository.Terminal, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid terminal ID: must be greater than 0")
	}

	terminal, err := s.terminalRepo.GetTerminalByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve terminal with ID %d: %w", id, err)
	}

	return terminal, nil
}

func (s *TerminalService) CreateTerminal(tokenString, serialNumber, address, name string) error {
	if err := ensureAdmin(s.authService, tokenString); err != nil {
		return err
	}

	if serialNumber == "" {
		return fmt.Errorf("serial number cannot be empty")
	}
	if address == "" {
		return fmt.Errorf("address cannot be empty")
	}
	if name == "" {
		return fmt.Errorf("terminal name cannot be empty")
	}

	err := s.terminalRepo.CreateTerminal(serialNumber, address, name)
	if err != nil {
		return fmt.Errorf("failed to create terminal: %w", err)
	}

	return nil
}

func (s *TerminalService) UpdateTerminal(tokenString string, id int, serialNumber, address, name string) error {
	if err := ensureAdmin(s.authService, tokenString); err != nil {
		return err
	}

	if id <= 0 {
		return fmt.Errorf("invalid terminal ID: must be greater than 0")
	}
	if serialNumber == "" {
		return fmt.Errorf("serial number cannot be empty")
	}
	if address == "" {
		return fmt.Errorf("address cannot be empty")
	}
	if name == "" {
		return fmt.Errorf("terminal name cannot be empty")
	}

	err := s.terminalRepo.UpdateTerminal(id, serialNumber, address, name)
	if err != nil {
		return fmt.Errorf("failed to update terminal with ID %d: %w", id, err)
	}

	return nil
}

func (s *TerminalService) DeleteTerminal(tokenString string, id int) error {
	if err := ensureAdmin(s.authService, tokenString); err != nil {
		return err
	}

	if id <= 0 {
		return fmt.Errorf("invalid terminal ID: must be greater than 0")
	}

	err := s.terminalRepo.DeleteTerminalByID(id)
	if err != nil {
		return fmt.Errorf("failed to delete terminal with ID %d: %w", id, err)
	}

	return nil
}
