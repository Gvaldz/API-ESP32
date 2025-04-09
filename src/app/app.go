package app

import (
	"database/sql"
	"esp32/src/core"
	consumer_amqp		 "esp32/src/consumer_amqp"
	cages 				 "esp32/src/internal/sensores/cages/infrastructure"
	food 				 "esp32/src/internal/sensores/food/infrastructure"
	humidity 			 "esp32/src/internal/sensores/humidity/infrastructure"
	motion 				 "esp32/src/internal/sensores/motion/infrastructure"
	temperature 		 "esp32/src/internal/sensores/temperatura/infrastructure"
	websocketapp 		 "esp32/src/internal/services/websocket/application"
	websocketinfra 		 "esp32/src/internal/services/websocket/infrastructure"
	websocketc 			 "esp32/src/internal/services/websocket/infrastructure/controllers"
	"esp32/src/server"
)

type Application struct {
	DB            *sql.DB
	AMQPConn      *core.AMQPConnection
	Server        *server.Server
	AMQPConsumer  *consumer_amqp.RabbitMQConsumer
	Hasher		  *core.BcryptHasher
	tokenService  *core.JWTService
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
	wsService := websocketapp.NewWebSocketService()
	tokenService := core.NewJWTService()
    userRepo := core.NewUserRepository(db).(*core.UserRepository)
	

	wsHandler := websocketc.NewWebSocketController(wsService, *tokenService)
	wsRoutes := websocketinfra.NewWebSocketRoutes(wsHandler)
	cageDeps := cages.NewCageDependencies(db)
	tempDeps := temperature.NewTemperatureDependencies(db, amqpConn, wsService, userRepo)
	motionDeps := motion.NewMotionDependencies(db, amqpConn, wsService, userRepo)
	humidityDeps := humidity.NewHumidityDependencies(db, amqpConn, wsService, userRepo)
	foodDeps := food.NewFoodDependencies(db, amqpConn, wsService, userRepo)

    
    server := server.NewServer( 
        tempDeps.GetRoutes(),
        motionDeps.GetRoutes(),
        humidityDeps.GetRoutes(),
        foodDeps.GetRoutes(),
        cageDeps.GetRoutes(),
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
		DB:            db,
		AMQPConn:      amqpConn,
		Server:        server,
		AMQPConsumer:  consumer,
		Hasher:        hasher,
		tokenService:  tokenService,
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