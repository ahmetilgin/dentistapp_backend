package server

import (
	"backend/auth"
	authhttp "backend/auth/delivery/http"
	authmongo "backend/auth/repository/mongo"
	authusecase "backend/auth/usecase"
	emailservice "backend/email_service"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"backend/job"
	jobhttp "backend/job/delivery/http"
	jobmongo "backend/job/repository/mongo"
	jobusecase "backend/job/usecase"

	"backend/country"
	regionhttp "backend/country/delivery/http"
	regionmongo "backend/country/repository/mongo"
	regionusecase "backend/country/usecase"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type App struct {
	httpServer   *http.Server
	jobUC        job.UseCase
	authUC       auth.UseCase
	regionUC     country.UseCase
	emailService *emailservice.EmailService
}

func NewApp(isProduction bool) *App {
	db := initDB(isProduction)

	userRepo := authmongo.NewUserRepository(db,
		viper.GetString("mongo.normal_user_collection"),
		viper.GetString("mongo.business_user_collection"),
		viper.GetString("mongo.password_reset_token_collection"))

	jobRepo := jobmongo.NewJobRepository(db,
		viper.GetString("mongo.profession_collection"),
		viper.GetString("mongo.job_collection"),
		userRepo)

	regionRepo := regionmongo.NewRegionRepository(db,
		viper.GetString("mongo.region_collection"))

	emailservice := emailservice.NewEmailService(
		viper.GetString("email_service.smtpServer"),
		viper.GetString("email_service.smtpPort"),
		viper.GetString("email_service.from"),
		viper.GetString("email_service.password"),
	)

	return &App{
		jobUC: jobusecase.NewJobUseCase(jobRepo),
		authUC: authusecase.NewAuthUseCase(
			userRepo,
			viper.GetString("auth.hash_salt"),
			[]byte(viper.GetString("auth.signing_key")),
			viper.GetDuration("auth.token_ttl"),
			emailservice,
		),
		regionUC:     regionusecase.NewRegionUseCase(regionRepo),
		emailService: emailservice,
	}
}

func (a *App) Run(port string) error {
	// Init gin handler
	router := gin.Default()
	router.Use(
		gin.Recovery(),
		gin.Logger(),
	)

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // fix me
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"}, // fix me
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           6 * time.Hour,
	}))

	// Set up http handlers
	// SignUp/SignIn endpoints
	authhttp.RegisterHTTPEndpoints(router, a.authUC)

	// API endpoints
	authMiddleware := authhttp.NewAuthMiddleware(a.authUC)
	api := router.Group("/api")
	jobhttp.RegisterHTTPEndpoints(api, a.jobUC, authMiddleware)
	regionhttp.RegisterHTTPEndpoints(api, a.regionUC, authMiddleware)

	// HTTP Server
	a.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil {
			log.Fatalf("Failed to listen and serve: %+v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return a.httpServer.Shutdown(ctx)
}

func initDB(isProduction bool) *mongo.Database {
	mongo_uri := ""
	if isProduction {
		mongo_uri = viper.GetString("mongo.production_database")
	} else {
		mongo_uri = viper.GetString("mongo.development_database")
	}
	client, err := mongo.NewClient(options.Client().ApplyURI(mongo_uri))
	if err != nil {
		log.Fatalf("Error occured while establishing connection to mongoDB")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	return client.Database(viper.GetString("mongo.name"))
}
