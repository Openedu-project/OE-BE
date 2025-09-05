package users

type UserService struct {
	repo *UserRepository
}

func NewUserService(repo *UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(dto CreateUserDTO) (*User, error) {
	user := User{
		Name:  dto.Name,
		Email: dto.Email,
	}
	if err := s.repo.Create(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) GetUsers() ([]User, error) {
	return s.repo.FindAll()
}
