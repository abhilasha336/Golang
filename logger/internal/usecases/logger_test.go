package usecases

import (
	"context"
	"logger/internal/entities"
	"logger/internal/repo/mock"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gitlab.com/tuneverse/toolkit/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestAddLog(t *testing.T) {
	userID, partnerID := uuid.New().String(), uuid.New().String()
	logEntry := entities.Log{
		Method:       "POST",
		Endpoint:     "/api/resource",
		Service:      "my-service",
		UserIP:       "192.168.1.100",
		Filename:     "example.txt",
		Message:      "Request processed successfully",
		ResponseCode: 200,
		DurationMS:   "50",
		Day:          time.Now().Day(),
		Month:        int(time.Now().Month()),
		Year:         time.Now().Year(),
		Timestamp:    time.Now().Round(1 * time.Second),
		LogLevel:     "INFO",
		URL:          "https://example.com/api/resource",
		RequestDump:  utils.RequestDataDump{},
		ResponseDump: utils.ResponseData{},
		Data: entities.Data{
			UserID:    userID,
			PartnerID: partnerID,
		},
	}

	testCases := []struct {
		name          string
		log           entities.Log
		buildStubs    func(store *mock.MockLoggerRepoImply)
		checkResponse func(t *testing.T, err error)
	}{
		{
			name: "ok",
			log: entities.Log{
				Method:       "POST",
				Endpoint:     "/api/resource",
				Service:      "my-service",
				UserIP:       "192.168.1.100",
				Filename:     "example.txt",
				Message:      "Request processed successfully",
				ResponseCode: 200,
				DurationMS:   "50",
				Timestamp:    time.Now().Round(1 * time.Second),
				LogLevel:     "INFO",
				URL:          "https://example.com/api/resource",
				RequestDump:  utils.RequestDataDump{},
				ResponseDump: utils.ResponseData{},
				Data: entities.Data{
					UserID:    userID,
					PartnerID: partnerID,
				},
			},
			buildStubs: func(store *mock.MockLoggerRepoImply) {
				store.EXPECT().
					AddLog(gomock.Any(), gomock.Eq(logEntry)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, err error) {
				require.Nil(t, err)
			},
		},
		{
			name: "Internal server error",
			log:  logEntry,
			buildStubs: func(store *mock.MockLoggerRepoImply) {
				store.EXPECT().
					AddLog(gomock.Any(), gomock.Any()).
					Times(1).
					Return(mongo.ErrClientDisconnected)
			},
			checkResponse: func(t *testing.T, err error) {
				require.Error(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mock.NewMockLoggerRepoImply(ctrl)
			tc.buildStubs(store)
			loggerUseCases := NewLoggerUseCases(store)

			err := loggerUseCases.AddLog(context.Background(), tc.log)
			tc.checkResponse(t, err)
		})
	}
}

func TestGetLogs(t *testing.T) {

	var (
		logs []*entities.Log
		n    int = 5
	)
	userID, partnerID := uuid.New().String(), uuid.New().String()
	for i := 0; i < n; i++ {
		logs = append(logs, &entities.Log{
			Method:       "POST",
			Endpoint:     "/api/resource",
			Service:      "my-service",
			UserIP:       "192.168.1.100",
			Filename:     "example.txt",
			Message:      "Request processed successfully",
			ResponseCode: 200,
			DurationMS:   "50",
			Timestamp:    time.Now().Round(1 * time.Second),
			LogLevel:     "INFO",
			URL:          "https://example.com/api/resource",
			RequestDump:  utils.RequestDataDump{},
			ResponseDump: utils.ResponseData{},
			Data: entities.Data{
				UserID:    userID,
				PartnerID: partnerID,
			},
		})
	}

	testCases := []struct {
		name          string
		query         entities.LogParams
		buildStubs    func(store *mock.MockLoggerRepoImply)
		checkresponse func(t *testing.T, resp *entities.Response, err error)
	}{
		{
			name: "ok",
			query: entities.LogParams{
				Page:       1,
				Limit:      5,
				StartDate:  "21-11-2023",
				EndDate:    "21-11-2023",
				Service:    "utility",
				HTTPMethod: "GET",
				LogLevel:   "INFO",
				UserID:     "user123",
				PartnerID:  "partner123",
				Endpoint:   "/api/{:version}/resource/{:id}",
				UserIP:     "127.0.0.1",
			},
			buildStubs: func(store *mock.MockLoggerRepoImply) {
				store.EXPECT().
					GetLogs(gomock.Any(), gomock.Any(), gomock.Eq(int32(1)), gomock.Eq(int32(5))).
					Times(1).
					Return(logs, int64(len(logs)), nil)
			},
			checkresponse: func(t *testing.T, resp *entities.Response, err error) {
				require.Nil(t, err)
				require.Equal(t, int64(len(logs)), resp.MetaData.Total)
			},
		},
		{
			name: "invalid start date",
			query: entities.LogParams{
				StartDate: "invalid",
			},
			buildStubs: func(store *mock.MockLoggerRepoImply) {
				store.EXPECT().
					GetLogs(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkresponse: func(t *testing.T, resp *entities.Response, err error) {
				require.Error(t, err)
			},
		},
		{
			name: "invalid end date",
			query: entities.LogParams{
				EndDate: "invalid",
			},
			buildStubs: func(store *mock.MockLoggerRepoImply) {
				store.EXPECT().
					GetLogs(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkresponse: func(t *testing.T, resp *entities.Response, err error) {
				require.Error(t, err)
			},
		},
		{
			name: "internal server error",
			query: entities.LogParams{
				Page:  1,
				Limit: 5,
			},
			buildStubs: func(store *mock.MockLoggerRepoImply) {
				store.EXPECT().
					GetLogs(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).Return(nil, int64(0), mongo.ErrClientDisconnected)
			},
			checkresponse: func(t *testing.T, resp *entities.Response, err error) {
				require.Error(t, err)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mock.NewMockLoggerRepoImply(ctrl)
			tc.buildStubs(store)
			logger := NewLoggerUseCases(store)
			resp, err := logger.GetLogs(context.Background(), tc.query)
			tc.checkresponse(t, resp, err)
		})
	}
}
