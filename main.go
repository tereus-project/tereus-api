package main

import (
	"log"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"github.com/tereus-project/tereus-api/env"
	"github.com/tereus-project/tereus-api/handlers"
	"github.com/tereus-project/tereus-api/services"
	"github.com/tereus-project/tereus-go-std/logging"
)

func main() {
	err := env.LoadEnv()
	if err != nil {
		log.Fatal(err)
	}

	config := env.Get()

	sentryHook, err := logging.SetupLog(logging.LogConfig{
		Format:       config.LogFormat,
		LogLevel:     config.LogLevel,
		ShowFilename: true,
		ReportCaller: true,
		SentryDSN:    config.SentryDSN,
	})
	if err != nil {
		logrus.WithError(err).Fatal("Failed to set log configuration")
	}
	defer sentryHook.Flush()
	defer logging.RecoverAndLogPanic()

	// Echo instance
	e := echo.New()
	e.HideBanner = true

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.Validator = &services.CustomValidator{Validator: validator.New()}

	// Initialize S3 service
	logrus.Debugln("Initializing S3 service")
	s3Service, err := services.NewS3Service(config.S3Endpoint, config.S3AccessKey, config.S3SecretKey, config.S3Bucket)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to initialize S3 service")
	}

	if err := s3Service.MakeBucketIfNotExists(config.S3Bucket); err != nil {
		logrus.WithError(err).Fatalln("Failed to create bucket")
	}

	// Initialize Kafka service
	logrus.Debugln("Initializing Kafka service")
	kafkaService, err := services.NewKafkaService(config.KafkaEndpoint)
	if err != nil {
		logrus.WithError(err).Fatalln("Failed to initialize Kafka service")
	}
	defer kafkaService.CloseAllWriters()

	// Initialize database service
	logrus.Debugln("Initializing database service")
	databaseService, err := services.NewDatabaseService(config.DatabaseDriver, config.DatabaseEndpoint)
	if err != nil {
		logrus.WithError(err).Fatalln("Failed to initialize database service")
	}
	defer databaseService.Close()

	if err := databaseService.AutoMigrate(); err != nil {
		logrus.WithError(err).Fatalln("Failed to migrate database")
	}

	// Initialize GitHub service
	logrus.Debugln("Initializing GitHub service")
	githubService, err := services.NewGithubService(config.GithubOAuthClientId, config.GithubOAuthClientSecret)
	if err != nil {
		logrus.WithError(err).Fatalln("Failed to initialize RabbitMQ service")
	}

	// Initialize token service
	logrus.Debugln("Initializing token service")
	tokenService := services.NewTokenService(databaseService)
	if err != nil {
		logrus.WithError(err).Fatalln("Failed to initialize token service")
	}

	// Initialize Stripe service
	logrus.Debugln("Initializing stripe service")
	stripeService := services.NewStripeService(config.StripeSecretKey)
	if err != nil {
		logrus.WithError(err).Fatalln("Failed to initialize stripe service")
	}

	// Initialize subscription service
	logrus.Debugln("Initializing subscription service")
	subscriptionService := services.NewSubscriptionService(databaseService, stripeService)
	if err != nil {
		logrus.WithError(err).Fatalln("Failed to initialize subscription service")
	}

	logrus.Debugln("Starting submission completion listener")
	err = startSubmissionStatusListener(kafkaService, databaseService)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to start submission completion listener")
	}

	remixHandler, err := handlers.NewRemixHandler(s3Service, kafkaService, databaseService, tokenService)
	if err != nil {
		log.Fatal(err)
	}

	authHandler, err := handlers.NewAuthHandler(databaseService, githubService, tokenService)
	if err != nil {
		log.Fatal(err)
	}

	userHandler, err := handlers.NewUserHandler(databaseService, tokenService, subscriptionService)
	if err != nil {
		log.Fatal(err)
	}

	subscriptionHandler, err := handlers.NewSubscriptionHandler(databaseService, tokenService, subscriptionService)
	if err != nil {
		log.Fatal(err)
	}

	stripeWebhooksHandler, err := handlers.NewStripeWebhooksHandler(databaseService, subscriptionService, stripeService, config.StripeWebhookSecret)
	if err != nil {
		log.Fatal(err)
	}

	e.GET("/remix/:id/main", remixHandler.DownloadRemixedMain)
	e.GET("/remix/:id", remixHandler.DownloadRemixedFiles)
	e.POST("/remix/inline/:src/to/:target", remixHandler.RemixInline)
	e.POST("/remix/zip/:src/to/:target", remixHandler.RemixZip)
	e.POST("/remix/git/:src/to/:target", remixHandler.RemixGit)

	e.POST("/auth/signup/classic", authHandler.ClassicSignup)
	e.POST("/auth/login/github", authHandler.GithubLogin)

	e.GET("/users/me", userHandler.GetCurrentUser)
	e.GET("/users/me/submissions", userHandler.GetSubmissionsHistory)

	e.POST("/subscription/checkout", subscriptionHandler.CreateCheckoutSession)
	e.POST("/subscription/portal", subscriptionHandler.CreatePortalSession)

	e.POST("/stripe-webhooks", stripeWebhooksHandler.HandleWebhooks)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
