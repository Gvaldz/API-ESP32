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
	// Crear el repositorio de movimiento
	motionRepo := NewMotionRepo(d.DB, nil)

	// Crear los casos de uso de movimiento
	createMotionUseCase := application.NewCreateMotion(motionRepo)
	getByHamsterUseCase := application.NewGetByHamster(motionRepo)

	// Crear el controlador de movimiento
	createMotionController := controllers.NewCreateMotionController(createMotionUseCase)
	getByHamsterController := controllers.NewGetByHamsterController(getByHamsterUseCase)

	// Devolver las rutas de movimiento con el controlador necesario
	return NewMotionRoutes(createMotionController, getByHamsterController)
}
