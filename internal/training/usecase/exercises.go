package usecase

import (
	"context"

	"github.com/VladimirKholomyanskyy/gym-api/internal/training/model"
	"github.com/VladimirKholomyanskyy/gym-api/internal/training/repository"
)

type ExerciseUseCase interface {
	List(ctx context.Context, page, pageSize int) ([]model.Exercise, int64, error)
	GetByPrimaryMuscle(ctx context.Context, primaryMuscle string, page, pageSize int) ([]model.Exercise, int64, error)
	GetByID(ctx context.Context, id string) (*model.Exercise, error)
}

type exerciseUseCase struct {
	repo repository.ExerciseRepository
}

// NewExerciseUseCase creates a new instance of ExerciseUseCase
func NewExerciseUseCase(repo repository.ExerciseRepository) ExerciseUseCase {
	return &exerciseUseCase{
		repo: repo,
	}
}

func (s *exerciseUseCase) GetByPrimaryMuscle(ctx context.Context, primaryMuscle string, page, pageSize int) ([]model.Exercise, int64, error) {
	return s.repo.FindByPrimaryMuscle(ctx, primaryMuscle, page, pageSize)
}

// GetAllExercises retrieves all exercises.
func (s *exerciseUseCase) List(ctx context.Context, page, pageSize int) ([]model.Exercise, int64, error) {
	return s.repo.FindAll(ctx, page, pageSize)
}

// GetExercise retrieves a single exercise by ID.
func (s *exerciseUseCase) GetByID(ctx context.Context, id string) (*model.Exercise, error) {
	return s.repo.FindByID(ctx, id)
}
