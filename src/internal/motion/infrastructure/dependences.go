package infrastructure

import (
	"database/sql"
	"esp32/src/core"
	"esp32/src/internal/motion/application"
	"esp32/src/internal/motion/infrastructure/controllers"
	websocket "esp32/src/internal/websocket/application"
	cages	  "esp32/src/internal/cages/infrastructure"
)

type MotionDependencies struct {
	DB   *sql.DB
	AMQP *core.AMQPConnection
	WsService *websocket.WebSocketService
}

func NewMotionDependencies(db *sql.DB, amqp *core.AMQPConnection, wsService *websocket.WebSocketService) *MotionDependencies {
	return &MotionDependencies{DB: db, AMQP: amqp, WsService: wsService}
}

func (d *MotionDependencies) GetRoutes() *MotionRoutes {
	motionRepo := NewMotionRepo(d.DB, nil)
	cageRepo := cages.NewCageRepo(d.DB)

	createMotionUseCase := application.NewCreateMotion(motionRepo)
	getByHamsterUseCase := application.NewGetByHamster(motionRepo)

	createMotionController := controllers.NewCreateMotionController(createMotionUseCase, d.WsService, cageRepo)
	getByHamsterController := controllers.NewGetByHamsterController(getByHamsterUseCase)

	return NewMotionRoutes(createMotionController, getByHamsterController)
}
