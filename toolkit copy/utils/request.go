package utils

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.com/tuneverse/toolkit/consts"
)

func GetRequestIDFromRequest(r *http.Request) string {
	rid := r.Header.Get(consts.ContextRequestID)
	if len(rid) < 1 {
		rid = GenerateRequestID()
	}
	return rid
}

func GenerateRequestID() string {
	return uuid.New().String()
}

func GetRequestRoute(c *gin.Context) string {
	urlTemplate := c.Request.URL.Path
	for _, p := range c.Params {
		urlTemplate = strings.Replace(urlTemplate, p.Value, ":"+p.Key, 1)
	}
	return urlTemplate
}

func ConstructURL(req *http.Request) string {
	var scheme string
	if req.TLS != nil {
		scheme = "https"
	} else {
		scheme = "http"
	}
	// Construct the full URL
	return fmt.Sprintf("%s://%s%s?%s", scheme, req.Host, req.URL.Path, req.URL.RawQuery)
}

type RequestDataDump struct {
	Headers map[string][]string `json:"headers"`
	Body    string              `json:"body"`
}

// capture request data without modifying the request body
func GetRequestDump(req *http.Request) (*RequestDataDump, error) {
	requestData := RequestDataDump{
		Headers: req.Header,
	}

	buf, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	requestData.Body = string(buf)
	reader := io.NopCloser(bytes.NewBuffer(buf))
	req.Body = reader

	return &requestData, nil
}
