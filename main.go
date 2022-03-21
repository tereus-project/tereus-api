package main

import (
	"archive/zip"
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
	"github.com/tereus-project/tereus-api/ent"
	"github.com/tereus-project/tereus-api/env"
)

var (
	minioClient *minio.Client
	s           RabbitMQService
	client      *ent.Client
)

func main() {
	err := env.LoadEnv()
	if err != nil {
		log.Fatal(err)
	}

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Routes
	e.GET("/", hello)
	e.POST("/remix/:src/to/:target", remix)

	// Initialize minio client object.
	minioClient, err = minio.New(env.S3Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(env.S3AccessKey, env.S3SecretKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalln(err)
	}

	// Create tereus S3 bucket if it doesn't exist
	exists, err := minioClient.BucketExists(context.Background(), env.S3Bucket)
	if err != nil {
		log.Fatalln(err)
	}

	if !exists {
		err = minioClient.MakeBucket(context.Background(), env.S3Bucket, minio.MakeBucketOptions{})
		if err != nil {
			// Check to see if we already own this bucket (which happens if you run this twice)
			log.Fatalln(err)
		}
	}

	// Initialize RabbitMQ
	s, err = NewRabbitMQService(env.RabbitMQEndpoint)
	if err != nil {
		log.Fatalln(err)
	}

	// Connect to DB
	client, err = ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		logrus.WithError(err).Fatal("Failed to connect to database")

	}
	defer client.Close()

	// Run the auto migration tool
	if err := client.Schema.Create(context.Background()); err != nil {
		logrus.WithError(err).Fatal("failed creating schema resources")
	}

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

// Handler
func hello(c echo.Context) error {
	return c.JSON(http.StatusOK, "Hello, World!")
}

type remixResult struct {
	ID             string `json:"id"`
	SourceLanguage string `json:"source_language"`
	TargetLanguage string `json:"target_language"`
}

func remix(c echo.Context) error {
	srcLanguage := c.Param("src")
	targetLanguage := c.Param("target")

	jobID := uuid.New()

	// Open file and unzip it
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Missing file")
	}
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to open file")
	}
	defer src.Close()

	zipReader, err := zip.NewReader(src, int64(file.Size))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to unzip file")
	}

	// Upload files to minio
	for _, file := range zipReader.File {
		if file.FileInfo().IsDir() {
			continue
		}
		f, err := file.Open()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "Failed to open file")
		}
		defer f.Close()
		log.Println(file.Name)
		_, err = minioClient.PutObject(
			context.Background(),
			env.S3Bucket,
			fmt.Sprintf("remix/%s/%s", jobID, file.Name),
			f,
			file.FileInfo().Size(),
			minio.PutObjectOptions{},
		)
		if err != nil {
			c.Logger().Error(err)
			return c.JSON(http.StatusInternalServerError, "Failed to upload file to object storage")
		}
	}

	// Publish job to exchange
	err = s.publishJob(remixJob{
		ID:             jobID.String(),
		SourceLanguage: srcLanguage,
		TargetLanguage: targetLanguage,
	})
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, "Failed to publish job to RabbitMQ")
	}

	_, err = client.Submission.Create().
		SetID(jobID).
		SetSourceLanguage(srcLanguage).
		SetTargetLanguage(targetLanguage).
		Save(context.Background())
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, "Failed to create submission")
	}

	return c.JSON(http.StatusOK, remixResult{
		ID:             jobID.String(),
		SourceLanguage: srcLanguage,
		TargetLanguage: targetLanguage,
	})
}
