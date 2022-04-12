package handlers

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/tereus-project/tereus-api/ent/submission"
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

// POST /remix/:src/to/:target
func (h *RemixHandler) Remix(c echo.Context) error {
	srcLanguage := c.Param("src")
	targetLanguage := c.Param("target")

	jobID := uuid.New()

	// Open file and unzip it
	file, err := c.FormFile("file")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing file")
	}

	source, err := file.Open()
	if err != nil {
		logrus.WithError(err).Error("Failed to open file")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to open file")
	}
	defer source.Close()

	zipReader, err := zip.NewReader(source, int64(file.Size))
	if err != nil {
		logrus.WithError(err).Error("Failed to unzip file")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to unzip file")
	}

	// Upload files to minio
	for _, file := range zipReader.File {
		if file.FileInfo().IsDir() {
			continue
		}

		f, err := file.Open()
		if err != nil {
			logrus.WithError(err).Error("Failed to open file")
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf(`Failed to open file "%s"`, file.Name))
		}
		defer f.Close()

		log.Println(file.Name)
		_, err = h.S3Service.PutObject(env.S3Bucket, fmt.Sprintf("remix/%s/%s", jobID, file.Name), f, file.FileInfo().Size())
		if err != nil {
			logrus.WithError(err).Error("Failed to upload file to S3")
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf(`Failed to upload file "%s" to object storage`, file.Name))
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
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to publish job to RabbitMQ")
	}

	_, err = h.DatabaseService.Submission.Create().
		SetID(jobID).
		SetSourceLanguage(srcLanguage).
		SetTargetLanguage(targetLanguage).
		Save(context.Background())
	if err != nil {
		logrus.WithError(err).Error("Failed to save submission to database")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save submission to database")
	}

	return c.JSON(http.StatusOK, remixResult{
		ID:             jobID.String(),
		SourceLanguage: srcLanguage,
		TargetLanguage: targetLanguage,
	})
}

// GET /remix/:id
func (h *RemixHandler) DownloadRemixedFiles(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid job ID")
	}

	// Get job from database
	job, err := h.DatabaseService.Submission.Query().Where(submission.ID(id)).Only(context.Background())
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "This remixing job does not exist")
	}

	if job.Status != "done" {
		return echo.NewHTTPError(http.StatusNotFound, "This remixing job is not done yet")
	}

	objectStoragePath := fmt.Sprintf("%s/%s", env.SubmissionsFolder, job.ID)

	// Get files from S3
	paths, err := h.S3Service.GetObjects(env.S3Bucket, objectStoragePath)
	if err != nil {
		logrus.WithError(err).Error("Failed to get files from S3")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get files from S3")
	}

	c.Response().Header().Set("Content-Type", "application/zip")
	c.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.zip", job.ID))

	// Create zip file
	zipFile := zip.NewWriter(c.Response().Writer)

	for path := range paths {
		reader, err := h.S3Service.GetObject(env.S3Bucket, path)
		if err != nil {
			logrus.WithError(err).Error("Failed to get file from S3")
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get file from S3")
		}
		defer reader.Close()

		objectRelativePath := strings.TrimPrefix(path, objectStoragePath)
		zippedFilePath := fmt.Sprintf("%s/%s", job.ID, objectRelativePath)

		writer, err := zipFile.Create(zippedFilePath)
		if err != nil {
			logrus.WithError(err).Error("Failed to create file in zip")
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create file in zip")
		}

		_, err = io.Copy(writer, reader)
		if err != nil {
			logrus.WithError(err).Error("Failed to copy file to zip")
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to copy file to zip")
		}
	}

	err = zipFile.Close()
	if err != nil {
		logrus.WithError(err).Error("Failed to close zip file")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to close zip file")
	}

	return nil
}
