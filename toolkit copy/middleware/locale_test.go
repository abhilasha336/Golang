package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gitlab.com/tuneverse/toolkit/middleware"
)

func TestLocaleMiddleware(t *testing.T) {
	// Create a new Gin router
	router := gin.New()

	// Add the Locale middleware
	router.Use(middleware.Localize(middleware.LocaleOptions{
		HeaderLabel:  "Accept-Language",
		ContextLabel: "lang",
	}))

	// Define a route handler that retrieves the language from the context
	router.GET("/language", func(c *gin.Context) {
		lang, _ := c.Get("lang")
		c.JSON(http.StatusOK, gin.H{
			"language": lang,
		})
	})

	// Create a new HTTP request
	req, err := http.NewRequest("GET", "/language", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Set the Accept-Language header
	req.Header.Set("Accept-Language", "fr")

	// Create a new HTTP response recorder
	w := httptest.NewRecorder()

	// Perform the HTTP request
	router.ServeHTTP(w, req)

	// Assert that the response status code is 200 OK
	assert.Equal(t, http.StatusOK, w.Code)

	// Assert that the response body contains the expected language
	assert.JSONEq(t, `{"language": "fr"}`, w.Body.String())
}
