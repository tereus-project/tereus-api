package services

import (
	"context"
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

func (s *TokenService) GetUser(tokenId uuid.UUID) (*ent.User, error) {
	return s.databaseService.Token.Query().
		Where(token.ID(tokenId)).
		QueryUser().
		Only(context.Background())
}

func (s *TokenService) GetUserFromContext(c echo.Context) (*ent.User, error) {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "No authorization header provided")
	}

	token, err := uuid.Parse(strings.TrimPrefix(authHeader, "Bearer "))
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
	}

	user, err := s.GetUser(token)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "Expired or invalid token")
	}

	return user, nil
}
