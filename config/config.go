package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     	string
	DBPort     	string
	DBUser     	string
	DBPassword 	string
	DBName     	string
	DBTimezone 	string
	RedisHost  	string
	RedisPass  	string
	SMTPHost 	string
	SMTPPort 	string
	SenderEmail string
	SenderPass 	string
	SECRET_KEY 	string
	GinPort   	string
	UploadDir	string
}

var AppConfig *Config

func LoadConfig() {
	err := godotenv.Load()
    if err != nil {
        log.Println("Warning: .env file not found, using system env")
    }	

	AppConfig = &Config{
		DBHost: getEnv("POSTGRES_HOST", "localhost"),
		DBPort: getEnv("POSTGRES_PORT", "3306"),
		DBUser: getEnv("pOSTGRES_USER", "root"),
		DBPassword: getEnv("POSTGRES_PASSWORD", ""),
		DBName: getEnv("POSTGRES_DB", "test"),
		DBTimezone: getEnv("POSTGRES_TIMEZONE", "UTC"),
		RedisHost: getEnv("REDIS_ADDR", "localhost"),
		RedisPass: getEnv("REDIS_PASSWORD", ""),
		SMTPHost: getEnv("SMTP_HOST", "smtp.example.com"),
		SMTPPort: getEnv("SMTP_PORT", ""),
		SenderEmail: getEnv("SENDER_EMAIL", ""),
		SenderPass: getEnv("SENDER_PASS", ""),
		SECRET_KEY: getEnv("JWT_SECRET",""),
		GinPort: getEnv("GIN_PORT", "8080"),
		UploadDir: getEnv("UPLOAD_DIR",""),
	}
}

func getEnv(key string, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}