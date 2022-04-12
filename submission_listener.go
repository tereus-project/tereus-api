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
				delivery.Nack(false, true)
				continue
			}

			id, err := uuid.Parse(message.ID)
			if err != nil {
				logrus.WithError(err).Error("Failed to parse submission ID")
				delivery.Nack(false, false)
				continue
			}

			err = databaseService.Submission.
				UpdateOneID(id).
				SetStatus(submission.StatusDone).
				Exec(context.Background())
			if err != nil {
				logrus.WithError(err).Error("Failed to update submission status")
				delivery.Nack(false, true)
				continue
			}

			delivery.Ack(false)
		}
	}()

	return nil
}
