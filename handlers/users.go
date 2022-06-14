package handlers

import (
	"archive/zip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/tereus-project/tereus-api/ent"
	"github.com/tereus-project/tereus-api/ent/submission"
	"github.com/tereus-project/tereus-api/ent/user"
	"github.com/tereus-project/tereus-api/services"
)

type UserHandler struct {
	databaseService     *services.DatabaseService
	tokenService        *services.TokenService
	subscriptionService *services.SubscriptionService
	s3Service           *services.S3Service
}

func NewUserHandler(databaseService *services.DatabaseService, tokenService *services.TokenService, subscriptionService *services.SubscriptionService, s3Service *services.S3Service) (*UserHandler, error) {
	return &UserHandler{
		databaseService:     databaseService,
		tokenService:        tokenService,
		subscriptionService: subscriptionService,
		s3Service:           s3Service,
	}, nil
}

type getCurrentUserResultSubscription struct {
	Tier      string `json:"tier"`
	ExpiresAt string `json:"expires_at"`
	Cancelled bool   `json:"cancelled"`
}

type getCurrentUserResult struct {
	ID                string                            `json:"id"`
	Email             string                            `json:"email"`
	Subscription      *getCurrentUserResultSubscription `json:"subscription"`
	CurrentUsageBytes int64                             `json:"current_usage_bytes"`
}

// GET /users/me
func (h *UserHandler) GetCurrentUser(c echo.Context) error {
	loggedUser, err := h.tokenService.GetUserFromContext(c)
	if err != nil {
		return err
	}

	subscription, err := h.subscriptionService.GetCurrentUserSubscription(loggedUser.ID)
	if err != nil && err.(*ent.NotFoundError) == nil {
		logrus.WithError(err).Error("Failed to get current user subscription")
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var subscriptionResult *getCurrentUserResultSubscription
	if subscription != nil {
		subscriptionResult = &getCurrentUserResultSubscription{
			Tier:      subscription.Tier.String(),
			ExpiresAt: subscription.ExpiresAt.Format(time.RFC3339),
			Cancelled: subscription.Cancelled,
		}
	}

	currentUsageBytes, err := h.databaseService.Submission.Query().
		Where(submission.HasUserWith(user.ID(loggedUser.ID))).
		Modify(func(s *sql.Selector) {
			// COALESCE is so that we have a default value of 0 if there are no volumes in the DB
			// otherwise it will return a null row which will cause an error
			s.Select("COALESCE(SUM(submission_source_size_bytes + submission_target_size_bytes),0) as usage")
		}).
		Int(context.Background())

	return c.JSON(http.StatusOK, getCurrentUserResult{
		ID:                loggedUser.ID.String(),
		Email:             loggedUser.Email,
		Subscription:      subscriptionResult,
		CurrentUsageBytes: int64(currentUsageBytes),
	})
}

type submissionsHistoryItem struct {
	ID              string `json:"id"`
	SourceLanguage  string `json:"source_language"`
	TargetLanguage  string `json:"target_language"`
	IsInline        bool   `json:"is_inline"`
	IsPublic        bool   `json:"is_public"`
	Status          string `json:"status"`
	Reason          string `json:"reason"`
	CreatedAt       string `json:"created_at"`
	ShareID         string `json:"share_id"`
	SourceSizeBytes int    `json:"source_size_bytes"`
	TargetSizeBytes int    `json:"target_size_bytes"`
}

type submissionsHistory struct {
	Submissions []*submissionsHistoryItem `json:"submissions"`
}

// GET /users/me/submissions
func (h *UserHandler) GetSubmissionsHistory(c echo.Context) error {
	tereusUser, err := h.tokenService.GetUserFromContext(c)
	if err != nil {
		return err
	}

	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 10
	}

	submissions, err := h.databaseService.Submission.Query().
		Where(submission.HasUserWith(user.ID(tereusUser.ID))).
		Order(ent.Desc(submission.FieldCreatedAt)).
		Offset(limit * (page - 1)).
		Limit(limit).
		All(context.Background())
	if err != nil {
		return err
	}

	items := submissionsHistory{
		Submissions: make([]*submissionsHistoryItem, len(submissions)),
	}

	for i, s := range submissions {
		items.Submissions[i] = &submissionsHistoryItem{
			ID:              s.ID.String(),
			SourceLanguage:  s.SourceLanguage,
			TargetLanguage:  s.TargetLanguage,
			IsInline:        s.IsInline,
			IsPublic:        s.IsPublic,
			Status:          s.Status.String(),
			Reason:          s.Reason,
			CreatedAt:       s.CreatedAt.Format(time.RFC3339),
			ShareID:         s.ShareID,
			SourceSizeBytes: s.SubmissionSourceSizeBytes,
			TargetSizeBytes: s.SubmissionTargetSizeBytes,
		}
	}

	totalItems, err := h.databaseService.Submission.Query().
		Where(submission.HasUserWith(user.ID(tereusUser.ID))).
		Count(context.Background())
	if err != nil {
		return err
	}

	response := PaginatedResponse[*submissionsHistoryItem]{
		Items: items.Submissions,
		Meta: PaginatedMeta{
			ItemCount:    len(items.Submissions),
			TotalItems:   totalItems,
			ItemsPerPage: limit,
			TotalPages:   int(math.Ceil(float64(totalItems) / float64(limit))),
			CurrentPage:  page,
		},
	}

	return c.JSON(http.StatusOK, response)
}

// DELETE /users/me
func (h *UserHandler) DeleteCurrentUser(c echo.Context) error {
	loggedUser, err := h.tokenService.GetUserFromContext(c)
	if err != nil {
		return err
	}

	// Get all submissions
	submissions, err := h.databaseService.Submission.Query().
		Where(submission.HasUserWith(user.ID(loggedUser.ID))).
		All(context.Background())
	if err != nil {
		return err
	}

	logrus.WithField("submissions", len(submissions)).Info("Deleting submissions for user")

	// Delete all submissions
	for _, s := range submissions {
		err := h.s3Service.DeleteSubmission(s.ID.String())
		if err != nil {
			logrus.WithError(err).Error("Failed to delete submission")
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		err = h.databaseService.Submission.DeleteOneID(s.ID).Exec(context.Background())
		if err != nil {
			logrus.WithError(err).Error("Failed to delete submission")
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	// TODO cancel subscription if paid

	// Delete user (and cascade delete all relations)
	err = h.databaseService.User.DeleteOneID(loggedUser.ID).Exec(context.Background())
	if err != nil {
		logrus.WithError(err).Error("Failed to delete user")
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

// GET /users/me/export
func (h *UserHandler) GetExport(c echo.Context) error {
	loggedUser, err := h.tokenService.GetUserFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to export data")
	}

	file, err := ioutil.TempFile("", "export")
	if err != nil {
		logrus.WithError(err).Error("Failed to create temporary file")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to export data")
	}
	defer file.Close()
	defer os.Remove(file.Name())

	zipWriter := zip.NewWriter(file)

	userJSON, err := json.Marshal(loggedUser)
	if err != nil {
		logrus.WithError(err).Error("Failed to marshal user")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to export data")
	}

	userFile, err := zipWriter.Create("user.json")
	if err != nil {
		logrus.WithError(err).Error("Failed to create user file in zip")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to export data")
	}

	_, err = userFile.Write(userJSON)
	if err != nil {
		logrus.WithError(err).Error("Failed to write user file in zip")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to export data")
	}

	// Write submissions in zip file
	submissions, err := h.databaseService.Submission.Query().
		Where(submission.HasUserWith(user.ID(loggedUser.ID))).
		All(context.Background())
	if err != nil {
		logrus.WithError(err).Error("Failed to get submissions")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to export data")
	}

	submissionJSON, err := json.Marshal(submissions)
	if err != nil {
		logrus.WithError(err).Error("Failed to marshal submissions")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to export data")
	}

	submissionFile, err := zipWriter.Create("submissions.json")
	if err != nil {
		logrus.WithError(err).Error("Failed to create submissions file in zip")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to export data")
	}

	_, err = submissionFile.Write(submissionJSON)
	if err != nil {
		logrus.WithError(err).Error("Failed to write submissions file in zip")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to export data")
	}

	// Extract submissions source and results from S3
	for _, s := range submissions {
		objects := h.s3Service.ListSubmissionFiles(s.ID.String())
		if err != nil {
			logrus.WithError(err).Error("Failed to list submission files")
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to export data")
		}

		for object := range objects {
			o, err := h.s3Service.GetObject(object.Path)
			if err != nil {
				logrus.WithError(err).Error("Failed to get object")
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to export data")
			}

			codeFile, err := zipWriter.Create(object.Path)
			if err != nil {
				logrus.WithError(err).Error("Failed to create code file in zip")
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to export data")
			}

			_, err = io.Copy(codeFile, o)
			if err != nil {
				logrus.WithError(err).Error("Failed to copy object")
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to export data")
			}
		}
	}

	err = zipWriter.Close()
	if err != nil {
		logrus.WithError(err).Error("Failed to close zip writer")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to export data")
	}

	c.Response().Header().Set("Content-Type", "application/zip")
	c.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=export-user-%s.zip", loggedUser.ID.String()))

	_, err = file.Seek(0, 0)
	if err != nil {
		logrus.WithError(err).Error("Failed to seek file")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to export data")
	}

	_, err = io.Copy(c.Response(), file)
	if err != nil {
		logrus.WithError(err).Error("Failed to write to response")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to export data")
	}

	return nil
}
