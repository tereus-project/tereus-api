package services

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/tereus-project/tereus-api/services/internal"
)

type QueueService struct {
	kafkaService *internal.KafkaService
}

func NewQueueService(kafkaEndpoint string) (*QueueService, error) {
	logrus.Debugln("Initializing Kafka service")
	kafkaService, err := internal.NewKafkaService(kafkaEndpoint, "api")
	if err != nil {
		return nil, fmt.Errorf("Failed initializing kafka service: %s", err)
	}
	defer kafkaService.CloseAllWriters()

	return &QueueService{
		kafkaService: kafkaService,
	}, nil
}

func (s *QueueService) publish(topic string, key string, message interface{}) error {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return s.kafkaService.Publish(topic, key, messageBytes)
}

func (s *QueueService) consume(topic string) <-chan []byte {
	return s.kafkaService.Consume(topic)
}

type Queue[T any] struct {
	queueService *QueueService
	topic        string
}

func NewQueue[T any](topic string, queueService *QueueService) *Queue[T] {
	return &Queue[T]{
		queueService: queueService,
		topic:        topic,
	}
}

func (q *Queue[T]) Publish(key string, message *T) error {
	return q.queueService.publish(q.topic, key, message)
}

func (q *Queue[T]) Consume() <-chan T {
	ch := make(chan T)

	go func() {
		defer close(ch)

		for bytes := range q.queueService.consume(q.topic) {
			var message T
			err := json.Unmarshal(bytes, &message)
			if err != nil {
				logrus.WithError(err).Error("Error unmarshaling message")
				continue
			}

			ch <- message
		}
	}()

	return ch
}
