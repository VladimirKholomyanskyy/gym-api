package usecase

import (
	"context"
	"fmt"

	"github.com/VladimirKholomyanskyy/gym-api/internal/common"
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
	// Add other dependencies like validators, external services, etc.
}

// NewExerciseUseCase creates a new instance of ExerciseUseCase
func NewExerciseUseCase(repo repository.ExerciseRepository) ExerciseUseCase {
	return &exerciseUseCase{
		repo: repo,
	}
}

func (s *exerciseUseCase) GetByPrimaryMuscle(ctx context.Context, primaryMuscle string, page, pageSize int) ([]model.Exercise, int64, error) {
	exercises, totalCount, err := s.repo.FindByPrimaryMuscle(ctx, primaryMuscle, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve exercises: %w", err)
	}
	return exercises, totalCount, nil
}

// GetAllExercises retrieves all exercises.
func (s *exerciseUseCase) List(ctx context.Context, page, pageSize int) ([]model.Exercise, int64, error) {
	exercises, totalCount, err := s.repo.FindAll(ctx, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve exercises: %w", err)
	}
	return exercises, totalCount, nil
}

// GetExercise retrieves a single exercise by ID.
func (s *exerciseUseCase) GetByID(ctx context.Context, id string) (*model.Exercise, error) {
	exercise, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve exercise: %w", err)
	}
	if exercise == nil {
		return nil, common.NewNotFoundError("exercise not found")
	}
	return exercise, nil
}
