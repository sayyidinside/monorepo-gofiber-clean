package rabbitmq

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/infrastructure/config"
)

type RabbitMQClient struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

func Connect(cfg *config.Config) (*RabbitMQClient, error) {
	conn, err := amqp.Dial(cfg.RabbitMQURL)
	if err != nil {
		conn.Close()
		log.Fatalf("failed to connect database: %v", err)
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		ch.Close()
		log.Fatalf("failed to connect database: %v", err)
		return nil, err
	}

	return &RabbitMQClient{
		Connection: conn,
		Channel:    ch,
	}, err
}
