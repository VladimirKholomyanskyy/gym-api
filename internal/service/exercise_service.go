package service

import (
	"github.com/VladimirKholomyanskyy/gym-api/internal/models"
	"github.com/VladimirKholomyanskyy/gym-api/internal/repository"
)

type ExerciseService struct {
	Repo *repository.ExerciseRepository
}

func (s *ExerciseService) GetExercisesByPrimaryMuscle(primaryMuscle string) ([]models.Exercise, error) {
	return s.Repo.GetExercisesByPrimaryMuscle(primaryMuscle)
}

func (s *ExerciseService) GetAllExercises() ([]models.Exercise, error) {
	return s.Repo.GetAllExercises()
}
