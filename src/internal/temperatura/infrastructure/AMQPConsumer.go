package infrastructure

import (
    "encoding/json"
    "fmt"
    "log"
    "os"

    "esp32/src/core"
    "esp32/src/internal/temperatura/domain"
    "esp32/src/internal/temperatura/infrastructure/controllers"
    "github.com/joho/godotenv"
    amqp "github.com/rabbitmq/amqp091-go"
)

type AMQPConsumer struct {
    conn       *core.AMQPConnection
    createTemp *controllers.CreateTemperatureController
}

func NewAMQPConsumer(conn *core.AMQPConnection, createTemp *controllers.CreateTemperatureController) *AMQPConsumer {
    return &AMQPConsumer{
        conn:       conn,
        createTemp: createTemp,
    }
}

func (c *AMQPConsumer) Consume() {
    // Cargar variables de entorno
    if err := godotenv.Load(); err != nil {
        log.Fatalf("Error cargando .env: %v", err)
    }

    amqpServer := os.Getenv("AMQP_SERVER")
    exchangeName := "sensor_data"
    sensorType := "temperatura" // Solo nos interesa consumir los mensajes de temperatura

    // Conectar a RabbitMQ
    connRabbit, err := amqp.Dial(amqpServer)
    if err != nil {
        log.Fatalf("Error conectando a RabbitMQ: %v", err)
    }
    defer connRabbit.Close()

    chRabbit, err := connRabbit.Channel()
    if err != nil {
        log.Fatalf("Error abriendo canal en RabbitMQ: %v", err)
    }
    defer chRabbit.Close()

    // Declarar Exchange de tipo "direct"
    err = chRabbit.ExchangeDeclare(exchangeName, "direct", true, false, false, false, nil)
    if err != nil {
        log.Fatalf("Error declarando exchange: %v", err)
    }

    // Declarar la cola y enlazarla al exchange
    q, err := chRabbit.QueueDeclare(sensorType, true, false, false, false, nil)
    if err != nil {
        log.Fatalf("Error declarando cola: %v", err)
    }

    err = chRabbit.QueueBind(q.Name, sensorType, exchangeName, false, nil)
    if err != nil {
        log.Fatalf("Error enlazando la cola al exchange: %v", err)
    }

    // Consumir mensajes
    msgs, err := chRabbit.Consume(
        q.Name,
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

    fmt.Printf("Esperando mensajes del sensor: %s\n", sensorType)
    for msg := range msgs {
        var sensorData struct {
            Sensor      string  `json:"sensor"`
            IDHamster   int     `json:"idhamster"`
            Temperatura float64 `json:"temperatura,omitempty"`
        }

        if err := json.Unmarshal(msg.Body, &sensorData); err != nil {
            log.Printf("Error al deserializar el mensaje: %v", err)
            continue
        }

        // Solo procesar los mensajes de tipo "temperatura"
        if sensorData.Sensor != "temperatura" {
            continue
        }

        if sensorData.IDHamster == 0 || sensorData.Temperatura == 0 {
            log.Println("⚠️ Advertencia: Datos no válidos o mensaje incorrecto.")
            continue
        }

        fmt.Printf("📩 Mensaje de temperatura recibido: %+v\n", sensorData)

        // Procesar la temperatura en la capa de aplicación
        temperature := domain.Temperature{
            IDHamster: int32(sensorData.IDHamster),
            Temperatura:   sensorData.Temperatura, // Usar temperatura como "Humedad" por ejemplo
        }

        // Procesar humedad (o temperatura) en la capa de aplicación
        if err := c.createTemp.ProcessTemperature(temperature); err != nil {
            log.Printf("Error procesando temperatura: %v", err)
        }
    }
}
