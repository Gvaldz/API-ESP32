package cmd

import (
	temperatureDeps "esp32/src/internal/temperatura/infrastructure"
	"esp32/src/core"
	"esp32/src/server"
	"log"
)

func Init() {
	db, err := core.ConnectDB()
	if err != nil {
		log.Fatal("Error al conectar a la base de datos:", err)
	}

	mqttConn, err := core.NewMQTTConnection()
	if err != nil {
		log.Fatal("Error al conectar a MQTT:", err)
	}
	defer mqttConn.Close()

	temperatureDependencies := temperatureDeps.NewTemperatureDependencies(db, mqttConn)
	temperatureRoutes := temperatureDependencies.GetRoutes()

	server.Run(temperatureRoutes)
}