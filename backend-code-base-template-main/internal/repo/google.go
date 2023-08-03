package repo

import "database/sql"

type OauthGoogleRepo struct {
	repo *sql.DB
}

type OauthGoogleRepoImply interface {
	GetOauthCredentials()
}

func NewOauthGoogleRepo(repo *sql.DB) OauthGoogleRepoImply {
	return &OauthGoogleRepo{
		repo: repo,
	}
}

func (oauth *OauthGoogleRepo) GetOauthCredentials() {
	// query part comes here
}
