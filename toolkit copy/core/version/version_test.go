package version

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gitlab.com/tuneverse/toolkit/middleware"
)

type TestObject struct{}

func (t *TestObject) GetUsersv1(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "GetUsers_v1"})
}

func (t *TestObject) GetUsersv2(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "GetUsers_v2"})
}

func (t *TestObject) GetUsers(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "GetUsers"})
}

func TestRenderHandler(t *testing.T) {
	// Create a new Gin context with the necessary headers
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Define the test options for the middleware function
	object := &TestObject{}
	method := "GetUsers"
	args := []interface{}{}

	var expectedBody string

	// Define the test options for the middleware function
	options := middleware.VersionOptions{
		VersionParamLookup: func(c *gin.Context) string {
			return c.Query("version")
		},
		AcceptedVersions: []string{"v1", "v2"},
	}

	c.Request, _ = http.NewRequest("GET", "/test?version=v1", nil)
	c.Request.Header.Set("Accept-version", "v2")

	t.Run("test 001", func(t *testing.T) {

		// Call the middleware function with the test options and context
		middleware.APIVersionGuard(options)(c)

		// Call the middleware function with the test options and context
		RenderHandler(c, object, method, args...)

		// Assert that the middleware function behaves as expected
		expectedBody = "{\"message\":\"GetUsers\"}"
		if w.Body.String() != expectedBody {
			t.Errorf("Expected response body %s but got %s", expectedBody, w.Body.String())
		}

	})
	t.Run("test 002", func(t *testing.T) {
		// Call the middleware function with an invalid method and assert the result
		w = httptest.NewRecorder()
		method = "InvalidMethod"
		c.Request.Header.Set("Accept-version", "v1")

		// Call the middleware function with the test options and context
		middleware.APIVersionGuard(options)(c)

		defer func() {
			if r := recover(); r != nil {
				expectedBody = "unable to locate the method InvalidMethod"
				if r != expectedBody {
					t.Errorf("Expected response body %s but got %s", expectedBody, r)
				}
			}
		}()
		RenderHandler(c, object, method, args...)

	})
}
