package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"zionechainapi/configs"
	"zionechainapi/internal/services"
)

// Auth is the authentication middleware
func Auth(config *configs.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get auth header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
			c.Abort()
			return
		}

		// Check if header is in correct format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header format must be Bearer {token}"})
			c.Abort()
			return
		}

		// Get token
		token := parts[1]

		// Validate token
		authService := services.NewAuthService(config)
		claims, err := authService.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		// Set user ID and role in context
		c.Set("userID", claims.UserID)
		c.Set("userRole", claims.Role)

		c.Next()
	}
}

// RequireRole is the role-based access control middleware
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user role from context
		userRole, exists := c.Get("userRole")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			c.Abort()
			return
		}

		// Check if user has required role
		role := userRole.(string)
		for _, r := range roles {
			// Admin role has access to everything
			if role == "admin" || role == r {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
		c.Abort()
	}
}

// GetUserID gets the user ID from the context
func GetUserID(c *gin.Context) uint {
	userID, exists := c.Get("userID")
	if !exists {
		return 0
	}
	return userID.(uint)
}

// GetUserRole gets the user role from the context
func GetUserRole(c *gin.Context) string {
	userRole, exists := c.Get("userRole")
	if !exists {
		return ""
	}
	return userRole.(string)
} 