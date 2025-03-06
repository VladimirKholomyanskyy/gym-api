package usecase

import (
	"context"
	"fmt"

	openapi "github.com/VladimirKholomyanskyy/gym-api/internal/api/go"
	"github.com/VladimirKholomyanskyy/gym-api/internal/auth"
	"github.com/VladimirKholomyanskyy/gym-api/internal/common"
	"github.com/VladimirKholomyanskyy/gym-api/internal/training/model"
	"github.com/VladimirKholomyanskyy/gym-api/internal/training/repository"
)

type WorkoutExerciseUseCase interface {
	Create(ctx context.Context, profileId string, input openapi.CreateWorkoutExerciseRequest) (*model.WorkoutExercise, error)
	GetByID(ctx context.Context, profileId, workoutExerciseId string) (*model.WorkoutExercise, error)
	List(ctx context.Context, profileId, workoutId string, page, pageSize int) ([]model.WorkoutExercise, int64, error)
	Update(ctx context.Context, profileId, workoutExerciseId string, input openapi.PatchWorkoutExerciseRequest) (*model.WorkoutExercise, error)
	Delete(ctx context.Context, profileId, workoutExerciseId string) error
	Reorder(ctx context.Context, profileID, workoutExerciseId string, request openapi.ReorderWorkoutExerciseRequest) error
}

type workoutExerciseUseCase struct {
	repo          repository.WorkoutExerciseRepository
	authorization *auth.Authorization
}

// NewExerciseUseCase creates a new instance of ExerciseUseCase
func NewWorkoutExerciseUseCase(repo repository.WorkoutExerciseRepository, authorization *auth.Authorization) WorkoutExerciseUseCase {
	return &workoutExerciseUseCase{
		repo:          repo,
		authorization: authorization,
	}
}

// AddExerciseToWorkout adds an exercise to a workout, ensuring ownership of the underlying training program.
func (uc *workoutExerciseUseCase) Create(ctx context.Context, profileID string, input openapi.CreateWorkoutExerciseRequest) (*model.WorkoutExercise, error) {
	if err := common.ValidateUUIDs(profileID, input.WorkoutId, input.ExerciseId); err != nil {
		return nil, err
	}
	if err := uc.authorization.CanModifyWorkout(ctx, profileID, input.WorkoutId); err != nil {
		return nil, err
	}
	if input.Sets <= 0 {
		return nil, common.NewValidationError("sets must be greater then zero")
	}
	if input.Reps <= 0 {
		return nil, common.NewValidationError("reps must be greater then zero")
	}
	workoutExercise := &model.WorkoutExercise{
		WorkoutID:  input.WorkoutId,
		ExerciseID: input.ExerciseId,
		Sets:       int(input.Sets),
		Reps:       int(input.Reps),
	}
	if err := uc.repo.Create(ctx, workoutExercise); err != nil {
		return nil, fmt.Errorf("failed to create workout exercise: %w", err)
	}
	return workoutExercise, nil
}

// GetAllWorkoutExercisesByWorkout retrieves workout exercises, verifying that the workout belongs to the user's program.
func (uc *workoutExerciseUseCase) List(ctx context.Context, profileID, workoutID string, page, pageSize int) ([]model.WorkoutExercise, int64, error) {
	if err := uc.authorization.CanModifyWorkout(ctx, profileID, workoutID); err != nil {
		return nil, 0, err
	}

	workoutExercises, totalCount, err := uc.repo.FindByWorkoutID(ctx, workoutID, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve workout exercises: %w", err)
	}
	return workoutExercises, totalCount, nil

}

func (uc *workoutExerciseUseCase) GetByID(ctx context.Context, profileID, workoutExerciseID string) (*model.WorkoutExercise, error) {
	if err := common.ValidateUUIDs(profileID, workoutExerciseID); err != nil {
		return nil, err
	}
	workoutExercise, err := uc.repo.FindByID(ctx, workoutExerciseID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve workout exercise: %w", err)
	}
	if err := uc.authorization.CanModifyWorkout(ctx, profileID, workoutExercise.WorkoutID); err != nil {
		return nil, err
	}

	return workoutExercise, nil
}

// UpdateWorkoutExercise updates a workout exercise, ensuring it belongs to a workout under the user's program.
func (uc *workoutExerciseUseCase) Update(ctx context.Context, profileID, workoutExerciseID string, input openapi.PatchWorkoutExerciseRequest) (*model.WorkoutExercise, error) {
	if err := common.ValidateUUIDs(profileID, workoutExerciseID); err != nil {
		return nil, err
	}
	workoutExercise, err := uc.GetByID(ctx, profileID, workoutExerciseID)
	if err != nil {
		return nil, err
	}
	updated := false
	if input.Sets != nil && *input.Sets > 0 {
		workoutExercise.Sets = int(*input.Sets)
		updated = true
	}
	if input.Reps != nil && *input.Reps > 0 {
		workoutExercise.Reps = int(*input.Reps)
		updated = true
	}
	if !updated {
		return workoutExercise, nil
	}
	if err = uc.repo.Update(ctx, workoutExercise); err != nil {
		return nil, fmt.Errorf("failed to update workout exercise: %w", err)
	}
	return workoutExercise, nil
}

// DeleteWorkoutExercise deletes a workout exercise if its parent workout belongs to the profile.
func (uc *workoutExerciseUseCase) Delete(ctx context.Context, profileID, workoutExerciseID string) error {
	if err := common.ValidateUUIDs(profileID, workoutExerciseID); err != nil {
		return err
	}
	workoutExercise, err := uc.GetByID(ctx, profileID, workoutExerciseID)
	if err != nil {
		return err
	}
	if err = uc.repo.Delete(ctx, workoutExercise.ID); err != nil {
		return fmt.Errorf("failed to delete workout exercise: %w", err)
	}
	return nil
}

// Reorder implements WorkoutExerciseUseCase.
func (uc *workoutExerciseUseCase) Reorder(ctx context.Context, profileID, workoutExerciseID string, request openapi.ReorderWorkoutExerciseRequest) error {
	if err := common.ValidateUUIDs(profileID, workoutExerciseID); err != nil {
		return err
	}
	workoutExercise, err := uc.GetByID(ctx, profileID, workoutExerciseID)
	if err != nil {
		return err
	}
	if err = uc.repo.Reorder(ctx, workoutExercise.ID, int(request.Position)); err != nil {
		return fmt.Errorf("failed to delete workout exercise: %w", err)
	}
	return nil
}
