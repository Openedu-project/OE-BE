package services

import (
	"log"

	"communication-service/configs"

	"github.com/streadway/amqp"
)

type RabbitMQService struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMQService() *RabbitMQService {
	conn, err := amqp.Dial(configs.Env.RabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}
	return &RabbitMQService{
		conn:    conn,
		channel: ch,
	}
}

func (s *RabbitMQService) Consume(handler func(d amqp.Delivery)) {
	err := s.channel.ExchangeDeclare(
		"user_events", // name
		"topic",       // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // argument
	)
	if err != nil {
		log.Fatalf("Failed to declare an exchange: %s", err)
	}

	q, err := s.channel.QueueDeclare(
		"communication_queue", // name
		true,                  // durable
		false,                 // delete when unused
		false,                 // exclusive
		false,                 // no-wait
		nil,                   // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %s", err)
	}
	err = s.channel.QueueBind(
		q.Name,        // queue name
		"user.*",      // routing key
		"user_events", // exchange
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to bind a queue: %s", err)
	}
	msgs, err := s.channel.Consume(
		q.Name, // queue
		"",     // consume
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %s", err)
	}
	forever := make(chan bool)

	go func() {
		for d := range msgs {
			handler(d)
		}
	}()
	log.Print("[*] Waiting for messages on exchange 'user_events'. To exit press Ctrl+C")
	<-forever
}

func (s *RabbitMQService) Close() {
	s.channel.Close()
	s.conn.Close()
}
