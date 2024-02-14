package app

import (
	"context"
	"fmt"
	"log"
	"member/config"
	"member/internal/consts"
	"member/internal/controllers"
	"member/internal/entities"
	"member/internal/middlewares"

	"member/internal/repo"

	"member/internal/repo/driver"
	"member/internal/usecases"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"gitlab.com/tuneverse/toolkit/core/logger"
	"gitlab.com/tuneverse/toolkit/middleware"
)

// Run initializes and starts the member service.
func Run() {
	// Initialize environment configuration
	cfg, err := config.LoadConfig(consts.AppName)
	if err != nil {
		panic(err)
	}

	// Create a file logger configuration with specified settings.
	file := &logger.FileMode{
		LogfileName:  "member.log",      // Log file name.
		LogPath:      "log",             // Log file path. Add this folder to .gitignore as logs/
		LogMaxAge:    7,                 // Maximum log file age in days.
		LogMaxSize:   1024 * 1024 * 100, // Maximum log file size (10 MB).
		LogMaxBackup: 5,                 // Maximum number of log file backups to keep.
	}

	// Configure client options for the logger.
	clientOpt := &logger.ClientOptions{
		Service:             consts.AppName, // Service name.
		LogLevel:            "trace",        // Log level.
		IncludeRequestDump:  true,           // Include request data in logs.
		IncludeResponseDump: true,           // Include response data in logs.
	}
	var log *logger.Logger

	// Check if the application is in debug mode.
	if cfg.Debug {
		log = logger.InitLogger(clientOpt, file)
	} else {
		db := &logger.CloudMode{
			URL:    cfg.LoggerServiceURL,
			Secret: cfg.LoggerSecret,
		}

		log = logger.InitLogger(clientOpt, db, file)
	}

	// Connect to the database
	pgsqlDB, err := driver.ConnectDB(cfg.Db)
	if err != nil {
		log.Fatalf("unable to connect to the database")
		return
	}

	// Initialize the router
	router := initRouter()
	if !cfg.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	api := router.Group("/api")

	// Middleware initialization
	m := middlewares.NewMiddlewares(cfg, pgsqlDB)
	api.Use(m.PartnerID())
	//api.Use(m.JwtMiddleware())

	// Additional middleware specific to the API group
	api.Use(middleware.LogMiddleware(map[string]interface{}{}))
	api.Use(middleware.APIVersionGuard(middleware.VersionOptions{
		AcceptedVersions: cfg.AcceptedVersions,
	}))

	api.Use(middleware.Localize())

	api.Use(middleware.ErrorLocalization(
		middleware.ErrorLocaleOptions{
			Cache:                  cache.New(5*time.Minute, 10*time.Minute),
			CacheExpiration:        time.Duration(time.Hour * 24),
			CacheKeyLabel:          "ERROR_CACHE_KEY_LABEL",
			LocalisationServiceURL: fmt.Sprintf("%s/localization/error", cfg.LocalisationServiceURL),
		},
	))

	api.Use(middleware.EndpointExtraction(
		middleware.EndPointOptions{
			Cache:           cache.New(5*time.Minute, 10*time.Minute),
			CacheExpiration: time.Duration(time.Hour * 24),
			CacheKeyLabel:   "ENDPOINT_CACHE_KEY_LABEL",
			EndPointsURL:    fmt.Sprintf("%s/localization/endpointname", cfg.EndpointURL),
		},
	))

	// Initialize user-related components
	{
		// Initialize the repository
		memberRepo := repo.NewMemberRepo(pgsqlDB, cfg)
		// Initialize use cases
		memberUseCases := usecases.NewMemberUseCases(memberRepo)
		// Initialize controllers
		memberControllers := controllers.NewMemberController(api, memberUseCases)
		// Initialize the routes
		memberControllers.InitRoutes()
	}
	// Run the application
	launch(cfg, router)
}

// initRouter initializes the Gin router.
func initRouter() *gin.Engine {
	router := gin.Default()
	gin.SetMode(gin.DebugMode)

	// CORS settings
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "DELETE", "GET", "OPTIONS"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Add common middlewares here
	return router
}

// launch starts the HTTP server.
func launch(cfg *entities.EnvConfig, router *gin.Engine) {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%v", cfg.Port),
		Handler: router,
	}

	go func() {
		// Service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	fmt.Println("Server listening on port", cfg.Port)

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	// Check for ctx.Done() timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("Timeout of 5 seconds.")
	}

	log.Println("Server exiting")
}
