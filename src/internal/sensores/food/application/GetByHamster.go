package application

import (
	"esp32/src/internal/sensores/food/domain"
)

type GetByHamster struct {
	repo domain.FoodRepository
}

func NewGetByHamster(repo domain.FoodRepository) *GetByHamster {
	return &GetByHamster{repo: repo}
}

func (cp *GetByHamster) Execute(IDHamster string) ([]domain.Food, error) {
	return cp.repo.GetByHamster(IDHamster)
}
