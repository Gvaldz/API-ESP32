package infrastructure

import (
	"database/sql"
	"esp32/src/core"
	"esp32/src/internal/food/application"
	"esp32/src/internal/food/infrastructure/controllers"
)

type FoodDependencies struct {
	DB   *sql.DB
	AMQP *core.AMQPConnection
}

func NewFoodDependencies(db *sql.DB, amqp *core.AMQPConnection) *FoodDependencies {
	return &FoodDependencies{DB: db, AMQP: amqp}
}

func (d *FoodDependencies) GetRoutes() *FoodRoutes {
	foodRepo := NewFoodRepo(d.DB, nil)

	createFoodUseCase := application.NewCreateStatusFood(foodRepo)
	getByHamsterUseCase := application.NewGetByHamster(foodRepo)

	createFoodController := controllers.NewCreateStatusFoodController(createFoodUseCase)
	getByHamsterController := controllers.NewGetByHamsterController(getByHamsterUseCase)


	return NewFoodRoutes(createFoodController, getByHamsterController)
}
