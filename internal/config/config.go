package config

import (
	"os"
)

// AuthMode определяет способ авторизации amoCRM
type AuthMode string

const (
	AuthModeToken AuthMode = "token"
	AuthModeOAuth AuthMode = "oauth"
)

// Config holds all application configuration
type Config struct {
	TelegramToken string
	Debug         bool

	// Ollama settings
	OllamaURL   string
	OllamaModel string

	// amoCRM
	AmoCRMAuthMode     AuthMode
	AmoCRMBaseURL      string
	AmoCRMToken        string
	AmoCRMClientID     string
	AmoCRMClientSecret string
	AmoCRMRedirectURI  string
}

// Load loads configuration from environment variables
func Load() *Config {
	authMode := AuthMode(getEnvOrDefault("AMOCRM_AUTH_MODE", "token"))
	if authMode != AuthModeToken && authMode != AuthModeOAuth {
		authMode = AuthModeToken
	}

	return &Config{
		TelegramToken:      os.Getenv("TELEGRAM_BOT_TOKEN"),
		Debug:              os.Getenv("DEBUG") == "true" || os.Getenv("DEBUG") == "1",
		OllamaURL:          getEnvOrDefault("OLLAMA_URL", "http://localhost:11434"),
		OllamaModel:        getEnvOrDefault("OLLAMA_MODEL", "gpt-oss:120b-cloud"),
		AmoCRMAuthMode:     authMode,
		AmoCRMBaseURL:      os.Getenv("AMOCRM_BASE_URL"),
		AmoCRMToken:        os.Getenv("AMOCRM_ACCESS_TOKEN"),
		AmoCRMClientID:     os.Getenv("AMOCRM_CLIENT_ID"),
		AmoCRMClientSecret: os.Getenv("AMOCRM_CLIENT_SECRET"),
		AmoCRMRedirectURI:  os.Getenv("AMOCRM_REDIRECT_URI"),
	}
}

// IsAmoCRMConfigured checks if amoCRM credentials are set
func (c *Config) IsAmoCRMConfigured() bool {
	if c.AmoCRMBaseURL == "" {
		return false
	}

	if c.AmoCRMAuthMode == AuthModeOAuth {
		return c.AmoCRMClientID != "" && c.AmoCRMClientSecret != "" && c.AmoCRMRedirectURI != ""
	}

	// token mode
	return c.AmoCRMToken != ""
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
