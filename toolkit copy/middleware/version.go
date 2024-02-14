package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/tuneverse/toolkit/consts"
	"gitlab.com/tuneverse/toolkit/utils"
)

type optionVersionLookup func(*gin.Context) string

type VersionOptions struct {
	VersionParamLookup optionVersionLookup
	AcceptedVersions   []string
}

// APIVersionGuard
// Middleware function to check Accept-version from API Header
func APIVersionGuard(option VersionOptions) gin.HandlerFunc {
	return func(c *gin.Context) {
		var version string

		if option.VersionParamLookup != nil {
			version = option.VersionParamLookup(c)
		} else {
			version = c.Param("version")
		}

		if version == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Missing version parameter"})
			return
		}

		// get and prepare the version name
		apiVersion := utils.PrepareVersionName(version)
		apiVersion = strings.ToUpper(apiVersion)

		var formattedVersions []string

		for _, version := range option.AcceptedVersions {
			formattedVersion := utils.PrepareVersionName(version)
			formattedVersions = append(formattedVersions, formattedVersion)
		}

		// set the list of system accepting version in the context
		systemAcceptedVersionsList := formattedVersions
		c.Set(consts.ContextSystemAcceptedVersions, systemAcceptedVersionsList)

		// check the version exists in the accepted list
		// find index of version from Accepted versions
		var found bool
		for index, version := range systemAcceptedVersionsList {
			version = strings.ToUpper(version)
			if version == apiVersion {
				found = true
				c.Set(consts.ContextAcceptedVersionIndex, index)
			}

		}
		if !found {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Given version is not supported by the system"})
			return
		}

		c.Next()
	}
}
