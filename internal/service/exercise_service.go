package service

import (
	"github.com/VladimirKholomyanskyy/gym-api/internal/models"
	"github.com/VladimirKholomyanskyy/gym-api/internal/repository"
)

type ExerciseService struct {
	repo *repository.ExerciseRepository
}

func NewExerciseService(repo *repository.ExerciseRepository) *ExerciseService {
	return &ExerciseService{repo: repo}
}

func (s *ExerciseService) GetExercisesByPrimaryMuscle(primaryMuscle string) ([]models.Exercise, error) {
	return s.repo.FindByPrimaryMuscle(primaryMuscle)
}

func (s *ExerciseService) GetAllExercises() ([]models.Exercise, error) {
	return s.repo.FindAll()
}
