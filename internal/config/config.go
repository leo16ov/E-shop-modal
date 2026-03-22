package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBPort     string
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	Debug      string
	JWTSecret  string
}

func LoadConfig() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("No se encontro el archivo .env")
	}
	appConfig := &Config{
		DBPort:     getEnv("DB_PORT", ":5050"),
		DBHost:     getEnv("DB_HOST", "MySQL"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PWD", "12344321"),
		DBName:     getEnv("DB_NAME", "..."),
		Debug:      getEnv("DEBUG", "dev"),
		JWTSecret:  getEnv("JWT_SECRET", "mysecretkey"),
	}
	return appConfig
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
