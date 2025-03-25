package infrastructure

import (
	"encoding/json"
	"fmt"
	"log"

	"esp32/src/core"
	"esp32/src/internal/humidity/domain"
	"esp32/src/internal/humidity/infrastructure/controllers"
)

type AMQPConsumer struct {
	conn        *core.AMQPConnection
	createTempC *controllers.CreateHumidityController
}

func NewAMQPConsumer(conn *core.AMQPConnection, createTempC *controllers.CreateHumidityController) *AMQPConsumer {
	return &AMQPConsumer{
		conn:        conn,
		createTempC: createTempC,
	}
}

func (c *AMQPConsumer) Consume() {
	msgs, err := c.conn.Channel.Consume(
		"sensor_humedad",
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
			Humedad float64 `json:"humedad"`

		}

		if err := json.Unmarshal(msg.Body, &sensorData); err != nil {
			log.Printf("Error al deserializar el mensaje: %v", err)
			continue
		}

		if sensorData.IDHamster == 0 || sensorData.Humedad == 0 {
			log.Println(" Advertencia: Datos no válidos o mensaje incorrecto.")
			continue
		}

		humidity := domain.Humidity{
			IDHamster:   int32(sensorData.IDHamster),
			Humedad: sensorData.Humedad,
		}

		if err := c.createTempC.ProcessHumidity(humidity); err != nil {
			log.Printf("Error al procesar la humedad: %v", err)
		}
	}
}
