package infrastructure

import (
	"database/sql"
	"esp32/src/core"
	"esp32/src/internal/temperatura/application"
	"esp32/src/internal/temperatura/infrastructure/controllers"
	amqpConsumer "esp32/src/internal/consumer_amqp"
)

type TemperatureDependencies struct {
	DB   *sql.DB
	AMQP *core.AMQPConnection
}

func NewTemperatureDependencies(db *sql.DB, amqp *core.AMQPConnection) *TemperatureDependencies {
	return &TemperatureDependencies{DB: db, AMQP: amqp}
}

func (d *TemperatureDependencies) GetRoutes() *TemperatureRoutes {
	// Crear el repositorio de temperatura
	temperatureRepo := NewTemperatureRepo(d.DB, nil)

	// Crear los casos de uso
	createTemperatureUseCase := application.NewCreateTemperature(temperatureRepo)
	getByHamsterUseCase := application.NewGetByHamster(temperatureRepo)

	// Crear el controlador de temperatura
	createTemperatureController := controllers.NewCreateTemperatureController(createTemperatureUseCase)



	// Crear el consumidor AMQP y asociar los controladores
	amqpConsumer := amqpConsumer.NewRabbitMQConsumer(d.AMQP, nil, createTemperatureController, nil)

	// Asociar el consumidor AMQP al controlador de temperatura y humedad
	amqpConsumer.CreateTemp = createTemperatureController
	

	// Crear el controlador para obtener la temperatura por hámster
	getByHamsterController := controllers.NewGetByHamsterController(getByHamsterUseCase)

	// Iniciar el consumidor AMQP en una goroutine
	go amqpConsumer.Start()

	// Devolver las rutas de temperatura y humedad
	return NewTemperatureRoutes(createTemperatureController, getByHamsterController)
}
