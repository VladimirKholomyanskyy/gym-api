package service

import (
	"github.com/VladimirKholomyanskyy/gym-api/internal/models"
	"github.com/VladimirKholomyanskyy/gym-api/internal/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}
func (s *UserService) CreateUser(externalUserID string) (*models.User, error) {
	user := &models.User{
		ExternalID: externalUserID,
	}
	err := s.repo.CreateUser(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) FindUserByID(id uint) (*models.User, error) {
	return s.repo.FindByID(id)
}

func (s *UserService) FindUserByExternalID(id string) (*models.User, error) {
	return s.repo.FindByExternalID(id)
}

func (s *UserService) UpdateUser(user *models.User) error {
	return s.repo.UpdateUser(user)
}

func (s *UserService) DeleteUser(id uint) error {
	return s.repo.DeleteUser(id)
}
