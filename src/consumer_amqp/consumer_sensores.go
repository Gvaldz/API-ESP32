package consumeramqp

import (
	"encoding/json"
	"log"
	"os"
	"fmt"
	"esp32/src/core"
	"github.com/joho/godotenv"
	depenencesHumidity		 "esp32/src/internal/sensores/humidity/domain"
	controllersHumidity  	 "esp32/src/internal/sensores/humidity/infrastructure/controllers"
	dependencesTemperature 	 "esp32/src/internal/sensores/temperatura/domain"
	controllersTemperature   "esp32/src/internal/sensores/temperatura/infrastructure/controllers"
	dependencesMotion 		 "esp32/src/internal/sensores/motion/domain"
	controllersMotion 		 "esp32/src/internal/sensores/motion/infrastructure/controllers"
	dependencesFood			 "esp32/src/internal/sensores/food/domain"
	controllersFood			 "esp32/src/internal/sensores/food/infrastructure/controllers"
	dependencesCage			 "esp32/src/internal/sensores/cages/domain"
	controllersCage			 "esp32/src/internal/sensores/cages/infrastructure/controllers"
	amqp 					 "github.com/rabbitmq/amqp091-go"
)

type RabbitMQConsumer struct {
	conn       *core.AMQPConnection
	CreateHumidity *controllersHumidity.CreateHumidityController
	CreateTemp *controllersTemperature.CreateTemperatureController
	CreateMov  *controllersMotion.CreateMotionController
	CreateFood *controllersFood.CreateStatusFoodController
	CreateCage *controllersCage.CreateCageController
}

func NewRabbitMQConsumer(conn *core.AMQPConnection, CreateHumidity *controllersHumidity.CreateHumidityController, createTemp *controllersTemperature.CreateTemperatureController, createMov *controllersMotion.CreateMotionController, createFood *controllersFood.CreateStatusFoodController, createCage *controllersCage.CreateCageController) *RabbitMQConsumer {
	

	return &RabbitMQConsumer{
		conn:       conn,
		CreateHumidity: CreateHumidity,
		CreateTemp: createTemp,
		CreateMov: createMov,
		CreateFood: createFood,
		CreateCage: createCage,
	}
	
}

func (c *RabbitMQConsumer) Start() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error cargando .env: %v", err)
	}

	amqpServer := os.Getenv("AMQP_SERVER")

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
	defer ch.Close()

	q, err := ch.QueueDeclare("sensores", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Error declarando cola: %v", err)
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Error al consumir mensajes: %v", err)
	}
	log.Println("Esperando mensajes...")

	for msg := range msgs {
		var sensorData struct {
			Sensor      string  `json:"sensor"`
			IDHamster   string     `json:"idhamster"`
			Humedad     float64 `json:"humedad,omitempty"`
			Temperatura float64 `json:"temperatura,omitempty"`
			Movimiento  int   	`json:"movimiento,omitempty"`
			Alimento    int   	`json:"alimento,omitempty"`
			Porcentaje  float32 `json:"porcentaje,omitempty"`
			Token       string  `json:"token"` 
		}

		if err := json.Unmarshal(msg.Body, &sensorData); err != nil {
			log.Printf("Error deserializando el mensaje: %v", err)
			continue
		}

		log.Printf("Mensaje recibido: Sensor: %s, IDamster: %s\n", sensorData.Sensor, sensorData.IDHamster)


		if sensorData.Sensor == "idhamster"{

			if sensorData.IDHamster == "" {
				log.Println("Advertencia: ID no valido.")
				continue
			}
	
			jaula := dependencesCage.Cage{
				Idjaula: string(sensorData.IDHamster),
			}
	
			if err := c.CreateCage.ProcessCage(jaula); err != nil {
				log.Printf("Error al procesar jaula: %v", err)}
		}


		switch sensorData.Sensor {
		case "temperatura":
			if sensorData.Temperatura == 0 {
				log.Printf("Temperatura no válida para el hámster ID: %s\n", sensorData.IDHamster)
				continue
			}
			log.Printf("Procesando temperatura: %v", sensorData.Temperatura)
		
			temperature := dependencesTemperature.Temperature{
				IDHamster:   string(sensorData.IDHamster),
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
				log.Printf("Humedad no válida para el hámster ID: %s\n", sensorData.IDHamster)
				continue
			}
			log.Printf("Procesando humedad: %v", sensorData.Humedad)

			humidity := depenencesHumidity.Humidity{
				IDHamster: string(sensorData.IDHamster),
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
			if sensorData.IDHamster == "" {
				log.Println("Advertencia: Datos no válidos o mensaje incorrecto.")
				continue
			}
	
			mov := sensorData.Movimiento == 1  
			fmt.Printf("Mensaje de Movimiento recibido: %+v\n", sensorData)
	
			motion := dependencesMotion.Motion{
				IDHamster:  string(sensorData.IDHamster),
				Movimiento: mov,  
			}
	
			if err := c.CreateMov.ProcessMotion(motion); err != nil {
				log.Printf("Error al procesar el movimiento: %v", err)
			}
		case "alimento":
			if sensorData.Porcentaje == 0 {
				log.Printf("Porcentaje de alimento no válido para el hámster ID: %s\n", sensorData.IDHamster)
				continue
			}
		
			food := dependencesFood.Food{
				IDHamster:  string(sensorData.IDHamster),
				Alimento:   sensorData.Alimento,
				Porcentaje: sensorData.Porcentaje,
			}
		
			err := c.CreateFood.ProcessFood(food)
			if err != nil {
				log.Printf("Error procesando estado de alimento: %v", err)
			}
		
		}
}
}
