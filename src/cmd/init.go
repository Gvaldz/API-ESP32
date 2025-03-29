package cmd

import (
	"esp32/src/app"
	"log"
)

func Init() {
	app, err := app.NewApplication()
	if err != nil {
		log.Fatal("Error al inicializar la aplicación:", err)
	}
	defer app.Close()

	if err := app.Start(); err != nil {
		log.Fatal("Error al iniciar la aplicación:", err)
	}
}