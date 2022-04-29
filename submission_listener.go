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

func startSubmissionStatusListener(rabbitMQService *services.RabbitMQService, databaseService *services.DatabaseService) error {
	queue, err := rabbitMQService.NewQueue("submission_status_q", "submission_status_ex", "#")
	if err != nil {
		return err
	}

	deliveries, err := queue.Consume()
	if err != nil {
		return err
	}

	go func() {
		for delivery := range deliveries {
			logrus.Info("Received submission status message")

			var message submissionStatusMessage
			if err := json.Unmarshal(delivery.Body, &message); err != nil {
				logrus.WithError(err).Error("Failed to unmarshal submission status message")

				if err := delivery.Nack(false, true); err != nil {
					logrus.WithError(err).Error("Failed to nack submission status message")
				}

				continue
			}

			id, err := uuid.Parse(message.ID)
			if err != nil {
				logrus.WithError(err).Error("Failed to parse submission ID")

				if err := delivery.Nack(false, false); err != nil {
					logrus.WithError(err).Error("Failed to nack submission status message")
				}

				continue
			}

			err = databaseService.Submission.
				UpdateOneID(id).
				SetStatus(message.Status).
				SetReason(message.Reason).
				Exec(context.Background())
			if err != nil {
				logrus.WithError(err).Error("Failed to update submission status")

				if err := delivery.Nack(false, true); err != nil {
					logrus.WithError(err).Error("Failed to nack submission status message")
				}

				continue
			}

			if err := delivery.Ack(false); err != nil {
				logrus.WithError(err).Error("Failed to acknowledge submission status message")
			}
		}
	}()

	return nil
}
