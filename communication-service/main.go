package main

import (
	"log"

	"communication-service/config"
)

func main() {
	config.LoadConfig()
	log.Printf("Communication Service started in %s environment on port %s", config.Cfg.AppEnv, config.Cfg.Port)

	select {}
}
