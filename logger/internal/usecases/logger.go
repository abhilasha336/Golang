package usecases

import (
	"context"
	"fmt"
	"logger/internal/consts"
	"logger/internal/entities"
	"logger/internal/repo"
	"strings"
	"time"

	constants "gitlab.com/tuneverse/toolkit/consts"
	"gitlab.com/tuneverse/toolkit/models"
	"gitlab.com/tuneverse/toolkit/utils"
	"go.mongodb.org/mongo-driver/bson"
)

// LoggerUseCases represents the use cases for handling log-related operations.
type LoggerUseCases struct {
	repo repo.LoggerRepoImply
}

// LoggerUseCaseImply specifies the contract for the LoggerUseCases type.
type LoggerUseCaseImply interface {
	AddLog(ctx context.Context, log entities.Log) error
	GetLogs(ctx context.Context, params entities.LogParams) (*entities.Response, error)
	BuildFilterQuery(ctx context.Context, params entities.LogParams) (bson.M, error)
}

/*
NewLoggerUseCases creates a new instance of the LoggerUseCases type, initializing it
with the provided LoggerRepoImply repository for database interactions
*/
func NewLoggerUseCases(userRepo repo.LoggerRepoImply) LoggerUseCaseImply {
	return &LoggerUseCases{
		repo: userRepo,
	}
}

// AddLog insert log data based on the provided information.
func (logger *LoggerUseCases) AddLog(ctx context.Context, log entities.Log) error {
	log.Month = int(log.Timestamp.Month())
	log.Day = log.Timestamp.Day()
	log.Year = log.Timestamp.Year()
	err := logger.repo.AddLog(ctx, log)
	if err != nil {
		return err
	}
	return nil
}

/*
GetLogs retrieves log data based on the provided log parameters and returns the logs
as well as metadata such as pagination information.
*/
func (logger *LoggerUseCases) GetLogs(ctx context.Context, params entities.LogParams) (*entities.Response, error) {

	filters, err := logger.BuildFilterQuery(ctx, params)
	if err != nil {
		return nil, err
	}
	page, limit := utils.Paginate(params.Page, params.Limit, 0)
	logs, totalRecords, err := logger.repo.GetLogs(ctx, filters, page, limit)
	if err != nil {
		return nil, err
	}
	metaData := utils.MetaDataInfo(&models.MetaData{
		Total:       totalRecords,
		PerPage:     limit,
		CurrentPage: page,
	})

	resp := entities.Response{
		MetaData: metaData,
		Data:     logs,
	}
	return &resp, nil
}

/*
BuildFilterQuery constructs a MongoDB query filter based on the provided LogParams.
It takes LogParams as input and returns a BSON M (MongoDB filter) that can be used
to filter log entries in a MongoDB collection.
*/
func (logger *LoggerUseCases) BuildFilterQuery(ctx context.Context, params entities.LogParams) (bson.M, error) {
	// Initialize variables for parsed start and end dates and an error variable.
	var (
		parsedStartDate, parsedEndDate time.Time
		err                            error
	)

	// Create an empty BSON M to store the query filters.
	filters := bson.M{}

	// Check if the start date parameter is not empty.
	if !utils.IsEmpty(params.StartDate) {
		// Parse and validate the start date.
		params.StartDate, parsedStartDate, err = utils.ConvertDate(params.StartDate)
		if err != nil {
			return nil, fmt.Errorf("invalid date format for start date: %w", err)
		}
	}

	// Check if the end date parameter is not empty.
	if !utils.IsEmpty(params.EndDate) {
		// Parse and validate the end date.
		params.EndDate, parsedEndDate, err = utils.ConvertDate(params.EndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid date format for end date: %w", err)
		}
	}

	// Check if the UserIP parameter is not empty and validate its format.
	if !utils.IsEmpty(params.UserIP) {
		if !utils.IsValidIPAddress(params.UserIP) {
			return nil, fmt.Errorf("invalid IP address format")
		}
		// Add the UserIP filter to the BSON M.
		filters[constants.ContextRequestIP] = params.UserIP
	}

	// Check if the HTTPMethod parameter is not empty and convert it to uppercase.
	if !utils.IsEmpty(params.HTTPMethod) {
		// Add the HTTPMethod filter to the BSON M.
		filters[constants.ContextRequestMethod] = strings.ToUpper(params.HTTPMethod)
	}

	// Handle cases based on the presence of StartDate and EndDate parameters.
	switch {
	case !utils.IsEmpty(params.StartDate) && !utils.IsEmpty(params.EndDate):
		if parsedStartDate.After(parsedEndDate) {
			return nil, fmt.Errorf("start date cannot be greater than end date")
		}
		params.EndDate = utils.AddHours(parsedEndDate, consts.HoursInADay)
		filters[constants.ContextTimeStamp] = bson.M{"$gte": parsedStartDate, "$lt": parsedEndDate}
	case !utils.IsEmpty(params.StartDate):
		filters[constants.ContextTimeStamp] = bson.M{"$gte": parsedStartDate}
	case !utils.IsEmpty(params.EndDate):
		params.EndDate = utils.AddHours(parsedEndDate, consts.HoursInADay)
		filters[constants.ContextTimeStamp] = bson.M{"$lt": parsedEndDate}
	}

	// Check if the Service parameter is not empty and convert it to lowercase.
	if !utils.IsEmpty(params.Service) {
		// Add the Service filter to the BSON M.
		filters[constants.ContextService] = strings.ToLower(params.Service)
	}

	// Check if the LogLevel parameter is not empty and convert it to lowercase.
	if !utils.IsEmpty(params.LogLevel) {
		// Add the LogLevel filter to the BSON M.
		filters[constants.ContextLogLevel] = strings.ToLower(params.LogLevel)
	}

	// Check if the Endpoint parameter is not empty.
	if !utils.IsEmpty(params.Endpoint) {
		// Add the Endpoint filter to the BSON M.
		filters[constants.ContextRequestUriTemplate] = params.Endpoint
	}

	// Check if the UserID parameter is not empty.
	if !utils.IsEmpty(params.UserID) {
		// Add the UserID filter to the BSON M.
		filters[consts.UserID] = params.UserID
	}

	// Check if the PartnerID parameter is not empty.
	if !utils.IsEmpty(params.PartnerID) {
		// Add the PartnerID filter to the BSON M.
		filters[consts.PartnerID] = params.PartnerID
	}

	// Return the constructed filter and no error.
	return filters, nil
}
