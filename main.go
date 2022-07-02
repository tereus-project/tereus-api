package main

import (
	"log"

	"github.com/go-playground/validator"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	echoSwagger "github.com/swaggo/echo-swagger"
	_ "github.com/tereus-project/tereus-api/docs"
	"github.com/tereus-project/tereus-api/env"
	"github.com/tereus-project/tereus-api/handlers"
	"github.com/tereus-project/tereus-api/services"
	"github.com/tereus-project/tereus-api/workers"
	"github.com/tereus-project/tereus-go-std/logging"
	"github.com/tereus-project/tereus-go-std/queue"
)

// @title Tereus API
// @version 1.0
// @description The main API for the Tereus project.

// @contact.name Tereus Team
// @contact.url https://github.com/tereus-project

// @host api.tereus.dev
// @BasePath /
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

	// Initialize storage service
	logrus.Debugln("Initializing storage service")
	storageService, err := services.NewStorageService(config.S3Endpoint, config.S3AccessKey, config.S3SecretKey, config.S3Bucket, config.S3HTTPSEnabled)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to initialize storage service")
	}

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
		logrus.WithError(err).Fatalln("Failed to initialize GitHub service")
	}

	// Initialize GitLab service
	logrus.Debugln("Initializing GitLab service")
	gitlabService, err := services.NewGitlabService(config.GitlabOAuthClientId, config.GitlabOAuthClientSecret)
	if err != nil {
		logrus.WithError(err).Fatalln("Failed to initialize GitLab service")
	}

	// Initialize token service
	logrus.Debugln("Initializing token service")
	tokenService := services.NewTokenService(databaseService)

	// Initialize subscription service
	logrus.Debugln("Initializing subscription service")
	subscriptionService := services.NewSubscriptionService(
		config.StripeSecretKey,
		services.TierPrices{
			BasePriceId:    config.StripeTierProBase,
			MeteredPriceId: config.StripeTierProMetered,
		},
		services.TierPrices{
			BasePriceId:    config.StripeTierEnterpriseBase,
			MeteredPriceId: config.StripeTierEnterpriseMetered,
		},
		databaseService,
	)

	// Initialize queue service
	logrus.Debugln("Initializing queue service")
	queueService, err := queue.NewQueueService(config.NSQEndpoint, config.NSQLookupdEndpoint)
	if err != nil {
		logrus.WithError(err).Fatalln("Failed to initialize queue service")
	}
	defer queueService.Close()

	// Initialize submission service
	logrus.Debugln("Initializing submission service")
	submissionService := services.NewSubmissionService(queueService, databaseService, storageService)

	logrus.Debugln("Starting submission status consumer worker")
	err = workers.RegisterStatusConsumerWorker(submissionService, queueService)
	if err != nil {
		logrus.WithError(err).Fatalln("Failed to start submission status consumer worker")
	}

	logrus.Debugln("Starting subscription data usage reporting worker")
	go workers.SubscriptionDataUsageReportingWorker(subscriptionService, databaseService)

	logrus.Debugln("Starting retention worker")
	go workers.RetentionWorker(databaseService, storageService)

	transpilationHandler, err := handlers.NewTranspilationHandler(storageService, databaseService, tokenService, submissionService)
	if err != nil {
		log.Fatal(err)
	}

	authHandler, err := handlers.NewAuthHandler(databaseService, githubService, gitlabService, tokenService)
	if err != nil {
		log.Fatal(err)
	}

	userHandler, err := handlers.NewUserHandler(databaseService, tokenService, subscriptionService, storageService)
	if err != nil {
		log.Fatal(err)
	}

	submissionHandler, err := handlers.NewSubmissionsHandler(databaseService, tokenService, storageService)
	if err != nil {
		log.Fatal(err)
	}

	subscriptionHandler, err := handlers.NewSubscriptionHandler(databaseService, tokenService, subscriptionService)
	if err != nil {
		log.Fatal(err)
	}

	stripeWebhooksHandler, err := handlers.NewStripeWebhooksHandler(databaseService, subscriptionService, config.StripeWebhookSecret)
	if err != nil {
		log.Fatal(err)
	}

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.POST("/submissions/inline/:src/to/:target", transpilationHandler.TranspileInline)
	e.POST("/submissions/zip/:src/to/:target", transpilationHandler.TranspileZip)
	e.POST("/submissions/git/:src/to/:target", transpilationHandler.TranspileGit)

	e.DELETE("/submissions/:id", submissionHandler.DeleteSubmission)
	e.PATCH("/submissions/:id/visibility", submissionHandler.UpdateSubmissionVisibility)

	e.GET("/submissions/:id/download", transpilationHandler.DownloadTranspiledFiles)
	e.GET("/submissions/:id/inline/source", transpilationHandler.DownloadInlineTranspilationSource)
	e.GET("/submissions/:id/inline/output", transpilationHandler.DownloadInlineTranspiledOutput)

	e.POST("/auth/login/github", authHandler.LoginGithub)
	e.POST("/auth/revoke/github", authHandler.RevokeGithub)
	e.POST("/auth/login/gitlab", authHandler.LoginGitlab)
	e.POST("/auth/revoke/gitlab", authHandler.RevokeGitlab)
	e.POST("/auth/check", authHandler.Check)

	e.GET("/users/me", userHandler.GetCurrentUser)
	e.DELETE("/users/me", userHandler.DeleteCurrentUser)
	e.GET("/users/me/linked-accounts", userHandler.GetCurrentUserLinkedAccounts)
	e.GET("/users/me/submissions", userHandler.GetSubmissionsHistory)
	e.GET("/users/me/export", userHandler.GetExport)

	e.POST("/subscription/checkout", subscriptionHandler.CreateCheckoutSession)
	e.POST("/subscription/portal", subscriptionHandler.CreatePortalSession)

	e.POST("/stripe-webhooks", stripeWebhooksHandler.HandleWebhooks)

	p := prometheus.NewPrometheus("echo", nil)
	p.Use(e)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
