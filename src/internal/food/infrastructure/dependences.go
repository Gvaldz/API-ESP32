package infrastructure

import (
	"database/sql"
	"esp32/src/core"
	"esp32/src/internal/food/application"
	"esp32/src/internal/food/infrastructure/controllers"
	websocket "esp32/src/internal/websocket/application"
	cages	  "esp32/src/internal/cages/infrastructure"
)
type FoodDependencies struct {
	DB   *sql.DB
	AMQP *core.AMQPConnection
	WsService *websocket.WebSocketService
}

func NewFoodDependencies(db *sql.DB, amqp *core.AMQPConnection, wsService *websocket.WebSocketService) *FoodDependencies {
	return &FoodDependencies{DB: db, AMQP: amqp, WsService: wsService}
}

func (d *FoodDependencies) GetRoutes() *FoodRoutes {
	foodRepo := NewFoodRepo(d.DB, nil)
	cageRepo := cages.NewCageRepo(d.DB)

	createFoodUseCase := application.NewCreateStatusFood(foodRepo)
	getByHamsterUseCase := application.NewGetByHamster(foodRepo)

	createFoodController := controllers.NewCreateStatusFoodController(createFoodUseCase, d.WsService, cageRepo)
	getByHamsterController := controllers.NewGetByHamsterController(getByHamsterUseCase)


	return NewFoodRoutes(createFoodController, getByHamsterController)
}
