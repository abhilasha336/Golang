package usecases

import (
	"context"
	"oauth/internal/entities"
	"oauth/internal/repo"
)

// OauthUseCase holds OauthRepoImply interface
type OauthUseCase struct {
	useCase   repo.OauthRepoImply
	oauthData entities.OAuthData
}

// OuathUsecaseImply which implements functions
type OuathUsecaseImply interface {
	GetOauthCredentials(context.Context, string, string) (entities.OAuthCredentials, error)
	PostRefreshToken(context.Context, entities.Refresh, string, string, *string) error
	DeleteAndInsertRefreshToken(ctx context.Context, oldToken, newToken, newRefresh, partnerID string, memberID *string) error
	Logout(context.Context, entities.Refresh, string, string, string) error
	GetProviderName(ctx context.Context, id string) (string, error)
	GetPartnerId(ctx context.Context, clientID, clientSecret string) (string, string, error)
}

// NewOauthUseCase function assign values to OauthUseCase
func NewOauthUseCase(oauth repo.OauthRepoImply, oauthData entities.OAuthData) OuathUsecaseImply {
	return &OauthUseCase{
		useCase:   oauth,
		oauthData: oauthData,
	}
}

// function helps to implement businuss logics with data
func (oauth *OauthUseCase) GetOauthCredentials(ctx context.Context, provider, partnerID string) (entities.OAuthCredentials, error) {
	return oauth.useCase.GetOauthCredentials(ctx, provider, partnerID)
}

func (oauth *OauthUseCase) PostRefreshToken(ctx context.Context, refreshToken entities.Refresh, accessToken string, partnerID string, memberID *string) error {
	return oauth.useCase.PostRefreshToken(ctx, refreshToken, accessToken, partnerID, memberID)
}
func (oauth *OauthUseCase) DeleteAndInsertRefreshToken(ctx context.Context, oldToken, newToken, newRefresh, partnerID string, memberID *string) error {
	return oauth.useCase.DeleteAndInsertRefreshToken(ctx, oldToken, newToken, newRefresh, partnerID, memberID)
}
func (oauth *OauthUseCase) Logout(ctx context.Context, refreshToken entities.Refresh, accessToken string, partnerID, memberID string) error {
	return oauth.useCase.Logout(ctx, refreshToken, accessToken, partnerID, memberID)
}

func (oauth *OauthUseCase) GetProviderName(ctx context.Context, id string) (string, error) {
	return oauth.useCase.GetProviderName(ctx, id)
}

func (oauth *OauthUseCase) GetPartnerId(ctx context.Context, clientID, clientSecret string) (string, string, error) {
	return oauth.useCase.GetPartnerId(ctx, clientID, clientSecret)
}
