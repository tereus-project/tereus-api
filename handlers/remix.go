package handlers

import (
	"archive/zip"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
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
	"github.com/tereus-project/tereus-api/env"
	"github.com/tereus-project/tereus-api/services"
)

type RemixHandler struct {
	s3Service         *services.S3Service
	databaseService   *services.DatabaseService
	tokenService      *services.TokenService
	submissionService *services.SubmissionService
}

func NewRemixHandler(s3Service *services.S3Service, databaseService *services.DatabaseService, tokenService *services.TokenService, submissionService *services.SubmissionService) (*RemixHandler, error) {
	return &RemixHandler{
		s3Service:         s3Service,
		databaseService:   databaseService,
		tokenService:      tokenService,
		submissionService: submissionService,
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
	user, err := h.tokenService.GetUserFromContext(c)
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

	err = h.submissionService.CheckSupport(srcLanguage, targetLanguage)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	submissionId := uuid.New()

	switch remixType {
	case InlineRemixType:
		if body.SourceCode == "" {
			return c.JSON(http.StatusBadRequest, "Missing source code")
		}

		reader := strings.NewReader(body.SourceCode)
		_, err := h.s3Service.PutObject(fmt.Sprintf("remix/%s/%s", submissionId, "main.c"), reader, reader.Size())
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

			_, err = h.s3Service.PutObject(fmt.Sprintf("remix/%s/%s", submissionId, file.Name), f, file.FileInfo().Size())
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

			_, err = h.s3Service.PutObject(fmt.Sprintf("remix/%s/%s", submissionId, file.Name()), f, info.Size())
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

	err = h.submissionService.PublishSubmission(&services.SubmissionMessage{
		ID:             submissionId.String(),
		SourceLanguage: srcLanguage,
		TargetLanguage: targetLanguage,
	})
	if err != nil {
		logrus.WithError(err).Error("Failed to publish submission to Kafka")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to publish submission to Kafka")
	}

	submissionCreation := h.databaseService.Submission.Create().
		SetID(submissionId).
		SetSourceLanguage(srcLanguage).
		SetTargetLanguage(targetLanguage).
		SetIsInline(remixType == InlineRemixType).
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

// GET /submissions/:id/download
func (h *RemixHandler) DownloadRemixedFiles(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid submission ID")
	}

	// Get sub from database
	sub, err := h.databaseService.Submission.Query().Where(submission.ID(id)).Only(context.Background())
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "This submission does not exist")
	}

	if sub.Status != "done" {
		return echo.NewHTTPError(http.StatusNotFound, "This submission is not done yet")
	}

	config := env.Get()
	objectStoragePath := fmt.Sprintf("%s/%s", config.SubmissionsFolder, sub.ID)

	c.Response().Header().Set("Content-Type", "application/zip")
	c.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.zip", sub.ID))

	// Create zip file
	zipFile := zip.NewWriter(c.Response().Writer)

	objects := h.s3Service.GetObjects(objectStoragePath)
	for object := range objects {
		if object.Err != nil {
			logrus.WithError(object.Err).Error("Failed to get file from S3")
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get files from S3")
		}

		reader, err := h.s3Service.GetObject(object.Path)
		if err != nil {
			logrus.WithError(err).Error("Failed to get file from S3")
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get file from S3")
		}

		objectRelativePath := strings.TrimPrefix(object.Path, objectStoragePath)
		zippedFilePath := fmt.Sprintf("%s/%s", sub.ID, objectRelativePath)

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

type downloadInlineResponse struct {
	Data           string `json:"data"`
	SourceLanguage string `json:"source_language"`
	TargetLanguage string `json:"target_language"`
}

// GET /submissions/:id/inline/source
func (h *RemixHandler) DownloadInlineRemixSource(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid submissions ID")
	}

	// Get sub from database
	sub, err := h.databaseService.Submission.Query().Where(submission.ID(id)).Only(context.Background())
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "This submission does not exist")
	}

	if !sub.IsInline {
		return echo.NewHTTPError(http.StatusNotFound, "This submission is not inline")
	}

	if !sub.IsPublic {
		user, err := h.tokenService.GetUserFromContext(c)
		if err != nil {
			return err
		}

		owner, err := sub.QueryUser().OnlyID(context.Background())
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get owner of submission")
		}

		if user.ID != owner {
			return echo.NewHTTPError(http.StatusForbidden, "This submission is not public and you are not the owner")
		}
	}

	if sub.Status == submission.StatusCleaned {
		return echo.NewHTTPError(http.StatusNotFound, "This submission has been cleaned")
	}

	objectStoragePath := fmt.Sprintf("remix/%s/main.%s", sub.ID, sub.SourceLanguage)

	// Get files from S3
	object, err := h.s3Service.GetObject(objectStoragePath)
	if err != nil {
		logrus.WithError(err).Error("Failed to get files from S3")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get files from S3")
	}
	defer object.Close()

	if _, err := object.Stat(); err != nil {
		return echo.NewHTTPError(http.StatusNoContent)
	}

	data, err := ioutil.ReadAll(object)
	if err != nil {
		logrus.WithError(err).Error("Failed to read file from S3")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to read file from S3")
	}

	return c.JSON(http.StatusOK, downloadInlineResponse{
		Data:           base64.StdEncoding.EncodeToString(data),
		SourceLanguage: sub.SourceLanguage,
		TargetLanguage: sub.TargetLanguage,
	})
}

// GET /submissions/:id/inline/output
func (h *RemixHandler) DownloadInlineRemixdOutput(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid submission ID")
	}

	// Get sub from database
	sub, err := h.databaseService.Submission.Query().Where(submission.ID(id)).Only(context.Background())
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "This submission does not exist")
	}

	if !sub.IsInline {
		return echo.NewHTTPError(http.StatusNotFound, "This submission is not inline")
	}

	if !sub.IsPublic {
		user, err := h.tokenService.GetUserFromContext(c)
		if err != nil {
			return err
		}

		owner, err := sub.QueryUser().OnlyID(context.Background())
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get owner of submission")
		}

		if user.ID != owner {
			return echo.NewHTTPError(http.StatusForbidden, "This submission is not public and you are not the owner")
		}
	}

	if sub.Status == submission.StatusFailed {
		return echo.NewHTTPError(http.StatusOK, map[string][]string{
			"errors": {sub.Reason},
		})
	}

	if sub.Status != submission.StatusDone {
		return echo.NewHTTPError(http.StatusNotFound, "This submission is not done yet")
	}

	config := env.Get()
	objectStoragePath := fmt.Sprintf("%s/%s/main.%s", config.SubmissionsFolder, sub.ID, sub.TargetLanguage)

	// Get files from S3
	object, err := h.s3Service.GetObject(objectStoragePath)
	if err != nil {
		logrus.WithError(err).Error("Failed to get files from S3")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get files from S3")
	}
	defer object.Close()

	if _, err := object.Stat(); err != nil {
		return echo.NewHTTPError(http.StatusNoContent)
	}

	data, err := ioutil.ReadAll(object)
	if err != nil {
		logrus.WithError(err).Error("Failed to read file from S3")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to read file from S3")
	}

	return c.JSON(http.StatusOK, downloadInlineResponse{
		Data:           base64.StdEncoding.EncodeToString(data),
		SourceLanguage: sub.TargetLanguage,
		TargetLanguage: sub.TargetLanguage,
	})
}
