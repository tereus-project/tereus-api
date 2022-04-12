package handlers

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
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

	jobsQueues map[string]map[string]*services.RabbitMQQueue
}

func NewRemixHandler(s3Service *services.S3Service, rabbitMQService *services.RabbitMQService, databaseService *services.DatabaseService) (*RemixHandler, error) {
	var err error

	jobsQueues := map[string]map[string]*services.RabbitMQQueue{
		"c": {
			"go": nil,
		},
	}

	for sourceLanguage := range jobsQueues {
		for targetLanguage := range jobsQueues[sourceLanguage] {
			jobsQueues[sourceLanguage][targetLanguage], err = rabbitMQService.NewQueue("remix_jobs_q", "remix_jobs_ex", fmt.Sprintf("remix_jobs_%s_to_%s_rk", sourceLanguage, targetLanguage))
			if err != nil {
				return nil, err
			}
		}
	}

	return &RemixHandler{
		S3Service:       s3Service,
		RabbitMQService: rabbitMQService,
		DatabaseService: databaseService,

		jobsQueues: jobsQueues,
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

type remixReq struct {
	GitRepo string `json:"git_repo"`
}

type RemixType int64

const (
	UndefinedRemixType RemixType = iota
	InlineRemixType
	ZipRemixType
	GitRemixType
)

func (h *RemixHandler) RemixInline(c echo.Context) error {
	return h.Remix(c, InlineRemixType)
}

func (h *RemixHandler) RemixZip(c echo.Context) error {
	return h.Remix(c, ZipRemixType)
}

func (h *RemixHandler) RemixGit(c echo.Context) error {
	return h.Remix(c, GitRemixType)
}

func (h *RemixHandler) Remix(c echo.Context, remixType RemixType) error {
	srcLanguage := strings.ToLower(c.Param("src"))
	targetLanguage := strings.ToLower(c.Param("target"))

	if sourceMap, ok := h.jobsQueues[srcLanguage]; !ok {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("source language %s is not supported", srcLanguage))
	} else if _, ok := sourceMap[targetLanguage]; !ok {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("target language %s is not supported for source language %s", targetLanguage, srcLanguage))
	}

	jobID := uuid.New()

	switch remixType {
	case InlineRemixType:
		sourceCode := c.Param("code")
		if sourceCode == "" {
			return c.JSON(http.StatusBadRequest, "Missing source code")
		}

		_, err := h.S3Service.PutObject(env.S3Bucket, fmt.Sprintf("remix/%s/%s", jobID, "main.c"), strings.NewReader(sourceCode), strings.NewReader(sourceCode).Size())
		if err != nil {
			logrus.WithError(err).Error("Failed to upload file to S3")
			return c.JSON(http.StatusInternalServerError, "Failed to upload file to object storage")
		}
	case ZipRemixType:
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
			_, err = h.S3Service.PutObject(env.S3Bucket, fmt.Sprintf("remix/%s/%s", jobID, file.Name), f, file.FileInfo().Size())
			if err != nil {
				logrus.WithError(err).Error("Failed to upload file to S3")
				return c.JSON(http.StatusInternalServerError, fmt.Sprintf(`Failed to upload file "%s" to object storage`, file.Name))
			}
		}

	case GitRemixType:
		req := new(remixReq)
		err := c.Bind(req)
		if err != nil {
			return err
		}

		_, err = git.PlainClone("/tmp/"+jobID.String(), false, &git.CloneOptions{
			URL:      req.GitRepo,
			Progress: os.Stdout,
		})
		if err != nil {
			logrus.WithError(err).Error("Failed to clone git repo")
			return c.JSON(http.StatusInternalServerError, "Failed to clone git repo")
		}

		// List files in git repo
		files, err := ioutil.ReadDir("/tmp/" + jobID.String())
		if err != nil {
			logrus.WithError(err).Error("Failed to list files in git repo")
			return c.JSON(http.StatusInternalServerError, "Failed to list files in git repo")
		}
		for _, file := range files {
			if file.IsDir() {
				continue
			}

			f, err := os.Open("/tmp/" + jobID.String() + "/" + file.Name())
			if err != nil {
				logrus.WithError(err).Error("Failed to open file")
				return c.JSON(http.StatusInternalServerError, fmt.Sprintf(`Failed to open file "%s"`, file.Name()))
			}
			defer f.Close()
			_, err = h.S3Service.PutObject(env.S3Bucket, fmt.Sprintf("remix/%s/%s", jobID, file.Name()), f, file.Size())
			if err != nil {
				logrus.WithError(err).Error("Failed to upload file to S3")
				return c.JSON(http.StatusInternalServerError, fmt.Sprintf(`Failed to upload file "%s" to object storage`, file.Name()))
			}
		}

	default:
		return c.JSON(http.StatusBadRequest, "Invalid remix type")
	}

	// Publish job to exchange
	err := h.jobsQueues[srcLanguage][targetLanguage].Publish(remixJob{
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
