package middleware

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/tuneverse/toolkit/consts"
)

type LocaleOptions struct {
	HeaderLabel  string
	ContextLabel string
}

// Localize
func Localize(options ...LocaleOptions) gin.HandlerFunc {
	return func(c *gin.Context) {
		headerLabel := consts.HeaderLanguage
		contextLabel := consts.ContextLocaleLang

		if len(options) > 0 {
			opt := options[0]
			if opt.HeaderLabel != "" {
				headerLabel = opt.HeaderLabel
			}

			if opt.ContextLabel != "" {
				contextLabel = opt.ContextLabel
			}
		}

		lan := c.Request.Header.Get(headerLabel)
		if lan == "" {
			lan = "en"
		}

		// setting the language
		c.Set(contextLabel, lan)
		c.Next()
	}
}
