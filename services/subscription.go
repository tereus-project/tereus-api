package services

import (
	"context"
	"fmt"
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

func (s *SubscriptionService) GetOrCreateCustomer(subscribingUser *ent.User, lastUserSubscription *ent.Subscription) (*stripe.Customer, *ent.Subscription, error) {
	var err error
	var stripeCustomer *stripe.Customer
	var stripeSubscription *stripe.Subscription

	// Retrieve Stripe data from the local saved sbscription details
	if lastUserSubscription != nil {
		customerParams := &stripe.CustomerParams{}

		stripeCustomer, err = customer.Get(lastUserSubscription.StripeCustomerID, customerParams)
		if err != nil {
			return nil, lastUserSubscription, err
		}

		if lastUserSubscription.StripeSubscriptionID != "" {
			stripeSubscription, err = sub.Get(lastUserSubscription.StripeSubscriptionID, nil)
			if err != nil {
				return nil, lastUserSubscription, err
			}
		}
	}

	if stripeCustomer != nil && stripeSubscription != nil {
		return stripeCustomer, lastUserSubscription, nil
	}

	// If there is no local saved customer, then we need try to retrieve the
	// existing customer from Stripe. If there is no customer in Stripe, then
	// we need to create a new customer.
	if stripeCustomer == nil {
		searchParams := &stripe.CustomerSearchParams{
			SearchParams: stripe.SearchParams{
				Query: fmt.Sprintf("metadata['user_id']:'%s'", subscribingUser.ID.String()),
			},
		}
		searchParams.AddExpand("subscriptions")
		customerSearchIter := customer.Search(searchParams)

		if customerSearchIter.Next() {
			stripeCustomer = customerSearchIter.Customer()
		} else {
			customerParams := &stripe.CustomerParams{
				Email: stripe.String(subscribingUser.Email),
			}
			customerParams.AddMetadata("user_id", subscribingUser.ID.String())

			var err error
			stripeCustomer, err = customer.New(customerParams)
			if err != nil {
				return nil, lastUserSubscription, err
			}
		}
	}

	tier := "free"

	// If there is no local saved subscription but the customer has one,
	// then we need try to retrieve the existing subscription from Stripe.
	if stripeSubscription == nil && stripeCustomer.Subscriptions != nil {
		for _, stripeCustomerSubscription := range stripeCustomer.Subscriptions.Data {
			maybeTier := s.GetTierFromPriceId(stripeCustomerSubscription.Items.Data[0].Price.ID)
			if maybeTier == "" {
				continue
			}

			stripeSubscription = stripeCustomerSubscription
			tier = maybeTier
			break
		}
	}

	subscriptionCreate := s.databaseService.Subscription.Create().
		SetStripeCustomerID(stripeCustomer.ID).
		SetUser(subscribingUser)

	// Only save subscription details if there is an unsaved Stripe subscription
	if stripeSubscription != nil && stripeCustomer.Subscriptions != nil {
		subscriptionCreate.
			SetStripeSubscriptionID(stripeSubscription.ID).
			SetExpiresAt(time.Unix(stripeSubscription.CurrentPeriodEnd, 0)).
			SetCancelled(stripeSubscription.CancelAtPeriodEnd).
			SetTier(subscription.Tier(tier))
	}

	subscriptionId, err := subscriptionCreate.
		OnConflict(
			sql.ConflictColumns("user_subscription"),
		).
		UpdateNewValues().
		ID(context.Background())
	if err != nil {
		return stripeCustomer, nil, err
	}

	lastUserSubscription, err = s.databaseService.Subscription.Get(context.Background(), subscriptionId)
	if err != nil {
		return stripeCustomer, lastUserSubscription, err
	}

	return stripeCustomer, lastUserSubscription, nil
}

type CheckoutSessionConfig struct {
	User       *ent.User
	Tier       string
	SuccessURL string
	CancelURL  string
}

func (s *SubscriptionService) CreateCheckoutSession(stripeCustomer *stripe.Customer, config *CheckoutSessionConfig) (*stripe.CheckoutSession, error) {
	tierPrices := s.stripeTierPriceIds[config.Tier]

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
		SuccessURL:          stripe.String(config.SuccessURL),
		CancelURL:           stripe.String(config.CancelURL),
		Customer:            stripe.String(stripeCustomer.ID),
		AllowPromotionCodes: stripe.Bool(true),
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

	userId, err := uuid.Parse(stripeCustomer.Metadata["user_id"])
	if err != nil {
		return err
	}

	subscribingUser, err := s.databaseService.User.Get(context.Background(), userId)
	if err != nil {
		return err
	}

	_, err = s.databaseService.Subscription.Create().
		SetStripeCustomerID(stripeCustomer.ID).
		SetStripeSubscriptionID(stripeSubscription.ID).
		SetTier(subscription.Tier(s.GetTierFromPriceId(stripeSubscription.Items.Data[0].Price.ID))).
		SetExpiresAt(time.Unix(stripeSubscription.CurrentPeriodEnd, 0)).
		SetCancelled(stripeSubscription.CancelAtPeriodEnd).
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

func (s *SubscriptionService) CancelStripeSubscription(stripeSubscriptionId string) error {
	stripeSubscription, err := sub.Get(stripeSubscriptionId, nil)
	if err != nil {
		return err
	}

	_, err = sub.Update(stripeSubscription.ID, &stripe.SubscriptionParams{
		CancelAtPeriodEnd: stripe.Bool(true),
	})

	return err
}

func (s *SubscriptionService) CancelSubscription(userId uuid.UUID, expiresAt time.Time) error {
	_, err := s.databaseService.Subscription.Update().
		SetExpiresAt(expiresAt).
		SetCancelled(true).
		Where(
			subscription.HasUserWith(
				user.ID(userId),
			),
		).
		Save(context.Background())

	return err
}

func (s *SubscriptionService) CancelSubscriptionFromStripeCustomerId(stripeCustomerId string, expiresAt time.Time) error {
	subscribingUser, err := s.databaseService.User.Query().
		Where(
			user.HasSubscriptionWith(
				subscription.StripeCustomerID(stripeCustomerId),
			),
		).
		Only(context.Background())
	if err != nil {
		return err
	}

	return s.CancelSubscription(subscribingUser.ID, expiresAt)
}
