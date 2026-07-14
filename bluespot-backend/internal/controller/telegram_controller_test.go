package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"bluespot/internal/service"

	"github.com/gin-gonic/gin"
)

type telegramSenderStub struct {
	err  error
	text string
}

func (s *telegramSenderStub) SendMessage(_ context.Context, text string) error {
	s.text = text
	return s.err
}

func TestTelegramControllerSend(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name      string
		body      string
		senderErr error
		wantCode  MyCode
		wantText  string
	}{
		{name: "success", body: `{"text":"hello"}`, wantCode: CodeSuccess, wantText: "hello"},
		{name: "missing text", body: `{}`, wantCode: CodeInvalidParam},
		{name: "blank text", body: `{"text":"   "}`, wantCode: CodeInvalidParam},
		{name: "sender error", body: `{"text":"hello"}`, senderErr: errors.New("telegram unavailable"), wantCode: CodeServerBusy, wantText: "hello"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sender := &telegramSenderStub{err: tt.senderErr}
			telegramService := service.NewTelegramService(sender)
			telegramController := NewTelegramController(telegramService)

			router := gin.New()
			router.POST("/telegram", telegramController.Send)

			req := httptest.NewRequest(http.MethodPost, "/telegram", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()
			router.ServeHTTP(response, req)

			if response.Code != http.StatusOK {
				t.Fatalf("HTTP status = %d, want %d", response.Code, http.StatusOK)
			}

			var body ResponseData
			if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
				t.Fatalf("decode response: %v", err)
			}
			if body.Code != tt.wantCode {
				t.Fatalf("response code = %d, want %d", body.Code, tt.wantCode)
			}
			if sender.text != tt.wantText {
				t.Fatalf("sent text = %q, want %q", sender.text, tt.wantText)
			}
		})
	}
}
