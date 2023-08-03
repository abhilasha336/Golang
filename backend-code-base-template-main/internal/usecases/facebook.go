package usecases

import "backend-code-base-template/internal/repo"

type OauthFacebookUseCase struct {
	useCase repo.OauthFacebookRepoImply
}

type OuathFacebookUsecaseImply interface {
	GetOauthCredentials()
}

func NewOauthFacebookUseCase(oauth repo.OauthFacebookRepoImply) OuathGoogleUsecaseImply {
	return &OauthFacebookUseCase{
		useCase: oauth,
	}
}

func (oauth *OauthFacebookUseCase) GetOauthCredentials() {

}
