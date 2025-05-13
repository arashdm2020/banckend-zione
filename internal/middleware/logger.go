package middleware

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestLogger is a middleware that logs detailed information about HTTP requests
func RequestLogger() gin.HandlerFunc {
	// Create logs directory if it doesn't exist
	logsDir := "logs"
	if _, err := os.Stat(logsDir); os.IsNotExist(err) {
		if err := os.MkdirAll(logsDir, 0755); err != nil {
			fmt.Printf("Error creating logs directory: %v\n", err)
		}
	}

	// Create or open log file for appending
	logFile := filepath.Join(logsDir, time.Now().Format("2006-01-02")+".log")
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening log file: %v\n", err)
	} else {
		defer file.Close()
	}

	// Keep a console copy of logs
	multiWriter := gin.DefaultWriter

	// If log file opened successfully, write to both console and file
	if file != nil {
		multiWriter = io.MultiWriter(gin.DefaultWriter, file)
	}

	return func(c *gin.Context) {
		// Get the request body for logging
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			// Restore the body for the next middleware/handler
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get request details
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		userAgent := c.Request.UserAgent()
		
		// Format request parameters (if any)
		var params string
		if len(requestBody) > 0 {
			// Only log if body is not too large
			if len(requestBody) < 1024 { // Only log if less than 1KB
				params = string(requestBody)
			} else {
				params = fmt.Sprintf("[Body too large: %d bytes]", len(requestBody))
			}
		} else if query != "" {
			params = "?" + query
		}

		// Get response status
		responseStatus := "Success"
		if statusCode >= 400 {
			responseStatus = "Error"
		}
		
		// Format the log entry
		logEntry := fmt.Sprintf("[REQUEST] %v | %s | %s %s | %d | %v | %s | %s | User-Agent: %s | %s\n",
			time.Now().Format("2006/01/02 - 15:04:05"),
			responseStatus,
			method, path,
			statusCode,
			latency,
			clientIP,
			params,
			userAgent,
			c.Errors.String(),
		)
		
		// Write to multiWriter (console and file if available)
		fmt.Fprint(multiWriter, logEntry)
	}
} 