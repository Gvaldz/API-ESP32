package app

import (
	"database/sql"
	"esp32/src/core"
	consumer_amqp	"esp32/src/internal/consumer_amqp"
	humidity 		"esp32/src/internal/humidity/infrastructure"
	motion			"esp32/src/internal/motion/infrastructure"
	temperatura		"esp32/src/internal/temperatura/infrastructure"
	"esp32/src/server"
)

type Application struct {
	DB            *sql.DB
	AMQPConn      *core.AMQPConnection
	Server        *server.Server
	AMQPConsumer  *consumer_amqp.RabbitMQConsumer
}

func NewApplication() (*Application, error) {
	db, err := core.ConnectDB()
	if err != nil {
		return nil, err
	}

	amqpConn, err := core.NewAMQPConnection()
	if err != nil {
		return nil, err
	}

	tempDeps := temperatura.NewTemperatureDependencies(db, amqpConn)
	motionDeps := motion.NewMotionDependencies(db, amqpConn)
	humidityDeps := humidity.NewHumidityDependencies(db, amqpConn)

	server := server.NewServer(
		tempDeps.GetRoutes(),
		motionDeps.GetRoutes(),
		humidityDeps.GetRoutes(),
	)

	consumer := consumer_amqp.NewRabbitMQConsumer(
		amqpConn,
		humidityDeps.GetRoutes().CreateHumidityController,
		tempDeps.GetRoutes().CreateTemperatureController,
		motionDeps.GetRoutes().CreateMotionController,
	)

	return &Application{
		DB:           db,
		AMQPConn:     amqpConn,
		Server:       server,
		AMQPConsumer: consumer,
	}, nil
}

func (a *Application) Start() error {
	go a.AMQPConsumer.Start()
	return a.Server.Run()
}

func (a *Application) Close() {
	if a.AMQPConn != nil {
		a.AMQPConn.Close()
	}
}