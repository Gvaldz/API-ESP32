package core

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

// AMQPConnection estructura para gestionar la conexión con RabbitMQ.
type AMQPConnection struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

// NewAMQPConnection constructor para crear una nueva conexión y canal AMQP.
func NewAMQPConnection() (*AMQPConnection, error) {
	conn, err := amqp.Dial("amqp://ale:ale05@54.156.170.232:5672/")
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
