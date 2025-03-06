package usecase

import (
	"context"
	"fmt"
	"strings"

	openapi "github.com/VladimirKholomyanskyy/gym-api/internal/api/go"
	"github.com/VladimirKholomyanskyy/gym-api/internal/auth"
	"github.com/VladimirKholomyanskyy/gym-api/internal/common"
	"github.com/VladimirKholomyanskyy/gym-api/internal/training/model"
	"github.com/VladimirKholomyanskyy/gym-api/internal/training/repository"
)

type WorkoutUseCase interface {
	Create(ctx context.Context, profileID, programID string, input openapi.CreateWorkoutRequest) (*model.Workout, error)
	GetByProgramIDAndWorkoutID(ctx context.Context, profileID, programID, workoutID string) (*model.Workout, error)
	GetByWorkoutID(ctx context.Context, profileID, workoutID string) (*model.Workout, error)
	List(ctx context.Context, profileID, programID string, page, pageSize int) ([]model.Workout, int64, error)
	Update(ctx context.Context, profileID, programID, workoutID string, input openapi.PatchWorkoutRequest) (*model.Workout, error)
	Delete(ctx context.Context, profileID, programID, workoutID string) error
	Reorder(ctx context.Context, profileID, programID, workoutID string, request openapi.ReorderWorkoutRequest) error
}

type workoutUseCase struct {
	repo          repository.WorkoutRepository
	authorization *auth.Authorization
}

// NewExerciseUseCase creates a new instance of ExerciseUseCase
func NewWorkoutUseCase(repo repository.WorkoutRepository, authorization *auth.Authorization) WorkoutUseCase {
	return &workoutUseCase{
		repo:          repo,
		authorization: authorization,
	}
}

func (s *workoutUseCase) Create(ctx context.Context, profileID, programID string, input openapi.CreateWorkoutRequest) (*model.Workout, error) {
	if err := common.ValidateUUIDs(profileID, programID); err != nil {
		return nil, err
	}

	if err := s.authorization.CanModifyTrainingProgram(ctx, profileID, programID); err != nil {
		return nil, err
	}
	if !common.HasText(&input.Name) {
		return nil, common.NewValidationError("workout name cannot be empty")
	}

	workout := &model.Workout{
		Name:              strings.TrimSpace(input.Name),
		TrainingProgramID: programID,
	}
	if err := s.repo.Create(ctx, workout); err != nil {
		return nil, fmt.Errorf("failed to create workout: %w", err)
	}
	return workout, nil
}

// GetWorkout retrieves a workout. Ownership check is recommended; consider using a FindByIDAndProgramID method.
func (s *workoutUseCase) GetByProgramIDAndWorkoutID(ctx context.Context, profileID, programID, workoutID string) (*model.Workout, error) {
	if err := common.ValidateUUIDs(profileID, programID, workoutID); err != nil {
		return nil, err
	}
	if err := s.authorization.CanModifyTrainingProgram(ctx, profileID, programID); err != nil {
		return nil, err
	}

	workout, err := s.repo.FindByID(ctx, workoutID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve workout: %w", err)
	}
	if workout == nil {
		return nil, common.NewNotFoundError("workout not found")
	}
	if workout.TrainingProgramID != programID {
		return nil, common.NewForbiddenError("you do not have permission to modify this workout")
	}
	return workout, nil
}
func (s *workoutUseCase) GetByWorkoutID(ctx context.Context, profileID, workoutID string) (*model.Workout, error) {

	workout, err := s.repo.FindByID(ctx, workoutID)
	if err := s.authorization.CanModifyTrainingProgram(ctx, profileID, workout.TrainingProgramID); err != nil {
		return nil, err
	}
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve workout: %w", err)
	}
	if workout == nil {
		return nil, common.NewNotFoundError("workout not found")
	}

	return workout, nil
}

// GetAllWorkouts retrieves all workouts for a given training program.
func (s *workoutUseCase) List(ctx context.Context, profileID, programID string, page, pageSize int) ([]model.Workout, int64, error) {
	err := s.authorization.CanModifyTrainingProgram(ctx, profileID, programID)
	if err != nil {
		return nil, 0, err
	}

	workouts, totalCount, err := s.repo.FindByTrainingProgramID(ctx, programID, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve workouts: %w", err)
	}
	return workouts, totalCount, nil
}

// UpdateWorkout updates a workout's details, ensuring ownership.
func (s *workoutUseCase) Update(ctx context.Context, profileID, programID, workoutID string, input openapi.PatchWorkoutRequest) (*model.Workout, error) {
	if err := common.ValidateUUIDs(profileID, programID, workoutID); err != nil {
		return nil, err
	}
	workout, err := s.GetByProgramIDAndWorkoutID(ctx, profileID, programID, workoutID)
	if err != nil {
		return nil, err
	}

	updated := false

	if common.HasText(input.Name) {
		workout.Name = strings.TrimSpace(*input.Name)
		updated = true
	}

	if !updated {
		return workout, nil // No changes to apply
	}

	if err = s.repo.Update(ctx, workout); err != nil {
		return nil, fmt.Errorf("failed to update workout: %w", err)
	}

	return workout, nil
}

// DeleteWorkout deletes a workout if it belongs to the user's training program.
func (s *workoutUseCase) Delete(ctx context.Context, profileID, programID, workoutID string) error {
	if err := common.ValidateUUIDs(profileID, programID, workoutID); err != nil {
		return err
	}
	workout, err := s.GetByProgramIDAndWorkoutID(ctx, profileID, programID, workoutID)
	if err != nil {
		return err
	}

	if err = s.repo.Delete(ctx, workout.ID); err != nil {
		return fmt.Errorf("failed to delete workout: %w", err)
	}

	return nil
}

// Reorder implements WorkoutUseCase.
func (s *workoutUseCase) Reorder(ctx context.Context, profileID, programID, workoutID string, request openapi.ReorderWorkoutRequest) error {
	if err := common.ValidateUUIDs(profileID, programID, workoutID); err != nil {
		return err
	}
	_, err := s.GetByProgramIDAndWorkoutID(ctx, profileID, programID, workoutID)
	if err != nil {
		return err
	}
	if err = s.repo.Reorder(ctx, workoutID, int(request.Position)); err != nil {
		return fmt.Errorf("failed to reorder workout: %w", err)
	}

	return nil
}
