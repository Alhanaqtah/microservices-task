package rabbitmq

import (
	"fmt"
	"user-managment-service/internal/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Broker struct {
	ch *amqp.Channel
}

func New(cfg config.Broker) (*Broker, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s", cfg.User, cfg.Password, cfg.Host, cfg.Port))
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	_, err = ch.QueueDeclare(
		cfg.QueueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &Broker{ch: ch}, nil
}
