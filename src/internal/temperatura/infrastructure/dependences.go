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
	temperatureRepo := NewTemperatureRepo(d.DB, nil)

	createTemperatureUseCase := application.NewCreateTemperature(temperatureRepo)
	getByHamsterUseCase := application.NewGetByHamster(temperatureRepo)

	createTemperatureController := controllers.NewCreateTemperatureController(createTemperatureUseCase)
	getByHamsterController := controllers.NewGetByHamsterController(getByHamsterUseCase)

	return NewTemperatureRoutes(createTemperatureController, getByHamsterController)
}
