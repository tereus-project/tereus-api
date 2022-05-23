package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/sub"
	"github.com/tereus-project/tereus-api/services"
)

type StripeWebhooksHandler struct {
	databaseService     *services.DatabaseService
	subscriptionService *services.SubscriptionService
	stripeService       *services.StripeService

	endpointSecret string
}

func NewStripeWebhooksHandler(databaseService *services.DatabaseService, subscriptionService *services.SubscriptionService, stripeService *services.StripeService, endpointSecret string) (*StripeWebhooksHandler, error) {
	return &StripeWebhooksHandler{
		databaseService:     databaseService,
		subscriptionService: subscriptionService,
		endpointSecret:      endpointSecret,
	}, nil
}

func (h *StripeWebhooksHandler) HandleWebhooks(c echo.Context) error {
	event, err := h.stripeService.ConstructWebhookEvent(c.Response().Writer, c.Request(), h.endpointSecret)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	switch event.Type {
	case "invoice.paid":
		var invoice stripe.Invoice
		err := json.Unmarshal(event.Data.Raw, &invoice)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		stripeSubscription, err := sub.Get(invoice.Subscription.ID, nil)
		if err != nil {
			logrus.WithError(err).Error("Failed to fetch subscription from invoice")
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		err = h.subscriptionService.UpdateSubscription(stripeSubscription)
		if err != nil {
			logrus.WithError(err).Error("Failed to create subscription from invoice")
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	case "customer.subscription.deleted":
		var invoice stripe.Invoice
		err := json.Unmarshal(event.Data.Raw, &invoice)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		err = h.subscriptionService.ExpireSubscriptionFromInvoice(&invoice)
		if err != nil {
			logrus.WithError(err).Error("Failed to expire subscription from invoice")
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.NoContent(200)
}
