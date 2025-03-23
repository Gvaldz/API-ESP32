package infrastructure

import (
	"encoding/json"
	"fmt"
	"log"

	"esp32/src/core"
	"esp32/src/internal/temperatura/domain"
	"esp32/src/internal/temperatura/infrastructure/controllers"
)

// AMQPConsumer estructura para consumir mensajes de la cola.
type AMQPConsumer struct {
	conn        *core.AMQPConnection
	createTempC *controllers.CreateTemperatureController
}

// NewAMQPConsumer constructor
func NewAMQPConsumer(conn *core.AMQPConnection, createTempC *controllers.CreateTemperatureController) *AMQPConsumer {
	return &AMQPConsumer{
		conn:        conn,
		createTempC: createTempC,
	}
}

// Método para consumir mensajes
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
		fmt.Println("📩 Recibido desde RabbitMQ:", string(msg.Body))

		var sensorData struct {
			IDHamster   int     `json:"idhamster"`
			Temperatura float64 `json:"temperatura"`
			

		}

		if err := json.Unmarshal(msg.Body, &sensorData); err != nil {
			log.Printf("❌ Error al deserializar el mensaje: %v", err)
			continue
		}

		if sensorData.IDHamster == 0 || sensorData.Temperatura == 0 {
			log.Println("⚠️ Advertencia: Datos no válidos o mensaje incorrecto.")
			continue
		}

		temperature := domain.Temperature{
			IDHamster:   int32(sensorData.IDHamster),
			Temperatura: sensorData.Temperatura,
		}

		if err := c.createTempC.ProcessTemperature(temperature); err != nil {
			log.Printf("❌ Error al procesar la temperatura: %v", err)
		}
	}
}
