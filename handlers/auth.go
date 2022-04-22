package handlers

import (
	"context"
	"fmt"
	"net/http"

	"entgo.io/ent/dialect/sql"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/tereus-project/tereus-api/ent/user"
	"github.com/tereus-project/tereus-api/services"
)

type AuthHandler struct {
	DatabaseService *services.DatabaseService
	GithubService   *services.GithubService
	TokenService    *services.TokenService
}

func NewAuthHandler(databaseService *services.DatabaseService, githubService *services.GithubService, tokenService *services.TokenService) (*AuthHandler, error) {
	return &AuthHandler{
		DatabaseService: databaseService,
		GithubService:   githubService,
		TokenService:    tokenService,
	}, nil
}

type classicSignupBody struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type signupResult struct {
	Token string `json:"token"`
}

// /auth/signup/classic
func (h *AuthHandler) ClassicSignup(c echo.Context) error {
	body := new(classicSignupBody)
	if err := c.Bind(body); err != nil {
		return err
	}

	if err := c.Validate(body); err != nil {
		return err
	}

	user, err := h.DatabaseService.User.Create().
		SetEmail(body.Email).
		SetPassword(body.Password).
		Save(context.Background())
	if err != nil {
		return err
	}

	token, err := h.TokenService.GenerateToken(user.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create token")
	}

	return c.JSON(200, signupResult{
		Token: token.String(),
	})
}

type githubSignupBody struct {
	Code string `json:"code" validate:"required"`
}

// /auth/login/github
func (h *AuthHandler) GithubLogin(c echo.Context) error {
	body := new(githubSignupBody)
	if err := c.Bind(body); err != nil {
		logrus.Error(err)
		return err
	}

	if err := c.Validate(body); err != nil {
		logrus.Error(err)
		return err
	}

	githubAuth, err := h.GithubService.GenerateAccessTokenFromCode(body.Code)
	if err != nil {
		return err
	}

	githubClient := h.GithubService.NewClient(githubAuth.AccessToken)
	githubUser, err := githubClient.GetUser()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	email := githubUser.GetEmail()
	if email == "" {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Failed to retrieve GitHub user email. Make sure to enable user:email scope when authenticating with GitHub. You can revoke the access token at https://github.com/settings/connections/applications/%s and retry.", h.GithubService.ClientId))
	}

	userId, err := h.DatabaseService.User.Create().
		SetEmail(email).
		SetGithubAccessToken(githubAuth.AccessToken).
		OnConflict(
			sql.ConflictColumns(user.FieldEmail),
			sql.UpdateWhere(
				sql.NotNull(user.FieldGithubAccessToken),
			),
		).
		UpdateGithubAccessToken().
		ID(context.Background())
	if err != nil {
		logrus.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create user")
	}

	token, err := h.TokenService.GenerateToken(userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create token")
	}

	return c.JSON(200, signupResult{
		Token: token.String(),
	})
}
