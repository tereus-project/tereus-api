package handlers

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/tereus-project/tereus-api/ent"
	"github.com/tereus-project/tereus-api/ent/submission"
	"github.com/tereus-project/tereus-api/services"
)

type SubmissionsHandler struct {
	databaseService *services.DatabaseService
	tokenService    *services.TokenService
	s3Service       *services.S3Service
}

func NewSubmissionsHandler(databaseService *services.DatabaseService, tokenService *services.TokenService, s3Service *services.S3Service) (*SubmissionsHandler, error) {
	return &SubmissionsHandler{
		databaseService: databaseService,
		tokenService:    tokenService,
		s3Service:       s3Service,
	}, nil
}

// DELETE /submissions/:id
func (h *SubmissionsHandler) DeleteSubmission(c echo.Context) error {
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
	if sub.Status == submission.StatusCleaned {
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
		SetStatus(submission.StatusCleaned).
		Exec(context.Background())
	if err != nil {
		logrus.WithError(err).Error("Failed to set submission as deleted")
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

type updateSubmissionVisibilityBody struct {
	IsPublic bool `json:"is_public"`
}

type updateSubmissionVisibilityResponse struct {
	Id       string `json:"id"`
	IsPublic bool   `json:"is_public"`
}

// PATCH /submissions/:id/visibility
func (h *SubmissionsHandler) UpdateSubmissionVisibility(c echo.Context) error {
	tereusUser, err := h.tokenService.GetUserFromContext(c)
	if err != nil {
		return err
	}

	body := new(updateSubmissionVisibilityBody)

	if err := c.Bind(body); err != nil {
		return err
	}

	if err := c.Validate(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
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
		return echo.NewHTTPError(http.StatusForbidden, "You are not allowed to update this submission")
	}

	err = h.databaseService.Submission.
		UpdateOneID(sub.ID).
		SetIsPublic(body.IsPublic).
		Exec(context.Background())
	if err != nil {
		logrus.WithError(err).Error("Failed to update submission visibility")
		return err
	}

	return c.JSON(http.StatusOK, &updateSubmissionVisibilityResponse{
		Id:       sub.ID.String(),
		IsPublic: body.IsPublic,
	})
}
