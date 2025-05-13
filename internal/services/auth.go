package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"zionechainapi/configs"
	"zionechainapi/internal/database"
	"zionechainapi/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthService handles authentication and authorization
type AuthService struct {
	config *configs.Config
}

// NewAuthService creates a new auth service
func NewAuthService(config *configs.Config) *AuthService {
	return &AuthService{
		config: config,
	}
}

// LoginRequest represents the login request
type LoginRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest represents the register request
type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

// TokenResponse represents the token response
type TokenResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	User         UserResponse `json:"user"`
}

// UserResponse represents the user response
type UserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	Role  string `json:"role"`
}

// Claims represents the JWT claims
type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// RefreshTokenClaims represents the refresh token claims
type RefreshTokenClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

// Login authenticates a user and returns a JWT token
func (s *AuthService) Login(req LoginRequest) (*TokenResponse, error) {
	var user models.User
	
	// Find user by phone
	if err := database.DB.Preload("Role").Where("phone = ?", req.Phone).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid phone or password")
		}
		return nil, err
	}

	// Check password
	if !user.CheckPassword(req.Password) {
		return nil, errors.New("invalid phone or password")
	}

	// Generate tokens
	accessToken, expiresAt, err := s.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	// Return token response
	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		User: UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Phone: user.Phone,
			Role:  user.Role.Name,
		},
	}, nil
}

// Register creates a new user account
func (s *AuthService) Register(req RegisterRequest) (*TokenResponse, error) {
	// Check if user already exists
	var count int64
	if err := database.DB.Model(&models.User{}).Where("email = ? OR phone = ?", req.Email, req.Phone).Count(&count).Error; err != nil {
		return nil, err
	}

	if count > 0 {
		return nil, errors.New("user with this email or phone already exists")
	}

	// Create user
	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: req.Password,
		RoleID:   models.RoleUser, // Default to user role
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return nil, err
	}

	// Load role
	if err := database.DB.Preload("Role").First(&user, user.ID).Error; err != nil {
		return nil, err
	}

	// Generate tokens
	accessToken, expiresAt, err := s.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	// Return token response
	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		User: UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Phone: user.Phone,
			Role:  user.Role.Name,
		},
	}, nil
}

// RefreshToken refreshes the access token using a refresh token
func (s *AuthService) RefreshToken(refreshToken string) (*TokenResponse, error) {
	// Parse refresh token
	token, err := jwt.ParseWithClaims(refreshToken, &RefreshTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	// Validate token
	claims, ok := token.Claims.(*RefreshTokenClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	// Get user
	var user models.User
	if err := database.DB.Preload("Role").First(&user, claims.UserID).Error; err != nil {
		return nil, err
	}

	// Generate new tokens
	accessToken, expiresAt, err := s.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	// Return token response
	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    expiresAt,
		User: UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Phone: user.Phone,
			Role:  user.Role.Name,
		},
	}, nil
}

// ValidateToken validates a JWT token and returns the claims
func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// GetUserByID gets a user by ID
func (s *AuthService) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := database.DB.Preload("Role").First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// ChangePassword changes a user's password
func (s *AuthService) ChangePassword(userID uint, currentPassword, newPassword string) error {
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return err
	}

	// Verify current password
	if !user.CheckPassword(currentPassword) {
		return errors.New("current password is incorrect")
	}

	// Update password
	user.Password = newPassword
	return database.DB.Save(&user).Error
}

// Helper functions
func (s *AuthService) generateAccessToken(user models.User) (string, time.Time, error) {
	expiresAt := time.Now().Add(s.config.JWT.AccessTokenExpiry)

	claims := &Claims{
		UserID: user.ID,
		Role:   user.Role.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   fmt.Sprintf("%d", user.ID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.config.JWT.Secret))

	return tokenString, expiresAt, err
}

func (s *AuthService) generateRefreshToken(user models.User) (string, error) {
	expiresAt := time.Now().Add(s.config.JWT.RefreshTokenExpiry)

	claims := &RefreshTokenClaims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   fmt.Sprintf("%d", user.ID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWT.Secret))
} 