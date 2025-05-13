package controllers

import (
	"github.com/gin-gonic/gin"
	"zionechainapi/configs"
	"zionechainapi/internal/services"
	"zionechainapi/internal/utils"
)

// AuthController handles authentication-related routes
type AuthController struct {
	config      *configs.Config
	authService *services.AuthService
}

// NewAuthController creates a new auth controller
func NewAuthController(config *configs.Config) *AuthController {
	return &AuthController{
		config:      config,
		authService: services.NewAuthService(config),
	}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with name, email, phone, and password
// @Tags auth
// @Accept json
// @Produce json
// @Param body body services.RegisterRequest true "Register request"
// @Success 201 {object} utils.Response{data=services.TokenResponse} "User registered successfully"
// @Failure 400 {object} utils.Response "Bad request"
// @Failure 422 {object} utils.Response "Validation error"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /api/auth/register [post]
func (c *AuthController) Register(ctx *gin.Context) {
	var req services.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err.Error())
		return
	}

	token, err := c.authService.Register(req)
	if err != nil {
		utils.BadRequestResponse(ctx, "Failed to register user", err.Error())
		return
	}

	utils.CreatedResponse(ctx, "User registered successfully", token)
}

// Login godoc
// @Summary Login a user
// @Description Login a user with phone and password
// @Tags auth
// @Accept json
// @Produce json
// @Param body body services.LoginRequest true "Login request"
// @Success 200 {object} utils.Response{data=services.TokenResponse} "User logged in successfully"
// @Failure 400 {object} utils.Response "Bad request"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 422 {object} utils.Response "Validation error"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /api/auth/login [post]
func (c *AuthController) Login(ctx *gin.Context) {
	var req services.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err.Error())
		return
	}

	token, err := c.authService.Login(req)
	if err != nil {
		utils.UnauthorizedResponse(ctx, err.Error())
		return
	}

	utils.OKResponse(ctx, "User logged in successfully", token)
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Refresh access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param body body map[string]string true "Refresh token request"
// @Success 200 {object} utils.Response{data=services.TokenResponse} "Token refreshed successfully"
// @Failure 400 {object} utils.Response "Bad request"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 422 {object} utils.Response "Validation error"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /api/auth/refresh [post]
func (c *AuthController) RefreshToken(ctx *gin.Context) {
	var req map[string]string
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err.Error())
		return
	}

	refreshToken, ok := req["refresh_token"]
	if !ok || refreshToken == "" {
		utils.BadRequestResponse(ctx, "Refresh token is required", nil)
		return
	}

	token, err := c.authService.RefreshToken(refreshToken)
	if err != nil {
		utils.UnauthorizedResponse(ctx, err.Error())
		return
	}

	utils.OKResponse(ctx, "Token refreshed successfully", token)
}

// Me godoc
// @Summary Get current user
// @Description Get current authenticated user
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=services.UserResponse} "User retrieved successfully"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /api/auth/me [get]
func (c *AuthController) Me(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		utils.UnauthorizedResponse(ctx, "User not authenticated")
		return
	}

	user, err := c.authService.GetUserByID(userID.(uint))
	if err != nil {
		utils.InternalServerErrorResponse(ctx, err.Error())
		return
	}

	utils.OKResponse(ctx, "User retrieved successfully", services.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Phone: user.Phone,
		Role:  user.Role.Name,
	})
}

// Routes registers auth routes
func (c *AuthController) Routes(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", c.Register)
		auth.POST("/login", c.Login)
		auth.POST("/refresh", c.RefreshToken)
		auth.GET("/me", c.Me)
	}
} 