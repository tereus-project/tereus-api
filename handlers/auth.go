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
	databaseService *services.DatabaseService
	githubService   *services.GithubService
	tokenService    *services.TokenService
}

func NewAuthHandler(databaseService *services.DatabaseService, githubService *services.GithubService, tokenService *services.TokenService) (*AuthHandler, error) {
	return &AuthHandler{
		databaseService: databaseService,
		githubService:   githubService,
		tokenService:    tokenService,
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

	user, err := h.databaseService.User.Create().
		SetEmail(body.Email).
		SetPassword(body.Password).
		Save(context.Background())
	if err != nil {
		return err
	}

	token, err := h.tokenService.GenerateToken(user.ID)
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

	githubAuth, err := h.githubService.GenerateAccessTokenFromCode(body.Code)
	if err != nil {
		return err
	}

	githubClient := h.githubService.NewClient(githubAuth.AccessToken)
	githubUser, err := githubClient.GetUser()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	email := githubUser.GetEmail()
	if email == "" {
		emails := githubClient.GetEmails()

		if len(emails) == 0 {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Failed to retrieve GitHub user email. Make sure to enable user:email scope when authenticating with GitHub. You can revoke the access token at https://github.com/settings/connections/applications/%s and retry.", h.githubService.ClientId))
		}

		for _, e := range emails {
			if e.GetPrimary() {
				email = e.GetEmail()
				break
			}
		}
	}

	userId, err := h.databaseService.User.Create().
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

	token, err := h.tokenService.GenerateToken(userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create token")
	}

	return c.JSON(200, signupResult{
		Token: token.String(),
	})
}

// /auth/check
func (h *AuthHandler) Check(c echo.Context) error {
	valid, err := h.tokenService.VaidateTokenFromContext(c)
	if err != nil {
		return err
	}

	if !valid {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
	}

	return c.NoContent(http.StatusOK)
}
