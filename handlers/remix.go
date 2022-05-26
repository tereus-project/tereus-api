package handlers

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	transportHttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/tereus-project/tereus-api/ent/submission"
	"github.com/tereus-project/tereus-api/ent/subscription"
	"github.com/tereus-project/tereus-api/env"
	"github.com/tereus-project/tereus-api/services"
)

type RemixHandler struct {
	S3Service       *services.S3Service
	KafkaService    *services.KafkaService
	DatabaseService *services.DatabaseService
	TokenService    *services.TokenService

	supportedLanguages map[string]map[string]int
}

func NewRemixHandler(s3Service *services.S3Service, k *services.KafkaService, databaseService *services.DatabaseService, tokenService *services.TokenService) (*RemixHandler, error) {
	supportedLanguages := make(map[string]map[string]int)
	supportedLanguages["c"] = map[string]int{
		"go": 1,
	}

	return &RemixHandler{
		S3Service:          s3Service,
		KafkaService:       k,
		DatabaseService:    databaseService,
		TokenService:       tokenService,
		supportedLanguages: supportedLanguages,
	}, nil
}

type RemixResult struct {
	ID             string `json:"id"`
	SourceLanguage string `json:"source_language"`
	TargetLanguage string `json:"target_language"`
	Status         string `json:"status"`
	Reason         string `json:"reason"`
	CreatedAt      string `json:"created_at"`
}

type remixBody struct {
	GitRepo    string `json:"git_repo"`
	SourceCode string `json:"source_code"`
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
	user, err := h.TokenService.GetUserFromContext(c)
	if err != nil {
		return err
	}

	body := new(remixBody)

	if err := c.Bind(body); err != nil {
		return err
	}

	if err := c.Validate(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	srcLanguage := strings.ToLower(c.Param("src"))
	targetLanguage := strings.ToLower(c.Param("target"))

	if sourceMap, ok := h.supportedLanguages[srcLanguage]; !ok {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("source language %s is not supported", srcLanguage))
	} else if _, ok := sourceMap[targetLanguage]; !ok {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("target language %s is not supported for source language %s", targetLanguage, srcLanguage))
	}

	jobID := uuid.New()

	switch remixType {
	case InlineRemixType:
		if body.SourceCode == "" {
			return c.JSON(http.StatusBadRequest, "Missing source code")
		}

		reader := strings.NewReader(body.SourceCode)
		_, err := h.S3Service.PutObject(fmt.Sprintf("remix/%s/%s", jobID, "main.c"), reader, reader.Size())
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

			_, err = h.S3Service.PutObject(fmt.Sprintf("remix/%s/%s", jobID, file.Name), f, file.FileInfo().Size())
			if err != nil {
				logrus.WithError(err).Error("Failed to upload file to S3")
				return c.JSON(http.StatusInternalServerError, fmt.Sprintf(`Failed to upload file "%s" to object storage`, file.Name))
			}
		}

	case GitRemixType:
		if body.GitRepo == "" {
			return c.JSON(http.StatusBadRequest, "Missing git repository")
		}

		url, err := url.Parse(body.GitRepo)
		if err != nil {
			return c.JSON(http.StatusBadRequest, "Invalid git repository")
		}

		destination, err := os.MkdirTemp("", "tereus")
		if err != nil {
			logrus.WithError(err).Error("Failed to create temporary directory")
			return c.JSON(http.StatusInternalServerError, "Failed to clone git repository")
		}

		var auth transport.AuthMethod

		if user.GithubAccessToken != "" && url.Host == "github.com" {
			auth = &transportHttp.BasicAuth{
				Username: "tereus",
				Password: user.GithubAccessToken,
			}
		}

		_, err = git.PlainClone(destination, false, &git.CloneOptions{
			URL:      body.GitRepo,
			Progress: os.Stdout,
			Auth:     auth,
		})
		if err != nil {
			logrus.WithError(err).Error("Failed to clone git repository")
			return c.JSON(http.StatusInternalServerError, "Failed to clone git repository")
		}

		// List files in git repository
		files, err := os.ReadDir(destination)
		if err != nil {
			logrus.WithError(err).Error("Failed to list files in git repository")
			return c.JSON(http.StatusInternalServerError, "Failed to list files in git repository")
		}

		for _, file := range files {
			if file.IsDir() {
				continue
			}

			f, err := os.Open(destination + "/" + file.Name())
			if err != nil {
				logrus.WithError(err).Error("Failed to open file")
				return c.JSON(http.StatusInternalServerError, fmt.Sprintf(`Failed to open file "%s"`, file.Name()))
			}
			defer f.Close()

			info, err := f.Stat()
			if err != nil {
				logrus.WithError(err).Error("Failed to stat file")
				return c.JSON(http.StatusInternalServerError, fmt.Sprintf(`Failed to stat file "%s"`, file.Name()))
			}

			_, err = h.S3Service.PutObject(fmt.Sprintf("remix/%s/%s", jobID, file.Name()), f, info.Size())
			if err != nil {
				logrus.WithError(err).Error("Failed to upload file to S3")
				return c.JSON(http.StatusInternalServerError, fmt.Sprintf(`Failed to upload file "%s" to object storage`, file.Name()))
			}
		}

		err = os.RemoveAll(destination)
		if err != nil {
			logrus.WithError(err).Error("Failed to remove temporary directory")
		}
	default:
		return c.JSON(http.StatusBadRequest, "Invalid remix type")
	}

	err = h.KafkaService.PublishSubmission(services.RemixJob{
		ID:             jobID.String(),
		SourceLanguage: srcLanguage,
		TargetLanguage: targetLanguage,
	})
	if err != nil {
		logrus.WithError(err).Error("Failed to publish job to Kafka")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to publish job to Kafka")
	}

	submissionCreation := h.DatabaseService.Submission.Create().
		SetID(jobID).
		SetSourceLanguage(srcLanguage).
		SetTargetLanguage(targetLanguage).
		SetUserID(user.ID)

	if remixType == GitRemixType {
		submissionCreation.SetGitRepo(body.GitRepo)
	}

	s, err := submissionCreation.Save(context.Background())
	if err != nil {
		logrus.WithError(err).Error("Failed to save submission to database")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save submission to database")
	}

	return c.JSON(http.StatusOK, RemixResult{
		ID:             s.ID.String(),
		SourceLanguage: s.SourceLanguage,
		TargetLanguage: s.TargetLanguage,
		Status:         s.Status.String(),
		Reason:         s.Reason,
		CreatedAt:      s.CreatedAt.Format(time.RFC3339),
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

	config := env.Get()
	objectStoragePath := fmt.Sprintf("%s/%s", config.SubmissionsFolder, job.ID)

	c.Response().Header().Set("Content-Type", "application/zip")
	c.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.zip", job.ID))

	// Create zip file
	zipFile := zip.NewWriter(c.Response().Writer)

	objects := h.S3Service.GetObjects(objectStoragePath)
	for object := range objects {
		if object.Err != nil {
			logrus.WithError(object.Err).Error("Failed to get file from S3")
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get files from S3")
		}

		reader, err := h.S3Service.GetObject(object.Path)
		if err != nil {
			logrus.WithError(err).Error("Failed to get file from S3")
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get file from S3")
		}

		objectRelativePath := strings.TrimPrefix(object.Path, objectStoragePath)
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
		reader.Close()
	}

	err = zipFile.Close()
	if err != nil {
		logrus.WithError(err).Error("Failed to close zip file")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to close zip file")
	}

	return nil
}

// GET /remix/:id/main
func (h *RemixHandler) DownloadRemixedMain(c echo.Context) error {
	user, err := h.TokenService.GetUserFromContext(c)
	if err != nil {
		return err
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid job ID")
	}

	// Get job from database
	job, err := h.DatabaseService.Submission.Query().Where(submission.ID(id)).Only(context.Background())
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "This remixing job does not exist")
	}

	if job.Status == submission.StatusFailed {
		return echo.NewHTTPError(http.StatusOK, map[string][]string{
			"errors": {job.Reason},
		})
	}

	if job.Status != submission.StatusDone {
		return echo.NewHTTPError(http.StatusNotFound, "This remixing job is not done yet")
	}

	config := env.Get()
	objectStoragePath := fmt.Sprintf("%s/%s/main.%s", config.SubmissionsFolder, job.ID, job.TargetLanguage)

	// Get files from S3
	object, err := h.S3Service.GetObject(objectStoragePath)
	if err != nil {
		logrus.WithError(err).Error("Failed to get files from S3")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get files from S3")
	}
	defer object.Close()

	if _, err := object.Stat(); err != nil {
		return echo.NewHTTPError(http.StatusNoContent)
	}

	sub, err := user.QuerySubscription().First(context.Background())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get subscription")
	}
	if sub.Tier == subscription.TierFree {
		go func() {
			err := h.S3Service.ScheduleForDeletion(job.ID.String())
			if err != nil {
				logrus.WithError(err).Error("Failed to delete submission from S3")
			}
		}()
	}

	return c.Stream(http.StatusOK, "text/plain", object)
}
