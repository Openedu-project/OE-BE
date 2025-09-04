package configs

import (
	"fmt"
	"log"
	"strings"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Port string `env:"PORT,required"`
	// DBHost string `env:"DB_HOST,required"`
	// DBUser string `env:"DB_USER,required"`
	// DBPass string `env:"DB_PASS,required"`
	// DBName string `env:"DB_NAME,required"`
}

var Env Config

func InitEnv() {
	_ = godotenv.Load()
	if err := env.Parse(&Env); err != nil {
		printEnvError(err)
		log.Fatal("❌ Environment validation failed!")
	}
}

func printEnvError(err error) {
	for _, line := range strings.Split(err.Error(), ";") {
		if strings.TrimSpace(line) != "" {
			fmt.Println(strings.TrimSpace(line))
		}
	}

}
