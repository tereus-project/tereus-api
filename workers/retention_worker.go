package workers

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tereus-project/tereus-api/ent/submission"
	"github.com/tereus-project/tereus-api/ent/subscription"
	"github.com/tereus-project/tereus-api/ent/user"
	"github.com/tereus-project/tereus-api/services"
)

func RetentionWorker(databaseService *services.DatabaseService, s3Service *services.S3Service) {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		logrus.Info("New retention worker iteration")

		// Query all submissions where user is in free tier and updated date is > 1 day
		submissions, err := databaseService.Submission.Query().
			Where(submission.CreatedAtLT(time.Now().AddDate(0, 0, -1))).
			Where(submission.Or(
				submission.HasUserWith(user.HasSubscriptionWith(subscription.TierEQ(subscription.TierFree))),
				submission.Not(submission.HasUserWith(user.HasSubscription())),
			)).
			Where(submission.StatusEQ(submission.StatusDone)).
			All(context.Background())
		if err != nil {
			logrus.WithError(err).Errorln("Failed to get submissions")
		}

		logrus.WithFields(logrus.Fields{
			"count": len(submissions),
		}).Infoln("Found submissions to delete")

		for _, sub := range submissions {
			err = databaseService.Submission.
				UpdateOneID(sub.ID).
				SetStatus(submission.StatusCleaned).
				Exec(context.Background())
			if err != nil {
				logrus.WithError(err).Error("Failed to update submission")
			}

			// Delete submission from S3
			err = s3Service.DeleteSubmission(sub.ID.String())
			if err != nil {
				logrus.WithField("submission_id", sub.ID).WithError(err).Errorln("Failed to delete submission from S3")
			}
		}
	}
}
