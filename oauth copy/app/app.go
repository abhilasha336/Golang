package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"oauth/config"
	"oauth/internal/consts"
	"oauth/internal/controllers"
	"oauth/internal/entities"
	"oauth/internal/repo"
	"oauth/internal/repo/driver"
	"oauth/internal/usecases"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gitlab.com/tuneverse/toolkit/core/logger"
	"gitlab.com/tuneverse/toolkit/middleware"
)

// Run  method  function is used to start server
// env configuration
// logrus, zap
// use case intia
// repo initalization
// controller init
func Run() {

	// var lg *logger.Logger

	// init the env config
	cfg, err := config.LoadConfig(consts.AppName)
	if err != nil {
		panic(err)
	}

	// ################ For Logging to a File ################
	// Create a file logger configuration with specified settings.
	file := &logger.FileMode{
		LogfileName:  "OAuth.log",      // Log file name.
		LogPath:      "logs",           // Log file path. Add this folder to .gitignore as logs/
		LogMaxAge:    7,                // Maximum log file age in days.
		LogMaxSize:   1024 * 1024 * 10, // Maximum log file size (10 MB).
		LogMaxBackup: 5,                // Maximum number of log file backups to keep.
	}

	// ################ For Logging to a Database ################
	// Create a database logger configuration with the specified URL and secret.
	// For development purpose don't use cloudMode.

	// ############# Client Options #############
	// Configure client options for the logger.
	clientOpt := &logger.ClientOptions{
		Service:             consts.AppName, // Service name.
		LogLevel:            "info",         // Log level.
		IncludeRequestDump:  true,           // Include request data in logs.
		IncludeResponseDump: true,           // Include response data in logs.
	}

	//check gebug is true,if true log only in file &console
	if cfg.Debug {
		// Debug Mode: Logs will print to both the file and the console.

		// Initialize the logger with the specified configurations for file and console logging.
		logger.InitLogger(clientOpt, file)
	} else {
		// Release Mode: Logs will print to a database, file, and console.
		// Create a database logger configuration with the specified URL and secret.
		db := &logger.CloudMode{
			// Database API endpoint (for best practice, load this from an environment variable).
			URL: cfg.LoggerServiceURL,
			// Secret for authentication.
			Secret: cfg.LoggerSecret,
		}

		// 	// Initialize the logger with the specified configurations for database, file, and console logging.
		logger.InitLogger(clientOpt, db, file)
	}

	// database connection
	pgsqlDB, err := driver.ConnectDB(cfg.Db)
	if err != nil {
		log.Fatalf("unable to connect the database")
		return
	}

	// here initalizing the router
	router := initRouter()
	if !cfg.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	api := router.Group("/api")

	// m := middlewares.NewMiddlewares(cfg, pgsqlDB)
	// api.Use(m.JwtMiddleware())

	// middleware from toolkit/middleware  to intialise log middleware
	api.Use(middleware.LogMiddleware(map[string]interface{}{}))

	//middleware from toolkit/middleware  to intialise version middleware
	api.Use(middleware.APIVersionGuard(middleware.VersionOptions{
		AcceptedVersions: cfg.AcceptedVersions,
	}))

	// complete user related initialization
	{

		// repo initialization
		oauthRepo := repo.NewOauthRepo(pgsqlDB, cfg)
		// initilizing usecases
		oauthUseCases := usecases.NewOauthUseCase(oauthRepo, entities.OAuthData{})
		// initalizing controllers
		oauthControllers := controllers.NewOauthController(api, oauthUseCases, cfg)
		// init the routes
		oauthControllers.InitRoutes()

	}

	// run the app
	launch(cfg, router)
}

func initRouter() *gin.Engine {
	router := gin.Default()
	gin.SetMode(gin.DebugMode)

	// CORS
	// - PUT and PATCH methods
	// - Origin header
	// - Credentials share
	// - Preflight requests cached for 12 hours
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "DELETE", "GET", "OPTIONS"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		// AllowOriginFunc: func(origin string) bool {
		// },
		MaxAge: 12 * time.Hour,
	}))

	// common middlewares should be added here

	return router
}

// launch
func launch(cfg *entities.EnvConfig, router *gin.Engine) {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%v", cfg.Port),
		Handler: router,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	fmt.Println("Server listening in...", cfg.Port)
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1) // kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.

	<-ctx.Done()
	log.Println("timeout of 5 seconds.")

	log.Println("Server exiting")
}
