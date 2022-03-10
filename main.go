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
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	minioClient *minio.Client
	s           RabbitMQService
)

func main() {
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
	var err error
	minioClient, err = minio.New("localhost:9000", &minio.Options{
		Creds:  credentials.NewStaticV4("miniokey", "miniosecret", ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalln(err)
	}

	// Create tereus MinIO bucket
	err = minioClient.MakeBucket(context.Background(), "tereus", minio.MakeBucketOptions{})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, bucketExistsErr := minioClient.BucketExists(context.Background(), "tereus")
		if bucketExistsErr != nil {
			log.Fatalln(bucketExistsErr)
		}

		if !exists {
			log.Fatalln(err)
		}
	}

	// Initialize RabbitMQ
	s, err = NewRabbitMQService("amqp://admin:admin@localhost:5672/")
	if err != nil {
		log.Fatalln(err)
	}

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

type remixResult struct {
	ID             string `json:"id"`
	SourceLanguage string `json:"source_language"`
	TargetLanguage string `json:"target_language"`
}

func remix(c echo.Context) error {
	srcLanguage := c.Param("src")
	targetLanguage := c.Param("target")

	jobID := uuid.New().String()

	// Open file and unzip it
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	zipReader, err := zip.NewReader(src, int64(file.Size))
	if err != nil {
		log.Fatal(err)
	}

	// Upload files to minio
	for _, file := range zipReader.File {
		if file.FileInfo().IsDir() {
			continue
		}
		f, err := file.Open()
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		log.Println(file.Name)
		_, err = minioClient.PutObject(
			context.Background(),
			"tereus",
			fmt.Sprintf("remix/%s/%s", jobID, file.Name),
			f,
			file.FileInfo().Size(),
			minio.PutObjectOptions{},
		)
		if err != nil {
			c.Logger().Error(err)
			return c.String(http.StatusInternalServerError, "Internal Server Error")
		}
	}

	// Publish job to exchange
	err = s.publishJob(remixJob{
		ID:             jobID,
		SourceLanguage: srcLanguage,
		TargetLanguage: targetLanguage,
	})
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	return c.JSON(http.StatusOK, remixResult{
		ID:             jobID,
		SourceLanguage: srcLanguage,
		TargetLanguage: targetLanguage,
	})
}
