package service

import "context"

type TelegramSender interface {
	SendMessage(ctx context.Context, text string) error
}

type TelegramService struct {
	sender TelegramSender
}

func NewTelegramService(sender TelegramSender) *TelegramService {
	return &TelegramService{sender: sender}
}

func (s *TelegramService) SendMessage(ctx context.Context, text string) error {
	return s.sender.SendMessage(ctx, text)
}
