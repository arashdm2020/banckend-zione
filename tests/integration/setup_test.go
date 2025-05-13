package integration

import (
	"log"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"zionechainapi/configs"
	"zionechainapi/internal/controllers"
	"zionechainapi/internal/database"
	"zionechainapi/internal/middleware"
)

var (
	router *gin.Engine
	config *configs.Config
)

func TestMain(m *testing.M) {
	// Set test mode for Gin
	gin.SetMode(gin.TestMode)

	// Load test configuration
	var err error
	config, err = configs.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Override config with test values
	config.App.Env = "testing"
	config.Database.Name = "zione_test_db"

	// Setup database connection
	_, err = database.Connect(config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migrate the database
	if err := database.AutoMigrate(); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize router and routes
	router = gin.Default()
	setupRoutes(router)

	// Run tests
	code := m.Run()

	// Clean up
	if err := database.Close(); err != nil {
		log.Fatalf("Failed to close database connection: %v", err)
	}

	os.Exit(code)
}

func setupRoutes(router *gin.Engine) {
	// API base group
	api := router.Group("/api")

	// Initialize controllers
	authController := controllers.NewAuthController(config)
	projectController := controllers.NewProjectController(config)
	blogController := controllers.NewBlogController(config)
	categoryController := controllers.NewCategoryController(config)
	tagController := controllers.NewTagController(config)

	// Register routes
	authController.Routes(api)
	
	// Create auth middleware for protected routes
	authMiddleware := middleware.Auth(config)
	
	// Register controller routes
	projectController.Routes(api, authMiddleware)
	blogController.Routes(api, authMiddleware)
	categoryController.Routes(api, authMiddleware)
	tagController.Routes(api, authMiddleware)
}