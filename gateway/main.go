package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	routes "gateway/api"
	"gateway/configs"
)

func init() {
	configs.InitEnv()

	log.Println("Connecting to RabbitMQ at:", configs.Env.RabbitMQURL)
	configs.ConnectDatabase()
	configs.ConnectRabbitMQ()
}

func main() {
	routersInit := routes.InitRouter()
	port := configs.Env.Port

	endPoint := fmt.Sprintf("%s:%s", "0.0.0.0", port)
	maxHeaderBytes := 1 << 20

	server := &http.Server{
		Addr:           endPoint,
		Handler:        routersInit,
		MaxHeaderBytes: maxHeaderBytes,
	}

	go func() {
		log.Printf("✅ Start http server listening at: %s\n", endPoint)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Printf("Shutting down server...")
}
