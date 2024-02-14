package usecases

import (
	"context"

	"localization/internal/entities"
	"localization/internal/entities/db"
	"localization/internal/repo"
)

// ErrorCodesUseCases contains the use cases for error codes.
type ErrorCodesUseCases struct {
	repo repo.ErrorCodesRepoImply
}

// ErrorCodesUseCaseImply defines the interface for error code use cases.
type ErrorCodesUseCaseImply interface {
	GetError(string, string, string, string, string) (any, error)
	AppendError(context.Context, string, string, string, db.ErrorData, string) error
	AddTranslation(context.Context, []string) (any, error)
	AddEndpoint(context.Context, entities.RequestData) error
	GetEndpointName(context.Context) (entities.ResponseData, error)
}

// NewErrorCodesUseCases initializes a new ErrorCodesUseCases instance.
func NewErrorCodesUseCases(userRepo repo.ErrorCodesRepoImply) ErrorCodesUseCaseImply {
	return &ErrorCodesUseCases{
		repo: userRepo,
	}
}

// AppendError appends an error to the repository.
func (user *ErrorCodesUseCases) AppendError(ctx context.Context, errType string, endpoint string, lang string, code db.ErrorData, method string) error {
	return user.repo.AppendError(ctx, errType, endpoint, lang, code, method)
}

// GetError retrieves an error based on provided parameters.// GetError retrieves an error based on provided parameters.
func (user *ErrorCodesUseCases) GetError(errType string, endpoint string, lang string, field string, method string) (any, error) {
	resp, err := user.repo.GetError(errType, endpoint, lang, field, method)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// AddTranslation adds a translation for a language.
func (user *ErrorCodesUseCases) AddTranslation(ctx context.Context, language []string) (any, error) {
	output, err := user.repo.AddTranslation(ctx, language)
	if err != nil {
		return output, err
	}
	return output, nil
}

// AddEndpoint adds an endpoint to the repository.
func (user *ErrorCodesUseCases) AddEndpoint(ctx context.Context, endpoint entities.RequestData) error {
	err := user.repo.AddEndpoint(ctx, endpoint)
	if err != nil {
		return err
	}
	return nil
}

// GetEndpointName retrieves endpoint names.
func (user *ErrorCodesUseCases) GetEndpointName(ctx context.Context) (entities.ResponseData, error) {
	results, err := user.repo.GetEndpointName(ctx)
	if err != nil {
		return entities.ResponseData{}, err
	}
	return results, nil
}
