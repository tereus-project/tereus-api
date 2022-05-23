package services

import (
	"context"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v72"
	portalsession "github.com/stripe/stripe-go/v72/billingportal/session"
	"github.com/stripe/stripe-go/v72/checkout/session"
	"github.com/stripe/stripe-go/v72/customer"
	"github.com/stripe/stripe-go/v72/sub"
	"github.com/tereus-project/tereus-api/ent"
	"github.com/tereus-project/tereus-api/ent/subscription"
	"github.com/tereus-project/tereus-api/ent/user"
	"github.com/tereus-project/tereus-api/env"
)

type TierPrices struct {
	BasePriceId    string
	MeteredPriceId string
}

type SubscriptionService struct {
	databaseService *DatabaseService
	stripeService   *StripeService

	stripeTierPriceIds map[string]TierPrices
}

func NewSubscriptionService(databaseService *DatabaseService, stripeService *StripeService) *SubscriptionService {
	config := env.Get()

	return &SubscriptionService{
		databaseService: databaseService,
		stripeService:   stripeService,

		stripeTierPriceIds: map[string]TierPrices{
			"pro": {
				BasePriceId:    config.StripeTierProBase,
				MeteredPriceId: config.StripeTierProMetered,
			},
			"enterprise": {
				BasePriceId:    config.StripeTierEnterpriseBase,
				MeteredPriceId: config.StripeTierEnterpriseMetered,
			},
		},
	}
}

func (s *SubscriptionService) HasTier(tier string) bool {
	_, ok := s.stripeTierPriceIds[tier]
	return ok
}

func (s *SubscriptionService) GetTiers() []string {
	tiers := make([]string, 0, len(s.stripeTierPriceIds))
	for tier := range s.stripeTierPriceIds {
		tiers = append(tiers, tier)
	}

	return tiers
}

func (s *SubscriptionService) GetTierPrices(tier string) TierPrices {
	return s.stripeTierPriceIds[tier]
}

func (s *SubscriptionService) GetTierFromPriceId(priceId string) string {
	for tier, tierPrices := range s.stripeTierPriceIds {
		if tierPrices.BasePriceId == priceId || tierPrices.MeteredPriceId == priceId {
			return tier
		}
	}

	return ""
}

func (s *SubscriptionService) GetCurrentUserSubscription(userID uuid.UUID) (*ent.Subscription, error) {
	return s.databaseService.Subscription.Query().
		Where(
			subscription.HasUserWith(
				user.ID(userID),
			),
			subscription.ExpiresAtGTE(time.Now()),
		).
		Only(context.Background())
}

func (s *SubscriptionService) GetLastUserSubscription(userID uuid.UUID) (*ent.Subscription, error) {
	return s.databaseService.Subscription.Query().
		Where(
			subscription.HasUserWith(
				user.ID(userID),
			),
		).
		Only(context.Background())
}

func (s *SubscriptionService) getOrCreateCustomer(user *ent.User, lastUserSubscription *ent.Subscription) (*stripe.Customer, error) {
	if lastUserSubscription != nil {
		return customer.Get(lastUserSubscription.StripeCustomerID, nil)
	}

	customerParams := &stripe.CustomerParams{
		Email: stripe.String(user.Email),
	}

	customerParams.AddMetadata("user_id", user.ID.String())

	return customer.New(customerParams)
}

type CheckoutSessionConfig struct {
	User       *ent.User
	Tier       string
	SuccessURL string
	CancelURL  string
}

func (s *SubscriptionService) CreateCheckoutSession(lastUserSubscription *ent.Subscription, config *CheckoutSessionConfig) (*stripe.CheckoutSession, error) {
	tierPrices := s.stripeTierPriceIds[config.Tier]

	stripeCustomer, err := s.getOrCreateCustomer(config.User, lastUserSubscription)
	if err != nil {
		return nil, err
	}

	checkoutParams := &stripe.CheckoutSessionParams{
		Mode: stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(tierPrices.BasePriceId),
				Quantity: stripe.Int64(1),
			},
			{
				Price: stripe.String(tierPrices.MeteredPriceId),
			},
		},
		SuccessURL: stripe.String(config.SuccessURL),
		CancelURL:  stripe.String(config.CancelURL),
		Customer:   stripe.String(stripeCustomer.ID),
	}

	return session.New(checkoutParams)
}

type PortalSessionConfig struct {
	ReturnURL string
}

func (s *SubscriptionService) CreatePortalSession(currentUserSubscription *ent.Subscription, config *PortalSessionConfig) (*stripe.BillingPortalSession, error) {
	stripeCustomer, err := customer.Get(currentUserSubscription.StripeCustomerID, nil)
	if err != nil {
		return nil, err
	}

	return portalsession.New(&stripe.BillingPortalSessionParams{
		Customer:  stripe.String(stripeCustomer.ID),
		ReturnURL: stripe.String(config.ReturnURL),
	})
}

func (s *SubscriptionService) UpdateStripeSubscription(stripeSubscriptionId string, targetTier string) (*stripe.Subscription, error) {
	stripeSubscription, err := sub.Get(stripeSubscriptionId, nil)
	if err != nil {
		return nil, err
	}

	prices := s.GetTierPrices(targetTier)

	return sub.Update(stripeSubscription.ID, &stripe.SubscriptionParams{
		CancelAtPeriodEnd: stripe.Bool(false),
		ProrationBehavior: stripe.String(string(stripe.SubscriptionProrationBehaviorCreateProrations)),
		Items: []*stripe.SubscriptionItemsParams{
			{
				ID:    stripe.String(stripeSubscription.Items.Data[0].ID),
				Price: stripe.String(prices.BasePriceId),
			},
			{
				ID:    stripe.String(stripeSubscription.Items.Data[1].ID),
				Price: stripe.String(prices.MeteredPriceId),
			},
		},
	})
}

func (s *SubscriptionService) UpdateSubscription(stripeSubscription *stripe.Subscription) error {
	stripeCustomer, err := customer.Get(stripeSubscription.Customer.ID, nil)
	if err != nil {
		return err
	}

	id, err := uuid.Parse(stripeCustomer.Metadata["user_id"])
	if err != nil {
		return err
	}

	subscribingUser, err := s.databaseService.User.Get(context.Background(), id)
	if err != nil {
		return err
	}

	_, err = s.databaseService.Subscription.Create().
		SetStripeCustomerID(stripeCustomer.ID).
		SetStripeSubscriptionID(stripeSubscription.ID).
		SetTier(s.GetTierFromPriceId(stripeSubscription.Items.Data[0].Price.ID)).
		SetExpiresAt(time.Unix(stripeSubscription.CurrentPeriodEnd, 0)).
		SetUser(subscribingUser).
		OnConflict(
			sql.ConflictColumns("user_subscription"),
		).
		UpdateNewValues().
		ID(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (s *SubscriptionService) ExpireSubscriptionFromInvoice(invoice *stripe.Invoice) error {
	stripeCustomer, err := customer.Get(invoice.Customer.ID, nil)
	if err != nil {
		return err
	}

	id, err := uuid.Parse(stripeCustomer.Metadata["user_id"])
	if err != nil {
		return err
	}

	subscribedUser, err := s.databaseService.User.Get(context.Background(), id)
	if err != nil {
		return err
	}

	_, err = s.databaseService.Subscription.Update().
		SetExpiresAt(time.Now()).
		Where(
			subscription.HasUserWith(
				user.ID(subscribedUser.ID),
			),
		).
		Save(context.Background())

	return nil
}
