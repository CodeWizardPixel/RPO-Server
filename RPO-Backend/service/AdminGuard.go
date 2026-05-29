package service

import "fmt"

func ensureAdmin(authService *AuthService, tokenString string) error {
	if tokenString == "" {
		return fmt.Errorf("token is required")
	}
	if authService == nil {
		return fmt.Errorf("auth service is not configured")
	}

	claims, err := authService.ValidateJWT(tokenString)
	if err != nil {
		return fmt.Errorf("invalid token: %w", err)
	}
	if claims.IsAdmin != 1 {
		return fmt.Errorf("forbidden: admin access required")
	}

	return nil
}
