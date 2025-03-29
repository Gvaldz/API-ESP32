package infrastructure

import (
	"database/sql"
	"esp32/src/core"
	"esp32/src/internal/motion/application"
	"esp32/src/internal/motion/infrastructure/controllers"
	"log"
)

type MotionDependences struct {
	DB   *sql.DB
	AMQP *core.AMQPConnection
}

func NewMotionDependences(db *sql.DB, amqp *core.AMQPConnection) *MotionDependences {
	return &MotionDependences{DB: db, AMQP: amqp}
}

func (d *MotionDependences) GetRoutes() *MotionRoutes {
	// Inicialización del consumidor AMQP
	amqpConsumer := NewAMQPConsumer(d.AMQP, nil)

	// Repositorio de movimientos
	motionRepo := NewMotionRepo(d.DB, amqpConsumer)

	// Casos de uso
	createMotionUseCase := application.NewCreateMotion(motionRepo)
	getByHamsterUseCase := application.NewGetByHamster(motionRepo)

	// Verificar que los casos de uso no sean nil
	if createMotionUseCase == nil || getByHamsterUseCase == nil {
		log.Fatalf("Error al crear los casos de uso: uno o más casos de uso son nil")
	}

	// Controladores
	createMotionController := controllers.NewCreateMotionController(createMotionUseCase)
	getByHamsterController := controllers.NewGetByHamsterController(getByHamsterUseCase)

	// Verificar que los controladores no sean nil
	if createMotionController == nil || getByHamsterController == nil {
		log.Fatalf("Error al crear los controladores: uno o más controladores son nil")
	}

	// Asignar el controlador al consumidor AMQP
	amqpConsumer.createMotionC = createMotionController

	// Verificar que el controlador esté correctamente asignado
	if amqpConsumer.createMotionC == nil {
		log.Fatalf("Error al asignar el controlador de movimiento al consumidor AMQP")
	}

	// Iniciar el consumo en un goroutine
	go amqpConsumer.Consume()

	// Retornar las rutas de movimiento
	return NewMotionRoutes(createMotionController, getByHamsterController)
}
