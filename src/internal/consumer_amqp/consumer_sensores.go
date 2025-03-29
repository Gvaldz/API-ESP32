package consumeramqp

import (
	"encoding/json"
	"log"
	"os"
	"fmt"
	"esp32/src/core"
	dHum "esp32/src/internal/humidity/domain"
	cHum  "esp32/src/internal/humidity/infrastructure/controllers"
	"github.com/joho/godotenv"
	dTemp "esp32/src/internal/temperatura/domain"
	cTemp "esp32/src/internal/temperatura/infrastructure/controllers"
	dMov "esp32/src/internal/motion/domain"
	cMov "esp32/src/internal/motion/infrastructure/controllers"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQConsumer struct {
	conn       *core.AMQPConnection
	CreateHumC *cHum.CreateHumidityController
	CreateTemp *cTemp.CreateTemperatureController
	CreateMov  *cMov.CreateMotionController
}

func NewRabbitMQConsumer(conn *core.AMQPConnection, createHumC *cHum.CreateHumidityController, createTemp *cTemp.CreateTemperatureController, createMov *cMov.CreateMotionController) *RabbitMQConsumer {
	

	return &RabbitMQConsumer{
		conn:       conn,
		CreateHumC: createHumC,
		CreateTemp: createTemp,
		CreateMov: createMov,
	}
	
}

func (c *RabbitMQConsumer) Start() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error cargando .env: %v", err)
	}
	log.Println("Variables de entorno cargadas correctamente.")

	amqpServer := os.Getenv("AMQP_SERVER")
	log.Printf("Conectando a RabbitMQ en: %s\n", amqpServer)

	// Conectar a RabbitMQ
	connRabbit, err := amqp.Dial(amqpServer)
	if err != nil {
		log.Fatalf("Error conectando a RabbitMQ: %v", err)
	}
	log.Println("Conexión a RabbitMQ establecida.")
	defer connRabbit.Close()

	ch, err := connRabbit.Channel()
	if err != nil {
		log.Fatalf("Error abriendo canal en RabbitMQ: %v", err)
	}
	log.Println("Canal RabbitMQ abierto.")
	defer ch.Close()

	// Consumir los mensajes de la cola "sensores" existente
	q, err := ch.QueueDeclare("sensores", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Error declarando cola: %v", err)
	}
	log.Println("Cola 'sensores' declarada.")

	// Consumir los mensajes de la cola
	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Error al consumir mensajes: %v", err)
	}
	log.Println("Esperando mensajes...")

	for msg := range msgs {
		var sensorData struct {
			Sensor      string  `json:"sensor"`
			IDHamster   int     `json:"idhamster"`
			Humedad     float64 `json:"humedad,omitempty"`
			Temperatura float64 `json:"temperatura,omitempty"`
			Movimiento  int   	`json:"movimiento,omitempty"`
		}

		// Deserializar el mensaje recibido
		if err := json.Unmarshal(msg.Body, &sensorData); err != nil {
			log.Printf("Error deserializando el mensaje: %v", err)
			continue
		}

		log.Printf("Mensaje recibido: Sensor: %s, IDHamster: %d\n", sensorData.Sensor, sensorData.IDHamster)

		// Lógica de negocio según el tipo de sensor
		switch sensorData.Sensor {
		case "temperatura":
			// Validar y procesar datos de temperatura
			if sensorData.Temperatura == 0 {
				log.Printf("Temperatura no válida para el hámster ID: %d\n", sensorData.IDHamster)
				continue
			}
			log.Printf("Procesando temperatura: %v", sensorData.Temperatura)

			temperature := dTemp.Temperature{
				IDHamster:   int32(sensorData.IDHamster),
				Temperatura: sensorData.Temperatura,
			}

			log.Printf("Temperatura procesada: %v", temperature)

			if c.CreateTemp != nil {
				err := c.CreateTemp.ProcessTemperature(temperature)
				if err != nil {
					log.Printf("Error procesando temperatura: %v", err)
				} else {
					log.Println("Temperatura procesada exitosamente.")
				}
			} else {
				log.Println("El controlador de temperatura es nil, no se puede procesar.")
			}

		case "humedad":
			// Validar y procesar datos de humedad
			if sensorData.Humedad == 0 {
				log.Printf("Humedad no válida para el hámster ID: %d\n", sensorData.IDHamster)
				continue
			}
			log.Printf("Procesando humedad: %v", sensorData.Humedad)

			humidity := dHum.Humidity{
				IDHamster: int32(sensorData.IDHamster),
				Humedad:   sensorData.Humedad,
			}

			log.Printf("Humedad procesada: %v", humidity)

			if c.CreateHumC != nil {
				err := c.CreateHumC.ProcessHumidity(humidity)
				if err != nil {
					log.Printf("Error procesando humedad: %v", err)
				} else {
					log.Println("Humedad procesada exitosamente.")
				}
			} else {
				log.Println("El controlador de humedad es nil, no se puede procesar.")
			}
		case "movimiento":
			// Validar y procesar datos de humedad
			if sensorData.IDHamster == 0 {
				log.Println("Advertencia: Datos no válidos o mensaje incorrecto.")
				continue
			}
	
			// Convertir el valor de "movimiento" a bool
			mov := sensorData.Movimiento == 1  // 1 es true, 0 es false
	
			fmt.Printf("📩 Mensaje de Movimiento recibido: %+v\n", sensorData)
	
			motion := dMov.Motion{
				IDHamster:  int32(sensorData.IDHamster),
				Movimiento: mov,  // Asignar el valor booleano
			}
	
			if err := c.CreateMov.ProcessMotion(motion); err != nil {
				log.Printf("Error al procesar el movimiento: %v", err)
			}
		}
}
}
