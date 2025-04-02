package infrastructure

import (
	"database/sql"
	"esp32/src/core"
	"esp32/src/internal/temperatura/application"
	"esp32/src/internal/temperatura/infrastructure/controllers"
	websocket "esp32/src/internal/websocket/application"
	cages	  "esp32/src/internal/cages/infrastructure"
)

type TemperatureDependencies struct {
	DB   *sql.DB
	AMQP *core.AMQPConnection
	WsService *websocket.WebSocketService
}

func NewTemperatureDependencies(db *sql.DB, amqp *core.AMQPConnection, wsService *websocket.WebSocketService) *TemperatureDependencies {
	return &TemperatureDependencies{DB: db, AMQP: amqp, WsService: wsService}
}

func (d *TemperatureDependencies) GetRoutes() *TemperatureRoutes {
	temperatureRepo := NewTemperatureRepo(d.DB, nil)
	cageRepo := cages.NewCageRepo(d.DB)

	createTemperatureUseCase := application.NewCreateTemperature(temperatureRepo)
	getByHamsterUseCase := application.NewGetByHamster(temperatureRepo)

	createTemperatureController := controllers.NewCreateTemperatureController(createTemperatureUseCase, d.WsService, cageRepo)
	getByHamsterController := controllers.NewGetByHamsterController(getByHamsterUseCase)

	return NewTemperatureRoutes(createTemperatureController, getByHamsterController)
}
