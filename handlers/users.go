package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/tereus-project/tereus-api/ent"
	"github.com/tereus-project/tereus-api/ent/submission"
	"github.com/tereus-project/tereus-api/ent/user"
	"github.com/tereus-project/tereus-api/services"
)

type UserHandler struct {
	TokenService    *services.TokenService
	DatabaseService *services.DatabaseService
}

func NewUserHandler(tokenService *services.TokenService, databaseService *services.DatabaseService) (*UserHandler, error) {
	return &UserHandler{
		TokenService:    tokenService,
		DatabaseService: databaseService,
	}, nil
}

type getCurrentUserResult struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

// /users/me
func (h *UserHandler) GetCurrentUser(c echo.Context) error {
	user, err := h.TokenService.GetUserFromContext(c)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, getCurrentUserResult{
		ID:    user.ID.String(),
		Email: user.Email,
	})
}

type submissionsHistory struct {
	Submissions []RemixResult `json:"submissions"`
}

func (h *UserHandler) GetSubmissionsHistory(c echo.Context) error {
	tereusUser, err := h.TokenService.GetUserFromContext(c)
	if err != nil {
		return err
	}

	submissions, err := h.DatabaseService.Submission.Query().
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
			CreatedAt:      s.CreatedAt.Format(time.RFC3339),
			SourceLanguage: s.SourceLanguage,
			TargetLanguage: s.TargetLanguage,
			Status:         s.Status.String(),
		}
	}

	return c.JSON(http.StatusOK, response)
}
