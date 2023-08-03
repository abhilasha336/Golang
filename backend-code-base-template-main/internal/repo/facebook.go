package repo

import "database/sql"

type OauthFacebookRepo struct {
	repo *sql.DB
}

type OauthFacebookRepoImply interface {
	GetOauthCredentials()
}

func NewOauthFacebookRepo(repo *sql.DB) OauthFacebookRepoImply {
	return &OauthFacebookRepo{
		repo: repo,
	}
}

func (oauth *OauthFacebookRepo) GetOauthCredentials() {
	// query part comes here
}
