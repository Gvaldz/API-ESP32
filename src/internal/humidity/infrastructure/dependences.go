package infrastructure

import (
	"database/sql"
	"esp32/src/core"
	"esp32/src/internal/humidity/application"
	"esp32/src/internal/humidity/infrastructure/controllers"
)

type HumidityDependeces struct {
	DB   *sql.DB
	AMQP *core.AMQPConnection
}

func NewHumidityDependeces(db *sql.DB, amqp *core.AMQPConnection) *HumidityDependeces {
	return &HumidityDependeces{DB: db, AMQP: amqp}
}

func (d *HumidityDependeces) GetRoutes() *HumidityRoutes {
	amqpConsumer := NewAMQPConsumer(d.AMQP, nil)

	humidityRepo := NewHumidityRepo(d.DB, amqpConsumer)
	createHumidityUseCase := application.NewCreateHumidity(humidityRepo)
	getByHamsterUseCase := application.NewGetByHamster(humidityRepo)

	createHumidityController := controllers.NewCreateHumidityController(createHumidityUseCase)
	getByHamsterController := controllers.NewGetByHamsterController(getByHamsterUseCase)

	amqpConsumer.createHumC = createHumidityController
	go amqpConsumer.Consume()

	return NewHumidityRoutes(createHumidityController, getByHamsterController)
}
