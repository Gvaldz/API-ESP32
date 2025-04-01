package infrastructure

import (
	"database/sql"
	"esp32/src/core"
	"esp32/src/internal/motion/application"
	"esp32/src/internal/motion/infrastructure/controllers"
)

type MotionDependencies struct {
	DB   *sql.DB
	AMQP *core.AMQPConnection
}

func NewMotionDependencies(db *sql.DB, amqp *core.AMQPConnection) *MotionDependencies {
	return &MotionDependencies{DB: db, AMQP: amqp}
}

func (d *MotionDependencies) GetRoutes() *MotionRoutes {
	motionRepo := NewMotionRepo(d.DB, nil)

	createMotionUseCase := application.NewCreateMotion(motionRepo)
	getByHamsterUseCase := application.NewGetByHamster(motionRepo)

	createMotionController := controllers.NewCreateMotionController(createMotionUseCase)
	getByHamsterController := controllers.NewGetByHamsterController(getByHamsterUseCase)

	return NewMotionRoutes(createMotionController, getByHamsterController)
}
