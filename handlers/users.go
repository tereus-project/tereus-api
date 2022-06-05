package handlers

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

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
	ID           string                            `json:"id"`
	Email        string                            `json:"email"`
	Subscription *getCurrentUserResultSubscription `json:"subscription"`
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

	return c.JSON(http.StatusOK, getCurrentUserResult{
		ID:           loggedUser.ID.String(),
		Email:        loggedUser.Email,
		Subscription: subscriptionResult,
	})
}

type submissionsHistoryItem struct {
	ID             string `json:"id"`
	SourceLanguage string `json:"source_language"`
	TargetLanguage string `json:"target_language"`
	IsInline       bool   `json:"is_inline"`
	IsPublic       bool   `json:"is_public"`
	Status         string `json:"status"`
	Reason         string `json:"reason"`
	CreatedAt      string `json:"created_at"`
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

	submissions, err := h.databaseService.Submission.Query().
		Where(submission.HasUserWith(user.ID(tereusUser.ID))).
		Order(ent.Desc(submission.FieldCreatedAt)).
		All(context.Background())
	if err != nil {
		return err
	}

	response := submissionsHistory{
		Submissions: make([]*submissionsHistoryItem, len(submissions)),
	}

	for i, s := range submissions {
		response.Submissions[i] = &submissionsHistoryItem{
			ID:             s.ID.String(),
			SourceLanguage: s.SourceLanguage,
			TargetLanguage: s.TargetLanguage,
			IsInline:       s.IsInline,
			IsPublic:       s.IsPublic,
			Status:         s.Status.String(),
			Reason:         s.Reason,
			CreatedAt:      s.CreatedAt.Format(time.RFC3339),
		}
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
		return err
	}

	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	userJSON, err := json.Marshal(loggedUser)
	if err != nil {
		return err
	}

	userFile, err := zipWriter.Create("user.json")
	if err != nil {
		return err
	}

	_, err = userFile.Write(userJSON)
	if err != nil {
		return err
	}

	// Write submissions in zip file
	submissions, err := h.databaseService.Submission.Query().
		Where(submission.HasUserWith(user.ID(loggedUser.ID))).
		All(context.Background())
	if err != nil {
		return err
	}

	for _, s := range submissions {
		submissionJSON, err := json.Marshal(s)
		if err != nil {
			return err
		}

		submissionFile, err := zipWriter.Create(s.ID.String() + ".json")
		if err != nil {
			return err
		}

		_, err = submissionFile.Write(submissionJSON)
		if err != nil {
			return err
		}
	}

	err = zipWriter.Close()
	if err != nil {
		logrus.WithError(err).Error("Failed to close zip writer")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to export data")
	}

	c.Response().Header().Set("Content-Type", "application/zip")
	c.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=export-user-%s.zip", loggedUser.ID.String()))

	_, err = io.Copy(c.Response(), buf)
	if err != nil {
		logrus.WithError(err).Error("Failed to write to response")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to export data")
	}

	return nil
}
