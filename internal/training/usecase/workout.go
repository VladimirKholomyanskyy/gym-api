package usecase

import (
	"context"
	"strings"

	openapi "github.com/VladimirKholomyanskyy/gym-api/internal/api/go"
	"github.com/VladimirKholomyanskyy/gym-api/internal/auth"
	customerrors "github.com/VladimirKholomyanskyy/gym-api/internal/customErrors"
	"github.com/VladimirKholomyanskyy/gym-api/internal/training/model"
	"github.com/VladimirKholomyanskyy/gym-api/internal/training/repository"
	"github.com/VladimirKholomyanskyy/gym-api/internal/utils"
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
	if err := s.authorization.CanModifyTrainingProgram(ctx, profileID, programID); err != nil {
		return nil, err
	}
	workout := &model.Workout{
		Name:              strings.TrimSpace(input.Name),
		TrainingProgramID: programID,
	}
	if err := s.repo.Create(ctx, workout); err != nil {
		return nil, err
	}
	return workout, nil
}

// GetWorkout retrieves a workout. Ownership check is recommended; consider using a FindByIDAndProgramID method.
func (s *workoutUseCase) GetByProgramIDAndWorkoutID(ctx context.Context, profileID, programID, workoutID string) (*model.Workout, error) {
	if err := s.authorization.CanModifyTrainingProgram(ctx, profileID, programID); err != nil {
		return nil, err
	}
	workout, err := s.repo.GetByID(ctx, workoutID)
	if err != nil {
		return nil, err
	}
	if workout.TrainingProgramID != programID {
		return nil, customerrors.ErrAccessForbidden
	}
	return workout, nil
}
func (s *workoutUseCase) GetByWorkoutID(ctx context.Context, profileID, workoutID string) (*model.Workout, error) {
	workout, err := s.repo.GetByID(ctx, workoutID)
	if err != nil {
		return nil, err
	}
	if err := s.authorization.CanModifyTrainingProgram(ctx, profileID, workout.TrainingProgramID); err != nil {
		return nil, err
	}
	return workout, nil
}

// GetAllWorkouts retrieves all workouts for a given training program.
func (s *workoutUseCase) List(ctx context.Context, profileID, programID string, page, pageSize int) ([]model.Workout, int64, error) {
	err := s.authorization.CanModifyTrainingProgram(ctx, profileID, programID)
	if err != nil {
		return nil, 0, err
	}
	return s.repo.GetAllByTrainingProgramID(ctx, programID, page, pageSize)
}

// UpdateWorkout updates a workout's details, ensuring ownership.
func (s *workoutUseCase) Update(ctx context.Context, profileID, programID, workoutID string, input openapi.PatchWorkoutRequest) (*model.Workout, error) {
	workout, err := s.GetByProgramIDAndWorkoutID(ctx, profileID, programID, workoutID)
	if err != nil {
		return nil, err
	}
	updates := make(map[string]any)
	if utils.HasText(input.Name) {
		updates["name"] = strings.TrimSpace(*input.Name)
	}
	if len(updates) == 0 {
		return workout, nil // No changes to apply
	}
	return s.repo.UpdatePartial(ctx, workout.ID, updates)
}

// DeleteWorkout deletes a workout if it belongs to the user's training program.
func (s *workoutUseCase) Delete(ctx context.Context, profileID, programID, workoutID string) error {
	workout, err := s.GetByProgramIDAndWorkoutID(ctx, profileID, programID, workoutID)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, workout.ID)
}

// Reorder implements WorkoutUseCase.
func (s *workoutUseCase) Reorder(ctx context.Context, profileID, programID, workoutID string, request openapi.ReorderWorkoutRequest) error {
	_, err := s.GetByProgramIDAndWorkoutID(ctx, profileID, programID, workoutID)
	if err != nil {
		return err
	}
	return s.repo.Reorder(ctx, workoutID, int(request.Position))
}
