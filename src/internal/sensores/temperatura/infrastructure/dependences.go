package infrastructure

import (
	"database/sql"
	"esp32/src/core"
	cages "esp32/src/internal/sensores/cages/infrastructure"
	"esp32/src/internal/sensores/temperatura/application"
	"esp32/src/internal/sensores/temperatura/infrastructure/controllers"
	fcm "esp32/src/internal/services/fcm"
	websocket "esp32/src/internal/services/websocket/application"
)

type TemperatureDependencies struct {
    DB        *sql.DB
    AMQP      *core.AMQPConnection
    WsService *websocket.WebSocketService
    FCMSender *fcm.FCMSender 
    UserRepo  *core.UserRepository  
}

func NewTemperatureDependencies(
    db *sql.DB, 
    amqp *core.AMQPConnection, 
    wsService *websocket.WebSocketService, 
    fcmSender *fcm.FCMSender, 
    userRepo *core.UserRepository,  
) *TemperatureDependencies {
    return &TemperatureDependencies{
        DB:        db,
        AMQP:      amqp,
        WsService: wsService,
        FCMSender: fcmSender,
        UserRepo:  userRepo,
    }
}

func (d *TemperatureDependencies) GetRoutes() *TemperatureRoutes {
	temperatureRepo := NewTemperatureRepo(d.DB, nil)
	cageRepo := cages.NewCageRepo(d.DB)

	createTemperatureUseCase := application.NewCreateTemperature(temperatureRepo)
	getByHamsterUseCase := application.NewGetByHamster(temperatureRepo)

	createTemperatureController := controllers.NewCreateTemperatureController(
		createTemperatureUseCase, 
		d.WsService, 
		cageRepo,
		d.UserRepo,    
		d.FCMSender,   
	)
	getByHamsterController := controllers.NewGetByHamsterController(getByHamsterUseCase)

	return NewTemperatureRoutes(createTemperatureController, getByHamsterController)
}