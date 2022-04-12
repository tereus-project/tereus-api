package main

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/tereus-project/tereus-api/ent/submission"
	"github.com/tereus-project/tereus-api/services"
)

type submissionCompletionMessage struct {
	ID string `json:"id"`
}

func startSubmissionCompletionListener(rabbitMQService *services.RabbitMQService, databaseService *services.DatabaseService) error {
	queue, err := rabbitMQService.NewQueue("submission_completion_q", "submission_completion_ex", "#")
	if err != nil {
		return err
	}

	deliveries, err := queue.Consume()
	if err != nil {
		return err
	}

	go func() {
		for delivery := range deliveries {
			logrus.Info("Received submission completion message")

			var message submissionCompletionMessage
			if err := json.Unmarshal(delivery.Body, &message); err != nil {
				logrus.WithError(err).Error("Failed to unmarshal submission completion message")

				if err := delivery.Nack(false, true); err != nil {
					logrus.WithError(err).Error("Failed to nack submission completion message")
				}

				continue
			}

			id, err := uuid.Parse(message.ID)
			if err != nil {
				logrus.WithError(err).Error("Failed to parse submission ID")

				if err := delivery.Nack(false, false); err != nil {
					logrus.WithError(err).Error("Failed to nack submission completion message")
				}

				continue
			}

			err = databaseService.Submission.
				UpdateOneID(id).
				SetStatus(submission.StatusDone).
				Exec(context.Background())
			if err != nil {
				logrus.WithError(err).Error("Failed to update submission status")

				if err := delivery.Nack(false, true); err != nil {
					logrus.WithError(err).Error("Failed to nack submission completion message")
				}

				continue
			}

			if err := delivery.Ack(false); err != nil {
				logrus.WithError(err).Error("Failed to acknowledge submission completion message")
			}
		}
	}()

	return nil
}
