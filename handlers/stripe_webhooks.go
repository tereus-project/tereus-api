package handlers

import (
	"encoding/json"
	"net/http"
	"time"

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
	case "customer.subscription.updated":
		var subscription stripe.Subscription
		err := json.Unmarshal(event.Data.Raw, &subscription)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		if subscription.CancelAtPeriodEnd {
			err = h.subscriptionService.CancelSubscriptionFromStripeCustomerId(subscription.Customer.ID, time.Unix(subscription.CancelAt, 0))
			if err != nil {
				logrus.WithError(err).Error("Failed to cancel subscription")
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
		} else if subscription.Status == "active" {
			err = h.subscriptionService.UpdateSubscription(&subscription)
			if err != nil {
				logrus.WithError(err).Error("Failed to update subscription")
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
		}
	case "customer.subscription.deleted":
		var invoice stripe.Invoice
		err := json.Unmarshal(event.Data.Raw, &invoice)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		err = h.subscriptionService.CancelSubscriptionFromStripeCustomerId(invoice.Customer.ID, time.Unix(invoice.PeriodEnd, 0))
		if err != nil {
			logrus.WithError(err).Error("Failed to expire subscription from invoice")
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	case "customer.deleted":
		var customer stripe.Customer
		err := json.Unmarshal(event.Data.Raw, &customer)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		err = h.subscriptionService.RemoveSubscriptionsFromStripeCustomerId(customer.ID)
		if err != nil {
			logrus.WithError(err).Error("Failed to remove customer's subscriptions")
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.NoContent(200)
}
