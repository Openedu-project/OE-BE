package configs

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Environment struct {
	AppEnv string `env:"APP_ENV,required"`
	Port   string `env:"PORT,required"`
	DBHost string `env:"DB_HOST,required"`
	DBUser string `env:"DB_USER,required"`
	DBPass string `env:"DB_PASS,required"`
	DBName string `env:"DB_NAME,required"`
	DBPort string `env:"DB_PORT,required"`
	// JwtSecretAccess  string `env:"JWT_SECRET_ACCESS,required"`
	// JwtSecretRefresh string `env:"JWT_SECRET_REFRESH,required"`
	// JwtExpiredTime   string `env:"JWT_EXPIRED_TIME,required"`
	// AESSecret        string `env:"AES_SECRET,required"`
	SMTPHost    string
	SMTPPort    string
	SMTPUser    string
	SMTPPass    string
	RabbitMQURL string
}

var Env Environment

func InitEnv() {
	_ = godotenv.Load()
	if err := env.Parse(&Env); err != nil {
		printEnvError(err)
		log.Fatal("❌ Environment validation failed!")
	}

	Env.AppEnv = os.Getenv("APP_ENV")
	Env.Port = os.Getenv("PORT")
	Env.SMTPHost = os.Getenv("SMTP_HOST")
	Env.SMTPPort = os.Getenv("SMTP_PORT")
	Env.SMTPUser = os.Getenv("SMTP_USER")
	Env.SMTPPass = os.Getenv("SMTP_PASS")
	Env.RabbitMQURL = os.Getenv("RABBITMQ_URL")
}

func IsProduction() bool {
	isProd := Env.AppEnv == "production"
	fmt.Println("isProd: ", isProd)
	return isProd
}

func printEnvError(err error) {
	for _, line := range strings.Split(err.Error(), ";") {
		if strings.TrimSpace(line) != "" {
			fmt.Println(strings.TrimSpace(line))
		}
	}
}
