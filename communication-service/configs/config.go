package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EnvConfig struct {
	AppEnv      string
	Port        string
	SMTPHost    string
	SMTPPort    string
	SMTPUser    string
	SMTPPass    string
	RabbitMQURL string
}

var Env EnvConfig

func InitEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("⚠️ No .env file found, using environment variables")
	}

	Env.AppEnv = os.Getenv("APP_ENV")
	Env.Port = os.Getenv("PORT")
	Env.SMTPHost = os.Getenv("SMTP_HOST")
	Env.SMTPPort = os.Getenv("SMTP_PORT")
	Env.SMTPUser = os.Getenv("SMTP_USER")
	Env.SMTPPass = os.Getenv("SMTP_PASS")
	Env.RabbitMQURL = os.Getenv("RABBITMQ_URL")

	log.Println("✅ Environment loaded")
}
