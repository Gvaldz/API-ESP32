package application

import (
	"esp32/src/internal/sensores/humidity/domain"
)

type GetByHamster struct {
	repo domain.HumidityRepository
}

func NewGetByHamster(repo domain.HumidityRepository) *GetByHamster {
	return &GetByHamster{repo: repo}
}

func (cp *GetByHamster) Execute(IDHamster string) ([]domain.Humidity, error) {
	return cp.repo.GetByHamster(IDHamster)
}
