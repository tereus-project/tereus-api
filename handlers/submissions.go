package handlers

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/sirupsen/logrus"
	"github.com/tereus-project/tereus-api/ent"
	"github.com/tereus-project/tereus-api/ent/submission"
	"github.com/tereus-project/tereus-api/services"
)

type SubmissionsHandler struct {
	databaseService *services.DatabaseService
	tokenService    *services.TokenService
	storageService  *services.StorageService
}

func NewSubmissionsHandler(databaseService *services.DatabaseService, tokenService *services.TokenService, storageService *services.StorageService) (*SubmissionsHandler, error) {
	return &SubmissionsHandler{
		databaseService: databaseService,
		tokenService:    tokenService,
		storageService:  storageService,
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
	err = h.storageService.DeleteSubmission(sub.ID.String())
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
	ShareID  string `json:"share_id"`
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

	share_id := ""
	if body.IsPublic {
		share_id, err = gonanoid.Generate("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", 8)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate share ID")
		}

		err := h.databaseService.Submission.
			UpdateOneID(sub.ID).
			SetIsPublic(true).
			SetShareID(share_id).
			Exec(context.Background())
		if err != nil {
			logrus.WithError(err).Error("Failed to update submission visibility")
			return err
		}
	} else {
		err := h.databaseService.Submission.
			UpdateOneID(sub.ID).
			SetIsPublic(false).
			ClearShareID().
			Exec(context.Background())
		if err != nil {
			logrus.WithError(err).Error("Failed to update submission visibility")
			return err
		}
	}

	return c.JSON(http.StatusOK, &updateSubmissionVisibilityResponse{
		Id:       sub.ID.String(),
		IsPublic: body.IsPublic,
		ShareID:  share_id,
	})
}
