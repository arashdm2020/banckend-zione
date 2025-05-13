package services_test

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"zionechainapi/configs"
	"zionechainapi/internal/models"
	"zionechainapi/internal/services"
)

func TestNewAuthService(t *testing.T) {
	// Create a test config
	config := &configs.Config{
		JWT: configs.JWTConfig{
			Secret:             "test-secret",
			AccessTokenExpiry:  time.Minute * 15,
			RefreshTokenExpiry: time.Hour * 24 * 7,
		},
	}

	// Create a new auth service
	authService := services.NewAuthService(config)

	// Assert that the auth service is not nil
	assert.NotNil(t, authService)
}

func TestValidateToken(t *testing.T) {
	// This is more of an integration test and would need a mock database
	// For now, we'll just test token validation with a mocked token

	// Create a test config
	config := &configs.Config{
		JWT: configs.JWTConfig{
			Secret:             "test-secret",
			AccessTokenExpiry:  time.Minute * 15,
			RefreshTokenExpiry: time.Hour * 24 * 7,
		},
	}

	// Create a new auth service
	authService := services.NewAuthService(config)

	// Create a test claims
	claims := &services.Claims{
		UserID: 1,
		Role:   "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   "1",
		},
	}

	// Create a token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.JWT.Secret))
	assert.NoError(t, err)

	// Validate the token
	validatedClaims, err := authService.ValidateToken(tokenString)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), validatedClaims.UserID)
	assert.Equal(t, "admin", validatedClaims.Role)
}