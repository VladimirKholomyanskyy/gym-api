package auth

import (
	"context"
	"fmt"

	customerrors "github.com/VladimirKholomyanskyy/gym-api/internal/customErrors"
	"github.com/VladimirKholomyanskyy/gym-api/internal/training/repository"
)

type Authorization struct {
	trainingProgramRepo repository.TrainingProgramRepository
	workoutRepo         repository.WorkoutRepository
}

// NewExerciseUseCase creates a new instance of ExerciseUseCase
func NewAuthorization(trainingProgramRepo repository.TrainingProgramRepository, workoutRepo repository.WorkoutRepository) *Authorization {
	return &Authorization{
		trainingProgramRepo: trainingProgramRepo,
		workoutRepo:         workoutRepo,
	}
}

func (a *Authorization) CanModifyTrainingProgram(ctx context.Context, profileId, programId string) error {
	program, err := a.trainingProgramRepo.FindByID(ctx, programId)
	if err != nil {
		return fmt.Errorf("failed to retrieve training program: %w", err)
	}
	if program == nil {
		return customerrors.ErrEntityNotFound
	}
	if program.ProfileID != profileId {
		return customerrors.ErrAccessForbidden
	}
	return nil
}

func (a *Authorization) CanModifyWorkout(ctx context.Context, profileId, workoutId string) error {
	workout, err := a.workoutRepo.GetByID(ctx, workoutId)
	if err != nil {
		return fmt.Errorf("failed to retrieve workout: %w", err)
	}
	if workout == nil {
		return customerrors.ErrAccessForbidden
	}
	return a.CanModifyTrainingProgram(ctx, profileId, workout.TrainingProgramID)
}
