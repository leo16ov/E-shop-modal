package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DSN       string
	Debug     string
	JWTSecret string
	MPToken   string
}

func LoadConfig() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("No se encontro el archivo .env")
	}
	if getEnv("DEBUG", "dev") == "dev" {
		return &Config{
			DSN:       getEnv("DSN", "..."),
			Debug:     getEnv("DB_NAME", "Dev"),
			JWTSecret: getEnv("JWT_SECRET", "mysecretkey"),
			MPToken:   getEnv("TOKEN_MP_TEST", "..."),
		}
	}
	return &Config{
		Debug:     getEnv("DB_NAME", "Prod"),
		JWTSecret: getEnv("JWT_SECRET", "mysecretkey"),
		MPToken:   getEnv("TOKEN_MP_PROD", "..."),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
