package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tereus-project/tereus-api/services"
)

type UserHandler struct {
	TokenService *services.TokenService
}

func NewUserHandler(tokenService *services.TokenService) (*UserHandler, error) {
	return &UserHandler{
		TokenService: tokenService,
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
