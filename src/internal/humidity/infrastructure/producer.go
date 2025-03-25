package infrastructure

import (
	"esp32/src/internal/temperatura/domain"
	"fmt"
	"log"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type MQTTProducer struct {
	client MQTT.Client
}

func NewMQTTProducer(client MQTT.Client) *MQTTProducer {
	return &MQTTProducer{client: client}
}

func (p *MQTTProducer) PublishMessage(topic string, message string) error {
    log.Printf("Intentando publicar mensaje en el tópico: %s", topic)
	token := p.client.Publish(topic, 1, true, message)
	token.Wait()
	if token.Error() != nil {
		log.Printf("Error publicando mensaje en %s: %s", topic, token.Error())
		return token.Error()
	}
	log.Println("Mensaje publicado correctamente")
	
	return nil
}

func (p *MQTTProducer) SendTemperatureMessage(temperature domain.Temperature) error {
	topic := "temperatures/processed"
	message := fmt.Sprintf("IDHamster: %d, Temperatura: %.2f", temperature.IDHamster, temperature.Temperatura)
    log.Printf("Preparando mensaje para enviar: %s", message) 
	return p.PublishMessage(topic, message)
}