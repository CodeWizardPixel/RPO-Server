package handlers

import (
	"fmt"
	"net/http"
	"strings"
)

func extractBearerToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("authorization header is required")
	}

	const bearerScheme = "Bearer "
	if !strings.HasPrefix(authHeader, bearerScheme) {
		return "", fmt.Errorf("invalid authorization header format")
	}

	tokenString := strings.TrimPrefix(authHeader, bearerScheme)
	if tokenString == "" {
		return "", fmt.Errorf("token is required")
	}

	return tokenString, nil
}

func writeServiceError(w http.ResponseWriter, err error) {
	errText := err.Error()
	switch {
	case strings.Contains(errText, "forbidden"):
		http.Error(w, errText, http.StatusForbidden)
	case strings.Contains(errText, "token"):
		http.Error(w, errText, http.StatusUnauthorized)
	default:
		http.Error(w, errText, http.StatusBadRequest)
	}
}
