package infrastructure

import (
	"esp32/src/internal/temperatura/application"
	"esp32/src/internal/temperatura/domain"
	"fmt"
	"log"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type MQTTConsumer struct {
	client              MQTT.Client
	createTemperatureUseCase *application.CreateTemperature
}

func NewMQTTConsumer(client MQTT.Client, createTemperatureUseCase *application.CreateTemperature) *MQTTConsumer {
	return &MQTTConsumer{
		client:              client,
		createTemperatureUseCase: createTemperatureUseCase,
	}
}

func (c *MQTTConsumer) Start() {
	topic := "temperatures/pub" 
	token := c.client.Subscribe(topic, 1, func(client MQTT.Client, msg MQTT.Message) {
		payload := string(msg.Payload())
		log.Printf("Mensaje recibido en el tópico %s: %s", msg.Topic(), payload)
	
		temperature, err := parseTemperatureMessage(payload)
		if err != nil {
			log.Printf("Error parseando el mensaje: %s", err)
			return
		}
	
		if err := c.createTemperatureUseCase.Execute(temperature); err != nil {
			log.Printf("Error guardando la temperatura: %s", err)
		}
	})

	if token.Wait() && token.Error() != nil {
		log.Fatalf("Error suscribiéndose al tópico %s: %s", topic, token.Error())
	}

	log.Printf("Suscrito al tópico %s", topic)
}

func parseTemperatureMessage(payload string) (domain.Temperature, error) {
	var temperature domain.Temperature

	_, err := fmt.Sscanf(payload, "IDHamster: %d, Temperatura: %f", &temperature.IDHamster, &temperature.Temperatura)
	if err != nil {
		return domain.Temperature{}, fmt.Errorf("error parseando el mensaje: %s", err)
	}
	return temperature, nil
}