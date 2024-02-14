package middlewares

import (
	"database/sql"
	"member/internal/consts"
	"member/internal/entities"

	"github.com/gin-gonic/gin"
)

// Middlewares structure for storing middleware values
type Middlewares struct {
	Cfg  *entities.EnvConfig
	Repo *sql.DB
}

// NewMiddlewares
func NewMiddlewares(cfg *entities.EnvConfig, repo *sql.DB) *Middlewares {
	return &Middlewares{
		Repo: repo,
		Cfg:  cfg,
	}
}

// Middleware function to get Partner ID from the header and store it in the context
func (m Middlewares) PartnerID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract partner ID from the header.
		partnerID := c.GetHeader("partner_id")
		// Set partner ID in the context.
		c.Set(consts.ContextPartnerID, partnerID)
		c.Next()
	}
}

// func (m Middlewares) JwtMiddleware() gin.HandlerFunc {
// 	return func(ctx *gin.Context) {

// 		skipUrl := "/api/v1.0/members"
// 		path := ctx.Request.URL.Path
// 		method := ctx.Request.Method

// 		// Check if the path and method match the conditions to skip middleware
// 		if path == skipUrl && (method == "POST" || method == "GET") {
// 			ctx.Next()
// 			return
// 		}
// 		token := ctx.Request.Header.Get("Authorization")
// 		isValid := utilities.ValidateJwtToken(token, m.Cfg.JwtKey)

// 		_ = utilities.ValidateJwtToken(token, "jwtkey")

// 		repo := repo.NewMemberRepo(m.Repo)
// 		_, err := repo.Middleware(context.Background(), token)
// 		if isValid.Valid && err == nil {
// 			ctx.Set("memberId", isValid.MemberID)

// 			ctx.Set("partner_id", isValid.PartnerID)
// 			ctx.Set("roles", isValid.Roles)
// 			ctx.Set("memberType", isValid.MemberType)
// 			ctx.Set("partnerName", isValid.PartnerName)
// 			ctx.Set("email", isValid.MemberEmail)

// 			ctx.Next()

// 		} // else {
// 		// 	ctx.JSON(http.StatusUnauthorized, gin.H{"message": "invalid"})
// 		// 	return
// 		// }
// 	}
// }
