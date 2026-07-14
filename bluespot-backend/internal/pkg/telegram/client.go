package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	defaultBaseURL  = "https://api.telegram.org"
	maxResponseSize = 1 << 20
)

type Client struct {
	botToken   string
	chatID     string
	baseURL    string
	httpClient *http.Client
}

type sendMessageRequest struct {
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}

type telegramResponse struct {
	OK          bool   `json:"ok"`
	Description string `json:"description"`
}

func NewClient(botToken, chatID string) *Client {
	return &Client{
		botToken: botToken,
		chatID:   chatID,
		baseURL:  defaultBaseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) SendMessage(ctx context.Context, text string) error {
	if strings.TrimSpace(c.botToken) == "" {
		return errors.New("telegram bot token is empty")
	}
	if strings.TrimSpace(c.chatID) == "" {
		return errors.New("telegram chat id is empty")
	}
	if strings.TrimSpace(text) == "" {
		return errors.New("telegram message text is empty")
	}

	body, err := json.Marshal(sendMessageRequest{
		ChatID: c.chatID,
		Text:   text,
	})
	if err != nil {
		return fmt.Errorf("marshal telegram message: %w", err)
	}

	url := fmt.Sprintf("%s/bot%s/sendMessage", strings.TrimRight(c.baseURL, "/"), c.botToken)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return errors.New("create telegram request failed")
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return errors.New("send telegram request failed")
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("telegram returned http status %d", resp.StatusCode)
	}

	var result telegramResponse
	if err := json.NewDecoder(io.LimitReader(resp.Body, maxResponseSize)).Decode(&result); err != nil {
		return fmt.Errorf("decode telegram response: %w", err)
	}
	if !result.OK {
		if strings.TrimSpace(result.Description) == "" {
			return errors.New("telegram rejected message")
		}
		return fmt.Errorf("telegram rejected message: %s", result.Description)
	}

	return nil
}
