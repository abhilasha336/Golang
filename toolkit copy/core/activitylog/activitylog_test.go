package activitylog

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gitlab.com/tuneverse/toolkit/models"
)

func TestAddActivity(t *testing.T) {
	// Create a sample ActivityLogOptions and ActivityLog for testing
	activityLogOptions := &ActivityLogOptions{
		url: "https://example.com",
		// Add any other necessary fields for your testing
	}

	activityLogModel := models.ActivityLog{}

	// Mock the HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Verify the request method and URL
		if req.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", req.Method)
		}
		if req.URL.Path != "/activitylog" { // Update this based on your actual endpoint
			t.Errorf("Expected request to /activitylog, got %s", req.URL.Path)
		}

		// Parse the request body
		var receivedData models.ActivityLog
		err := json.NewDecoder(req.Body).Decode(&receivedData)
		if err != nil {
			t.Errorf("Error decoding request body: %v", err)
		}

		// Respond with a dummy success response
		response := map[string]interface{}{"status": "success"}
		responseJSON, _ := json.Marshal(response)
		rw.Write(responseJSON)
	}))
	defer server.Close()

	// Set the URL in the ActivityLogOptions to the mock server URL
	activityLogOptions.url = server.URL

	// Call the AddActivity function
	response, err := activityLogOptions.Log(activityLogModel)

	// Check for errors
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Add assertions based on your expectations for the response
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.StatusCode)
	}

}
