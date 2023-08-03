package middlewares

// import (
// 	"backend-code-base-template/config"
// 	"backend-code-base-template/internal/consts"
// 	"backend-code-base-template/utilities"
// 	"net/http"
// 	"strings"

// 	"github.com/gin-gonic/gin"
// )

// type Middlewares struct {
// }

// // NewMiddlewares
// func NewMiddlewares() *Middlewares {
// 	return &Middlewares{}
// }

// // Middleware function to check Accept-version from API Header
// func (m Middlewares) ApiVersioning() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		version := c.Param("version")
// 		// if version == "" {
// 		// 	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Missing version parameter"})
// 		// 	return
// 		// }

// 		apiVersion := utilities.PrepareVersionName(version)
// 		apiVersion = strings.ToUpper(apiVersion)

// 		// set the accepting version in the context
// 		c.Set(consts.AcceptedVersions, apiVersion)

// 		// init the system Accepted versions
// 		// init the env config
// 		cfg, err := config.LoadConfig(consts.AppName)
// 		if err != nil {
// 			panic(err)
// 		}

// 		// set the list of system accepting version in the context
// 		systemAcceptedVersionsList := cfg.AcceptedVersions
// 		c.Set(consts.ContextSystemAcceptedVersions, systemAcceptedVersionsList)

// 		// check the version exists in the accepted list
// 		// find index of version from Accepted versions
// 		var found bool
// 		for index, version := range systemAcceptedVersionsList {
// 			version = strings.ToUpper(version)
// 			if version == apiVersion {
// 				found = true
// 				c.Set(consts.ContextAcceptedVersionIndex, index)
// 			}

// 		}
// 		if !found {
// 			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Given version is not supported by the system"})
// 			return
// 		}

// 		c.Next()
// 	}
// }
