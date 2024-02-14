package usecases

import (
	"context"
	"database/sql"
	"net/http"
	"oauth/config"
	"oauth/internal/consts"
	"oauth/internal/entities"
	"oauth/internal/repo/mockdb"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"gitlab.com/tuneverse/toolkit/utils"
)

func TestGetOauthCredentials(t *testing.T) {

	cfg, err := config.LoadConfig(consts.AppName)
	if err != nil {
		panic(err)
	}

	account := randomAccount()
	accountWithValues := randomAccountWith()
	decryptedClientID, _ := utils.Decrypt(accountWithValues.ClientID, []byte(cfg.EncryptionKey))
	decryptedClientSecret, _ := utils.Decrypt(accountWithValues.ClientSecret, []byte(cfg.EncryptionKey))

	accountWithValues.ClientID = decryptedClientID
	accountWithValues.ClientSecret = decryptedClientSecret

	testCases := []struct {
		ProviderName  string
		PartnerID     string
		buildStubs    func(store *mockdb.MockOauthRepoImply)
		checkResponse func(t *testing.T, cred entities.OAuthCredentials, err error)
	}{
		{
			ProviderName: consts.GoogleProvider,

			buildStubs: func(store *mockdb.MockOauthRepoImply) {
				store.EXPECT().
					GetOauthCredentials(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(accountWithValues, nil)
			},
			checkResponse: func(t *testing.T, cred entities.OAuthCredentials, err error) {
				require.Nil(t, err)
				require.NotNil(t, cred)
			},
		}, {
			ProviderName: consts.FacebookProvider,

			buildStubs: func(store *mockdb.MockOauthRepoImply) {
				store.EXPECT().
					GetOauthCredentials(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(accountWithValues, nil)
			},
			checkResponse: func(t *testing.T, cred entities.OAuthCredentials, err error) {
				require.Nil(t, err)
				require.NotNil(t, cred)
			},
		},
		{
			ProviderName: "",
			buildStubs: func(store *mockdb.MockOauthRepoImply) {
				store.EXPECT().
					GetOauthCredentials(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(account, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, cred entities.OAuthCredentials, err error) {
				require.NotNil(t, http.StatusNotFound, err)
				require.Equal(t, cred, account)

			},
		},
	}
	for _, tc := range testCases {

		t.Run(tc.ProviderName, func(t *testing.T) {

			ctrl := gomock.NewController(t)

			store := mockdb.NewMockOauthRepoImply(ctrl)
			defer ctrl.Finish()

			storeUseCase := NewOauthUseCase(store, entities.OAuthData{})

			tc.buildStubs(store)

			output, err := storeUseCase.GetOauthCredentials(context.Background(), tc.ProviderName, tc.PartnerID)
			tc.checkResponse(t, output, err)

		})
	}

}

func TestGetLogOut(t *testing.T) {

	testCases := []struct {
		Name          string
		RefreshToken  entities.Refresh
		AccessToken   string
		PartnerID     string
		MemberID      string
		buildStubs    func(store *mockdb.MockOauthRepoImply)
		checkResponse func(t *testing.T, refreshToken entities.Refresh, accessToken, partnerID, memberIDstring string, err error)
	}{
		{
			Name: "logoutok",
			RefreshToken: entities.Refresh{
				RefreshToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiIiLCJQYXJ0bmVySUQiOiIiLCJQYXJ0bmVyTmFtZSI6IiIsIlVzZXJUeXBlIjoiIiwiUm9sZXMiOltdLCJVc2VyRW1haWwiOiIiLCJleHAiOjE2OTc3OTQ1NjJ9.isg4Fpnuf7J7gLFdnIRK5n-PfsqCMStd-kX5IkMtiqc",
			},
			AccessToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiIiLCJQYXJ0bmVySUQiOiIiLCJQYXJ0bmVyTmFtZSI6IiIsIlVzZXJUeXBlIjoiIiwiUm9sZXMiOltdLCJVc2VyRW1haWwiOiIiLCJleHAiOjE2OTc3MTE3NjJ9.EBaphXVE6ZnZ8xJMZubZZfrlefwql1vI0T-sqTVx3LQ",
			PartnerID:   "614608f2-6538-4733-aded-96f902007254",
			MemberID:    "1f29a442-0f64-455a-a557-7b792713de80",
			buildStubs: func(store *mockdb.MockOauthRepoImply) {
				store.EXPECT().
					Logout(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, refreshToken entities.Refresh, accessToken, partnerID, memberIDstring string, err error) {
				require.Nil(t, err)
			},
		},
	}
	for _, tc := range testCases {

		t.Run(tc.Name, func(t *testing.T) {

			ctrl := gomock.NewController(t)

			store := mockdb.NewMockOauthRepoImply(ctrl)
			defer ctrl.Finish()

			storeUseCase := NewOauthUseCase(store, entities.OAuthData{})

			tc.buildStubs(store)

			err := storeUseCase.Logout(context.Background(), tc.RefreshToken, tc.AccessToken, tc.PartnerID, tc.MemberID)
			tc.checkResponse(t, tc.RefreshToken, tc.AccessToken, tc.PartnerID, tc.MemberID, err)

		})
	}
}

func TestPostRefreshToken(t *testing.T) {

	testCases := []struct {
		Name          string
		RefreshToken  entities.Refresh
		AccessToken   string
		PartnerID     string
		MemberID      string
		buildStubs    func(store *mockdb.MockOauthRepoImply)
		checkResponse func(t *testing.T, refreshToken entities.Refresh, accessToken, partnerID string, memberID *string, err error)
	}{
		{
			Name: "RefreshTokenPostok",
			RefreshToken: entities.Refresh{
				RefreshToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiIiLCJQYXJ0bmVySUQiOiIiLCJQYXJ0bmVyTmFtZSI6IiIsIlVzZXJUeXBlIjoiIiwiUm9sZXMiOltdLCJVc2VyRW1haWwiOiIiLCJleHAiOjE2OTc3OTQ1NjJ9.isg4Fpnuf7J7gLFdnIRK5n-PfsqCMStd-kX5IkMtiqc",
			},
			MemberID:    "614608f2-6538-4733-aded-96f902007254",
			AccessToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiIiLCJQYXJ0bmVySUQiOiIiLCJQYXJ0bmVyTmFtZSI6IiIsIlVzZXJUeXBlIjoiIiwiUm9sZXMiOltdLCJVc2VyRW1haWwiOiIiLCJleHAiOjE2OTc3MTE3NjJ9.EBaphXVE6ZnZ8xJMZubZZfrlefwql1vI0T-sqTVx3LQ",
			PartnerID:   "614608f2-6538-4733-aded-96f902007254",
			buildStubs: func(store *mockdb.MockOauthRepoImply) {
				store.EXPECT().
					PostRefreshToken(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, refreshToken entities.Refresh, accessToken, partnerID string, memberID *string, err error) {
				require.Nil(t, err)
			},
		},
	}
	for _, tc := range testCases {

		t.Run(tc.Name, func(t *testing.T) {

			ctrl := gomock.NewController(t)

			store := mockdb.NewMockOauthRepoImply(ctrl)
			defer ctrl.Finish()

			storeUseCase := NewOauthUseCase(store, entities.OAuthData{})

			tc.buildStubs(store)

			err := storeUseCase.PostRefreshToken(context.Background(), tc.RefreshToken, tc.AccessToken, tc.PartnerID, &tc.MemberID)
			tc.checkResponse(t, tc.RefreshToken, tc.AccessToken, tc.PartnerID, &tc.MemberID, err)

		})
	}
}

func TestDeleteAndInsertRefreshToken(t *testing.T) {

	testCases := []struct {
		Name            string
		RefreshToken    entities.Refresh
		NewToken        string
		NewRefreshToken string
		AccessToken     string
		PartnerID       string
		MemberID        string
		buildStubs      func(store *mockdb.MockOauthRepoImply)
		checkResponse   func(t *testing.T, refereshToken entities.Refresh, NewToken, NewRefreshToken, PartnerID string, memberID *string, err error)
	}{
		{
			Name: "DeleteInsertRefreshToken",
			RefreshToken: entities.Refresh{
				RefreshToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiIiLCJQYXJ0bmVySUQiOiIiLCJQYXJ0bmVyTmFtZSI6IiIsIlVzZXJUeXBlIjoiIiwiUm9sZXMiOltdLCJVc2VyRW1haWwiOiIiLCJleHAiOjE2OTc3OTQ1NjJ9.isg4Fpnuf7J7gLFdnIRK5n-PfsqCMStd-kX5IkMtiqc",
			},
			AccessToken:     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiIiLCJQYXJ0bmVySUQiOiIiLCJQYXJ0bmVyTmFtZSI6IiIsIlVzZXJUeXBlIjoiIiwiUm9sZXMiOltdLCJVc2VyRW1haWwiOiIiLCJleHAiOjE2OTc3MTE3NjJ9.EBaphXVE6ZnZ8xJMZubZZfrlefwql1vI0T-sqTVx3LQ",
			PartnerID:       "614608f2-6538-4733-aded-96f902007254",
			MemberID:        "1f29a442-0f64-455a-a557-7b792713de80",
			NewToken:        "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiIxZjI5YTQ0Mi0wZjY0LTQ1NWEtYTU1Ny03Yjc5MjcxM2RlODAiLCJQYXJ0bmVySUQiOiJhc2Rhc2RzYWRzZHMiLCJQYXJ0bmVyTmFtZSI6IiIsIlVzZXJUeXBlIjoiYWRtaW4iLCJSb2xlcyI6WyJtYW5hZ2VyIiwiYWRtaW4iXSwiVXNlckVtYWlsIjoiYXZsYXNoYWJoaTMzNkBnbWFpbC5jb20iLCJleHAiOjE2OTc3OTQ2NTV9.w2E77ayx7BqtqNLXGufXUV-EadyySN9dSYRnzNycUak",
			NewRefreshToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiIxZjI5YTQ0Mi0wZjY0LTQ1NWEtYTU1Ny03Yjc5MjcxM2RlODAiLCJQYXJ0bmVySUQiOiJhc2Rhc2RzYWRzZHMiLCJQYXJ0bmVyTmFtZSI6IiIsIlVzZXJUeXBlIjoiYWRtaW4iLCJSb2xlcyI6WyJtYW5hZ2VyIiwiYWRtaW4iXSwiVXNlckVtYWlsIjoiYXZsYXNoYWJoaTMzNkBnbWFpbC5jb20iLCJleHAiOjE2OTc3OTQ2Mjd9.bGvuV552M32a7N1Hn_vutEl7knMNJZWb707r_a_ds9Y",
			buildStubs: func(store *mockdb.MockOauthRepoImply) {
				store.EXPECT().
					DeleteAndInsertRefreshToken(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, refreshToken entities.Refresh, NewToken, NewRefreshToken, PartnerID string, MemberID *string, err error) {
				require.Nil(t, err)
			},
		},
	}
	for _, tc := range testCases {

		t.Run(tc.Name, func(t *testing.T) {

			ctrl := gomock.NewController(t)

			store := mockdb.NewMockOauthRepoImply(ctrl)
			defer ctrl.Finish()

			storeUseCase := NewOauthUseCase(store, entities.OAuthData{})

			tc.buildStubs(store)

			err := storeUseCase.DeleteAndInsertRefreshToken(context.Background(), tc.RefreshToken.RefreshToken, tc.NewToken, tc.NewRefreshToken, tc.PartnerID, &tc.MemberID)
			tc.checkResponse(t, tc.RefreshToken, tc.NewToken, tc.NewRefreshToken, tc.PartnerID, &tc.MemberID, err)

		})
	}
}

func TestGetProviderName(t *testing.T) {

	testCases := []struct {
		Name          string
		Id            string
		buildStubs    func(store *mockdb.MockOauthRepoImply)
		checkResponse func(t *testing.T, id string, err error)
	}{
		{
			Name: "internal",
			Id:   "542b9379-b418-404e-8cec-a90ea265cb2e",
			buildStubs: func(store *mockdb.MockOauthRepoImply) {
				store.EXPECT().
					GetProviderName(gomock.Any(), gomock.Any()).
					Times(1).
					Return("internal", nil)
			},
			checkResponse: func(t *testing.T, name string, err error) {
				require.Nil(t, err)
				require.Equal(t, "internal", name)

			},
		},
		{
			Name: "Internal",
			Id:   "542b9379-b418-404e-8cec-a90ea265cb2e",
			buildStubs: func(store *mockdb.MockOauthRepoImply) {
				store.EXPECT().
					GetProviderName(gomock.Any(), gomock.Any()).
					Times(1).
					Return("internal", nil)
			},
			checkResponse: func(t *testing.T, name string, err error) {
				require.Nil(t, err)
				require.Equal(t, "internal", name)

			},
		},
		{
			Name: "Google",
			Id:   "542b9379-b418-404e-8cec-a90ea265cb2e",
			buildStubs: func(store *mockdb.MockOauthRepoImply) {
				store.EXPECT().
					GetProviderName(gomock.Any(), gomock.Any()).
					Times(1).
					Return("google", nil)
			},
			checkResponse: func(t *testing.T, name string, err error) {
				require.Nil(t, err)
				require.Equal(t, "google", name)

			},
		},
	}
	for _, tc := range testCases {

		t.Run(tc.Name, func(t *testing.T) {

			ctrl := gomock.NewController(t)

			store := mockdb.NewMockOauthRepoImply(ctrl)
			defer ctrl.Finish()

			storeUseCase := NewOauthUseCase(store, entities.OAuthData{})

			tc.buildStubs(store)

			providerName, err := storeUseCase.GetProviderName(context.Background(), tc.Id)
			tc.checkResponse(t, providerName, err)
		})
	}
}

func TestGetGetPartnerId(t *testing.T) {

	testCases := []struct {
		Name          string
		Id            string
		Secret        string
		buildStubs    func(store *mockdb.MockOauthRepoImply)
		checkResponse func(t *testing.T, id, url string, err error)
	}{
		{
			Name:   "internal",
			Id:     "542b9379-b418-404e-8cec-a90ea265cb2e",
			Secret: "abc",
			buildStubs: func(store *mockdb.MockOauthRepoImply) {
				store.EXPECT().
					GetPartnerId(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return("542b9379-b418-404e-8cec-a90ea265cb2e", "abc.com", nil)
			},
			checkResponse: func(t *testing.T, id, url string, err error) {
				require.Nil(t, err)
				require.Equal(t, "542b9379-b418-404e-8cec-a90ea265cb2e", id)
				require.Equal(t, "abc.com", url)

			},
		},
	}
	for _, tc := range testCases {

		t.Run(tc.Name, func(t *testing.T) {

			ctrl := gomock.NewController(t)

			store := mockdb.NewMockOauthRepoImply(ctrl)
			defer ctrl.Finish()

			storeUseCase := NewOauthUseCase(store, entities.OAuthData{})

			tc.buildStubs(store)

			id, sec, err := storeUseCase.GetPartnerId(context.Background(), tc.Id, tc.Secret)
			tc.checkResponse(t, id, sec, err)
		})
	}
}

// Import necessary packages and modules

func randomAccount() entities.OAuthCredentials {
	return entities.OAuthCredentials{}
}

func randomAccountWith() entities.OAuthCredentials {
	return entities.OAuthCredentials{
		ClientID:     "OUaGvGXjkZSC0T3vYrn85OHmyAEUaaGiBH4/2wLMsbqqnrGBjc4InLy5FhHEC5G4ocOEbJLyph/bH0Bb3IKUTFZ9jaHkFfh+oG+9VRkdYucOJyu7WkDSOQ==",
		ClientSecret: "NsiqnuwX9IheDy6hzQ2M4RIFX7BRubzD/qZlowvd+L0eKXyvNNVCyTaPSwUO2DLlsiK6",
		RedirectURL:  "http://localhost:8080/api/v1.0/callback",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		TokenURL:     "https://www.googleapis.com/oauth2/v3/userinfo?access_token=",
	}
}
