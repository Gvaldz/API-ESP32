package cmd

import (
	"esp32/src/app"
	"esp32/src/core"
	"log"
)

func Init() {
	core.LoadConfig()

	log.Printf("Ruta de credenciales FCM: %s\n", core.Config.FCM.CredentialsPath)

	app, err := app.NewApplication()
	if err != nil {
		log.Fatal("Error al inicializar la aplicación:", err)
	}
	defer app.Close()

	if err := app.Start(); err != nil {
		log.Fatal("Error al iniciar la aplicación:", err)
	}
}
