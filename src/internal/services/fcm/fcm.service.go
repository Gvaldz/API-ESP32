package fcm

import (
	"context"
	"os"
	"fmt"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
	"esp32/src/core"
)

type FCMSender struct {
	client *messaging.Client
}

func NewFCMSender(ctx context.Context, cfg core.FCMConfig) (*FCMSender, error) {
	if _, err := os.Stat(cfg.CredentialsPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("archivo de credenciales FCM no encontrado: %s", cfg.CredentialsPath)
	}

	opt := option.WithCredentialsFile(cfg.CredentialsPath)
	app, err := firebase.NewApp(ctx, &firebase.Config{
		ProjectID: cfg.ProjectID,
	}, opt)
	if err != nil {
		return nil, err
	}

	client, err := app.Messaging(ctx)
	if err != nil {
		return nil, err
	}

	return &FCMSender{client: client}, nil
}

func (s *FCMSender) SendNotification(ctx context.Context, deviceToken string, payload NotificationPayload) error {
	message := &messaging.Message{
		Token: deviceToken,
		Notification: &messaging.Notification{
			Title: payload.Title,
			Body:  payload.Body,
		},
		Data: payload.Data,
	}

	_, err := s.client.Send(ctx, message)
	return err
}

type NotificationPayload struct {
	Title string
	Body  string
	Data  map[string]string
}