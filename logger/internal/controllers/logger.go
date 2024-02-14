package controllers

import (
	"net/http"

	"logger/internal/entities"
	"logger/internal/middlewares"
	"logger/internal/usecases"

	"github.com/gin-gonic/gin"
	"gitlab.com/tuneverse/toolkit/core/logger"
	"gitlab.com/tuneverse/toolkit/core/version"
)

// Options represents functional options for configuring a LoggerController.
type Options func(cntrl *LoggerController)

// LoggerController is responsible for handling logger-related operations.
type LoggerController struct {
	router      *gin.RouterGroup
	useCases    usecases.LoggerUseCaseImply
	middlewares *middlewares.Middlewares
}

// NewLoggerController creates a new LoggerController with the provided options.
func NewLoggerController(options ...Options) *LoggerController {
	controll := &LoggerController{}
	for _, opt := range options {
		opt(controll)
	}

	return controll
}

// WithRouteInit sets the router for the LoggerController.
func WithRouteInit(router *gin.RouterGroup) Options {
	return func(cntrl *LoggerController) {
		cntrl.router = router
	}
}

// WithUseCaseInit sets the useCases for the LoggerController.
func WithUseCaseInit(usecases usecases.LoggerUseCaseImply) Options {
	return func(cntrl *LoggerController) {
		cntrl.useCases = usecases
	}
}

// WithMiddlewareInit sets the middlewares for the LoggerController.
func WithMiddlewareInit(middlewares *middlewares.Middlewares) Options {
	return func(cntrl *LoggerController) {
		cntrl.middlewares = middlewares
	}
}

// InitRoutes initializes the routes for the LoggerController.
func (controller *LoggerController) InitRoutes() {
	// Define routes and handlers here.
	controller.router.GET("/:version/health", func(ctx *gin.Context) {
		version.RenderHandler(ctx, controller, "HealthHandler")
	})
	controller.router.POST("/:version/logs", func(ctx *gin.Context) {
		version.RenderHandler(ctx, controller, "AddLog")
	})
	controller.router.GET("/:version/logs", func(ctx *gin.Context) {
		version.RenderHandler(ctx, controller, "GetLogs")
	})

}

// HealthHandler is a handler for health checks.
func (controller *LoggerController) HealthHandler(ctx *gin.Context) {
	// HealthHandler handles health check requests and responds with a status message.
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "server running with base version",
	})
}

// AddLog is a handler for inserting logs.
func (controller *LoggerController) AddLog(ctx *gin.Context) {
	// AddLog handles requests to insert logs into the system.
	var log entities.Log
	ctxt := ctx.Request.Context()
	if err := ctx.BindJSON(&log); err != nil {
		logger.Error(ctxt, "AddLog failed: Ivalid request data, err=%s", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request data. Please check your input and try again.",
			"error":   err.Error(),
		})
		return
	}
	err := controller.useCases.AddLog(ctx, log)
	if err != nil {
		logger.Error(ctxt, "AddLog failed: server error, err=%s", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong on the server. Please try again later.",
			"error":   err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Logged successfully",
	})
}

// GetLogs is a handler for retrieving logs.
func (controller *LoggerController) GetLogs(ctx *gin.Context) {

	var params entities.LogParams
	ctxt := ctx.Request.Context()
	if err := ctx.BindQuery(&params); err != nil {
		logger.Error(ctxt, "GetLogs failed: Ivalid request data, err=%s", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Params",
		})
		return
	}
	resp, err := controller.useCases.GetLogs(ctx, params)
	if err != nil {
		logger.Error(ctxt, "GetLogs failed: err=%s", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	resp.StatusCode = 200
	ctx.JSON(http.StatusOK, resp)
}
