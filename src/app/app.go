package app

import (
	"context"
	"database/sql"
	consumer_amqp "esp32/src/consumer_amqp"
	"esp32/src/core"
	login "esp32/src/internal/auth/infrastructure"
	cages "esp32/src/internal/cages/infrastructure"
	fcm "esp32/src/internal/fcm"
	food "esp32/src/internal/food/infrastructure"
	humidity "esp32/src/internal/humidity/infrastructure"
	motion "esp32/src/internal/motion/infrastructure"
	temperature "esp32/src/internal/temperatura/infrastructure"
	users "esp32/src/internal/users/infrastructure"
	websocketapp "esp32/src/internal/websocket/application"
	websocketinfra "esp32/src/internal/websocket/infrastructure"
	websocketc "esp32/src/internal/websocket/infrastructure/controllers"
	"esp32/src/server"
	"fmt"
)

type Application struct {
	DB            *sql.DB
	AMQPConn      *core.AMQPConnection
	Server        *server.Server
	AMQPConsumer  *consumer_amqp.RabbitMQConsumer
	Hasher		  *core.BcryptHasher
	tokenService  *core.JWTService
	fcmSender     *fcm.FCMSender
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
	fcmSender, err := fcm.NewFCMSender(context.Background(), core.Config.FCM)
	if err != nil {
    return nil, fmt.Errorf("error inicializando FCM: %v", err)
	}

	hasher := core.NewBcryptHasher(12)
	wsService := websocketapp.NewWebSocketService()
	tokenService := core.NewJWTService()
	usersDeps := users.NewUserDependencies(db, amqpConn, hasher, fcmSender)

	wsHandler := websocketc.NewWebSocketController(wsService, *tokenService)
	wsRoutes := websocketinfra.NewWebSocketRoutes(wsHandler)
	cageDeps := cages.NewCageDependencies(db)
	tempDeps := temperature.NewTemperatureDependencies(db, amqpConn, wsService, fcmSender, usersDeps.UserRepo)
	motionDeps := motion.NewMotionDependencies(db, amqpConn, wsService)
	humidityDeps := humidity.NewHumidityDependencies(db, amqpConn, wsService)
	foodDeps := food.NewFoodDependencies(db, amqpConn, wsService)
	loginDeps := login.NewAuthDependencies(db, hasher, usersDeps.UserRepo)

	server := server.NewServer(
		tempDeps.GetRoutes(),
		motionDeps.GetRoutes(),
		humidityDeps.GetRoutes(),
		foodDeps.GetRoutes(),
		usersDeps.GetRoutes(),
		cageDeps.GetRoutes(),
		loginDeps.GetRoutes(),
		wsRoutes,
	)

	consumer := consumer_amqp.NewRabbitMQConsumer(
		amqpConn,
		humidityDeps.GetRoutes().CreateHumidityController,
		tempDeps.GetRoutes().CreateTemperatureController,
		motionDeps.GetRoutes().CreateMotionController,
		foodDeps.GetRoutes().CreateStatusFoodController,
		cageDeps.GetRoutes().CreateCageController,
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