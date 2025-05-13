package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
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
	// Define simple handlers
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			handleNotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"message": "Zione API is running!"}`)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status": "OK"}`)
	})

	// API routes handler
	http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		method := r.Method

		// API welcome endpoint
		if path == "/api" || path == "/api/" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `{"message": "Welcome to Zione API", "version": "1.0.0"}`)
			return
		}

		// Auth endpoints
		if strings.HasPrefix(path, "/api/auth/") {
			if path == "/api/auth/login" && method == "POST" {
				handleLogin(w, r)
				return
			}
			if path == "/api/auth/register" && method == "POST" {
				handleRegister(w, r)
				return
			}
		}

		// Projects endpoints
		if strings.HasPrefix(path, "/api/projects") {
			if path == "/api/projects" && method == "GET" {
				handleGetProjects(w, r)
				return
			}
			if path == "/api/projects" && method == "POST" {
				handleCreateProject(w, r)
				return
			}
		}

		// Blog endpoints
		if strings.HasPrefix(path, "/api/blog") {
			if path == "/api/blog" && method == "GET" {
				handleGetBlogPosts(w, r)
				return
			}
			if path == "/api/blog" && method == "POST" {
				handleCreateBlogPost(w, r)
				return
			}
		}

		// Categories endpoints
		if strings.HasPrefix(path, "/api/categories/") {
			if path == "/api/categories/projects" && method == "GET" {
				handleGetProjectCategories(w, r)
				return
			}
			if path == "/api/categories/blog" && method == "GET" {
				handleGetBlogCategories(w, r)
				return
			}
		}
		
		// Resume endpoints
		if strings.HasPrefix(path, "/api/resume/") {
			handleResumeRoutes(w, r)
			return
		}

		// If we get here, route was not found
		handleNotFound(w, r)
	})

	// Get port from environment or use default
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}

	// Create server
	srv := &http.Server{
		Addr: ":" + port,
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
}

// Handler functions for each endpoint
func handleNotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, `{"error": "Not Found", "message": "The requested resource was not found"}`)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message": "Login endpoint", "status": "Not implemented"}`)
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message": "Register endpoint", "status": "Not implemented"}`)
}

func handleGetProjects(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message": "Get projects endpoint", "projects": []}`)
}

func handleCreateProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message": "Create project endpoint", "status": "Not implemented"}`)
}

func handleGetBlogPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message": "Get blog posts endpoint", "posts": []}`)
}

func handleCreateBlogPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message": "Create blog post endpoint", "status": "Not implemented"}`)
}

func handleGetProjectCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message": "Get project categories endpoint", "categories": []}`)
}

func handleGetBlogCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message": "Get blog categories endpoint", "categories": []}`)
}

// Resume related handlers
func handleResumeRoutes(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	method := r.Method
	
	w.Header().Set("Content-Type", "application/json")
	
	// Personal info
	if path == "/api/resume/personal" {
		if method == "GET" {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `{"message": "Get personal info endpoint", "data": []}`)
			return
		} else if method == "POST" {
			w.WriteHeader(http.StatusCreated)
			fmt.Fprintf(w, `{"message": "Create personal info endpoint", "status": "Not implemented"}`)
			return
		}
	}
	
	if strings.HasPrefix(path, "/api/resume/personal/") {
		if method == "PUT" {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `{"message": "Update personal info endpoint", "status": "Not implemented"}`)
			return
		} else if method == "DELETE" {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `{"message": "Delete personal info endpoint", "status": "Not implemented"}`)
			return
		}
	}
	
	// Skills
	if path == "/api/resume/skills" {
		if method == "GET" {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `{"message": "Get skills endpoint", "data": []}`)
			return
		} else if method == "POST" {
			w.WriteHeader(http.StatusCreated)
			fmt.Fprintf(w, `{"message": "Create skill endpoint", "status": "Not implemented"}`)
			return
		}
	}
	
	if strings.HasPrefix(path, "/api/resume/skills/") {
		if method == "PUT" {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `{"message": "Update skill endpoint", "status": "Not implemented"}`)
			return
		} else if method == "DELETE" {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `{"message": "Delete skill endpoint", "status": "Not implemented"}`)
			return
		}
	}
	
	// Experience
	if path == "/api/resume/experience" {
		if method == "GET" {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `{"message": "Get experience endpoint", "data": []}`)
			return
		} else if method == "POST" {
			w.WriteHeader(http.StatusCreated)
			fmt.Fprintf(w, `{"message": "Create experience endpoint", "status": "Not implemented"}`)
			return
		}
	}
	
	if strings.HasPrefix(path, "/api/resume/experience/") {
		if method == "PUT" {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `{"message": "Update experience endpoint", "status": "Not implemented"}`)
			return
		} else if method == "DELETE" {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `{"message": "Delete experience endpoint", "status": "Not implemented"}`)
			return
		}
	}
	
	// Education
	if path == "/api/resume/education" {
		if method == "GET" {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `{"message": "Get education endpoint", "data": []}`)
			return
		} else if method == "POST" {
			w.WriteHeader(http.StatusCreated)
			fmt.Fprintf(w, `{"message": "Create education endpoint", "status": "Not implemented"}`)
			return
		}
	}
	
	if strings.HasPrefix(path, "/api/resume/education/") {
		if method == "PUT" {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `{"message": "Update education endpoint", "status": "Not implemented"}`)
			return
		} else if method == "DELETE" {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `{"message": "Delete education endpoint", "status": "Not implemented"}`)
			return
		}
	}
	
	// Certificates
	if path == "/api/resume/certificates" {
		if method == "GET" {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `{"message": "Get certificates endpoint", "data": []}`)
			return
		} else if method == "POST" {
			w.WriteHeader(http.StatusCreated)
			fmt.Fprintf(w, `{"message": "Create certificate endpoint", "status": "Not implemented"}`)
			return
		}
	}
	
	if strings.HasPrefix(path, "/api/resume/certificates/") {
		if method == "PUT" {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `{"message": "Update certificate endpoint", "status": "Not implemented"}`)
			return
		} else if method == "DELETE" {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `{"message": "Delete certificate endpoint", "status": "Not implemented"}`)
			return
		}
	}
	
	// Languages
	if path == "/api/resume/languages" {
		if method == "GET" {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `{"message": "Get languages endpoint", "data": []}`)
			return
		} else if method == "POST" {
			w.WriteHeader(http.StatusCreated)
			fmt.Fprintf(w, `{"message": "Create language endpoint", "status": "Not implemented"}`)
			return
		}
	}
	
	if strings.HasPrefix(path, "/api/resume/languages/") {
		if method == "PUT" {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `{"message": "Update language endpoint", "status": "Not implemented"}`)
			return
		} else if method == "DELETE" {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `{"message": "Delete language endpoint", "status": "Not implemented"}`)
			return
		}
	}
	
	// Publications
	if path == "/api/resume/publications" {
		if method == "GET" {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `{"message": "Get publications endpoint", "data": []}`)
			return
		} else if method == "POST" {
			w.WriteHeader(http.StatusCreated)
			fmt.Fprintf(w, `{"message": "Create publication endpoint", "status": "Not implemented"}`)
			return
		}
	}
	
	if strings.HasPrefix(path, "/api/resume/publications/") {
		if method == "PUT" {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `{"message": "Update publication endpoint", "status": "Not implemented"}`)
			return
		} else if method == "DELETE" {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `{"message": "Delete publication endpoint", "status": "Not implemented"}`)
			return
		}
	}
	
	// Complete resume
	if path == "/api/resume/complete" {
		if method == "GET" {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `{
				"message": "Get complete resume endpoint",
				"data": {
					"personal_info": [],
					"skills": [],
					"experience": [],
					"education": [],
					"certificates": [],
					"languages": [],
					"publications": []
				}
			}`)
			return
		}
	}
	
	handleNotFound(w, r)
} 
