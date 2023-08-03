package repo

import "database/sql"

type OauthSpotifyRepo struct {
	repo *sql.DB
}

type OauthSpotifyRepoImply interface {
	GetOauthCredentials()
}

func NewOauthSpotifyRepo(repo *sql.DB) OauthSpotifyRepoImply {
	return &OauthSpotifyRepo{
		repo: repo,
	}
}

func (oauth *OauthSpotifyRepo) GetOauthCredentials() {
	// query part comes here
}
