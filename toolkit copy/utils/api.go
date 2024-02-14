package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

var HTTPClient = &http.Client{}

// For api request
func APIRequest(method string, url string, headers map[string]interface{},
	body map[string]interface{}) (*http.Response, error) {

	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	// Create a strings.Reader from the JSON string
	reader := strings.NewReader(string(jsonData))

	request, err := http.NewRequest(method, url, reader)
	if err != nil {
		log.Errorf("unable to connect API server %v %v", url, err)
		return nil, err
	}
	request.Header.Add("Content-Type", "application/json")
	SETHeaders(*request, headers)

	response, err := HTTPClient.Do(request)
	if err != nil {
		log.Error("[HTTPClient.REQUEST]: Error occur on HTTClient.Do()", err)
		return nil, err
	}
	log.Infof("[HTTPClient.REQUEST] SERVICE Response Status: %v", response.Status)

	return response, nil
}

// set the headers
func SETHeaders(request http.Request, headers map[string]interface{}) http.Request {

	for key, value := range headers {
		request.Header.Add(key, fmt.Sprintf("%s", value))
	}
	return request
}
