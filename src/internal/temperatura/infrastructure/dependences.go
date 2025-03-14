package infrastructure

import (
	"database/sql"
	"esp32/src/internal/temperatura/application"
	"esp32/src/internal/temperatura/infrastructure/controllers"
	"esp32/src/core"
)

type TemperatureDependencies struct {
	DB   *sql.DB
	MQTT *core.MQTTConnection
}

func NewTemperatureDependencies(db *sql.DB, mqtt *core.MQTTConnection) *TemperatureDependencies {
	return &TemperatureDependencies{
		DB:   db,
		MQTT: mqtt,
	}
}

func (d *TemperatureDependencies) GetRoutes() *TemperatureRoutes {
	mqttProducer := NewMQTTProducer(d.MQTT.Client)
	temperatureRepo := NewTemperatureRepo(d.DB, mqttProducer)

	createTemperatureUseCase := application.NewCreateTemperature(temperatureRepo)
	getByHamsterUseCase := application.NewGetByHamster(temperatureRepo)


	createTemperatureController := controllers.NewCreateTemperatureController(createTemperatureUseCase)
	getByHamsterController := controllers.NewGetByHamsterController(getByHamsterUseCase)

	return NewTemperatureRoutes(
		createTemperatureController,
		getByHamsterController,
	)
}