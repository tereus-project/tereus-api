package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/tereus-project/tereus-api/ent/user"
	"github.com/tereus-project/tereus-api/services"
)

type AuthHandler struct {
	databaseService *services.DatabaseService
	githubService   *services.GithubService
	gitlabService   *services.GitlabService
	tokenService    *services.TokenService
}

func NewAuthHandler(
	databaseService *services.DatabaseService,
	githubService *services.GithubService,
	gitlabService *services.GitlabService,
	tokenService *services.TokenService,
) (*AuthHandler, error) {
	return &AuthHandler{
		databaseService: databaseService,
		githubService:   githubService,
		gitlabService:   gitlabService,
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

	hasConflictingUser, err := h.databaseService.User.Query().
		Where(
			user.And(
				user.EmailEQ(email),
				user.Or(
					user.GithubUserIDIsNil(),
					user.GithubUserIDNEQ(githubUser.GetID()),
				),
			),
		).
		Exist(context.Background())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to check if there is a conflicting user")
	}

	if hasConflictingUser {
		return echo.NewHTTPError(http.StatusBadRequest, "There is already a user with this email, you can link your GitHub account in your account settings")
	}

	userId, err := h.databaseService.User.Create().
		SetEmail(email).
		SetGithubUserID(githubUser.GetID()).
		SetGithubAccessToken(githubAuth.AccessToken).
		OnConflict(
			sql.ConflictColumns(user.FieldEmail),
		).
		UpdateNewValues().
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

type gitlabSignupBody struct {
	Code        string `json:"code" validate:"required"`
	RedirectUri string `json:"redirect_uri" validate:"required"`
}

// /auth/login/gitlab
func (h *AuthHandler) GitlabLogin(c echo.Context) error {
	body := new(gitlabSignupBody)
	if err := c.Bind(body); err != nil {
		logrus.Error(err)
		return err
	}

	if err := c.Validate(body); err != nil {
		logrus.Error(err)
		return err
	}

	gitlabAuth, err := h.gitlabService.GenerateAccessTokenFromCode(body.Code, body.RedirectUri)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	gitlabClient, err := h.gitlabService.NewClient(gitlabAuth.AccessToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	gitlabUser, err := gitlabClient.GetUser()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	email := gitlabUser.PublicEmail
	if email == "" {
		email = gitlabUser.Email
	}

	if email == "" {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Failed to retrieve GitLab user email. You can try to revoke the access token at https://gitlab.com/settings/connections/applications/%s and retry.", h.githubService.ClientId))
	}

	hasConflictingUser, err := h.databaseService.User.Query().
		Where(
			user.And(
				user.EmailEQ(email),
				user.Or(
					user.GitlabUserIDIsNil(),
					user.GitlabUserIDNEQ(gitlabUser.ID),
				),
			),
		).
		Exist(context.Background())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to check if there is a conflicting user")
	}

	if hasConflictingUser {
		return echo.NewHTTPError(http.StatusBadRequest, "There is already a user with this email, you can link your GitLab account in your account settings")
	}

	userId, err := h.databaseService.User.Create().
		SetEmail(email).
		SetGitlabUserID(gitlabUser.ID).
		SetGitlabAccessToken(gitlabAuth.AccessToken).
		SetGitlabRefreshToken(gitlabAuth.RefreshToken).
		SetGitlabAccessTokenExpiresAt(time.UnixMilli((gitlabAuth.CreatedAt + gitlabAuth.ExpiresIn) * 1000)).
		OnConflict(
			sql.ConflictColumns(user.FieldEmail),
		).
		UpdateNewValues().
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
