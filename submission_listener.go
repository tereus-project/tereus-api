package main

import (
	"context"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/tereus-project/tereus-api/ent/submission"
	"github.com/tereus-project/tereus-api/services"
)

func startSubmissionStatusListener(k *services.KafkaService, databaseService *services.DatabaseService) error {
	go func() {
		ch := k.ConsumeSubmissionStatus(context.Background())

		for {
			msg := <-ch

			logrus.Info("Received submission status message")

			id, err := uuid.Parse(msg.ID)
			if err != nil {
				logrus.WithError(err).Error("Failed to parse submission ID")
				continue
			}

			err = databaseService.Submission.
				Update().
				Where(
					submission.ID(id),
					submission.StatusIn(submission.StatusPending, submission.StatusProcessing),
				).
				SetStatus(msg.Status).
				SetReason(msg.Reason).
				Exec(context.Background())
			if err != nil {
				logrus.WithError(err).Error("Failed to update submission status")
				continue
			}
		}
	}()

	return nil
}
