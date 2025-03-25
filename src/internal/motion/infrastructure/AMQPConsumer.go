package infrastructure

import (
	"encoding/json"
	"fmt"
	"log"

	"esp32/src/core"
	"esp32/src/internal/motion/domain"
	"esp32/src/internal/motion/infrastructure/controllers"
)

type AMQPConsumer struct {
	conn        *core.AMQPConnection
	createMotionC *controllers.CreateMotionController
}

func NewAMQPConsumer(conn *core.AMQPConnection, createMotionC *controllers.CreateMotionController) *AMQPConsumer {
	return &AMQPConsumer{
		conn:        conn,
		createMotionC: createMotionC,
	}
}

func (c *AMQPConsumer) Consume() {
	msgs, err := c.conn.Channel.Consume(
		"sensor_data",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Error consumiendo cola: %v", err)
	}

	for msg := range msgs {
		fmt.Println("Recibido desde RabbitMQ:", string(msg.Body))

		var sensorData struct {
			IDHamster   int     `json:"idhamster"`
			Movimiento bool `json:"motion"`

		}

		if err := json.Unmarshal(msg.Body, &sensorData); err != nil {
			log.Printf("Error al deserializar el mensaje: %v", err)
			continue
		}

		if sensorData.IDHamster == 0 {
			log.Println("Advertencia: Datos no válidos o mensaje incorrecto.")
			continue
		}

		Motionerature := domain.Motion{
			IDHamster:   int32(sensorData.IDHamster),
			Movimiento: sensorData.Movimiento,
		}

		if err := c.createMotionC.ProcessMotion(Motionerature); err != nil {
			log.Printf("Error al procesar el movimiento: %v", err)
		}
	}
}
