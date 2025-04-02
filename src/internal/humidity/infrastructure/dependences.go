package infrastructure

import (
	"database/sql"
	"esp32/src/core"
	"esp32/src/internal/humidity/application"
	"esp32/src/internal/humidity/infrastructure/controllers"
	websocket "esp32/src/internal/websocket/application"
	cages	  "esp32/src/internal/cages/infrastructure"
)

type HumidityDependencies struct {
	DB   *sql.DB
	AMQP *core.AMQPConnection
	WsService *websocket.WebSocketService
}

func NewHumidityDependencies(db *sql.DB, amqp *core.AMQPConnection, wsService *websocket.WebSocketService) *HumidityDependencies {
	return &HumidityDependencies{DB: db, AMQP: amqp, WsService: wsService}
}

func (d *HumidityDependencies) GetRoutes() *HumidityRoutes {

	humidityRepo := NewHumidityRepo(d.DB, nil)
	cageRepo := cages.NewCageRepo(d.DB)

	createHumidityUseCase := application.NewCreateHumidity(humidityRepo)
	getByHamsterUseCase := application.NewGetByHamster(humidityRepo)

	createHumidityController := controllers.NewCreateHumidityController(createHumidityUseCase, d.WsService, cageRepo)
	getByHamsterController := controllers.NewGetByHamsterController(getByHamsterUseCase)

	return NewHumidityRoutes(createHumidityController, getByHamsterController)
}
