package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DSN               string
	Debug             string
	JWTSecret         []byte
	MPToken           string
	WebhookSecret     string
	NotificationURL   string
	OAuthRedirectURL  string
	OAuthIDClient     string
	OAuthSecretClient string
}

func LoadConfig() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("No se encontro el archivo .env")
	}
	if getEnv("DEBUG", "dev") == "dev" {
		return &Config{
			DSN:               getEnv("DSN_Cloud", "..."),
			Debug:             getEnv("DEBUG", "Dev"),
			JWTSecret:         []byte(getEnv("JWT_SECRET", "mysecretkey")),
			MPToken:           getEnv("TOKEN_MP_TEST", "..."),
			NotificationURL:   getEnv("NOTIFICATION_URL", "..."),
			WebhookSecret:     getEnv("MP_WEBHOOK_SECRET", "JLE02020"),
			OAuthRedirectURL:  getEnv("OAUTH_REDIRECT_URL", ".."),
			OAuthIDClient:     getEnv("OAUTH_ID_CLIENT", "..."),
			OAuthSecretClient: getEnv("OAUTH_SECRET_CLIENT", "..."),
		}
	}
	return &Config{
		DSN:               getEnv("DSN_Cloud", "..."),
		Debug:             getEnv("DEBUG", "Prod"),
		JWTSecret:         []byte(getEnv("JWT_SECRET", "mysecretkey")),
		MPToken:           getEnv("TOKEN_MP_PROD", "..."),
		NotificationURL:   getEnv("NOTIFICATION_URL", "..."),
		WebhookSecret:     getEnv("MP_WEBHOOK_SECRET", "JLE02020"),
		OAuthRedirectURL:  getEnv("OAUTH_REDIRECT_URL", ".."),
		OAuthIDClient:     getEnv("OAUTH_ID_CLIENT", "..."),
		OAuthSecretClient: getEnv("OAUTH_SECRET_CLIENT", "..."),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
