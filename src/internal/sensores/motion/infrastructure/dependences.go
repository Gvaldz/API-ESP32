package infrastructure

import (
	"database/sql"
	"esp32/src/core"
	cages "esp32/src/internal/sensores/cages/infrastructure"
	"esp32/src/internal/sensores/motion/application"
	"esp32/src/internal/sensores/motion/infrastructure/controllers"
	fcm "esp32/src/internal/services/fcm"
	websocket "esp32/src/internal/services/websocket/application"
)

type MotionDependencies struct {
	DB        *sql.DB
	AMQP      *core.AMQPConnection
	WsService *websocket.WebSocketService
	FCMSender *fcm.FCMSender
	UserRepo  *core.UserRepository
}

func NewMotionDependencies(
	db *sql.DB, 
	amqp *core.AMQPConnection, 
	wsService *websocket.WebSocketService, 
	fcmSender *fcm.FCMSender, 
	userRepo *core.UserRepository,
) *MotionDependencies {
	return &MotionDependencies{
		DB:        db,
		AMQP:      amqp,
		WsService: wsService,
		FCMSender: fcmSender,
		UserRepo:  userRepo,
	}
}

func (d *MotionDependencies) GetRoutes() *MotionRoutes {
	motionRepo := NewMotionRepo(d.DB, nil)
	cageRepo := cages.NewCageRepo(d.DB)

	createMotionUseCase := application.NewCreateMotion(motionRepo)
	getByHamsterUseCase := application.NewGetByHamster(motionRepo)

	createMotionController := controllers.NewCreateMotionController(
		createMotionUseCase, 
		d.WsService, 
		cageRepo,
		d.UserRepo,
		d.FCMSender,
	)
	getByHamsterController := controllers.NewGetByHamsterController(getByHamsterUseCase)

	return NewMotionRoutes(createMotionController, getByHamsterController)
}