package infrastructure

import (
	"database/sql"
	"esp32/src/core"
	"esp32/src/internal/temperatura/application"
	"esp32/src/internal/temperatura/infrastructure/controllers"
)

type TemperatureDependencies struct {
	DB   *sql.DB
	AMQP *core.AMQPConnection
}

func NewTemperatureDependencies(db *sql.DB, amqp *core.AMQPConnection) *TemperatureDependencies {
	return &TemperatureDependencies{DB: db, AMQP: amqp}
}

func (d *TemperatureDependencies) GetRoutes() *TemperatureRoutes {
	amqpConsumer := NewAMQPConsumer(d.AMQP, nil)

	temperatureRepo := NewTemperatureRepo(d.DB, amqpConsumer)
	createTemperatureUseCase := application.NewCreateTemperature(temperatureRepo)
	getByHamsterUseCase := application.NewGetByHamster(temperatureRepo)

	createTemperatureController := controllers.NewCreateTemperatureController(createTemperatureUseCase)
	getByHamsterController := controllers.NewGetByHamsterController(getByHamsterUseCase)

	amqpConsumer.createTempC = createTemperatureController
	go amqpConsumer.Consume()

	return NewTemperatureRoutes(createTemperatureController, getByHamsterController)
}
