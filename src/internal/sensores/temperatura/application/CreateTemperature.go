package application

import (
	"esp32/src/internal/sensores/temperatura/domain"
	"fmt"
)

type CreateTemperature struct {
	repo domain.TemperatureRepository
}

func NewCreateTemperature(repo domain.TemperatureRepository) *CreateTemperature {
	return &CreateTemperature{repo: repo}
}

func (c *CreateTemperature) Execute(temperature domain.Temperature) error {
	fmt.Printf("Guardando temperatura en la base de datos: %+v\n", temperature)
	return c.repo.CreateTemperature(temperature)
}
