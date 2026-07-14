package config

import (
	"testing"

	"github.com/spf13/viper"
)

func TestTelegramConfigFromEnvironment(t *testing.T) {
	t.Setenv("BLUESPOT_TG_BOT_TOKEN", "test-token")
	t.Setenv("BLUESPOT_TG_CHAT_ID", "-100123456")

	loader := viper.New()
	bindEnvVars(loader)

	var cfg Config
	if err := loader.Unmarshal(&cfg); err != nil {
		t.Fatalf("unmarshal config: %v", err)
	}

	if cfg.TG.BotToken != "test-token" {
		t.Fatalf("BotToken = %q, want %q", cfg.TG.BotToken, "test-token")
	}
	if cfg.TG.ChatID != "-100123456" {
		t.Fatalf("ChatID = %q, want %q", cfg.TG.ChatID, "-100123456")
	}
}
