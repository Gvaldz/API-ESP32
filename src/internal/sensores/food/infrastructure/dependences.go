package infrastructure

import (
	"database/sql"
	"esp32/src/core"
	cages "esp32/src/internal/sensores/cages/infrastructure"
	"esp32/src/internal/sensores/food/application"
	"esp32/src/internal/sensores/food/infrastructure/controllers"
	websocket "esp32/src/internal/services/websocket/application"
)

type FoodDependencies struct {
	DB        *sql.DB
	AMQP      *core.AMQPConnection
	WsService *websocket.WebSocketService
	UserRepo  *core.UserRepository
}

func NewFoodDependencies(
	db *sql.DB,
	amqp *core.AMQPConnection,
	wsService *websocket.WebSocketService,
	userRepo *core.UserRepository,
) *FoodDependencies {
	return &FoodDependencies{
		DB:        db,
		AMQP:      amqp,
		WsService: wsService,
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
	)
	getByHamsterController := controllers.NewGetByHamsterController(getByHamsterUseCase)

	return NewFoodRoutes(createFoodController, getByHamsterController)
}
