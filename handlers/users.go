package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
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

type submissionsHistory struct {
	Submissions []RemixResult `json:"submissions"`
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
		Submissions: make([]RemixResult, len(submissions)),
	}

	for i, s := range submissions {
		response.Submissions[i] = RemixResult{
			ID:             s.ID.String(),
			SourceLanguage: s.SourceLanguage,
			TargetLanguage: s.TargetLanguage,
			Status:         s.Status.String(),
			Reason:         s.Reason,
			CreatedAt:      s.CreatedAt.Format(time.RFC3339),
		}
	}

	return c.JSON(http.StatusOK, response)
}

// DELETE /submissions/:id
func (h *UserHandler) DeleteSubmission(c echo.Context) error {
	tereusUser, err := h.tokenService.GetUserFromContext(c)
	if err != nil {
		return err
	}

	submissionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid submission ID")
	}

	sub, err := h.databaseService.Submission.Get(context.Background(), submissionID)
	if err != nil {
		if ent.IsNotFound(err) {
			return echo.NewHTTPError(http.StatusNotFound, "Submission not found")
		}
		logrus.WithError(err).Error("Failed to get submission")
		return err
	}

	owner, err := sub.QueryUser().FirstID(c.Request().Context())
	if err != nil {
		logrus.WithError(err).Error("Failed to get submission owner")
		return err
	}

	if owner != tereusUser.ID {
		return echo.NewHTTPError(http.StatusForbidden, "You are not allowed to delete this submission")
	}

	// Already deleted, skip S3/DB deletion
	if sub.Status == submission.StatusDeleted {
		return c.NoContent(http.StatusNoContent)
	}

	// Delete from object storage
	err = h.s3Service.DeleteSubmission(sub.ID.String())
	if err != nil {
		logrus.WithError(err).Error("Failed to delete submission from object storage")
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Set submission as deleted
	err = h.databaseService.Submission.
		Update().
		Where(
			submission.ID(sub.ID),
		).
		SetStatus(submission.StatusDeleted).
		Exec(context.Background())
	if err != nil {
		logrus.WithError(err).Error("Failed to set submission as deleted")
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
