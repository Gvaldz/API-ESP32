package infrastructure

import (
	"database/sql"
	"esp32/src/core"
	cages "esp32/src/internal/sensores/cages/infrastructure"
	"esp32/src/internal/sensores/humidity/application"
	"esp32/src/internal/sensores/humidity/infrastructure/controllers"
	fcm "esp32/src/internal/services/fcm"
	websocket "esp32/src/internal/services/websocket/application"
)

type HumidityDependencies struct {
	DB        *sql.DB
	AMQP      *core.AMQPConnection
	WsService *websocket.WebSocketService
	FCMSender *fcm.FCMSender
	UserRepo  *core.UserRepository
}

func NewHumidityDependencies(
	db *sql.DB, 
	amqp *core.AMQPConnection, 
	wsService *websocket.WebSocketService, 
	fcmSender *fcm.FCMSender, 
	userRepo *core.UserRepository,
) *HumidityDependencies {
	return &HumidityDependencies{
		DB:        db,
		AMQP:      amqp,
		WsService: wsService,
		FCMSender: fcmSender,
		UserRepo:  userRepo,
	}
}

func (d *HumidityDependencies) GetRoutes() *HumidityRoutes {
	humidityRepo := NewHumidityRepo(d.DB, nil)
	cageRepo := cages.NewCageRepo(d.DB)

	createHumidityUseCase := application.NewCreateHumidity(humidityRepo)
	getByHamsterUseCase := application.NewGetByHamster(humidityRepo)

	createHumidityController := controllers.NewCreateHumidityController(
		createHumidityUseCase, 
		d.WsService, 
		cageRepo,
		d.UserRepo,
		d.FCMSender,
	)
	getByHamsterController := controllers.NewGetByHamsterController(getByHamsterUseCase)

	return NewHumidityRoutes(createHumidityController, getByHamsterController)
}