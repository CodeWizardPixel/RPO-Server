package service

import (
	"fmt"
	"go-back/repository"
)

type UserService struct {
	userRepo    *repository.UserRepository
	authService *AuthService
}

func NewUserService(userRepo *repository.UserRepository, authService *AuthService) *UserService {
	return &UserService{
		userRepo:    userRepo,
		authService: authService,
	}
}

func (s *UserService) GetAllUsers() ([]repository.User, error) {
	users, err := s.userRepo.GetAllUsers()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve all users: %w", err)
	}
	return users, nil
}

func (s *UserService) GetUserByID(id int) (*repository.User, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid user ID: must be greater than 0")
	}

	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user with ID %d: %w", id, err)
	}
	return user, nil
}

func (s *UserService) CreateUser(tokenString, login, name, passwordHash string, isAdmin int) error {
	if err := ensureAdmin(s.authService, tokenString); err != nil {
		return err
	}
	if login == "" {
		return fmt.Errorf("login cannot be empty")
	}
	if name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if passwordHash == "" {
		return fmt.Errorf("password hash cannot be empty")
	}
	if isAdmin != 0 && isAdmin != 1 {
		return fmt.Errorf("is_admin must be 0 or 1")
	}

	if err := s.userRepo.CreateUser(login, name, passwordHash, isAdmin); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (s *UserService) UpdateUser(tokenString string, id int, name, passwordHash string, isAdmin int) error {
	if err := ensureAdmin(s.authService, tokenString); err != nil {
		return err
	}
	if id <= 0 {
		return fmt.Errorf("invalid user ID: must be greater than 0")
	}
	if name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if passwordHash == "" {
		return fmt.Errorf("password hash cannot be empty")
	}
	if isAdmin != 0 && isAdmin != 1 {
		return fmt.Errorf("is_admin must be 0 or 1")
	}

	if err := s.userRepo.UpdateUser(id, name, passwordHash, isAdmin); err != nil {
		return fmt.Errorf("failed to update user with ID %d: %w", id, err)
	}
	return nil
}

func (s *UserService) DeleteUser(tokenString string, id int) error {
	if err := ensureAdmin(s.authService, tokenString); err != nil {
		return err
	}
	if id <= 0 {
		return fmt.Errorf("invalid user ID: must be greater than 0")
	}

	if err := s.userRepo.DeleteUserByID(id); err != nil {
		return fmt.Errorf("failed to delete user with ID %d: %w", id, err)
	}
	return nil
}

