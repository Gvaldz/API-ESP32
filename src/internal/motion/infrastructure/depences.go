package infrastructure

import (
	"database/sql"
	"esp32/src/core"
	"esp32/src/internal/motion/application"
	"esp32/src/internal/motion/infrastructure/controllers"
)

type MotionDependences struct {
	DB   *sql.DB
	AMQP *core.AMQPConnection
}

func NewMotionDependences(db *sql.DB, amqp *core.AMQPConnection) *MotionDependences {
	return &MotionDependences{DB: db, AMQP: amqp}
}

func (d *MotionDependences) GetRoutes() *MotionRoutes {
	amqpConsumer := NewAMQPConsumer(d.AMQP, nil)

	motionRepo := NewMotionRepo(d.DB, amqpConsumer)
	createMotionUseCase := application.NewCreateMotion(motionRepo)
	getByHamsterUseCase := application.NewGetByHamster(motionRepo)

	createMotionController := controllers.NewCreateMotionController(createMotionUseCase)
	getByHamsterController := controllers.NewGetByHamsterController(getByHamsterUseCase)

	amqpConsumer.createMotionC = createMotionController
	go amqpConsumer.Consume()

	return NewMotionRoutes(createMotionController, getByHamsterController)
}
