package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// Config holds the application settings.
type Config struct {
	BotToken    string `json:"botToken"`
	ChatID      string `json:"chatId"`
	DeviceAlias string `json:"deviceAlias"`
}

// GetConfigDir returns the directory path for PingSystem configuration (%APPDATA%\PingSystem).
func GetConfigDir() (string, error) {
	appData := os.Getenv("APPDATA")
	if appData == "" {
		userConfigDir, err := os.UserConfigDir()
		if err != nil {
			return "", fmt.Errorf("failed to get user config dir: %w", err)
		}
		appData = userConfigDir
	}
	return filepath.Join(appData, "PingSystem"), nil
}

// GetConfigPath returns the full path to config.json.
func GetConfigPath() (string, error) {
	dir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.json"), nil
}

// Load loads the configuration file. If it does not exist, it creates a default configuration.
func Load() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Auto-create default config file
		defaultConfig := &Config{
			BotToken:    "",
			ChatID:      "",
			DeviceAlias: "My Laptop",
		}
		if err := Save(defaultConfig); err != nil {
			return nil, fmt.Errorf("failed to create default config: %w", err)
		}
		return defaultConfig, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config JSON: %w", err)
	}

	// Default deviceAlias if blank
	if cfg.DeviceAlias == "" {
		cfg.DeviceAlias = "My Laptop"
	}

	return &cfg, nil
}

// Save writes the configuration to disk.
func Save(cfg *Config) error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Validate checks if the minimal required config values are set.
func (c *Config) Validate() error {
	if c.BotToken == "" {
		return errors.New("botToken is empty in config.json")
	}
	if c.ChatID == "" {
		return errors.New("chatId is empty in config.json")
	}
	return nil
}
