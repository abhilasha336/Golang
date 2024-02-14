package activitylog

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"gitlab.com/tuneverse/toolkit/utils"

	log "github.com/sirupsen/logrus"
	"gitlab.com/tuneverse/toolkit/consts"
	"gitlab.com/tuneverse/toolkit/models"
)

// ActivityLogOptions represents options for interacting with the activity log.
type ActivityLogOptions struct {
	url string
}

var activitylogObject *ActivityLogOptions

// Init initializes the ActivityLogOptions with the provided URL.
func Init(url string) (*ActivityLogOptions, error) {
	err := initTransportOptions(url)
	if err != nil {
		log.Error("Invalid transport options")

		return nil, err
	}
	setActivityLog(url)
	return &ActivityLogOptions{url: url}, err
}

// setActivityLog sets the singleton instance of ActivityLogOptions.
func setActivityLog(url string) {
	activitylogObject = &ActivityLogOptions{}
	activitylogObject.url = url
}

// initTransportOptions initializes transport options and checks health.
func initTransportOptions(url string) error {

	if utils.IsEmpty(url) {

		return fmt.Errorf("Invalid transport options")
	}
	err := ping(fmt.Sprintf("%s/%s", url, "health"))
	if err != nil {
		log.Error("Error occured while checking health ")
		return err
	}

	return nil

}

// ping checks the health of the activity log service.
func ping(url string) error {

	resp, err := utils.APIRequest(http.MethodGet, url, nil, nil)
	if err != nil {
		log.Error("Error occured while retrieving the response from API request")

		return fmt.Errorf("Ping request failed: %s", err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		log.Error("Error occured while checking status code")
		return fmt.Errorf("Ping request failed: invalid transport Options")
	}
	return nil
}

// AddActivity sends an activity log to the activity log service.
func (activitylog *ActivityLogOptions) Log(activitylogModel models.ActivityLog) (models.ActivityLogResponse, error) {
	var apiResponse map[string]interface{}
	var logResponse models.ActivityLogResponse
	jsonData, err := json.Marshal(activitylogModel)
	if err != nil {
		log.Errorf("error while parsing %v ", err)
		return logResponse, err
	}

	// Create a strings.Reader from the JSON string
	reader := strings.NewReader(string(jsonData))
	url := fmt.Sprintf("%s/%s", activitylog.url, consts.LogUrl)

	// Use http.Client from the standard library
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodPost, url, reader)

	if err != nil {
		log.Errorf("unable to connect API server %v %v", url, err)
		return logResponse, err
	}

	request.Header.Add("Content-Type", "application/json")

	response, err := client.Do(request)

	if err != nil {
		log.Error("[HTTPClient.REQUEST]: Error occur on client.Do()", err)
		return logResponse, err
	}

	resBody, err := io.ReadAll(response.Body)
	if err != nil {
		log.Errorf("unable to connect Activity log service %v", err)
		return logResponse, errors.New("unable to read data from activity log response")
	}
	defer response.Body.Close()

	err = json.Unmarshal(resBody, &apiResponse)

	if err != nil {
		log.Error("error occurred while unmarshaling the data", err)
		return logResponse, errors.New("error occurred while unmarshaling the data")
	}

	if _, ok := apiResponse["error"]; ok {
		return logResponse, errors.New(apiResponse["error"].(string))
	}

	logResponse.Body = string(resBody)
	logResponse.StatusCode = response.StatusCode
	return logResponse, nil
}
