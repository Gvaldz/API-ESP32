package application

import (
	"esp32/src/internal/sensores/humidity/domain"
	"fmt"
)

type CreateHumidity struct {
	repo domain.HumidityRepository
}

func NewCreateHumidity(repo domain.HumidityRepository) *CreateHumidity {
	return &CreateHumidity{repo: repo}
}

func (c *CreateHumidity) Execute(humidity domain.Humidity) error {
	fmt.Printf("Guardando humedad en la base de datos: %+v\n", humidity)
	return c.repo.CreateHumidity(humidity)
}
