package main

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/sub"
	"github.com/stripe/stripe-go/v72/usagerecord"
	"github.com/tereus-project/tereus-api/ent/submission"
	"github.com/tereus-project/tereus-api/ent/user"
	"github.com/tereus-project/tereus-api/env"
	"github.com/tereus-project/tereus-api/services"
)

const mb = 1024 * 1024

func startSubscriptionDataUsageReportingWorker(subscriptionService *services.SubscriptionService, databaseService *services.DatabaseService, s3Service *services.S3Service) {
	config := env.Get()

	go func() {
		subscriptionsBatchSize := 100
		ticker := time.NewTicker(time.Hour)
		defer ticker.Stop()

		for {
			for page := 0; ; page++ {
				subscriptions, err := subscriptionService.GetActiveSubscriptions(page*subscriptionsBatchSize, (page+1)*subscriptionsBatchSize)
				if err != nil {
					logrus.WithError(err).Errorln("Failed to get active subscriptions")
					continue
				}

				if len(subscriptions) == 0 {
					break
				}

				for _, subscription := range subscriptions {
					userId, err := subscription.QueryUser().OnlyID(context.Background())
					if err != nil {
						logrus.WithError(err).Errorf("Failed to get user ID for subscription %s", subscription.ID)
						continue
					}

					submissions, err := databaseService.Submission.Query().
						Where(
							submission.HasUserWith(
								user.ID(userId),
							),
						).
						All(context.Background())
					if err != nil {
						logrus.WithError(err).Errorf("Failed to get submission IDs for subscription %s", subscription.ID)
						continue
					}

					var totalBytes int64

					for _, submissionRow := range submissions {
						submissionBytesCount := s3Service.SizeofObjects(fmt.Sprintf("remix/%s/", submissionRow.ID))

						var submissionResultBytesCount int64
						if submissionRow.Status == submission.StatusDone {
							submissionResultBytesCount = s3Service.SizeofObjects(fmt.Sprintf("remix-results/%s/", submissionRow.ID.String()))
						}

						totalBytes += submissionBytesCount + submissionResultBytesCount
					}

					totalMegaBytes := totalBytes / mb

					logrus.Infof("Reporting usage for subscription %s: %dMB", subscription.ID, totalMegaBytes)

					stripeSubscription, err := sub.Get(subscription.StripeSubscriptionID, nil)
					if err != nil {
						logrus.WithError(err).Errorf("Failed to get stripe subscription %s", subscription.StripeSubscriptionID)
						continue
					}

					var dataRententionSubscriptionItem *stripe.SubscriptionItem
					for _, item := range stripeSubscription.Items.Data {
						if item.Price.ID == config.StripeTierProMetered || item.Price.ID == config.StripeTierEnterpriseMetered {
							dataRententionSubscriptionItem = item
							break
						}
					}

					if dataRententionSubscriptionItem == nil {
						logrus.Errorf("Failed to find data rentention subscription item for subscription %s", subscription.ID)
						continue
					}

					_, err = usagerecord.New(&stripe.UsageRecordParams{
						SubscriptionItem: stripe.String(dataRententionSubscriptionItem.ID),
						Quantity:         stripe.Int64(totalMegaBytes),
						Timestamp:        stripe.Int64(time.Now().Unix()),
						Action:           stripe.String(string(stripe.UsageRecordActionSet)),
					})

					if err != nil {
						logrus.WithError(err).Errorf("Failed to report usage for subscription %s", subscription.ID)
						continue
					}
				}
			}

			<-ticker.C
		}
	}()
}
