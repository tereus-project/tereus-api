package handlers

import (
	"archive/zip"
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/tereus-project/tereus-api/env"
	"github.com/tereus-project/tereus-api/services"
)

type RemixHandler struct {
	S3Service       *services.S3Service
	RabbitMQService *services.RabbitMQService
	DatabaseService *services.DatabaseService

	jobsQueue *services.RabbitMQQueue
}

func NewRemixHandler(s3Service *services.S3Service, rabbitMQService *services.RabbitMQService, databaseService *services.DatabaseService) (*RemixHandler, error) {
	jobsQueue, err := rabbitMQService.NewQueue("remix_jobs_q", "remix_jobs_ex", "remix_jobs_rk")
	if err != nil {
		return nil, err
	}

	return &RemixHandler{
		S3Service:       s3Service,
		RabbitMQService: rabbitMQService,
		DatabaseService: databaseService,
		jobsQueue:       jobsQueue,
	}, nil
}

type remixResult struct {
	ID             string `json:"id"`
	SourceLanguage string `json:"source_language"`
	TargetLanguage string `json:"target_language"`
}

type remixJob struct {
	ID             string `json:"id"`
	SourceLanguage string `json:"source_language"`
	TargetLanguage string `json:"target_language"`
}

func (h *RemixHandler) Remix(c echo.Context) error {
	srcLanguage := c.Param("src")
	targetLanguage := c.Param("target")

	jobID := uuid.New()

	// Open file and unzip it
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Missing file")
	}

	source, err := file.Open()
	if err != nil {
		logrus.WithError(err).Error("Failed to open file")
		return c.JSON(http.StatusInternalServerError, "Failed to open file")
	}
	defer source.Close()

	zipReader, err := zip.NewReader(source, int64(file.Size))
	if err != nil {
		logrus.WithError(err).Error("Failed to unzip file")
		return c.JSON(http.StatusInternalServerError, "Failed to unzip file")
	}

	// Upload files to minio
	for _, file := range zipReader.File {
		if file.FileInfo().IsDir() {
			continue
		}

		f, err := file.Open()
		if err != nil {
			logrus.WithError(err).Error("Failed to open file")
			return c.JSON(http.StatusInternalServerError, fmt.Sprintf(`Failed to open file "%s"`, file.Name))
		}
		defer f.Close()

		log.Println(file.Name)
		_, err = h.S3Service.PutObject(env.S3Bucket, fmt.Sprintf("remix/%s/%s", jobID, file.Name), f, file.FileInfo().Size())
		if err != nil {
			logrus.WithError(err).Error("Failed to upload file to S3")
			return c.JSON(http.StatusInternalServerError, fmt.Sprintf(`Failed to upload file "%s" to object storage`, file.Name))
		}
	}

	// Publish job to exchange
	err = h.jobsQueue.PublishJob(remixJob{
		ID:             jobID.String(),
		SourceLanguage: srcLanguage,
		TargetLanguage: targetLanguage,
	})
	if err != nil {
		logrus.WithError(err).Error("Failed to publish job to RabbitMQ")
		return c.JSON(http.StatusInternalServerError, "Failed to publish job to RabbitMQ")
	}

	_, err = h.DatabaseService.Submission.Create().
		SetID(jobID).
		SetSourceLanguage(srcLanguage).
		SetTargetLanguage(targetLanguage).
		Save(context.Background())
	if err != nil {
		logrus.WithError(err).Error("Failed to save submission to database")
		return c.JSON(http.StatusInternalServerError, "Failed to save submission to database")
	}

	return c.JSON(http.StatusOK, remixResult{
		ID:             jobID.String(),
		SourceLanguage: srcLanguage,
		TargetLanguage: targetLanguage,
	})
}
