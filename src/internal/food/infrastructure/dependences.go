package infrastructure

import (
	"database/sql"
	"esp32/src/core"
	"esp32/src/internal/food/application"
	"esp32/src/internal/food/infrastructure/controllers"
	websocket "esp32/src/internal/websocket/application"
	cages "esp32/src/internal/cages/infrastructure"
	fcm "esp32/src/internal/fcm"
)

type FoodDependencies struct {
	DB        *sql.DB
	AMQP      *core.AMQPConnection
	WsService *websocket.WebSocketService
	FCMSender *fcm.FCMSender
	UserRepo  *core.UserRepository
}

func NewFoodDependencies(
	db *sql.DB, 
	amqp *core.AMQPConnection, 
	wsService *websocket.WebSocketService, 
	fcmSender *fcm.FCMSender, 
	userRepo *core.UserRepository,
) *FoodDependencies {
	return &FoodDependencies{
		DB:        db,
		AMQP:      amqp,
		WsService: wsService,
		FCMSender: fcmSender,
		UserRepo:  userRepo,
	}
}

func (d *FoodDependencies) GetRoutes() *FoodRoutes {
	foodRepo := NewFoodRepo(d.DB, nil)
	cageRepo := cages.NewCageRepo(d.DB)

	createFoodUseCase := application.NewCreateStatusFood(foodRepo)
	getByHamsterUseCase := application.NewGetByHamster(foodRepo)

	createFoodController := controllers.NewCreateStatusFoodController(
		createFoodUseCase, 
		d.WsService, 
		cageRepo,
		d.UserRepo,
		d.FCMSender,
	)
	getByHamsterController := controllers.NewGetByHamsterController(getByHamsterUseCase)

	return NewFoodRoutes(createFoodController, getByHamsterController)
}