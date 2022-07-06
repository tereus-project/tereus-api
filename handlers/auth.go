package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

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

type signupResult struct {
	Token string `json:"token"`
}

type githubSignupBody struct {
	Code string `json:"code" validate:"required"`
}

// POST /auth/login/github
func (h *AuthHandler) LoginGithub(c echo.Context) error {
	tereusUser, _ := h.tokenService.GetUserFromContext(c)

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
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	githubClient := h.githubService.NewClient(githubAuth.AccessToken)
	githubUser, err := githubClient.GetUser()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Update user if logged in
	if tereusUser != nil {
		// Check if another user if using the same GitHub account
		hasConflictingUser, err := h.databaseService.User.Query().
			Where(
				user.And(
					user.GithubUserIDEQ(githubUser.GetID()),
					user.IDNEQ(tereusUser.ID),
				),
			).
			Exist(context.Background())
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to check if there is a conflicting user")
		}

		if hasConflictingUser {
			return echo.NewHTTPError(http.StatusBadRequest, "There is already a user with this GitHub account")
		}

		// Update the user
		_, err = h.databaseService.User.UpdateOneID(tereusUser.ID).
			SetGithubUserID(githubUser.GetID()).
			SetGithubAccessToken(githubAuth.AccessToken).
			Save(context.Background())
		if err != nil {
			logrus.Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update user")
		}

		token, _ := h.tokenService.GetTokenFromContext(c)
		return c.JSON(http.StatusOK, signupResult{
			Token: token.String(),
		})
	}

	// Update the user if they are already registered with this GitHub account
	existingUser, err := h.databaseService.User.Query().
		Where(user.GithubUserIDEQ(githubUser.GetID())).
		First(context.Background())
	if err == nil {
		_, err = h.databaseService.User.UpdateOneID(existingUser.ID).
			SetGithubAccessToken(githubAuth.AccessToken).
			Save(context.Background())
		if err != nil {
			logrus.Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update user")
		}

		token, err := h.tokenService.GenerateToken(existingUser.ID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create token")
		}

		return c.JSON(http.StatusOK, signupResult{
			Token: token.String(),
		})
	}

	email := githubUser.GetEmail()
	if email == "" {
		emails := githubClient.GetEmails()

		if len(emails) == 0 {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Failed to retrieve GitHub user email. Make sure to enable user:email scope when authenticating with GitHub. You can revoke the access token at https://github.com/settings/connections/applications/%s and retry", h.githubService.ClientId))
		}

		for _, e := range emails {
			if e.GetPrimary() {
				email = e.GetEmail()
				break
			}
		}
	}

	// Check if there is a user with this email
	hasEmailConflictingUser, err := h.databaseService.User.Query().
		Where(user.EmailEQ(email)).
		Exist(context.Background())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to check if there is a conflicting user")
	}

	if hasEmailConflictingUser {
		return echo.NewHTTPError(http.StatusBadRequest, "There is already a user with this email, you can link your GitHub account in your account settings")
	}

	// Save the new user
	newUser, err := h.databaseService.User.Create().
		SetEmail(email).
		SetGithubUserID(githubUser.GetID()).
		SetGithubAccessToken(githubAuth.AccessToken).
		Save(context.Background())
	if err != nil {
		logrus.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create user")
	}

	token, err := h.tokenService.GenerateToken(newUser.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create token")
	}

	return c.JSON(http.StatusOK, signupResult{
		Token: token.String(),
	})
}

type revokeResult struct {
	Success bool `json:"success"`
}

// POST /auth/revoke/github
func (h *AuthHandler) RevokeGithub(c echo.Context) error {
	tereusUser, err := h.tokenService.GetUserFromContext(c)
	if err != nil {
		return err
	}

	if tereusUser.GithubAccessToken == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "You don't have a GitHub account linked to your account")
	}

	if tereusUser.GitlabAccessToken == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "You don't have another provider linked to your account")
	}

	_, err = h.databaseService.User.UpdateOneID(tereusUser.ID).
		ClearGithubUserID().
		ClearGithubAccessToken().
		Save(context.Background())
	if err != nil {
		logrus.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update user")
	}

	return c.JSON(http.StatusOK, revokeResult{
		Success: true,
	})
}

type gitlabSignupBody struct {
	Code        string `json:"code" validate:"required"`
	RedirectUri string `json:"redirect_uri" validate:"required"`
}

// POST /auth/login/gitlab
func (h *AuthHandler) LoginGitlab(c echo.Context) error {
	tereusUser, _ := h.tokenService.GetUserFromContext(c)

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

	// Update user if logged in
	if tereusUser != nil {
		// Check if another user if using the same GitLab account
		hasConflictingUser, err := h.databaseService.User.Query().
			Where(
				user.And(
					user.GitlabUserIDEQ(gitlabUser.ID),
					user.IDNEQ(tereusUser.ID),
				),
			).
			Exist(context.Background())
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to check if there is a conflicting user")
		}

		if hasConflictingUser {
			return echo.NewHTTPError(http.StatusBadRequest, "There is already a user with this GitLab account")
		}

		// Update the user
		_, err = h.databaseService.User.UpdateOneID(tereusUser.ID).
			SetGitlabUserID(gitlabUser.ID).
			SetGitlabAccessToken(gitlabAuth.AccessToken).
			SetGitlabRefreshToken(gitlabAuth.RefreshToken).
			SetGitlabAccessTokenExpiresAt(time.UnixMilli((gitlabAuth.CreatedAt + gitlabAuth.ExpiresIn) * 1000)).
			Save(context.Background())
		if err != nil {
			logrus.Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update user")
		}

		token, _ := h.tokenService.GetTokenFromContext(c)
		return c.JSON(http.StatusOK, signupResult{
			Token: token.String(),
		})
	}

	// Update the user if they are already registered with this GitLab account
	existingUser, err := h.databaseService.User.Query().
		Where(user.GitlabUserIDEQ(gitlabUser.ID)).
		First(context.Background())
	if err == nil {
		_, err = h.databaseService.User.UpdateOneID(existingUser.ID).
			SetGitlabAccessToken(gitlabAuth.AccessToken).
			SetGitlabRefreshToken(gitlabAuth.RefreshToken).
			SetGitlabAccessTokenExpiresAt(time.UnixMilli((gitlabAuth.CreatedAt + gitlabAuth.ExpiresIn) * 1000)).
			Save(context.Background())
		if err != nil {
			logrus.Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update user")
		}

		token, err := h.tokenService.GenerateToken(existingUser.ID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create token")
		}

		return c.JSON(http.StatusOK, signupResult{
			Token: token.String(),
		})
	}

	email := gitlabUser.PublicEmail
	if email == "" {
		email = gitlabUser.Email
	}

	if email == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to retrieve GitLab user email. You can try to revoke the access token at https://gitlab.com/-/profile/applications and retry")
	}

	// Check if there is a user with this email
	hasEmailConflictingUser, err := h.databaseService.User.Query().
		Where(user.EmailEQ(email)).
		Exist(context.Background())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to check if there is a conflicting user")
	}

	if hasEmailConflictingUser {
		return echo.NewHTTPError(http.StatusBadRequest, "There is already a user with this email, you can link your GitLab account in your account settings")
	}

	// Save the new user
	newUser, err := h.databaseService.User.Create().
		SetEmail(email).
		SetGitlabUserID(gitlabUser.ID).
		SetGitlabAccessToken(gitlabAuth.AccessToken).
		SetGitlabRefreshToken(gitlabAuth.RefreshToken).
		SetGitlabAccessTokenExpiresAt(time.UnixMilli((gitlabAuth.CreatedAt + gitlabAuth.ExpiresIn) * 1000)).
		Save(context.Background())
	if err != nil {
		logrus.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create user")
	}

	token, err := h.tokenService.GenerateToken(newUser.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create token")
	}

	return c.JSON(http.StatusOK, signupResult{
		Token: token.String(),
	})
}

// POST /auth/revoke/gitlab
func (h *AuthHandler) RevokeGitlab(c echo.Context) error {
	tereusUser, err := h.tokenService.GetUserFromContext(c)
	if err != nil {
		return err
	}

	if tereusUser.GitlabAccessToken == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "You don't have a GitLab account linked to your account")
	}

	if tereusUser.GithubAccessToken == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "You don't have another provider linked to your account")
	}

	_, err = h.databaseService.User.UpdateOneID(tereusUser.ID).
		ClearGitlabUserID().
		ClearGitlabAccessToken().
		ClearGitlabRefreshToken().
		ClearGitlabAccessTokenExpiresAt().
		Save(context.Background())
	if err != nil {
		logrus.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update user")
	}

	return c.JSON(http.StatusOK, revokeResult{
		Success: true,
	})
}

type checkResult struct {
	Valid bool `json:"valid"`
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

	return c.JSON(http.StatusOK, checkResult{
		Valid: true,
	})
}
