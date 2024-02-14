package repo

import (
	"context"
	"database/sql"
	"errors"
	"oauth/internal/entities"

	log "gitlab.com/tuneverse/toolkit/core/logger"

	"gitlab.com/tuneverse/toolkit/utils"
)

// OauthRepo holds db and config
type OauthRepo struct {
	repo *sql.DB
	cfg  *entities.EnvConfig
}

// OauthRepoImply which implements functions
type OauthRepoImply interface {
	GetOauthCredentials(context.Context, string, string) (entities.OAuthCredentials, error)
	PostRefreshToken(context.Context, entities.Refresh, string, string, *string) error
	DeleteAndInsertRefreshToken(ctx context.Context, oldToken, newToken, newRefresh, partnerID string, memberID *string) error
	Logout(context.Context, entities.Refresh, string, string, string) error
	Middleware(context.Context, string) (string, error)
	GetProviderName(ctx context.Context, id string) (string, error)
	GetPartnerId(ctx context.Context, clientID, clientSecret string) (string, string, error)
}

// NewOauthRepo used to assign values to both database and config
func NewOauthRepo(repo *sql.DB, cfg *entities.EnvConfig) OauthRepoImply {
	return &OauthRepo{
		repo: repo,
		cfg:  cfg,
	}
}

// fn which retrieves partnerid with client id and client secret
func (oauth *OauthRepo) GetPartnerId(ctx context.Context, clientID, clientSecret string) (string, string, error) {

	var (
		partnerID, redirectUri string
		err                    error
		log                    = log.Log().WithContext(ctx)
	)

	GetCredentials := `
	SELECT 
		partner_id,redirect_uri 
	FROM partner_api_credential
	WHERE
	client_id=$1
	AND
	client_secret=$2
	`
	row := oauth.repo.QueryRowContext(
		ctx,
		GetCredentials,
		clientID,
		clientSecret,
	)

	err = row.Scan(
		&partnerID,
		&redirectUri,
	)
	if err != nil {
		log.Errorf("GetPartnerId- scan error:%v", err)
		return "", "", err
	}

	return partnerID, redirectUri, nil

}

// GetOauthCredentials function used to fetch oauth credentials from database
func (oauth *OauthRepo) GetOauthCredentials(ctx context.Context, provider, partnerID string) (entities.OAuthCredentials, error) {

	var (
		credential                               entities.OAuthCredentials
		decryptedClientID, decryptedClientSecret string
		log                                      = log.Log().WithContext(ctx)
		err                                      error
	)
	encryptionKey := oauth.cfg.EncryptionKey
	key := []byte(encryptionKey)

	GetCredentials := `
	SELECT 
		o1.name, 
		p1.client_id, 
		p1.client_secret, 
		p1.redirect_uri, 
		p1.scope, 
		p1.access_token_endpoint
	FROM partner_oauth_credential p1
	INNER JOIN oauth_provider o1 
	ON p1.oauth_provider_id = o1.id
	WHERE o1.name = $1
	AND p1.partner_id=$2
	`
	row := oauth.repo.QueryRowContext(
		ctx,
		GetCredentials,
		provider,
		partnerID,
	)

	err = row.Scan(
		&credential.ProviderName,
		&credential.ClientID,
		&credential.ClientSecret,
		&credential.RedirectURL,
		&credential.Scopes,
		&credential.TokenURL,
	)
	if err != nil {
		log.Errorf("GetOauthCredentials- scan error:%v", err)
		return credential, err
	}

	if decryptedClientID, err = utils.Decrypt(credential.ClientID, key); err != nil {
		log.Printf("GetOauthCredentials-decryptClientIderror: %v", err)
		return credential, err
	}
	if decryptedClientSecret, err = utils.Decrypt(credential.ClientSecret, key); err != nil {
		log.Printf("GetOauthCredentials-decryptClientSecreterror: %v", err)
		return credential, err
	}
	credential.ClientID = decryptedClientID
	credential.ClientSecret = decryptedClientSecret

	return credential, nil
}

// fn logout by changing status of _isrevoked field to true and add token in blacklisted field
func (oauth *OauthRepo) Logout(ctx context.Context, refreshToken entities.Refresh, accessToken, partnerID, memberID string) error {

	var (
		log = log.Log().WithContext(ctx)
	)

	query := `UPDATE refresh_token
	SET is_revoked = true, blacklisted = $1
	WHERE member_id = $2
	AND token=$3
	AND partner_id = $4`

	res, err := oauth.repo.ExecContext(ctx, query, accessToken, memberID, refreshToken.RefreshToken, partnerID)
	if err != nil {
		log.Errorf("LogOut-log out token update failed: %v", err)
		return err
	}
	// Check the number of rows affected by the update
	numRowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Errorf("LogOut-no rows affected: %v", err)
		return err
	}

	// Handle the result as needed
	if numRowsAffected == 0 {
		log.Errorf("LogOut-no rows affected")
		return errors.New("no rows updated")
	}
	return nil
}

// fn insert token in refreshtoken table
func (oauth *OauthRepo) PostRefreshToken(ctx context.Context, token entities.Refresh, active, partnerID string, memberID *string) error {

	var (
		log = log.Log().WithContext(ctx)
	)

	query := `INSERT INTO public.refresh_token (token, member_id, partner_id, active_token)
	VALUES ($1,$2,$3,$4)
	`
	res, err := oauth.repo.ExecContext(ctx, query, token.RefreshToken, &memberID, partnerID, active)
	if err != nil {
		log.Errorf("PostRefreshToken-token entry in db failed: %v", err)
		return err
	}
	// Check the number of rows affected by the update
	numRowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Errorf("PostRefreshToken-rows affected: %v", err)
		return err
	}

	// Handle the result as needed
	if numRowsAffected == 0 {
		log.Errorf("PostRefreshToken-no rows affected:")
		return errors.New("no rows updated")
	}
	return nil
}

// fn delete old token and insert new token in refresh token table
func (oauth *OauthRepo) DeleteAndInsertRefreshToken(ctx context.Context, oldToken, newToken, newRefresh, partnerID string, memberID *string) error {

	var (
		log = log.Log().WithContext(ctx)
	)

	// Start a transaction
	tx, err := oauth.repo.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			// Handle the rollback error, e.g., log it or take appropriate action
			log.Errorf("DeleteAndInsertRefreshToken-Error during rollback: %v", err)
		}
	}()

	// Step 1: Delete the old refresh token record
	deleteQuery := `DELETE FROM public.refresh_token WHERE token = $1 AND member_id = $2 AND partner_id = $3`
	_, err = tx.ExecContext(ctx, deleteQuery, oldToken, &memberID, partnerID)
	if err != nil {
		log.Errorf("DeleteAndInsertRefreshToken-Error during deletion:%v", err)
		return err
	}

	// Step 2: Insert the new refresh token record
	insertQuery := `INSERT INTO public.refresh_token (token, member_id, partner_id, active_token)
        VALUES ($1, $2, $3, $4)`
	_, err = tx.ExecContext(ctx, insertQuery, newRefresh, &memberID, partnerID, newToken)
	if err != nil {
		log.Errorf("DeleteAndInsertRefreshToken-Error during insertion:%v", err)
		return err
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		log.Errorf("DeleteAndInsertRefreshToken-Error during transaction commit:%v", err)
		return err
	}

	return nil
}

// fn to validate the existance of active token,used by middleware
func (oauth *OauthRepo) Middleware(ctx context.Context, token string) (string, error) {

	var (
		blacklistedToken string
		log              = log.Log().WithContext(ctx)
	)

	tokenBlacklisted := `
	SELECT  active_token
	FROM public.refresh_token
	WHERE is_revoked = false 
	AND active_token=$1;`

	row := oauth.repo.QueryRowContext(
		ctx,
		tokenBlacklisted,
		token,
	)

	err := row.Scan(&blacklistedToken)
	if err != nil {
		log.Printf("scan:%v", err)
		return blacklistedToken, err
	}

	return blacklistedToken, nil
}

// fn fetch provider name with provider table
func (oauth *OauthRepo) GetProviderName(ctx context.Context, id string) (string, error) {

	var (
		blacklistedToken string
		log              = log.Log().WithContext(ctx)
		data             entities.BasicMemberDataResponse
	)

	tokenBlacklisted := `
	select name from oauth_provider where id=$1;`

	row := oauth.repo.QueryRowContext(
		ctx,
		tokenBlacklisted,
		id,
	)

	err := row.Scan(&data.Data.Name)
	if err != nil {
		log.Printf("GetProviderName-scan:%v", err)
		return blacklistedToken, err
	}

	return data.Data.Name, nil
}
