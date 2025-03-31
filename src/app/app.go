package app

import (
	"database/sql"
	"esp32/src/core"
	consumer_amqp	"esp32/src/consumer_amqp"
	humidity 		"esp32/src/internal/humidity/infrastructure"
	motion			"esp32/src/internal/motion/infrastructure"
	temperature		"esp32/src/internal/temperatura/infrastructure"
	food			"esp32/src/internal/food/infrastructure"	
	users			"esp32/src/internal/users/infrastructure"
	cages 			"esp32/src/internal/cages/infrastructure"

	"esp32/src/server"
)

type Application struct {
	DB            *sql.DB
	AMQPConn      *core.AMQPConnection
	Server        *server.Server
	AMQPConsumer  *consumer_amqp.RabbitMQConsumer
	Hasher		  *core.BcryptHasher
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

	hasher := core.NewBcryptHasher(12)

	tempDeps := temperature.NewTemperatureDependencies(db, amqpConn)
	motionDeps := motion.NewMotionDependencies(db, amqpConn)
	humidityDeps := humidity.NewHumidityDependencies(db, amqpConn)
	foodDeps := food.NewFoodDependencies(db, amqpConn)
	usersDeps := users.NewUserDependencies(db, amqpConn, hasher)
	cageDeps := cages.NewCageDependencies(db)

	server := server.NewServer(
		tempDeps.GetRoutes(),
		motionDeps.GetRoutes(),
		humidityDeps.GetRoutes(),
		foodDeps.GetRoutes(),
		usersDeps.GetRoutes(),
		cageDeps.GetRoutes(),
	)

	consumer := consumer_amqp.NewRabbitMQConsumer(
		amqpConn,
		humidityDeps.GetRoutes().CreateHumidityController,
		tempDeps.GetRoutes().CreateTemperatureController,
		motionDeps.GetRoutes().CreateMotionController,
		foodDeps.GetRoutes().CreateStatusFoodController,
	)

	return &Application{
		DB:           db,
		AMQPConn:     amqpConn,
		Server:       server,
		AMQPConsumer: consumer,
		Hasher: 	  hasher,	
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