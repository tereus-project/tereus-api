package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/tereus-project/tereus-api/ent"
	"github.com/tereus-project/tereus-api/services"
)

type tierPrices struct {
	BasePriceId    string
	MeteredPriceId string
}

type SubscriptionHandler struct {
	databaseService     *services.DatabaseService
	tokenService        *services.TokenService
	subscriptionService *services.SubscriptionService
}

func NewSubscriptionHandler(databaseService *services.DatabaseService, tokenService *services.TokenService, subscriptionService *services.SubscriptionService) (*SubscriptionHandler, error) {
	return &SubscriptionHandler{
		databaseService:     databaseService,
		tokenService:        tokenService,
		subscriptionService: subscriptionService,
	}, nil
}

type checkoutBody struct {
	Tier       string `json:"tier" validate:"required"`
	SuccessURL string `json:"success_url" validate:"required"`
	CancelURL  string `json:"cancel_url" validate:"required"`
}

type checkoutResponse struct {
	RedirectURL string `json:"redirect_url"`
}

// POST /subscription/checkout
func (h *SubscriptionHandler) CreateCheckoutSession(c echo.Context) error {
	user, err := h.tokenService.GetUserFromContext(c)
	if err != nil {
		return err
	}

	body := new(checkoutBody)

	if err := c.Bind(body); err != nil {
		return err
	}

	if err := c.Validate(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if !h.subscriptionService.HasTier(body.Tier) {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			fmt.Sprintf("Invalid tier '%s', must be one of %s", body.Tier, strings.Join(h.subscriptionService.GetTiers(), ", ")),
		)
	}

	lastUserSubscription, err := h.subscriptionService.GetLastUserSubscription(user.ID)
	if err != nil && err.(*ent.NotFoundError) == nil {
		logrus.WithError(err).Error("Failed to get last user subscription")
		return err
	}

	if lastUserSubscription != nil && lastUserSubscription.ExpiresAt.After(time.Now()) {
		if lastUserSubscription.Tier == body.Tier {
			return echo.NewHTTPError(http.StatusBadRequest, "You already are subscribed to this tier")
		}

		stripeSubscription, err := h.subscriptionService.UpdateStripeSubscription(lastUserSubscription.StripeSubscriptionID, body.Tier)
		if err != nil {
			logrus.WithError(err).Error("Failed to update subscription")
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update subscription")
		}

		h.subscriptionService.UpdateSubscription(stripeSubscription)

		return c.JSON(http.StatusOK, checkoutResponse{
			RedirectURL: "",
		})
	}

	checkoutSession, err := h.subscriptionService.CreateCheckoutSession(lastUserSubscription, &services.CheckoutSessionConfig{
		User:       user,
		Tier:       body.Tier,
		SuccessURL: body.SuccessURL,
		CancelURL:  body.CancelURL,
	})
	if err != nil {
		logrus.WithError(err).Error("failed to create checkout session")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create checkout session")
	}

	return c.JSON(http.StatusOK, checkoutResponse{
		RedirectURL: checkoutSession.URL,
	})
}

type portalBody struct {
	ReturnURL string `json:"return_url" validate:"required"`
}

type portalResponse struct {
	RedirectURL string `json:"redirect_url"`
}

// POST /subscription/portal
func (h *SubscriptionHandler) CreatePortalSession(c echo.Context) error {
	user, err := h.tokenService.GetUserFromContext(c)
	if err != nil {
		return err
	}

	body := new(portalBody)

	if err := c.Bind(body); err != nil {
		return err
	}

	if err := c.Validate(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	subscription, err := h.subscriptionService.GetCurrentUserSubscription(user.ID)
	if err != nil {
		if err.(*ent.NotFoundError) != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "You are not subscribed to any tier")
		}

		logrus.WithError(err).Error("Failed to get current user subscription")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get current user subscription")
	}

	portalSession, err := h.subscriptionService.CreatePortalSession(subscription, &services.PortalSessionConfig{
		ReturnURL: body.ReturnURL,
	})
	if err != nil {
		logrus.WithError(err).Error("Failed to create portal session")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create portal session")
	}

	return c.JSON(http.StatusOK, portalResponse{
		RedirectURL: portalSession.URL,
	})
}
