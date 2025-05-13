package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"zionechainapi/internal/services"
)

func TestRegister(t *testing.T) {
	// Create a test user
	registerRequest := services.RegisterRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Phone:    "+1234567890",
		Password: "password123",
	}

	// Convert to JSON
	jsonData, err := json.Marshal(registerRequest)
	assert.NoError(t, err)

	// Create a request
	req, err := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(jsonData))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	w := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusCreated, w.Code)

	// Parse the response
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Assert that the response contains the expected fields
	assert.Equal(t, true, response["success"])
	assert.Equal(t, "User registered successfully", response["message"])
	assert.NotNil(t, response["data"])

	// Extract the token from the response
	data := response["data"].(map[string]interface{})
	assert.NotNil(t, data["access_token"])
	assert.NotNil(t, data["refresh_token"])
	assert.NotNil(t, data["expires_at"])
	assert.NotNil(t, data["user"])

	// Extract the user from the response
	user := data["user"].(map[string]interface{})
	assert.Equal(t, "Test User", user["name"])
	assert.Equal(t, "test@example.com", user["email"])
	assert.Equal(t, "+1234567890", user["phone"])
	assert.Equal(t, "user", user["role"])
}

func TestLogin(t *testing.T) {
	// Create a login request
	loginRequest := services.LoginRequest{
		Phone:    "+1234567890",
		Password: "password123",
	}

	// Convert to JSON
	jsonData, err := json.Marshal(loginRequest)
	assert.NoError(t, err)

	// Create a request
	req, err := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonData))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	w := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse the response
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Assert that the response contains the expected fields
	assert.Equal(t, true, response["success"])
	assert.Equal(t, "User logged in successfully", response["message"])
	assert.NotNil(t, response["data"])

	// Extract the token from the response
	data := response["data"].(map[string]interface{})
	assert.NotNil(t, data["access_token"])
	assert.NotNil(t, data["refresh_token"])
	assert.NotNil(t, data["expires_at"])
	assert.NotNil(t, data["user"])

	// Extract the user from the response
	user := data["user"].(map[string]interface{})
	assert.Equal(t, "Test User", user["name"])
	assert.Equal(t, "test@example.com", user["email"])
	assert.Equal(t, "+1234567890", user["phone"])
	assert.Equal(t, "user", user["role"])
}