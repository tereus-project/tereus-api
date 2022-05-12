package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"github.com/tereus-project/tereus-api/ent/submission"
)

type KafkaService struct {
	endpoint        string
	consumerGroupID string
	writers         map[string]*kafka.Writer
}

func NewKafkaService(endpoint string) (*KafkaService, error) {
	return &KafkaService{
		endpoint:        endpoint,
		consumerGroupID: "api",
		writers:         make(map[string]*kafka.Writer),
	}, nil
}

type RemixJob struct {
	ID             string `json:"id"`
	SourceLanguage string `json:"source_language"`
	TargetLanguage string `json:"target_language"`
}

func (k *KafkaService) PublishSubmission(remixJob RemixJob) error {
	msgBytes, err := json.Marshal(&remixJob)
	if err != nil {
		return err
	}

	topicName := fmt.Sprintf("remix_jobs_%s_to_%s", remixJob.SourceLanguage, remixJob.TargetLanguage)
	writer, err := k.getWriterForTopic(topicName)
	if err != nil {
		return err
	}

	err = writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(remixJob.ID),
			Value: msgBytes,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (k *KafkaService) getWriterForTopic(topicName string) (*kafka.Writer, error) {
	w, ok := k.writers[topicName]
	if !ok {
		w := &kafka.Writer{
			Addr:     kafka.TCP(k.endpoint),
			Topic:    topicName,
			Balancer: &kafka.LeastBytes{},
		}
		k.writers[topicName] = w

		return w, nil
	}

	return w, nil
}

type SubmissionStatusMessage struct {
	ID     string            `json:"id"`
	Status submission.Status `json:"status"`
	Reason string            `json:"reason"`
}

func (k *KafkaService) ConsumeSubmissionStatus(ctx context.Context) <-chan SubmissionStatusMessage {

	ch := make(chan SubmissionStatusMessage)

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{k.endpoint},
		GroupID:     k.consumerGroupID,
		Topic:       "submission_status",
		MaxWait:     1 * time.Second,
		StartOffset: kafka.LastOffset,
	})

	go func() {
		defer r.Close()

		for {
			msg, err := r.ReadMessage(ctx)
			if err != nil {
				logrus.WithError(err).Error("Error reading message")
			}

			var status SubmissionStatusMessage
			err = json.Unmarshal(msg.Value, &status)
			if err != nil {
				logrus.WithError(err).Error("Failed to unmarshal submission status message")
				continue
			}

			ch <- status
		}
	}()

	return ch
}

func (k *KafkaService) CloseAllWriters() {
	for _, w := range k.writers {
		err := w.Close()
		if err != nil {
			logrus.WithError(err).Error("Failed to close kafka writer")
		}
	}
}
