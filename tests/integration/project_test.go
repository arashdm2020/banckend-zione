package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"zionechainapi/internal/services"
)

var projectID uint
var accessToken string

func TestCreateProject(t *testing.T) {
	// First, login to get an access token
	loginAndGetToken(t)

	// Create a test project
	createRequest := services.CreateProjectRequest{
		Title:       "Test Project",
		Description: "This is a test project description",
		Content:     "This is the content of the test project",
		CategoryID:  1, // Assumes category ID 1 exists
		TagIDs:      []uint{1}, // Assumes tag ID 1 exists
		Featured:    true,
		Published:   true,
	}

	// Convert to JSON
	jsonData, err := json.Marshal(createRequest)
	assert.NoError(t, err)

	// Create a request
	req, err := http.NewRequest("POST", "/api/projects", bytes.NewBuffer(jsonData))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

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
	assert.Equal(t, "Project created successfully", response["message"])
	assert.NotNil(t, response["data"])

	// Extract the project from the response
	project := response["data"].(map[string]interface{})
	assert.Equal(t, "Test Project", project["title"])
	assert.Equal(t, "This is a test project description", project["description"])
	assert.Equal(t, "This is the content of the test project", project["content"])
	assert.Equal(t, float64(1), project["category_id"])
	assert.Equal(t, true, project["featured"])
	assert.Equal(t, true, project["published"])

	// Store the project ID for later tests
	projectID = uint(project["id"].(float64))
}

func TestGetProject(t *testing.T) {
	// Create a request
	req, err := http.NewRequest("GET", fmt.Sprintf("/api/projects/%d", projectID), nil)
	assert.NoError(t, err)

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
	assert.Equal(t, "Project retrieved successfully", response["message"])
	assert.NotNil(t, response["data"])

	// Extract the project from the response
	project := response["data"].(map[string]interface{})
	assert.Equal(t, float64(projectID), project["id"])
	assert.Equal(t, "Test Project", project["title"])
	assert.Equal(t, "This is a test project description", project["description"])
	assert.Equal(t, "This is the content of the test project", project["content"])
}

func TestUpdateProject(t *testing.T) {
	// Create an update request
	updateRequest := services.UpdateProjectRequest{
		Title:       "Updated Test Project",
		Description: "This is an updated test project description",
		Content:     "This is the updated content of the test project",
	}

	// Convert to JSON
	jsonData, err := json.Marshal(updateRequest)
	assert.NoError(t, err)

	// Create a request
	req, err := http.NewRequest("PUT", fmt.Sprintf("/api/projects/%d", projectID), bytes.NewBuffer(jsonData))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

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
	assert.Equal(t, "Project updated successfully", response["message"])
	assert.NotNil(t, response["data"])

	// Extract the project from the response
	project := response["data"].(map[string]interface{})
	assert.Equal(t, float64(projectID), project["id"])
	assert.Equal(t, "Updated Test Project", project["title"])
	assert.Equal(t, "This is an updated test project description", project["description"])
	assert.Equal(t, "This is the updated content of the test project", project["content"])
}

func TestDeleteProject(t *testing.T) {
	// Create a request
	req, err := http.NewRequest("DELETE", fmt.Sprintf("/api/projects/%d", projectID), nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	// Create a response recorder
	w := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusNoContent, w.Code)
}

func loginAndGetToken(t *testing.T) {
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

	// Extract the token from the response
	data := response["data"].(map[string]interface{})
	accessToken = data["access_token"].(string)
}