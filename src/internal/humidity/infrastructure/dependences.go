package infrastructure

import (
	"database/sql"
	"esp32/src/core"
	"esp32/src/internal/humidity/application"
	"esp32/src/internal/humidity/infrastructure/controllers"
	amqpConsumer "esp32/src/internal/consumer_amqp"
)

type HumidityDependencies struct {
	DB   *sql.DB
	AMQP *core.AMQPConnection
}

func NewHumidityDependencies(db *sql.DB, amqp *core.AMQPConnection) *HumidityDependencies {
	return &HumidityDependencies{DB: db, AMQP: amqp}
}

func (d *HumidityDependencies) GetRoutes() *HumidityRoutes {
	// Crear el repositorio de humedad
	humidityRepo := NewHumidityRepo(d.DB, nil)

	// Crear los casos de uso de humedad
	createHumidityUseCase := application.NewCreateHumidity(humidityRepo)
	getByHamsterUseCase := application.NewGetByHamster(humidityRepo)

	// Crear el controlador de humedad
	createHumidityController := controllers.NewCreateHumidityController(createHumidityUseCase)

	// Crear el consumidor AMQP y asociar el controlador de humedad
	amqpConsumer := amqpConsumer.NewRabbitMQConsumer(d.AMQP, createHumidityController,nil, nil )
	// Crear el controlador para obtener humedad por hámster
	getByHamsterController := controllers.NewGetByHamsterController(getByHamsterUseCase)

	// Iniciar el consumidor AMQP en una goroutine
	go amqpConsumer.Start()


	// Devolver las rutas de humedad
	return NewHumidityRoutes(createHumidityController, getByHamsterController)
}
