package infrastructure

import (
	"database/sql"
	"esp32/src/core"
	"esp32/src/internal/humidity/application"
	"esp32/src/internal/humidity/infrastructure/controllers"
)

type HumidityDependencies struct {
	DB   *sql.DB
	AMQP *core.AMQPConnection
}

func NewHumidityDependencies(db *sql.DB, amqp *core.AMQPConnection) *HumidityDependencies {
	return &HumidityDependencies{DB: db, AMQP: amqp}
}

func (d *HumidityDependencies) GetRoutes() *HumidityRoutes {

	humidityRepo := NewHumidityRepo(d.DB, nil)

	createHumidityUseCase := application.NewCreateHumidity(humidityRepo)
	getByHamsterUseCase := application.NewGetByHamster(humidityRepo)

	createHumidityController := controllers.NewCreateHumidityController(createHumidityUseCase)
	getByHamsterController := controllers.NewGetByHamsterController(getByHamsterUseCase)

	return NewHumidityRoutes(createHumidityController, getByHamsterController)
}
