package infrastructure

import (
	"database/sql"
	"esp32/src/core"
	"esp32/src/internal/temperatura/application"
	"esp32/src/internal/temperatura/infrastructure/controllers"
	websocket "esp32/src/internal/websocket/application"
	cages	  "esp32/src/internal/cages/infrastructure"
	fcm		  "esp32/src/internal/fcm"
	users      "esp32/src/internal/users/infrastructure"
)

type TemperatureDependencies struct {
    DB        *sql.DB
    AMQP      *core.AMQPConnection
    WsService *websocket.WebSocketService
    FCMSender *fcm.FCMSender 
    UserRepo  *users.UsersRepo 
}

func NewTemperatureDependencies(db *sql.DB, amqp *core.AMQPConnection, wsService *websocket.WebSocketService, fcmSender *fcm.FCMSender, userRepo *users.UsersRepo) *TemperatureDependencies {
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
		*d.UserRepo,    
		d.FCMSender,   
	)
	getByHamsterController := controllers.NewGetByHamsterController(getByHamsterUseCase)

	return NewTemperatureRoutes(createTemperatureController, getByHamsterController)
}