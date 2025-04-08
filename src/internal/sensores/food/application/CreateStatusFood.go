package application

import (
	"esp32/src/internal/sensores/food/domain"
	"fmt"
)

type CreateStatusFood struct {
	repo domain.FoodRepository
}

func NewCreateStatusFood(repo domain.FoodRepository) *CreateStatusFood {
	return &CreateStatusFood{repo: repo}
}

func (c *CreateStatusFood) Execute(food domain.Food) error {
	fmt.Printf("Guardando estatus de alimento en la base de datos: %+v\n", food)
	return c.repo.CreateStatusFood(food)
}
