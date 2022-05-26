package services

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/tereus-project/tereus-api/ent"
	"github.com/tereus-project/tereus-api/ent/token"
)

type TokenService struct {
	databaseService *DatabaseService
}

func NewTokenService(databaseService *DatabaseService) *TokenService {
	return &TokenService{
		databaseService: databaseService,
	}
}

func (s *TokenService) GenerateToken(userId uuid.UUID) (uuid.UUID, error) {
	token, err := s.databaseService.Token.Create().
		SetUserID(userId).
		Save(context.Background())

	if err != nil {
		return uuid.Nil, err
	}

	return token.ID, nil
}

func (s *TokenService) ValidateToken(tokenId uuid.UUID) (bool, error) {
	return s.databaseService.Token.Query().
		Where(
			token.ID(tokenId),
			token.IsActive(true),
		).
		Exist(context.Background())
}

func (s *TokenService) GetUser(tokenId uuid.UUID) (*ent.User, error) {
	return s.databaseService.Token.Query().
		Where(
			token.ID(tokenId),
			token.IsActive(true),
		).
		QueryUser().
		Only(context.Background())
}

func (s *TokenService) GetTokenFromContext(c echo.Context) (uuid.UUID, error) {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		authHeader = c.Request().URL.Query().Get("token")
	}

	if authHeader == "" {
		return uuid.UUID{}, fmt.Errorf("No authorization header provided")
	}

	token, err := uuid.Parse(strings.TrimPrefix(authHeader, "Bearer "))
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("Invalid token: %s", err.Error())
	}

	return token, nil
}

// Return a *echo.HTTPError if failing
func (s *TokenService) GetUserFromContext(c echo.Context) (*ent.User, error) {
	token, err := s.GetTokenFromContext(c)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	user, err := s.GetUser(token)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "Expired or invalid token")
	}

	return user, nil
}

// Return a *echo.HTTPError if failing
func (s *TokenService) VaidateTokenFromContext(c echo.Context) (bool, error) {
	token, err := s.GetTokenFromContext(c)
	if err != nil {
		return false, echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	valid, err := s.ValidateToken(token)
	if err != nil {
		return false, echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed check token validity: %s", err.Error()))
	}

	return valid, nil
}
