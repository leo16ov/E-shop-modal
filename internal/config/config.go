package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DSN           string
	Debug         string
	JWTSecret     string
	MPToken       string
	WebhookSecret string
}

func LoadConfig() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("No se encontro el archivo .env")
	}
	if getEnv("DEBUG", "dev") == "dev" {
		return &Config{
			DSN:           getEnv("DSN_Cloud", "..."),
			Debug:         getEnv("DEBUG", "Dev"),
			JWTSecret:     getEnv("JWT_SECRET", "mysecretkey"),
			MPToken:       getEnv("TOKEN_MP_TEST", "..."),
			WebhookSecret: getEnv("MP_WEBHOOK_SECRET", "JLE02020"),
		}
	}
	return &Config{
		DSN:           getEnv("DSN_Cloud", "..."),
		Debug:         getEnv("DEBUG", "Prod"),
		JWTSecret:     getEnv("JWT_SECRET", "mysecretkey"),
		MPToken:       getEnv("TOKEN_MP_PROD", "..."),
		WebhookSecret: getEnv("MP_WEBHOOK_SECRET", "JLE02020"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
