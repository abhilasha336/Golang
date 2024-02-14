package middlewares

import (
	"database/sql"
	"net/http"

	"localization/internal/consts"
	"localization/internal/entities"
	"localization/utilities"

	"github.com/gin-gonic/gin"
)

// Middleware structure
type Middlewares struct {
	Cfg    *entities.EnvConfig
	psqldb *sql.DB
}

// NewMiddlewares function
func NewMiddlewares(cfg *entities.EnvConfig, psqldb *sql.DB) *Middlewares {
	return &Middlewares{
		Cfg:    cfg,
		psqldb: psqldb,
	}
}

// SetLanguageInContextMiddleware is a middleware function that sets the language in the Gin context.
func (m Middlewares) SetLanguageInContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve the language using the provided Postgres database connection.
		language, err := utilities.GetLanguage(m.psqldb)
		if err != nil {
			// If there's an error fetching the language, abort the request with an error response.
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "An unexpected error occured"})
			return
		}

		// Set the language slice in the Gin context for further use in the request.
		c.Set(consts.Language, language)
		// Proceed to the next middleware or route handler.
		c.Next()
	}
}
