package core

import (
	"log"
	"os"
)

type FCMConfig struct {
	CredentialsPath string
	ProjectID       string
}

type AppConfig struct {
	FCM FCMConfig
}

var Config AppConfig

func LoadConfig() {
	Config.FCM.CredentialsPath = os.Getenv("FCM_CREDENTIALS_PATH")
	Config.FCM.ProjectID = os.Getenv("FCM_PROJECT_ID")

	if Config.FCM.CredentialsPath == "" {
		Config.FCM.CredentialsPath = "src/internal/fcm/serviceAccountKey.json"
	}

	if Config.FCM.ProjectID == "" {
		Config.FCM.ProjectID = "rodismart-77622"
	}

	log.Printf("Ruta de credenciales FCM: %s\n", Config.FCM.CredentialsPath)
}
