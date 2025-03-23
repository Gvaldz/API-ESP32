package cmd

import (
	temperatureDeps "esp32/src/internal/temperatura/infrastructure"
	"esp32/src/core"
	"esp32/src/server"
	"log"
)

func Init() {
	// Conectar a la base de datos
	db, err := core.ConnectDB()
	if err != nil {
		log.Fatal("Error al conectar a la base de datos:", err)
	}

	// Conectar a RabbitMQ (AMQP)
	amqpConn, err := core.NewAMQPConnection() // Aquí usas AMQPConnection en lugar de MQTT
	if err != nil {
		log.Fatal("Error al conectar a RabbitMQ:", err)
	}
	defer amqpConn.Close()

	// Crear las dependencias de temperatura pasando la conexión AMQP
	temperatureDependencies := temperatureDeps.NewTemperatureDependencies(db, amqpConn)
	temperatureRoutes := temperatureDependencies.GetRoutes()

	// Ejecutar el servidor con las rutas de temperatura
	server.Run(temperatureRoutes)
}
