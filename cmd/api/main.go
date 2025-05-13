package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"zionechainapi/configs"
	"zionechainapi/internal/controllers"
	"zionechainapi/internal/database"
	"zionechainapi/internal/middleware"
)

// Define available routes for better documentation
var availableRoutes = []struct {
	Method string
	Path   string
	Desc   string
	Access string
}{
	{"GET", "/", "API Status - Check if API is running", "Public"},
	{"GET", "/health", "Health Check - Server health status", "Public"},
	{"GET", "/api", "API Welcome - Welcome message and version info", "Public"},
	{"POST", "/api/auth/login", "Login via phone/password", "Public"},
	{"POST", "/api/auth/register", "Register new user", "Public"},
	{"GET", "/api/projects", "Get list of projects", "Public"},
	{"POST", "/api/projects", "Create project", "Admin"},
	{"GET", "/api/blog", "Get blog posts", "Public"},
	{"POST", "/api/blog", "Create blog post", "Admin"},
	{"GET", "/api/categories/projects", "Get project categories", "Public"},
	{"GET", "/api/categories/blog", "Get blog categories", "Public"},
	
	// Resume endpoints
	{"GET", "/api/resume/personal", "Get personal information", "Public"},
	{"POST", "/api/resume/personal", "Create personal information", "Admin"},
	{"PUT", "/api/resume/personal/:id", "Update personal information", "Admin"},
	{"DELETE", "/api/resume/personal/:id", "Delete personal information", "Admin"},
	
	{"GET", "/api/resume/skills", "Get skills", "Public"},
	{"POST", "/api/resume/skills", "Create skill", "Admin"},
	{"PUT", "/api/resume/skills/:id", "Update skill", "Admin"},
	{"DELETE", "/api/resume/skills/:id", "Delete skill", "Admin"},
	
	{"GET", "/api/resume/experience", "Get work experience", "Public"},
	{"POST", "/api/resume/experience", "Create work experience", "Admin"},
	{"PUT", "/api/resume/experience/:id", "Update work experience", "Admin"},
	{"DELETE", "/api/resume/experience/:id", "Delete work experience", "Admin"},
	
	{"GET", "/api/resume/education", "Get education details", "Public"},
	{"POST", "/api/resume/education", "Create education detail", "Admin"},
	{"PUT", "/api/resume/education/:id", "Update education detail", "Admin"},
	{"DELETE", "/api/resume/education/:id", "Delete education detail", "Admin"},
	
	{"GET", "/api/resume/certificates", "Get certificates", "Public"},
	{"POST", "/api/resume/certificates", "Create certificate", "Admin"},
	{"PUT", "/api/resume/certificates/:id", "Update certificate", "Admin"},
	{"DELETE", "/api/resume/certificates/:id", "Delete certificate", "Admin"},
	
	{"GET", "/api/resume/languages", "Get languages", "Public"},
	{"POST", "/api/resume/languages", "Create language", "Admin"},
	{"PUT", "/api/resume/languages/:id", "Update language", "Admin"},
	{"DELETE", "/api/resume/languages/:id", "Delete language", "Admin"},
	
	{"GET", "/api/resume/publications", "Get publications", "Public"},
	{"POST", "/api/resume/publications", "Create publication", "Admin"},
	{"PUT", "/api/resume/publications/:id", "Update publication", "Admin"},
	{"DELETE", "/api/resume/publications/:id", "Delete publication", "Admin"},
	
	{"GET", "/api/resume/complete", "Get complete resume", "Public"},
}

func main() {
	// Load configuration
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Setup database connection
	db, err := database.Connect(config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migrate the database
	if err := database.AutoMigrate(); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Set Gin mode based on environment
	if config.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize Gin router
	router := gin.Default()

	// Add basic middleware
	router.Use(gin.Recovery())
	router.Use(middleware.RequestLogger())

	// Root endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Zione API is running!"})
	})

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	// API base group
	api := router.Group("/api")
	
	// API welcome endpoint
	api.GET("", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to Zione API",
			"version": "1.0.0",
		})
	})

	// Initialize controllers
	authController := controllers.NewAuthController(config)
	projectController := controllers.NewProjectController(config)
	blogController := controllers.NewBlogController(config)
	categoryController := controllers.NewCategoryController(config)
	tagController := controllers.NewTagController(config)
	
	// Initialize resume controller with the database connection
	resumeController := controllers.NewResumeController(db)

	// Register auth routes (no middleware needed for these)
	authController.Routes(api)
	
	// Create auth middleware for protected routes
	authMiddleware := middleware.Auth(config)
	
	// Register controller routes that need auth for some endpoints
	projectController.Routes(api, authMiddleware)
	blogController.Routes(api, authMiddleware)
	categoryController.Routes(api, authMiddleware)
	tagController.Routes(api, authMiddleware)
	
	// Register resume routes
	resumeController.Routes(api)

	// Get port from environment or use default
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}

	// Create server with configured router
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Print available routes
	fmt.Println("\n=== Available API Routes ===")
	fmt.Println("Server will start on http://localhost:" + port)
	fmt.Println()
	
	fmt.Printf("%-7s %-40s %-35s %s\n", "Method", "Route", "Description", "Access")
	fmt.Println(strings.Repeat("-", 100))
	for _, route := range availableRoutes {
		fmt.Printf("%-7s %-40s %-35s %s\n", route.Method, "http://localhost:"+port+route.Path, route.Desc, route.Access)
	}
	fmt.Println("\nPress Ctrl+C to stop the server")
	fmt.Println("=============================")

	// Start server in a goroutine
	go func() {
		fmt.Printf("\nServer is running on port %s...\n", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	fmt.Println("Shutting down server...")
	
	// Create context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	
	// Close database connection
	if err := database.Close(); err != nil {
		log.Fatalf("Failed to close database connection: %v", err)
	}
	
	fmt.Println("Server exited properly")
} 
