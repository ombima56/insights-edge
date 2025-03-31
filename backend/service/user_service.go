package service

import (
	"github.com/ombima56/insights-edge/backend/models"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserByWallet(walletAddr string) (*models.User, error)
	GetUserByID(id int64) (*models.User, error)
}

type UserService struct {
	userRepository UserRepository
}

func NewUserService(userRepository UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (s *UserService) RegisterUser(user *models.User) error {
	return s.userRepository.CreateUser(user)
}

func (s *UserService) LoginUser(email string, password string) (string, error) {
	user, err := s.userRepository.GetUserByEmail(email)
	if err != nil {
		return "", models.ErrInvalidCredentials
	}

	// TODO: Implement password verification
	if user.Password != password {
		return "", models.ErrInvalidCredentials
	}

	// TODO: Implement JWT token generation
	return "token", nil
}

func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	return s.userRepository.GetUserByEmail(email)
}

func (s *UserService) GetUserByWallet(walletAddr string) (*models.User, error) {
	return s.userRepository.GetUserByWallet(walletAddr)
}

func (s *UserService) GetUserByID(id int64) (*models.User, error) {
	return s.userRepository.GetUserByID(id)
}
