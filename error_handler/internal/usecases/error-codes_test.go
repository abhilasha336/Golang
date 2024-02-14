package usecases_test

// import (
// 	"context"
// 	"errors"
// 	"localization/internal/entities/db"
// 	"localization/internal/repo/mock"
// 	"localization/internal/usecases"
// 	"testing"

// 	"github.com/golang/mock/gomock"
// 	"github.com/stretchr/testify/assert"
// 	// "your-package-path/errorcodes" // Update with the correct package path
// )

// func TestAppendError(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockRepo := mock.NewMockErrorCodesRepoImply(ctrl)
// 	errorCodesUseCase := usecases.NewErrorCodesUseCases(mockRepo)

// 	// Define test data
// 	errType := "testError"
// 	endpoint := "testEndpoint"
// 	lang := "en"
// 	method := "GET"

// 	// Mock successful case
// 	mockRepo.EXPECT().AppendError(gomock.Any(), errType, endpoint, lang, gomock.Any(), method).Return(nil)

// 	err := errorCodesUseCase.AppendError(context.Background(), errType, endpoint, lang, db.ErrorData{}, method)
// 	assert.NoError(t, err)

// 	// Mock an error case
// 	mockRepo.EXPECT().AppendError(gomock.Any(), errType, endpoint, lang, gomock.Any(), method).Return(errors.New("mock error"))

// 	err = errorCodesUseCase.AppendError(context.Background(), errType, endpoint, lang, db.ErrorData{}, method)
// 	assert.Error(t, err)
// }
