package rabbitmq_infrastructure

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/submission"
	entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/submission"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQAdapter struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	queues     map[string]string
}

func NewRabbitMQAdapter() (*RabbitMQAdapter, error) {
	queues := map[string]string{
		string(entities.LanguagePython): getEnvOrDefault("PYTHON_QUEUE", "python.queue"),
		string(entities.LanguageJava):   getEnvOrDefault("JAVA_QUEUE", "java.queue"),
		string(entities.LanguageCPP):    getEnvOrDefault("CPP_QUEUE", "cpp.queue"),
	}

	amqpURL := getEnvOrDefault("RABBITMQ_URL", "amqp://guest:guest@rabbitmq:5672/")

	conn, err := connectWithRetry(amqpURL, 10, 2*time.Second)
	if err != nil {
		return nil, fmt.Errorf("connect to rabbitmq: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("open rabbitmq channel: %w", err)
	}

	for language, queueName := range queues {
		_, err = ch.QueueDeclare(
			queueName,
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			_ = ch.Close()
			_ = conn.Close()
			return nil, fmt.Errorf("declare rabbitmq queue for language %q (%q): %w", language, queueName, err)
		}
	}

	return &RabbitMQAdapter{
		connection: conn,
		channel:    ch,
		queues:     queues,
	}, nil
}

func connectWithRetry(amqpURL string, attempts int, delay time.Duration) (*amqp.Connection, error) {
	var lastErr error

	for i := 0; i < attempts; i++ {
		conn, err := amqp.Dial(amqpURL)
		if err == nil {
			return conn, nil
		}

		lastErr = err
		if i < attempts-1 {
			time.Sleep(delay)
		}
	}

	return nil, fmt.Errorf("after %d attempts: %w", attempts, lastErr)
}

func (a *RabbitMQAdapter) PublishSubmission(dto dtos.SubmissionResultPublishedDTO) error {
	queueName, ok := a.queues[dto.Language]
	if !ok || queueName == "" {
		return fmt.Errorf("no queue configured for language %q", dto.Language)
	}

	payload, err := json.Marshal(dto)
	if err != nil {
		return fmt.Errorf("marshal submission payload: %w", err)
	}

	err = a.channel.Publish(
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         payload,
		},
	)
	if err != nil {
		return fmt.Errorf("publish submission to queue %q: %w", queueName, err)
	}

	return nil
}

func getEnvOrDefault(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
