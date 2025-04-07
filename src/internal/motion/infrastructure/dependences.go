package infrastructure

import (
	"database/sql"
	"esp32/src/core"
	"esp32/src/internal/motion/application"
	"esp32/src/internal/motion/infrastructure/controllers"
	websocket "esp32/src/internal/websocket/application"
	cages "esp32/src/internal/cages/infrastructure"
	fcm "esp32/src/internal/fcm"
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