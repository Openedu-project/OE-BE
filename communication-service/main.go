package main

import (
	"log"
	"os"
	"os/signal"

	"communication-service/configs"
)

func init() {
	// Load config, DB, RabbitMQ trước khi main chạy
	configs.InitEnv()
	configs.ConnectDatabase()
}

func main() {
	// Start listener nhận message từ RabbitMQ

	log.Println("✅ Communication service started...")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutting down communication service...")
}
