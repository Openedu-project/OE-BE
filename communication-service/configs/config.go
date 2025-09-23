package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv      string
	Port        string
	SMTPHost    string
	SMTPPort    string
	SMTPUser    string
	SMTPPass    string
	RabbitMQURL string
}

var Env *Config

func InitEnv() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found, using system env")
	}

	Env = &Config{
		AppEnv:      os.Getenv("APP_ENV"),
		Port:        os.Getenv("PORT"),
		SMTPHost:    os.Getenv("SMTP_HOST"),
		SMTPPort:    os.Getenv("SMTP_PORT"),
		SMTPUser:    os.Getenv("SMTP_USER"),
		SMTPPass:    os.Getenv("SMTP_PASS"),
		RabbitMQURL: os.Getenv("RABBITMQ_URL"),
	}

	log.Printf("✅ Environment loaded: %s", Env.AppEnv)
}
