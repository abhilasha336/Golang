package middlewares

import (
	"context"
	"database/sql"
	"net/http"
	"oauth/internal/entities"
	"oauth/internal/repo"
	"oauth/utilities"

	"github.com/gin-gonic/gin"
	log "gitlab.com/tuneverse/toolkit/core/logger"
)

type Middlewares struct {
	repo *sql.DB
	cfg  *entities.EnvConfig
}

// NewMiddlewares
func NewMiddlewares(cfg *entities.EnvConfig, repo *sql.DB) *Middlewares {
	return &Middlewares{
		repo: repo,
		cfg:  cfg,
	}
}

func (m *Middlewares) JwtMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log := log.Log().WithContext(ctx)
		skipUrl := "/api/v1.0/members"
		skipMethod := "POST"
		path := ctx.Request.URL.Path
		if path == skipUrl && ctx.Request.Method == skipMethod {
			ctx.Next()
			return
		}
		token := ctx.Request.Header.Get("Authorization")
		isValid := utilities.ValidateJwtToken(token, m.cfg.JwtKey)

		repo := repo.NewOauthRepo(m.repo, m.cfg)
		_, err := repo.Middleware(context.Background(), token)
		if isValid.Valid && err == nil {
			ctx.Set("memberId", isValid.MemberID)
			ctx.Set("partnerId", isValid.PartnerID)
			ctx.Set("roles", isValid.Roles)
			ctx.Set("memberType", isValid.MemberType)
			ctx.Set("partnerName", isValid.PartnerName)
			ctx.Set("email", isValid.MemberEmail)
			_, ok := ctx.Get("memberId")
			if ok {
				log.Printf("middleware context set failed")
			}
			ctx.Next()

		} else {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "invalid,unauthorised"})
			return
		}
	}
}
