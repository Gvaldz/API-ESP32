package core

import (
	"context"
	"fmt"
	"log"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

type FCMClient struct {
	app    *firebase.App
	client *messaging.Client
}

// Inicializa Firebase
func NewFCMClient() (*FCMClient, error) {
	opt := option.WithCredentialsFile("esp32/src/config/firebase_key.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, fmt.Errorf("error inicializando Firebase: %v", err)
	}

	client, err := app.Messaging(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error obteniendo cliente de FCM: %v", err)
	}

	return &FCMClient{app: app, client: client}, nil
}

// Envía una notificación a un dispositivo
func (f *FCMClient) SendNotification(token string, title string, body string) error {
	msg := &messaging.Message{
		Token: token, // Token del dispositivo
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
	}

	// Envía la notificación
	resp, err := f.client.Send(context.Background(), msg)
	if err != nil {
		return fmt.Errorf("error enviando notificación: %v", err)
	}

	log.Printf("Notificación enviada con éxito: %s\n", resp)
	return nil
}
