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
	// Crear el repositorio de temperatura
	temperatureRepo := NewTemperatureRepo(d.DB, nil)

	// Crear los casos de uso de temperatura
	createTemperatureUseCase := application.NewCreateTemperature(temperatureRepo)
	getByHamsterUseCase := application.NewGetByHamster(temperatureRepo)


	// Crear el controlador de temperatura
	createTemperatureController := controllers.NewCreateTemperatureController(createTemperatureUseCase)
	getByHamsterController := controllers.NewGetByHamsterController(getByHamsterUseCase)


	// Devolver las rutas de temperatura con el controlador necesario
	return NewTemperatureRoutes(createTemperatureController, getByHamsterController)
}
