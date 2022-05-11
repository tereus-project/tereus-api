package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaService struct {
	SubmissionStatusConsumer *kafka.Consumer
	Producer                 *kafka.Producer
	Admin                    *kafka.AdminClient
	Topics                   map[string]kafka.TopicPartition
}

func NewKafkaService(endpoint string) (*KafkaService, error) {
	config := kafka.ConfigMap{
		"bootstrap.servers": endpoint,
		"group.id":          "api",
		"auto.offset.reset": "earliest",
	}

	submissionStatusConsumer, err := kafka.NewConsumer(&config)
	if err != nil {
		return nil, err
	}
	admin, err := kafka.NewAdminClient(&config)
	if err != nil {
		return nil, err
	}

	submissionStatusConsumer.SubscribeTopics([]string{"submission_status"}, nil)

	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": endpoint,
	})
	if err != nil {
		return nil, err
	}

	// Delivery report handler for produced messages
	go func() {
		for e := range producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	return &KafkaService{
		SubmissionStatusConsumer: submissionStatusConsumer,
		Producer:                 producer,
		Admin:                    admin,
		Topics:                   make(map[string]kafka.TopicPartition),
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
	topic, err := k.getTopic(topicName)
	if err != nil {
		return err
	}

	k.Producer.Produce(&kafka.Message{
		TopicPartition: topic,
		Value:          msgBytes,
	}, nil)
	if err != nil {
		return err
	}

	return nil
}

func (k *KafkaService) getTopic(name string) (kafka.TopicPartition, error) {
	topic, ok := k.Topics[name]
	if !ok {
		// Create topic
		_, err := k.Admin.CreateTopics(context.Background(), []kafka.TopicSpecification{
			{
				Topic:         name,
				NumPartitions: 1,
				Config:        map[string]string{},
			},
		})
		if err != nil {
			return kafka.TopicPartition{}, err
		}

		// Save it and return it
		topic = kafka.TopicPartition{Topic: &name, Partition: kafka.PartitionAny}
		k.Topics[name] = topic

		return topic, nil
	}

	return topic, nil
}
