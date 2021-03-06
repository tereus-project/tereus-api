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
	"github.com/tereus-project/tereus-api/ent"
	"github.com/tereus-project/tereus-api/ent/submission"
	"github.com/tereus-project/tereus-api/env"
	"github.com/tereus-project/tereus-api/services"
)

type TranspilationHandler struct {
	storageService    *services.StorageService
	databaseService   *services.DatabaseService
	tokenService      *services.TokenService
	submissionService *services.SubmissionService
}

func NewTranspilationHandler(storageService *services.StorageService, databaseService *services.DatabaseService, tokenService *services.TokenService, submissionService *services.SubmissionService) (*TranspilationHandler, error) {
	return &TranspilationHandler{
		storageService:    storageService,
		databaseService:   databaseService,
		tokenService:      tokenService,
		submissionService: submissionService,
	}, nil
}

type TranspilationResult struct {
	ID             string `json:"id"`
	SourceLanguage string `json:"source_language"`
	TargetLanguage string `json:"target_language"`
	Status         string `json:"status"`
	Reason         string `json:"reason"`
	CreatedAt      string `json:"created_at"`
}

type transpilationBody struct {
	GitRepo    string `json:"git_repo"`
	SourceCode string `json:"source_code"`
}

type TranspilationType int64

const (
	UndefinedTranspilationType TranspilationType = iota
	InlineTranspilationType
	ZipTranspilationType
	GitTranspilationType
)

func (h *TranspilationHandler) TranspileInline(c echo.Context) error {
	return h.Transpile(c, InlineTranspilationType)
}

func (h *TranspilationHandler) TranspileZip(c echo.Context) error {
	return h.Transpile(c, ZipTranspilationType)
}

func (h *TranspilationHandler) TranspileGit(c echo.Context) error {
	return h.Transpile(c, GitTranspilationType)
}

func (h *TranspilationHandler) Transpile(c echo.Context, transpilationType TranspilationType) error {
	user, err := h.tokenService.GetUserFromContext(c)
	if err != nil {
		return err
	}

	body := new(transpilationBody)

	if err := c.Bind(body); err != nil {
		return err
	}

	if err := c.Validate(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	srcLanguage := strings.ToLower(c.Param("src"))
	targetLanguage := strings.ToLower(c.Param("target"))

	languagePairDetails, err := h.submissionService.GetLanguagePairDetails(srcLanguage, targetLanguage)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	submissionId := uuid.New()
	submissionSourceSize := 0

	switch transpilationType {
	case InlineTranspilationType:
		if body.SourceCode == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Missing source code")
		}

		reader := strings.NewReader(body.SourceCode)
		submissionSourceSize = int(reader.Size())
		_, err := h.storageService.PutSubmissionObject(
			submissionId.String(),
			fmt.Sprintf("main%s", languagePairDetails.SourceLanguageFileExtension),
			reader,
			reader.Size(),
		)
		if err != nil {
			logrus.WithError(err).Error("Failed to upload file to S3")
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to upload file to object storage")
		}
	case ZipTranspilationType:
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

			submissionSourceSize += int(file.UncompressedSize64)

			_, err = h.storageService.PutSubmissionObject(submissionId.String(), file.Name, f, file.FileInfo().Size())
			if err != nil {
				logrus.WithError(err).Error("Failed to upload file to S3")
				return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf(`Failed to upload file "%s" to object storage`, file.Name))
			}
		}

	case GitTranspilationType:
		if body.GitRepo == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Missing git repository")
		}

		url, err := url.Parse(body.GitRepo)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid git repository")
		}

		destination, err := os.MkdirTemp("", "tereus")
		if err != nil {
			logrus.WithError(err).Error("Failed to create temporary directory")
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to clone git repository")
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
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to clone git repository: %s", err.Error()))
		}

		// List files in git repository
		files, err := os.ReadDir(destination)
		if err != nil {
			logrus.WithError(err).Error("Failed to list files in git repository")
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to list files in git repository")
		}

		for _, file := range files {
			if file.IsDir() {
				continue
			}

			f, err := os.Open(destination + "/" + file.Name())
			if err != nil {
				logrus.WithError(err).Error("Failed to open file")
				return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf(`Failed to open file "%s"`, file.Name()))
			}
			defer f.Close()

			info, err := f.Stat()
			if err != nil {
				logrus.WithError(err).Error("Failed to stat file")
				return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf(`Failed to stat file "%s"`, file.Name()))
			}

			submissionSourceSize += int(info.Size())

			_, err = h.storageService.PutSubmissionObject(submissionId.String(), file.Name(), f, info.Size())
			if err != nil {
				logrus.WithError(err).Error("Failed to upload file to S3")
				return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf(`Failed to upload file "%s" to object storage`, file.Name()))
			}
		}

		err = os.RemoveAll(destination)
		if err != nil {
			logrus.WithError(err).Error("Failed to remove temporary directory")
		}
	default:
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid transpilation type")
	}

	newSubmission := services.SubmissionMessage{
		ID:             submissionId.String(),
		SourceLanguage: srcLanguage,
		TargetLanguage: targetLanguage,
	}
	err = h.submissionService.PublishSubmissionToTranspile(newSubmission)
	if err != nil {
		logrus.WithError(err).Error("Failed to publish submission")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to process submission")
	}

	submissionCreation := h.databaseService.Submission.Create().
		SetID(submissionId).
		SetSourceLanguage(srcLanguage).
		SetTargetLanguage(targetLanguage).
		SetSubmissionSourceSizeBytes(submissionSourceSize).
		SetIsInline(transpilationType == InlineTranspilationType).
		SetUserID(user.ID).
		SetProcessingStartedAt(time.Now())

	if transpilationType == GitTranspilationType {
		submissionCreation.SetGitRepo(body.GitRepo)
	}

	s, err := submissionCreation.Save(context.Background())
	if err != nil {
		logrus.WithError(err).Error("Failed to save submission to database")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save submission to database")
	}

	return c.JSON(http.StatusOK, TranspilationResult{
		ID:             s.ID.String(),
		SourceLanguage: s.SourceLanguage,
		TargetLanguage: s.TargetLanguage,
		Status:         s.Status.String(),
		Reason:         s.Reason,
		CreatedAt:      s.CreatedAt.Format(time.RFC3339Nano),
	})
}

// GET /submissions/:id/download
func (h *TranspilationHandler) DownloadTranspiledFiles(c echo.Context) error {
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

	objects := h.storageService.GetObjects(objectStoragePath)
	for object := range objects {
		if object.Err != nil {
			logrus.WithError(object.Err).Error("Failed to get file from S3")
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get files from S3")
		}

		reader, err := h.storageService.GetObject(object.Path)
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
	Data                 string `json:"data"`
	Status               string `json:"status"`
	SourceLanguage       string `json:"source_language"`
	TargetLanguage       string `json:"target_language"`
	SourceSizeBytes      int    `json:"source_size_bytes"`
	TargetSizeBytes      int    `json:"target_size_bytes"`
	ProcessingStartedAt  string `json:"processing_started_at"`
	ProcessingFinishedAt string `json:"processing_finished_at"`
	Reason               string `json:"reason"`
}

// GET /submissions/:id/inline/source
func (h *TranspilationHandler) DownloadInlineTranspilationSource(c echo.Context) error {
	var subID uuid.UUID
	var shareID string

	subID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		if len(c.Param("id")) != 8 {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid submission ID")
		}
		shareID = c.Param("id")
	}

	// Get sub from database
	var sub *ent.Submission
	if subID != uuid.Nil {
		sub, err = h.databaseService.Submission.Query().
			Where(submission.ID(subID)).
			Only(context.Background())
		if err != nil {
			if ent.IsNotFound(err) {
				return echo.NewHTTPError(http.StatusNotFound, "This submission does not exist")
			}

			logrus.WithError(err).Error("Failed to get submission")
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get submission")
		}
	} else {
		sub, err = h.databaseService.Submission.Query().
			Where(submission.ShareID(shareID)).
			Only(context.Background())
		if err != nil {
			if ent.IsNotFound(err) {
				return echo.NewHTTPError(http.StatusNotFound, "This submission does not exist")
			}

			logrus.WithError(err).Error("Failed to get submission")
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get submission")
		}
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

	languagePairDetails, err := h.submissionService.GetLanguagePairDetails(sub.SourceLanguage, sub.TargetLanguage)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get language pair details")
	}
	objectStoragePath := fmt.Sprintf("transpilations/%s/main%s", sub.ID, languagePairDetails.SourceLanguageFileExtension)

	// Get files from S3
	object, err := h.storageService.GetObject(objectStoragePath)
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
		Data:                 base64.StdEncoding.EncodeToString(data),
		Status:               string(sub.Status),
		SourceLanguage:       sub.SourceLanguage,
		TargetLanguage:       sub.TargetLanguage,
		SourceSizeBytes:      sub.SubmissionSourceSizeBytes,
		TargetSizeBytes:      sub.SubmissionTargetSizeBytes,
		ProcessingStartedAt:  sub.ProcessingStartedAt.Format(time.RFC3339Nano),
		ProcessingFinishedAt: sub.ProcessingFinishedAt.Format(time.RFC3339Nano),
	})
}

// GET /submissions/:id/inline/output
func (h *TranspilationHandler) DownloadInlineTranspiledOutput(c echo.Context) error {
	var subID uuid.UUID
	var shareID string

	subID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		if len(c.Param("id")) != 8 {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid submission ID")
		}
		shareID = c.Param("id")
	}

	// Get sub from database
	var sub *ent.Submission
	if subID != uuid.Nil {
		sub, err = h.databaseService.Submission.Query().
			Where(submission.ID(subID)).
			Only(context.Background())
		if err != nil {
			if ent.IsNotFound(err) {
				return echo.NewHTTPError(http.StatusNotFound, "This submission does not exist")
			}

			logrus.WithError(err).Error("Failed to get submission")
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get submission")
		}
	} else {
		sub, err = h.databaseService.Submission.Query().
			Where(submission.ShareID(shareID)).
			Only(context.Background())
		if err != nil {
			if ent.IsNotFound(err) {
				return echo.NewHTTPError(http.StatusNotFound, "This submission does not exist")
			}

			logrus.WithError(err).Error("Failed to get submission")
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get submission")
		}
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

	if sub.Status != submission.StatusDone {
		return c.JSON(http.StatusOK, downloadInlineResponse{
			Data:                 "",
			Status:               string(sub.Status),
			SourceLanguage:       sub.TargetLanguage,
			TargetLanguage:       sub.TargetLanguage,
			SourceSizeBytes:      sub.SubmissionSourceSizeBytes,
			TargetSizeBytes:      0,
			ProcessingStartedAt:  sub.ProcessingStartedAt.Format(time.RFC3339Nano),
			ProcessingFinishedAt: sub.ProcessingFinishedAt.Format(time.RFC3339Nano),
			Reason:               sub.Reason,
		})
	}

	config := env.Get()
	languagePairDetails, err := h.submissionService.GetLanguagePairDetails(sub.SourceLanguage, sub.TargetLanguage)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get language pair details")
	}
	objectStoragePath := fmt.Sprintf("%s/%s/main%s", config.SubmissionsFolder, sub.ID, languagePairDetails.TargetLanguageFileExtension)

	// Get files from S3
	object, err := h.storageService.GetObject(objectStoragePath)
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
		Data:                 base64.StdEncoding.EncodeToString(data),
		Status:               string(sub.Status),
		SourceLanguage:       sub.TargetLanguage,
		TargetLanguage:       sub.TargetLanguage,
		SourceSizeBytes:      sub.SubmissionSourceSizeBytes,
		TargetSizeBytes:      sub.SubmissionTargetSizeBytes,
		ProcessingStartedAt:  sub.ProcessingStartedAt.Format(time.RFC3339Nano),
		ProcessingFinishedAt: sub.ProcessingFinishedAt.Format(time.RFC3339Nano),
	})
}
