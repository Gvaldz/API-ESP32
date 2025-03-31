package core

import (
	"fmt"
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

// AMQPConnection estructura para gestionar la conexión con RabbitMQ.
type AMQPConnection struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

func NewAMQPConnection() (*AMQPConnection, error) {
	rabbitMQURL := os.Getenv("AMQP_SERVER")
    if rabbitMQURL == "" {
        log.Fatal("La variable de entorno RABBITMQ_URL no está configurada")
    }

	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		return nil, fmt.Errorf("error conectando a RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("error abriendo canal en RabbitMQ: %v", err)
	}

	return &AMQPConnection{
		Connection: conn,
		Channel:    ch,
	}, nil
}

// Close cierra la conexión y el canal de RabbitMQ.
func (c *AMQPConnection) Close() {
	if c.Channel != nil {
		c.Channel.Close()
	}
	if c.Connection != nil {
		c.Connection.Close()
	}
}
