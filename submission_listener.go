package main

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/tereus-project/tereus-api/ent/submission"
	"github.com/tereus-project/tereus-api/services"
)

type submissionStatusMessage struct {
	ID     string            `json:"id"`
	Status submission.Status `json:"status"`
	Reason string            `json:"reason"`
}

func startSubmissionStatusListener(k *services.KafkaService, databaseService *services.DatabaseService) error {
	go func() {
		for {
			msg, err := k.SubmissionStatusConsumer.ReadMessage(-1)
			if err != nil {
				logrus.WithError(err).Error("Failed to read message")
				continue
			}

			logrus.Info("Received submission status message")

			var message submissionStatusMessage
			err = json.Unmarshal(msg.Value, &message)
			if err != nil {
				logrus.WithError(err).Error("Failed to unmarshal submission status message")
				continue
			}

			id, err := uuid.Parse(message.ID)
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
				SetStatus(message.Status).
				SetReason(message.Reason).
				Exec(context.Background())
			if err != nil {
				logrus.WithError(err).Error("Failed to update submission status")
				continue
			}
		}
	}()

	return nil
}
