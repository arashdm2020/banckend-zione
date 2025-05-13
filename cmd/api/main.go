package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"zionechainapi/configs"
	"zionechainapi/internal/controllers"
	"zionechainapi/internal/database"
	"zionechainapi/internal/middleware"
)

func main() {
	// Load configuration
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Set Gin mode
	if config.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize router
	router := gin.Default()

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     config.CORS.AllowedOrigins,
		AllowMethods:     config.CORS.AllowedMethods,
		AllowHeaders:     config.CORS.AllowedHeaders,
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Connect to database
	db, err := database.Connect(config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migrate database
	if err := database.AutoMigrate(); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize API routes
	initializeRoutes(router, config)

	// Create server
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", config.App.Host, config.App.Port),
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on %s...\n", server.Addr)
		var err error
		if config.TLS.Enabled {
			err = server.ListenAndServeTLS(config.TLS.CertFile, config.TLS.KeyFile)
		} else {
			err = server.ListenAndServe()
		}
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown server
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	// Close database connection
	if err := database.Close(); err != nil {
		log.Fatalf("Failed to close database connection: %v", err)
	}

	log.Println("Server exiting")
}

// Initialize routes for the API
func initializeRoutes(router *gin.Engine, config *configs.Config) {
	// API base group
	api := router.Group("/api")

	// Initialize controllers
	authController := controllers.NewAuthController(config)
	projectController := controllers.NewProjectController(config)
	blogController := controllers.NewBlogController(config)
	categoryController := controllers.NewCategoryController(config)

	// Register routes
	authController.Routes(api)
	
	// Create auth middleware for protected routes
	authMiddleware := middleware.Auth(config)
	
	// Register controller routes
	projectController.Routes(api, authMiddleware)
	blogController.Routes(api, authMiddleware)
	categoryController.Routes(api, authMiddleware)
} 