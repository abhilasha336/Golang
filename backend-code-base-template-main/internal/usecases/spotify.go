package usecases

import "backend-code-base-template/internal/repo"

type OauthSpotifyUseCase struct {
	useCase repo.OauthSpotifyRepoImply
}

type OuathSpotifyUsecaseImply interface {
	GetOauthCredentials()
}

func NewOauthSpotifyUseCase(oauth repo.OauthSpotifyRepoImply) OuathSpotifyUsecaseImply {
	return &OauthSpotifyUseCase{
		useCase: oauth,
	}
}

func (oauth *OauthSpotifyUseCase) GetOauthCredentials() {

}
