package model

type SendTelegramMessageRequest struct {
	Text string `json:"text" binding:"required"`
}
