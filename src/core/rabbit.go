package core

import (
    "log"
    "os"
    "time"

    MQTT "github.com/eclipse/paho.mqtt.golang"
)

type MQTTConnection struct {
    Client MQTT.Client
}

func NewMQTTConnection() (*MQTTConnection, error) {
    mqttBrokerURL := os.Getenv("MQTT_BROKER_URL")
    if mqttBrokerURL == "" {
        log.Fatal("La variable de entorno MQTT_BROKER_URL no está configurada")
    }

    opts := MQTT.NewClientOptions().AddBroker(mqttBrokerURL)
    opts.SetClientID("go_mqtt_client")
    opts.SetAutoReconnect(true)
    opts.SetConnectTimeout(30 * time.Second)

    client := MQTT.NewClient(opts)
    if token := client.Connect(); token.Wait() && token.Error() != nil {
        return nil, token.Error()
    }

    return &MQTTConnection{Client: client}, nil
}

func (m *MQTTConnection) Close() {
    if m.Client != nil {
        m.Client.Disconnect(250) 
    }
}

func (m *MQTTConnection) Publish(topic string, qos byte, retained bool, payload interface{}) error {
    token := m.Client.Publish(topic, qos, retained, payload)
    token.Wait()
    return token.Error()
}

func (m *MQTTConnection) Subscribe(topic string, qos byte, callback MQTT.MessageHandler) error {
    if token := m.Client.Subscribe(topic, qos, callback); token.Wait() && token.Error() != nil {
        return token.Error()
    }
    return nil
}