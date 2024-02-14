package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gitlab.com/tuneverse/toolkit/middleware"
)

func TestAPIVersionGuard(t *testing.T) {
	t.Run("should fail", func(t *testing.T) {
		// Create a new Gin context with the necessary headers
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/test", nil)
		c.Request.Header.Set("Accept-version", "v1")

		// Define the test options for the middleware function
		options := middleware.VersionOptions{
			VersionParamLookup: nil,
			AcceptedVersions:   []string{"v1", "v2"},
		}

		// Call the middleware function with the test options and context
		middleware.APIVersionGuard(options)(c)

		// Assert that the middleware function behaves as expected
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d but got %d", http.StatusBadRequest, w.Code)
		}
		expectedBody := "{\"error\":\"Missing version parameter\"}"
		if w.Body.String() != expectedBody {
			t.Errorf("Expected response body %s but got %s", expectedBody, w.Body.String())
		}
	})

	t.Run("shouldn't fail", func(t *testing.T) {
		// Create a new Gin context with the necessary headers
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/test?version=v1", nil)
		c.Request.Header.Set("Accept-version", "v1")

		// Define the test options for the middleware function
		options := middleware.VersionOptions{
			VersionParamLookup: func(c *gin.Context) string {
				return c.Query("version")
			},
			AcceptedVersions: []string{"v1", "v2"},
		}

		// Call the middleware function with the test options and context
		middleware.APIVersionGuard(options)(c)

		// Assert that the middleware function behaves as expected
		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d but got %d", http.StatusOK, w.Code)
		}
	})
}
