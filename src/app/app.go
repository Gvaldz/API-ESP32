package app

import (
	"context"
	"database/sql"
	"esp32/src/core"
	consumer_amqp		 "esp32/src/consumer_amqp"
	login 				 "esp32/src/internal/services/auth/infrastructure"
	cages 				 "esp32/src/internal/sensores/cages/infrastructure"
	fcm 				 "esp32/src/internal/services/fcm"
	food 				 "esp32/src/internal/sensores/food/infrastructure"
	humidity 			 "esp32/src/internal/sensores/humidity/infrastructure"
	motion 				 "esp32/src/internal/sensores/motion/infrastructure"
	temperature 		 "esp32/src/internal/sensores/temperatura/infrastructure"
	users 				 "esp32/src/internal/users/infrastructure"
	websocketapp 		 "esp32/src/internal/services/websocket/application"
	websocketinfra 		 "esp32/src/internal/services/websocket/infrastructure"
	middleware 			 "esp32/src/server/middleware"
	websocketc 			 "esp32/src/internal/services/websocket/infrastructure/controllers"
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
    userRepo := core.NewUserRepository(db).(*core.UserRepository)
    authRepo := &core.AuthRepository{DB: db}
	
	usersDeps := users.NewUserDependencies(
		db,
		amqpConn,
		hasher,
		fcmSender,
		tokenService,
		authRepo, 
		userRepo, 
	)
	

	loginDeps := login.NewAuthDependencies(
		db,
		hasher,
		userRepo, 
	)

	authMiddleware := middleware.AuthMiddleware(
		tokenService,
		authRepo,
		"usuario",
	)

	wsHandler := websocketc.NewWebSocketController(wsService, *tokenService)
	wsRoutes := websocketinfra.NewWebSocketRoutes(wsHandler)
	cageDeps := cages.NewCageDependencies(db)
	tempDeps := temperature.NewTemperatureDependencies(db, amqpConn, wsService, fcmSender, userRepo)
	motionDeps := motion.NewMotionDependencies(db, amqpConn, wsService, fcmSender, userRepo)
	humidityDeps := humidity.NewHumidityDependencies(db, amqpConn, wsService, fcmSender, userRepo)
	foodDeps := food.NewFoodDependencies(db, amqpConn, wsService, fcmSender, userRepo)

    
    server := server.NewServer( 
        tempDeps.GetRoutes(),
        motionDeps.GetRoutes(),
        humidityDeps.GetRoutes(),
        foodDeps.GetRoutes(),
        usersDeps.GetRoutes(),
        cageDeps.GetRoutes(),
        loginDeps.GetRoutes(),
        wsRoutes,
        authMiddleware,
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
		DB:            db,
		AMQPConn:      amqpConn,
		Server:        server,
		AMQPConsumer:  consumer,
		Hasher:        hasher,
		tokenService:  tokenService,
		fcmSender:     fcmSender,
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