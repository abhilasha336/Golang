package usecases

import "backend-code-base-template/internal/repo"

type OauthGooglUseCase struct {
	useCase repo.OauthGoogleRepoImply
}

type OuathGoogleUsecaseImply interface {
	GetOauthCredentials()
}

func NewOauthGoogleUseCase(oauth repo.OauthGoogleRepoImply) OuathGoogleUsecaseImply {
	return &OauthGooglUseCase{
		useCase: oauth,
	}
}

func (oauth *OauthGooglUseCase) GetOauthCredentials() {

}
