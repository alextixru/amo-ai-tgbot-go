package config

import (
	"fmt"
	"os"
)

// PromptMissingCredentials checks for missing amoCRM credentials in debug mode
// and shows what needs to be configured
func PromptMissingCredentials(cfg *Config) error {
	if cfg.IsAmoCRMConfigured() || !cfg.Debug {
		return nil
	}

	fmt.Println()
	fmt.Println("⚠️  DEBUG MODE: amoCRM credentials not configured")
	fmt.Println()

	if cfg.AmoCRMBaseURL == "" {
		fmt.Println("   Missing: AMOCRM_BASE_URL")
	}

	if cfg.AmoCRMAuthMode == AuthModeOAuth {
		if cfg.AmoCRMClientID == "" {
			fmt.Println("   Missing: AMOCRM_CLIENT_ID")
		}
		if cfg.AmoCRMClientSecret == "" {
			fmt.Println("   Missing: AMOCRM_CLIENT_SECRET")
		}
		if cfg.AmoCRMRedirectURI == "" {
			fmt.Println("   Missing: AMOCRM_REDIRECT_URI")
		}
	} else {
		if cfg.AmoCRMToken == "" {
			fmt.Println("   Missing: AMOCRM_ACCESS_TOKEN")
		}
	}

	fmt.Println()
	fmt.Println("📝 Please fill in .env file and restart.")
	fmt.Println("   See .env.example for reference.")
	fmt.Println()

	// Create .env with template if it doesn't exist
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		template := `# amoCRM Config (fill in and restart)
AMOCRM_AUTH_MODE=token
AMOCRM_BASE_URL=https://your-domain.amocrm.ru
AMOCRM_ACCESS_TOKEN=your_token_here
`
		os.WriteFile(".env", []byte(template), 0644)
		fmt.Println("   Created .env template for you!")
	}

	return nil
}
