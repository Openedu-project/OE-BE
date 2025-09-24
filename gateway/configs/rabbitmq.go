package configs

import (
	"log"

	"github.com/streadway/amqp"
)

var (
	RabbitConn    *amqp.Connection
	RabbitChannel *amqp.Channel
)

func ConnectRabbitMQ() {
	log.Println("Connecting to RabbitMQ at:", Env.RabbitMQURL)
	var err error
	RabbitConn, err = amqp.Dial(Env.RabbitMQURL)
	if err != nil {
		log.Fatalf("❌ Failed to connect RabbitMQ: %v", err)
	}

	RabbitChannel, err = RabbitConn.Channel()
	if err != nil {
		log.Fatalf("❌ Failed to open RabbitMQ channel: %v", err)
	}
	log.Println("✅ RabbitMQ connected")
}
