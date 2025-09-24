package handler

import (
	"encoding/json"
	"log"

	"communication-service/services"

	"github.com/streadway/amqp"
)

type RabbitHandler struct{}

func NewRabbitHandler() *RabbitHandler {
	return &RabbitHandler{}
}

func (h *RabbitHandler) HandlerMessage(d amqp.Delivery) {
	switch d.RoutingKey {
	case "user.registered":
		h.handleUserRegistered(d)
	default:
		log.Printf("Unknown routing key: %s", d.RoutingKey)
	}
}

type UserRegisteredPayload struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (h *RabbitHandler) handleUserRegistered(d amqp.Delivery) {
	log.Printf("Received a message for user.registered: %s", d.Body)

	var payload UserRegisteredPayload
	err := json.Unmarshal(d.Body, &payload)
	if err != nil {
		log.Printf("Error devoding JSON: %s", err)
		return
	}

	err = services.SendWelcomeEmail(payload.Email, payload.Name)
	if err != nil {
		log.Printf("Error sending welcome email: %s", err)
		return
	}
	log.Printf("Welcome email sent to %s", payload.Email)
}
