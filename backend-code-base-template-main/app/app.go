package app

import (
	"backend-code-base-template/config"
	"backend-code-base-template/internal/controllers"
	"backend-code-base-template/internal/entities"
	"sync"

	// "backend-code-base-template/internal/middlewares"
	"backend-code-base-template/internal/repo"
	"backend-code-base-template/internal/repo/driver"
	"backend-code-base-template/internal/usecases"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// method run
// env configuration
// logrus, zap
// use case intia
// repo initalization
// controller init

var (
	appName = "backend_code_base"
)

func Run() {
	// init the env config
	cfg, err := config.LoadConfig(appName)
	if err != nil {
		panic(err)
	}

	// logrus init
	log := logrus.New()

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

	api := router.Group("")
	// middleware initialization
	// m := middlewares.NewMiddlewares()
	// api.Use(m.ApiVersioning())

	// complete user related initialization
	{

		// repo initialization
		oauthGoogleRepo := repo.NewOauthGoogleRepo(pgsqlDB)
		oauthFacebookRepo := repo.NewOauthFacebookRepo(pgsqlDB)
		oauthSpotifyRepo := repo.NewOauthSpotifyRepo(pgsqlDB)

		// initilizing usecases

		oauthGoogleUseCases := usecases.NewOauthGoogleUseCase(oauthGoogleRepo)
		oauthFacebookUseCases := usecases.NewOauthFacebookUseCase(oauthFacebookRepo)
		oauthSpotifyUseCases := usecases.NewOauthSpotifyUseCase(oauthSpotifyRepo)

		// initalizin controllers
		oauthGoogleControllers := controllers.NewOauthGoogleController(api, oauthGoogleUseCases)
		oauthFacebookControllers := controllers.NewOauthFacebookController(api, oauthFacebookUseCases)
		oauthSpotifyControllers := controllers.NewOauthSpotifyController(api, oauthSpotifyUseCases)
		tokenClaimController := controllers.NewTokenClaimsController(api)

		// init the routes
		oauthGoogleControllers.InitRoutes()
		oauthFacebookControllers.InitRoutes()
		oauthSpotifyControllers.InitRoutes()
		tokenClaimController.InitRoutes()

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
// launch
func launch(cfg *entities.EnvConfig, router *gin.Engine) {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%v", cfg.Port),
		Handler: router,
	}

	// Use WaitGroup to wait for the server goroutine to finish
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		// Signal the WaitGroup when the goroutine is finished
		defer wg.Done()

		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	fmt.Println("Server listening in...", cfg.Port)

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	// Wait for the server goroutine to finish
	wg.Wait()

	log.Println("Server exiting")
}
