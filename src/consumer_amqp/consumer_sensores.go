package consumeramqp

import (
	"encoding/json"
	"log"
	"os"
	"fmt"
	"esp32/src/core"
	"github.com/joho/godotenv"
	depenencesHumidity		 "esp32/src/internal/humidity/domain"
	controllersHumidity  	 "esp32/src/internal/humidity/infrastructure/controllers"
	dependencesTemperature 	 "esp32/src/internal/temperatura/domain"
	controllersTemperature   "esp32/src/internal/temperatura/infrastructure/controllers"
	dependencesMotion 		 "esp32/src/internal/motion/domain"
	controllersMotion 		 "esp32/src/internal/motion/infrastructure/controllers"
	dependencesFood			 "esp32/src/internal/food/domain"
	controllersFood			 "esp32/src/internal/food/infrastructure/controllers"
	amqp 					 "github.com/rabbitmq/amqp091-go"
)

type RabbitMQConsumer struct {
	conn       *core.AMQPConnection
	CreateHumidity *controllersHumidity.CreateHumidityController
	CreateTemp *controllersTemperature.CreateTemperatureController
	CreateMov  *controllersMotion.CreateMotionController
	CreateFood *controllersFood.CreateStatusFoodController
}

func NewRabbitMQConsumer(conn *core.AMQPConnection, CreateHumidity *controllersHumidity.CreateHumidityController, createTemp *controllersTemperature.CreateTemperatureController, createMov *controllersMotion.CreateMotionController, createFood *controllersFood.CreateStatusFoodController) *RabbitMQConsumer {
	

	return &RabbitMQConsumer{
		conn:       conn,
		CreateHumidity: CreateHumidity,
		CreateTemp: createTemp,
		CreateMov: createMov,
		CreateFood: createFood,
	}
	
}

func (c *RabbitMQConsumer) Start() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error cargando .env: %v", err)
	}
	log.Println("Variables de entorno cargadas correctamente.")

	amqpServer := os.Getenv("AMQP_SERVER")
	log.Printf("Conectando a RabbitMQ en: %s\n", amqpServer)

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

	q, err := ch.QueueDeclare("sensores", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Error declarando cola: %v", err)
	}
	log.Println("Cola 'sensores' declarada.")

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
			Alimento    int   	`json:"alimento,omitempty"`
			Porcentaje  float32   	`json:"porcentaje,omitempty"`
		}

		if err := json.Unmarshal(msg.Body, &sensorData); err != nil {
			log.Printf("Error deserializando el mensaje: %v", err)
			continue
		}

		log.Printf("Mensaje recibido: Sensor: %s, IDHamster: %d\n", sensorData.Sensor, sensorData.IDHamster)

		switch sensorData.Sensor {
		case "temperatura":
			if sensorData.Temperatura == 0 {
				log.Printf("Temperatura no válida para el hámster ID: %d\n", sensorData.IDHamster)
				continue
			}
			log.Printf("Procesando temperatura: %v", sensorData.Temperatura)

			temperature := dependencesTemperature.Temperature{
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
			if sensorData.Humedad == 0 {
				log.Printf("Humedad no válida para el hámster ID: %d\n", sensorData.IDHamster)
				continue
			}
			log.Printf("Procesando humedad: %v", sensorData.Humedad)

			humidity := depenencesHumidity.Humidity{
				IDHamster: int32(sensorData.IDHamster),
				Humedad:   sensorData.Humedad,
			}

			log.Printf("Humedad procesada: %v", humidity)

			if c.CreateHumidity != nil {
				err := c.CreateHumidity.ProcessHumidity(humidity)
				if err != nil {
					log.Printf("Error procesando humedad: %v", err)
				} else {
					log.Println("Humedad procesada exitosamente.")
				}
			} else {
				log.Println("El controlador de humedad es nil, no se puede procesar.")
			}
		case "movimiento":
			if sensorData.IDHamster == 0 {
				log.Println("Advertencia: Datos no válidos o mensaje incorrecto.")
				continue
			}
	
			mov := sensorData.Movimiento == 1  
			fmt.Printf("Mensaje de Movimiento recibido: %+v\n", sensorData)
	
			motion := dependencesMotion.Motion{
				IDHamster:  int32(sensorData.IDHamster),
				Movimiento: mov,  
			}
	
			if err := c.CreateMov.ProcessMotion(motion); err != nil {
				log.Printf("Error al procesar el movimiento: %v", err)
			}
		case "alimento":
			if sensorData.IDHamster == 0 {
				log.Println("Advertencia: Datos no válidos o mensaje incorrecto.")
				continue
			}
	
			fd := sensorData.Alimento == 1 
			fmt.Printf("Mensaje de alimento recibido: %+v\n", sensorData)
	
			food := dependencesFood.Food{
				IDHamster: int32(sensorData.IDHamster),
				Alimento:  fd,
				Porcentaje: sensorData.Porcentaje,
			}
	
			if err := c.CreateFood.ProcessFood(food); err != nil {
				log.Printf("Error al procesar estatus de alimento: %v", err)
			}
		}
}
}
