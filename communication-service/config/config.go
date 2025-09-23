package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv      string
	Port        string
	SMTPHost    string
	SMTPPort    int
	SMPTUser    string
	SMTPPass    string
	RabbitMQURL string
}

var Cfg *Config

func LoadConfig() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	Cfg = &Config{}

	Cfg.AppEnv = os.Getenv("APP_ENV")
	Cfg.Port = os.Getenv("PORT")
	Cfg.SMTPHost = os.Getenv("SMTP_HOST")
	smtpPort, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		log.Fatalf("Invalid SMTP_PORT: %v", err)
	}
	Cfg.SMPTUser = os.Getenv("SMTP_USER")
	Cfg.SMTPPass = os.Getenv("SMTP_PASS")

	Cfg.RabbitMQURL = os.Getenv("RABBITMQ_URL")

	Cfg.SMTPPort = smtpPort
}
