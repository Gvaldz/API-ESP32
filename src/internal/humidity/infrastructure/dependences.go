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
	// Crear el repositorio de humedad
	humidityRepo := NewHumidityRepo(d.DB, nil)

	// Crear los casos de uso de humedad
	createHumidityUseCase := application.NewCreateHumidity(humidityRepo)
	getByHamsterUseCase := application.NewGetByHamster(humidityRepo)

	// Crear los controladores de humedad
	createHumidityController := controllers.NewCreateHumidityController(createHumidityUseCase)
	getByHamsterController := controllers.NewGetByHamsterController(getByHamsterUseCase)

	// Devolver las rutas de humedad con el controlador necesario
	return NewHumidityRoutes(createHumidityController, getByHamsterController)
}
