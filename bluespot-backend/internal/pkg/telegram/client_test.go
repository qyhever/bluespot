package telegram

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestClientSendMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/bottest-token/sendMessage" {
			t.Errorf("path = %s, want /bottest-token/sendMessage", r.URL.Path)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", r.Header.Get("Content-Type"))
		}

		var request sendMessageRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if request.ChatID != "test-chat" {
			t.Errorf("ChatID = %q, want test-chat", request.ChatID)
		}
		if request.Text != "hello" {
			t.Errorf("Text = %q, want hello", request.Text)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{}}`))
	}))
	defer server.Close()

	client := NewClient("test-token", "test-chat")
	client.baseURL = server.URL
	client.httpClient = server.Client()

	if err := client.SendMessage(context.Background(), "hello"); err != nil {
		t.Fatalf("SendMessage() error = %v", err)
	}
}

func TestClientSendMessageErrors(t *testing.T) {
	tests := []struct {
		name       string
		botToken   string
		chatID     string
		text       string
		statusCode int
		response   string
		wantError  string
	}{
		{name: "missing bot token", chatID: "chat", text: "hello", wantError: "bot token is empty"},
		{name: "missing chat id", botToken: "token", text: "hello", wantError: "chat id is empty"},
		{name: "empty text", botToken: "token", chatID: "chat", text: " ", wantError: "message text is empty"},
		{name: "http error", botToken: "token", chatID: "chat", text: "hello", statusCode: http.StatusBadGateway, response: `{"ok":false}`, wantError: "http status 502"},
		{name: "telegram error", botToken: "token", chatID: "chat", text: "hello", statusCode: http.StatusOK, response: `{"ok":false,"description":"chat not found"}`, wantError: "chat not found"},
		{name: "invalid response", botToken: "token", chatID: "chat", text: "hello", statusCode: http.StatusOK, response: `{`, wantError: "decode telegram response"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := NewClient(tt.botToken, tt.chatID)
			client.baseURL = server.URL
			client.httpClient = server.Client()

			err := client.SendMessage(context.Background(), tt.text)
			if err == nil {
				t.Fatal("SendMessage() error = nil, want error")
			}
			if !strings.Contains(err.Error(), tt.wantError) {
				t.Fatalf("SendMessage() error = %q, want substring %q", err, tt.wantError)
			}
		})
	}
}

func TestClientSendMessageDoesNotExposeBotTokenInTransportError(t *testing.T) {
	const botToken = "secret-bot-token"

	client := NewClient(botToken, "test-chat")
	client.httpClient = &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("request failed for " + req.URL.String())
		}),
	}

	err := client.SendMessage(context.Background(), "hello")
	if err == nil {
		t.Fatal("SendMessage() error = nil, want error")
	}
	if strings.Contains(err.Error(), botToken) {
		t.Fatalf("SendMessage() error exposed bot token: %q", err)
	}
}
