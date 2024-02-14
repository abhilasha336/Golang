package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gitlab.com/tuneverse/toolkit/middleware"
)

type mockCache struct {
	data map[string]interface{}
}

func (c *mockCache) Get(k string) (interface{}, bool) {
	val, ok := c.data[k]
	return val, ok
}

func (c *mockCache) Set(k string, x interface{}, d time.Duration) {
	c.data[k] = x
}

func TestErrorLocalizationMiddleware(t *testing.T) {
	// Create a new Gin router
	router := gin.New()

	// Create a mock cache
	cache := &mockCache{
		data: make(map[string]interface{}),
	}

	router.Use(middleware.Localize(middleware.LocaleOptions{
		HeaderLabel:  "Accept-Language",
		ContextLabel: "lan",
	}))

	// Add the ErrorLocalization middleware
	router.Use(middleware.ErrorLocalization(middleware.ErrorLocaleOptions{
		Cache:                  cache,
		CacheExpiration:        time.Minute,
		CacheKeyLabel:          "error_data",
		LocalisationServiceURL: "https://jsonplaceholder.typicode.com/todos/1",
	}))

	// Define a route handler that retrieves the error data from the context
	router.GET("/error", func(c *gin.Context) {
		errorData, _ := c.Get("error_responses")
		c.JSON(http.StatusOK, errorData)
	})

	// Create a new HTTP request
	req, err := http.NewRequest("GET", "/error", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Set the Accept-Language header
	req.Header.Set("Accept-Language", "en")

	// Create a new HTTP response recorder
	w := httptest.NewRecorder()

	// Perform the HTTP request
	router.ServeHTTP(w, req)

	// Assert that the response status code is 200 OK
	assert.Equal(t, http.StatusOK, w.Code)
}
