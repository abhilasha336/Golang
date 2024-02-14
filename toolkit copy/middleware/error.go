package middleware

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gitlab.com/tuneverse/toolkit/consts"
	"gitlab.com/tuneverse/toolkit/models"
	"gitlab.com/tuneverse/toolkit/utils"

	"github.com/go-playground/validator/v10"
)

var (
	errorUnableToreadResponse = "unable to read data from localisation response"
)

type cache interface {
	Get(k string) (interface{}, bool)
	Set(k string, x interface{}, d time.Duration)
}

type UnimplementedCache struct{}

func (uimp *UnimplementedCache) Get(k string) (interface{}, bool) {
	panic("Get method id not implemented")
}

func (uimp *UnimplementedCache) Set(k string, x interface{}, d time.Duration) {
	panic("Set method id not implemented")
}

type ErrorLocaleOptions struct {
	Cache                  cache         `validate:"required"`
	CacheExpiration        time.Duration `validate:"required"`
	CacheKeyLabel          string        `validate:"required"`
	ContextErrorResponse   string
	LocalisationServiceURL string `validate:"required"`
	HeaderLanguage         string
}

type EndPointOptions struct {
	Cache            cache         `validate:"required"`
	CacheExpiration  time.Duration `validate:"required"`
	CacheKeyLabel    string        `validate:"required"`
	ContextEndPoints string
	EndPointsURL     string `validate:"required"`
	HeaderLanguage   string
}

// ErrorLocalization
func ErrorLocalization(option ...ErrorLocaleOptions) gin.HandlerFunc {
	if len(option) <= 0 {
		log.Fatal("please provide the error localization options")
	}

	opt := option[0]

	// Create a new validator instance
	validate := validator.New()

	// Validate the user struct
	err := validate.Struct(opt)
	if err != nil {
		log.Fatalf("localization options validation failed : %v", err)
	}

	return func(c *gin.Context) {

		cacheKeyLabel := consts.CacheErrorData
		if opt.CacheKeyLabel != "" {
			cacheKeyLabel = opt.CacheKeyLabel
		}

		contextErrorResponse := consts.ContextErrorResponses
		if opt.ContextErrorResponse != "" {
			contextErrorResponse = opt.ContextErrorResponse
		}

		localeLang := consts.ContextLocaleLang
		if opt.HeaderLanguage != "" {
			localeLang = opt.HeaderLanguage
		}

		//for storing error response data from context
		var errorData = make(map[string]interface{})
		cacheData, isFound := opt.Cache.Get(cacheKeyLabel)
		//check whether data is in cache or not
		if isFound {
			log.Infof("found data in cache")
			c.Set(contextErrorResponse, cacheData)
		} else {
			log.Infof("[ErrorLocalization] Cache not found, calling downstream API")

			language, ok := utils.GetContext[string](c, localeLang)
			if !ok {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "Unable to find the localization language",
				})
				return
			}

			headers := map[string]interface{}{
				consts.ContextLocaleLang: language,
			}

			resp, err := utils.APIRequest(http.MethodGet, opt.LocalisationServiceURL, headers, nil)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "An unexpected error occured",
				})
				return
			}

			//checking the sCacheErrorDatatatus code
			if resp.StatusCode != http.StatusOK {
				log.Errorf("%v %v", errorUnableToreadResponse, err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": errorUnableToreadResponse,
				})
				return
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Errorf("%v %v", errorUnableToreadResponse, err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": fmt.Sprintf("%v %v", errorUnableToreadResponse, err),
				})
				return
			}
			defer resp.Body.Close()

			err = json.Unmarshal(body, &errorData)
			if err != nil {
				log.Errorf("%v %v", errorUnableToreadResponse, err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": errorUnableToreadResponse,
				})
				return
			}

			// set error data in context
			c.Set(contextErrorResponse, errorData)

			// set error data in cache with an expiration time
			opt.Cache.Set(consts.CacheErrorData, errorData, opt.CacheExpiration)
		}
		c.Next()
	}
}

func EndpointExtraction(option ...EndPointOptions) gin.HandlerFunc {
	if len(option) <= 0 {
		log.Fatal("please provide the error localization options")
	}

	opt := option[0]

	// Create a new validator instance
	validate := validator.New()

	// Validate the user struct
	err := validate.Struct(opt)
	if err != nil {
		log.Fatalf("endpoint options validation failed : %v", err)
	}

	return func(c *gin.Context) {

		cacheKeyLabel := consts.CacheEndPointData
		if opt.CacheKeyLabel != "" {
			cacheKeyLabel = opt.CacheKeyLabel
		}

		ContextEndPoints := consts.ContextEndPoints
		if opt.ContextEndPoints != "" {
			ContextEndPoints = opt.ContextEndPoints
		}

		localeLang := consts.ContextLocaleLang
		if opt.HeaderLanguage != "" {
			localeLang = opt.HeaderLanguage
		}
		var endpoints models.ResponseData

		cacheData, isFound := opt.Cache.Get(cacheKeyLabel)
		//check whether data is in cache or not
		if isFound {
			log.Infof("found data in cache")
			c.Set(ContextEndPoints, cacheData)
		} else {

			language, ok := utils.GetContext[string](c, localeLang)
			if !ok {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "Unable to find endpoints",
				})
				return
			}

			headers := map[string]interface{}{
				consts.ContextLocaleLang: language,
			}

			resp, err := utils.APIRequest(http.MethodGet, opt.EndPointsURL, headers, nil)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "An unexpected error occured",
				})
				return
			}

			//checking the status code
			if resp.StatusCode != http.StatusOK {
				log.Errorf("unable to read data from localisation response %v", err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "unable to read data from localisation response",
				})
				return
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Errorf("unable to read data from localisation response %v", err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": fmt.Sprintf("unable to read data from localisation response %v", err),
				})
				return
			}
			defer resp.Body.Close()

			err = json.Unmarshal(body, &endpoints)
			if err != nil {
				log.Errorf("unable to read data from localisation response %v", err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "unable to read data from localisation response",
				})
				return
			}

			// set error data in context
			c.Set(ContextEndPoints, endpoints)
			c.Next()

		}
	}
}
