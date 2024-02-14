package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"localization/config"
	"localization/internal/consts"
	"localization/internal/controllers"
	"localization/internal/entities"
	"localization/internal/middlewares"
	"localization/internal/repo"
	"localization/internal/repo/driver"
	"localization/internal/usecases"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gitlab.com/tuneverse/toolkit/middleware"
)

// method run
// env configuration
// logrus, zap
// use case intia
// repo initalization
// controller init

func Run() {
	// init the env config
	cfg, err := config.LoadConfig(consts.AppName)
	if err != nil {
		panic(err)
	}

	// logrus init
	log := logrus.New()

	// database connection
	mongoDB, err := driver.ConnectDB(cfg.Db)
	if err != nil {
		log.Fatalf("unable to connect the database : %v", err)
		return
	}

	// database connection
	psqlDB, err := driver.ConnectPsqlDB(cfg.PsqlDb)
	if err != nil {
		log.Fatalf("unable to connect the database : %v", err)
		return
	}

	// here initalizing the router
	router := initRouter()
	if !cfg.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	// here initalizing the router
	// router := initRouter()
	// if !cfg.Debug {
	// 	gin.SetMode(gin.ReleaseMode)
	// }

	api := router.Group("/api")
	// middleware initialization
	m := middlewares.NewMiddlewares(cfg, psqlDB)
	// api.Use(m.ApiVersioning())
	api.Use(middleware.APIVersionGuard(middleware.VersionOptions{AcceptedVersions: cfg.AcceptedVersions}))
	api.Use(m.SetLanguageInContextMiddleware())

	// complete user related initialization
	{

		// repo initialization
		repo := repo.NewErrorCodesRepo(mongoDB, cfg)

		// initilizing usecases
		useCases := usecases.NewErrorCodesUseCases(repo)

		// initalizin controllers
		controller := controllers.NewErrorCodesController(
			controllers.WithRouteInit(api),
			controllers.WithUseCaseInit(useCases),
		)

		// init the routes
		controller.InitRoutes()
	}

	// runn the app
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
	quit := make(chan os.Signal)
	// kill (no param) default send syscanll.SIGTERM
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
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")
}
