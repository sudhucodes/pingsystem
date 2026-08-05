package config

import (
	"testing"
)

func TestConfigValidation(t *testing.T) {
	cfg := &Config{
		BotToken:    "",
		ChatID:      "",
		DeviceAlias: "Test Device",
	}

	if err := cfg.Validate(); err == nil {
		t.Error("expected error for empty botToken and chatId, got nil")
	}

	cfg.BotToken = "123456:ABC"
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for empty chatId, got nil")
	}

	cfg.ChatID = "987654321"
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected valid config, got error: %v", err)
	}
}
