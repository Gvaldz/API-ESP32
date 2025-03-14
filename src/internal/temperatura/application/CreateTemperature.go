package application

import (
	"esp32/src/internal/temperatura/domain"
)
type CreateTemperature struct {
	repo domain.TemperatureRepository
}

func NewCreateTemperature(repo domain.TemperatureRepository) *CreateTemperature {
	return &CreateTemperature{repo: repo}
}

func (cp *CreateTemperature) Execute(temperature domain.Temperature) error{
	
	return cp.repo.CreateTemperature(temperature)
}