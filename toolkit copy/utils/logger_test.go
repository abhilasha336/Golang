package utils

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gitlab.com/tuneverse/toolkit/consts"
)

func TestGetRequestUrl(t *testing.T) {
	tests := []struct {
		name           string
		request        *http.Request
		expectedResult string
	}{
		{
			name: "HTTP Request",
			request: &http.Request{
				Host: "example.com",
				URL: &url.URL{
					Path:     "/test",
					RawQuery: "param=value",
				},
			},
			expectedResult: "http://example.com/test?param=value",
		},
		{
			name: "HTTPS Request",
			request: &http.Request{
				Host: "example.com",
				TLS:  &tls.ConnectionState{},
				URL: &url.URL{
					Path:     "/secure",
					RawQuery: "param=value",
				},
			},
			expectedResult: "https://example.com/secure?param=value",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := ConstructURL(test.request)
			require.Equal(t, test.expectedResult, result)
		})
	}
}

func TestGetRequestRoute(t *testing.T) {
	tests := []struct {
		name           string
		context        *gin.Context
		expectedResult string
	}{
		{
			name: "Simple Route",
			context: &gin.Context{
				Request: &http.Request{
					URL: &url.URL{
						Path: "/user/123",
					},
				},
				Params: gin.Params{{Key: "id", Value: "123"}},
			},
			expectedResult: "/user/:id",
		},
		{
			name: "Complex Route",
			context: &gin.Context{
				Request: &http.Request{
					URL: &url.URL{
						Path: "/users/123/orders/456",
					},
				},
				Params: gin.Params{{Key: "user_id", Value: "123"}, {Key: "order_id", Value: "456"}},
			},
			expectedResult: "/users/:user_id/orders/:order_id",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := GetRequestRoute(test.context)
			require.Equal(t, test.expectedResult, result)
		})
	}
}

func TestTempDir(t *testing.T) {
	expectedResult := consts.DefaultDirName
	result := TempDir()
	require.Equal(t, expectedResult, result)
}

func TestGenerateRequestID(t *testing.T) {
	result := GenerateRequestID()
	_, err := uuid.Parse(result)
	require.NoError(t, err)
}
