package infrastructure

import (
	"encoding/json"
	"log"
	"fmt"
	"esp32/src/core"
	"github.com/joho/godotenv"
	"os"
	 amqp "github.com/rabbitmq/amqp091-go"
	"esp32/src/internal/motion/domain"
	"esp32/src/internal/motion/infrastructure/controllers"
)

type AMQPConsumer struct {
	conn           *core.AMQPConnection
	createMotionC  *controllers.CreateMotionController
}

func NewAMQPConsumer(conn *core.AMQPConnection, createMotionC *controllers.CreateMotionController) *AMQPConsumer {
	return &AMQPConsumer{
		conn:          conn,
		createMotionC: createMotionC,
	}
}
func (c *AMQPConsumer) Consume() {
    // Cargar variables de entorno
    if err := godotenv.Load(); err != nil {
        log.Fatalf("Error cargando .env: %v", err)
    }

    amqpServer := os.Getenv("AMQP_SERVER")
    exchangeName := "sensor_data"
    sensorType := "movimiento" // Solo nos interesa consumir los mensajes de temperatura

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
		log.Println("Recibido desde RabbitMQ:", string(msg.Body))

		var sensorData struct {
			Sensor      string  `json:"sensor"`
			IDHamster   int   `json:"idhamster"`
            Movimiento  int   `json:"movimiento,omitempty"`  // Recibimos un int desde RabbitMQ
		}

		if err := json.Unmarshal(msg.Body, &sensorData); err != nil {
			log.Printf("Error al deserializar el mensaje: %v", err)
			continue
		}

		if sensorData.IDHamster == 0 {
			log.Println("Advertencia: Datos no válidos o mensaje incorrecto.")
			continue
		}

		// Convertir el valor de "movimiento" a bool
		mov := sensorData.Movimiento == 1  // 1 es true, 0 es false

		fmt.Printf("📩 Mensaje de Movimiento recibido: %+v\n", sensorData)

		motion := domain.Motion{
			IDHamster:  int32(sensorData.IDHamster),
			Movimiento: mov,  // Asignar el valor booleano
		}

		if err := c.createMotionC.ProcessMotion(motion); err != nil {
			log.Printf("Error al procesar el movimiento: %v", err)
		}
	}
}
