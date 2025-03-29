package cmd

import (
	temperatureDeps "esp32/src/internal/temperatura/infrastructure"
	motionDeps	 	"esp32/src/internal/motion/infrastructure"
	humidityDeps 	"esp32/src/internal/humidity/infrastructure"
	"esp32/src/core"
	"esp32/src/server"
	"log"
)

func Init() {
	db, err := core.ConnectDB()
	if err != nil {
		log.Fatal("Error al conectar a la base de datos:", err)
	}

	amqpConn, err := core.NewAMQPConnection() 
	if err != nil {
		log.Fatal("Error al conectar a RabbitMQ:", err)
	}
	defer amqpConn.Close()

	temperatureDependencies := temperatureDeps.NewTemperatureDependencies(db, amqpConn)
	temperatureRoutes := temperatureDependencies.GetRoutes()

	motionDependences := motionDeps.NewMotionDependences(db, amqpConn)
	motionRoutes := motionDependences.GetRoutes()

	humidityDependences := humidityDeps.NewHumidityDependencies(db, amqpConn)
	humidityRoutes := humidityDependences.GetRoutes()

	server.Run(temperatureRoutes, motionRoutes, humidityRoutes)
}
