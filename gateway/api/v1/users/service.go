package users

import (
	"fmt"
	"gateway/utils"
)

type UserService struct {
	repo *UserRepository
}

func NewUserService(repo *UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(dto CreateUserDTO) (*User, error) {
	hashedPassword, err := utils.HashPassword(dto.Password)
	if err != nil {
		return nil, err
	}
	user := User{
		Name:     dto.Name,
		Email:    dto.Email,
		Password: hashedPassword,
	}
	if err := s.repo.Create(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) GetUsers() ([]User, error) {
	return s.repo.FindAll()
}

func (s *UserService) ValidateUser(email, password string) (*User, error) {
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
