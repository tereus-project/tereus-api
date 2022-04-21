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
)

func main() {
	err := env.LoadEnv()
	if err != nil {
		log.Fatal(err)
	}

	logrus.SetLevel(logrus.DebugLevel)

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
	s3Service, err := services.NewS3Service(env.S3Endpoint, env.S3AccessKey, env.S3SecretKey)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to initialize S3 service")
	}

	if err := s3Service.MakeBucketIfNotExists(env.S3Bucket); err != nil {
		logrus.WithError(err).Fatalln("Failed to create bucket")
	}

	// Initialize RabbitMQ service
	logrus.Debugln("Initializing RabbitMQ service")
	rabbitMQService, err := services.NewRabbitMQService(env.RabbitMQEndpoint)
	if err != nil {
		logrus.WithError(err).Fatalln("Failed to initialize RabbitMQ service")
	}
	defer rabbitMQService.Close()

	// Initialize database service
	logrus.Debugln("Initializing database service")
	databaseService, err := services.NewDatabaseService("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		logrus.WithError(err).Fatalln("Failed to initialize database service")
	}
	defer databaseService.Close()

	if err := databaseService.AutoMigrate(); err != nil {
		logrus.WithError(err).Fatalln("Failed to migrate database")
	}

	// Initialize GitHub service
	logrus.Debugln("Initializing GitHub service")
	githubService, err := services.NewGithubService(env.GithubOAuthClientId, env.GithubOAuthClientSecret)
	if err != nil {
		logrus.WithError(err).Fatalln("Failed to initialize RabbitMQ service")
	}

	// Initialize token service
	logrus.Debugln("Initializing token service")
	tokenService := services.NewTokenService(databaseService)
	if err != nil {
		logrus.WithError(err).Fatalln("Failed to initialize token service")
	}

	logrus.Debugln("Starting submission completion listener")
	err = startSubmissionCompletionListener(rabbitMQService, databaseService)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to start submission completion listener")
	}

	remixHandler, err := handlers.NewRemixHandler(s3Service, rabbitMQService, databaseService, tokenService)
	if err != nil {
		log.Fatal(err)
	}

	authHandler, err := handlers.NewAuthHandler(databaseService, githubService, tokenService)
	if err != nil {
		log.Fatal(err)
	}

	e.GET("/remix/:id", remixHandler.DownloadRemixedFiles)
	e.POST("/remix/inline/:src/to/:target", remixHandler.RemixInline)
	e.POST("/remix/zip/:src/to/:target", remixHandler.RemixZip)
	e.POST("/remix/git/:src/to/:target", remixHandler.RemixGit)

	e.POST("/auth/signup/classic", authHandler.ClassicSignup)
	e.POST("/auth/login/github", authHandler.GithubLogin)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
