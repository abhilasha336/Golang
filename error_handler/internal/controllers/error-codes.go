package controllers

import (
	"net/http"

	"localization/internal/consts"
	"localization/internal/entities"
	"localization/internal/entities/db"
	"localization/internal/middlewares"
	"localization/internal/usecases"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gitlab.com/tuneverse/toolkit/core/version"
	"gitlab.com/tuneverse/toolkit/utils"
)

// Options defines functional options for configuring ErrorCodesController.
type Options func(cntrl *ErrorCodesController)

// ErrorCodesController handles error codes and related functionality
type ErrorCodesController struct {
	router      *gin.RouterGroup
	useCases    usecases.ErrorCodesUseCaseImply
	middlewares *middlewares.Middlewares
}

// NewErrorCodesController creates a new ErrorCodesController.
func NewErrorCodesController(options ...Options) *ErrorCodesController {
	controll := &ErrorCodesController{}
	for _, opt := range options {
		opt(controll)
	}

	return controll
}

// WithRouteInit configures the router for ErrorCodesController.
func WithRouteInit(router *gin.RouterGroup) Options {
	return func(cntrl *ErrorCodesController) {
		cntrl.router = router
	}
}

// WithUseCaseInit configures the use cases for ErrorCodesController.
func WithUseCaseInit(usecases usecases.ErrorCodesUseCaseImply) Options {
	return func(cntrl *ErrorCodesController) {
		cntrl.useCases = usecases
	}
}

// WithMiddlewareInit configures the middlewares for ErrorCodesController.
func WithMiddlewareInit(middlewares *middlewares.Middlewares) Options {
	return func(cntrl *ErrorCodesController) {
		cntrl.middlewares = middlewares
	}
}

// InitRoutes initializes the routes for ErrorCodesController.
func (controller *ErrorCodesController) InitRoutes() {

	controller.router.GET("/:version/health", func(ctx *gin.Context) {
		version.RenderHandler(ctx, controller, "Health")
	})
	controller.router.PUT("/:version/localization/error", func(ctx *gin.Context) {
		version.RenderHandler(ctx, controller, "AddError")
	})
	controller.router.GET("/:version/localization/error", func(ctx *gin.Context) {
		version.RenderHandler(ctx, controller, "GetError")
	})
	controller.router.POST("/:version/localization/translation", func(ctx *gin.Context) {
		version.RenderHandler(ctx, controller, "AddTranslation")
	})
	controller.router.POST("/:version/localization/endpoint", func(ctx *gin.Context) {
		version.RenderHandler(ctx, controller, "AddEndpoint")
	})
	controller.router.GET("/:version/localization/endpointname", func(ctx *gin.Context) {
		version.RenderHandler(ctx, controller, "GetEndpointName")
	})

}

// HealthHandler handles health endpoint.
func (controller *ErrorCodesController) Health(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "server run with base version",
	})
}

// AddErrorHandler() is used to add/ update new errors
func (controller *ErrorCodesController) AddError(ctx *gin.Context) {
	var newData db.ErrorData

	if err := ctx.ShouldBindJSON(&newData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": consts.FailedJsonErr,
		})
		return
	}
	errType, endpoint, lang, method := ctx.Query("type"), ctx.Query("endpoint"), ctx.Request.Header.Get("Accept-Language"), ctx.Query("method")
	if lang == "" {
		lang = "en"
	}

	err := controller.useCases.AppendError(ctx, errType, endpoint, lang, newData, method)

	if err != nil {
		logrus.Errorf("[AppendCodes] %v", err)

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message":   consts.UnexpectedErr,
			"errorCode": consts.CodeL001,
			"errors":    consts.MissingErr,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": consts.Success,
	})

}

// GetErrorHandler handles getting error codes.
func (controller *ErrorCodesController) GetError(ctx *gin.Context) {

	errType, endpoint, lang, field, method := ctx.Query("type"), ctx.Query("endpoint"), ctx.Request.Header.Get("Accept-Language"), ctx.Query("field"), ctx.Query("method")

	if lang == "" {
		lang = "en"
	}
	resp, err := controller.useCases.GetError(errType, endpoint, lang, field, method)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"errorCode": http.StatusBadRequest,
			"message":   consts.VerifyMessage,
			"errors":    err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// AddTranslationHandler handles adding translations.
func (controller *ErrorCodesController) AddTranslation(ctx *gin.Context) {

	result, _ := utils.GetContext[[]string](ctx, consts.Language)
	_, err := controller.useCases.AddTranslation(ctx, result)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"errorCode": http.StatusBadRequest,
			"message":   consts.VerifyMessage,
			"errors":    err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": consts.Success,
	})
}

// AddEndpointHandler handles adding endpoints.
func (controller *ErrorCodesController) AddEndpoint(ctx *gin.Context) {

	var endpoint entities.RequestData
	if err := ctx.ShouldBindJSON(&endpoint); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": consts.FailedJsonErr,
		})
		return
	}
	err := controller.useCases.AddEndpoint(ctx, endpoint)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"errorCode": http.StatusBadRequest,
			"message":   consts.VerifyMessage,
			"errors":    err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": consts.Success,
	})
}

// GetEndpointNameHandler handles getting endpoint names.
func (controller *ErrorCodesController) GetEndpointName(ctx *gin.Context) {
	output, err := controller.useCases.GetEndpointName(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"errorCode": http.StatusBadRequest,
			"message":   consts.VerifyMessage,
			"errors":    err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, output)
}
