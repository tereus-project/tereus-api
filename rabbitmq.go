package main

import (
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQService struct {
	ch *amqp.Channel
}

func NewRabbitMQService(endpoint string) (RabbitMQService, error) {
	conn, err := amqp.Dial(endpoint)
	if err != nil {
		return RabbitMQService{}, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return RabbitMQService{}, err
	}

	// Ensure exchange exists
	err = ch.ExchangeDeclare(
		"remix_jobs_ex", // name
		"direct",        // type
		true,            // durable
		false,           // auto-deleted
		false,           // internal
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		return RabbitMQService{}, err
	}

	// Ensure queue exists
	jobsQ, err := ch.QueueDeclare(
		"remix_jobs_q", // name
		true,           // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		return RabbitMQService{}, err
	}

	// Bind queue to exchange
	err = ch.QueueBind(
		jobsQ.Name,      // queue name
		"remix_jobs_rk", // routing key
		"remix_jobs_ex", // exchange
		false,
		nil)
	if err != nil {
		return RabbitMQService{}, err
	}

	return RabbitMQService{ch: ch}, nil
}

type remixJob struct {
	ID             string `json:"id"`
	SourceLanguage string `json:"source_language"`
	TargetLanguage string `json:"target_language"`
}

// Publish a job to the exchange
func (s RabbitMQService) publishJob(job remixJob) error {
	b, err := json.Marshal(job)
	if err != nil {
		return err
	}
	err = s.ch.Publish(
		"remix_jobs_ex", // exchange
		"remix_jobs_rk", // routing key
		false,           // mandatory
		false,           // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        b,
		})
	return err
}
