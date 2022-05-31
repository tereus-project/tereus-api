package handlers

import (
	"context"
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
