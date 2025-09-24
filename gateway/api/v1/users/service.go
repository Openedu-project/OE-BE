package users

import (
	"encoding/json"
	"fmt"
	"log"

	"gateway/configs"
	"gateway/models"
	"gateway/utils"

	"github.com/streadway/amqp"
)

type UserService struct {
	repo *UserRepository
}

func NewUserService(repo *UserRepository) *UserService {
	return &UserService{repo: repo}
}

type UserRegisteredPayload struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (s *UserService) publishUserRegister(email, name string) {
	payload := UserRegisteredPayload{
		Email: email,
		Name:  name,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling user registered payload: %v", err)
		return
	}

	err = configs.RabbitChannel.Publish(
		"user_events",
		"user.registered",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Printf("Failed to publish user.registered event: %v", err)
	} else {
		log.Printf("Published user.registered event for user %s", email)
	}
}

func (s *UserService) CreateUser(dto CreateUserDTO) (*models.User, error) {
	hashedPassword, err := utils.HashPassword(dto.Password)
	if err != nil {
		return nil, err
	}
	user := models.User{
		Name:     dto.Name,
		Email:    dto.Email,
		Password: hashedPassword,
		Role:     "learner",
	}
	if err := s.repo.Create(&user); err != nil {
		return nil, err
	}
	go s.publishUserRegister(user.Email, user.Name)
	return &user, nil
}

func (s *UserService) GetUsers() ([]models.User, error) {
	return s.repo.FindAll()
}

func (s *UserService) ValidateUser(email, password string) (*models.User, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return nil, fmt.Errorf("invalid credentials")
	}

	return user, nil
}

func (s *UserService) UpdateRole(id uint, role string) (*models.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	user.Role = role
	if err := s.repo.db.Save(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}
