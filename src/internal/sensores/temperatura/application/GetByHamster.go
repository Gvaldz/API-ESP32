package application

import (
	"esp32/src/internal/sensores/temperatura/domain"
)

type GetByHamster struct {
	repo domain.TemperatureRepository
}

func NewGetByHamster(repo domain.TemperatureRepository) *GetByHamster {
	return &GetByHamster{repo: repo}
}

func (cp *GetByHamster) Execute(IDHamster string) ([]domain.Temperature, error) {
	return cp.repo.GetByHamster(IDHamster)
}
