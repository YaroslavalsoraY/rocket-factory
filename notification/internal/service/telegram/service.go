package telegram

import (
	"bytes"
	"context"
	"embed"
	"text/template"

	"go.uber.org/zap"
	"notification/internal/client/http"
	"notification/internal/model"
	"platform/pkg/logger"
)

const chatID = 723424205

//go:embed templates/paid_notification.tmpl
var paidNotification embed.FS

//go:embed templates/assembled_notification.tmpl
var assembledNotification embed.FS

type paidTemplateData struct {
	OrderUUID       string
	PaymentMethod   string
	TransactionUUID string
}

type assembledTemplateData struct {
	OrderUUID    string
	UserUUID     string
	BuildTimeSec int64
}

var paidTemplate = template.Must(template.ParseFS(paidNotification, "templates/paid_notification.tmpl"))

var assembledTemplate = template.Must(template.ParseFS(assembledNotification, "templates/assembled_notification.tmpl"))

type service struct {
	telegramClient http.TelegramClient
}

func NewService(telegramClient http.TelegramClient) *service {
	return &service{
		telegramClient: telegramClient,
	}
}

func (s *service) SendOrderPaidNotification(ctx context.Context, event model.OrderPaid) error {
	message, err := s.BuildOrderPaidMessage(event)
	if err != nil {
		return err
	}

	err = s.telegramClient.SendMessage(ctx, chatID, message)
	if err != nil {
		return err
	}

	logger.Info(ctx, "Telegram message sent to chat", zap.Int("chat_id", chatID), zap.String("message", message))
	return nil
}

func (s *service) BuildOrderPaidMessage(event model.OrderPaid) (string, error) {
	data := paidTemplateData{
		OrderUUID:       event.OrderUUID,
		PaymentMethod:   event.PaymentMethod,
		TransactionUUID: event.TransactionUUID,
	}

	var buf bytes.Buffer
	err := paidTemplate.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (s *service) SendShipAssembledNotification(ctx context.Context, event model.ShipAssembled) error {
	message, err := s.BuildShipAssembledMessage(event)
	if err != nil {
		return err
	}

	err = s.telegramClient.SendMessage(ctx, chatID, message)
	if err != nil {
		return err
	}

	logger.Info(ctx, "Telegram message sent to chat", zap.Int("chat_id", chatID), zap.String("message", message))
	return nil
}

func (s *service) BuildShipAssembledMessage(event model.ShipAssembled) (string, error) {
	data := assembledTemplateData{
		OrderUUID:    event.OrderUUID,
		UserUUID:     event.UserUUID,
		BuildTimeSec: event.BuildTimeSec,
	}

	var buf bytes.Buffer
	err := assembledTemplate.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
