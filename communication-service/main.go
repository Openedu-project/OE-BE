package main

import (
	"log"
	"os"
	"os/signal"

	"communication-service/configs"
	"communication-service/handlers"
)

func init() {
	configs.InitEnv()
}

func main() {
	log.Println("✅ Communication service started...")

	// Listener RabbitMQ (sẽ implement ở issue #47)
	go handlers.ListenRegisterUserSuccess()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutting down communication service...")
}
