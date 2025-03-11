package usecase

import (
	"context"

	openapi "github.com/VladimirKholomyanskyy/gym-api/internal/api/go"
	"github.com/VladimirKholomyanskyy/gym-api/internal/auth"
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
	if err := uc.authorization.CanModifyWorkout(ctx, profileID, input.WorkoutId); err != nil {
		return nil, err
	}

	workoutExercise := &model.WorkoutExercise{
		WorkoutID:  input.WorkoutId,
		ExerciseID: input.ExerciseId,
		Sets:       int(input.Sets),
		Reps:       int(input.Reps),
	}
	if err := uc.repo.Create(ctx, workoutExercise); err != nil {
		return nil, err
	}
	return workoutExercise, nil
}

// GetAllWorkoutExercisesByWorkout retrieves workout exercises, verifying that the workout belongs to the user's program.
func (uc *workoutExerciseUseCase) List(ctx context.Context, profileID, workoutID string, page, pageSize int) ([]model.WorkoutExercise, int64, error) {
	if err := uc.authorization.CanModifyWorkout(ctx, profileID, workoutID); err != nil {
		return nil, 0, err
	}
	return uc.repo.GetAllByWorkoutID(ctx, workoutID, page, pageSize)
}

func (uc *workoutExerciseUseCase) GetByID(ctx context.Context, profileID, workoutExerciseID string) (*model.WorkoutExercise, error) {
	workoutExercise, err := uc.repo.GetByID(ctx, workoutExerciseID)
	if err != nil {
		return nil, err
	}
	if err := uc.authorization.CanModifyWorkout(ctx, profileID, workoutExercise.WorkoutID); err != nil {
		return nil, err
	}
	return workoutExercise, nil
}

// UpdateWorkoutExercise updates a workout exercise, ensuring it belongs to a workout under the user's program.
func (uc *workoutExerciseUseCase) Update(ctx context.Context, profileID, workoutExerciseID string, input openapi.PatchWorkoutExerciseRequest) (*model.WorkoutExercise, error) {
	workoutExercise, err := uc.GetByID(ctx, profileID, workoutExerciseID)
	if err != nil {
		return nil, err
	}
	updates := make(map[string]any)
	if input.Sets != nil && *input.Sets > 0 {
		updates["sets"] = int(*input.Sets)
	}
	if input.Reps != nil && *input.Reps > 0 {
		updates["reps"] = int(*input.Reps)
	}
	if len(updates) != 0 {
		return workoutExercise, nil
	}
	return uc.repo.UpdatePartial(ctx, workoutExercise.ID, updates)
}

// DeleteWorkoutExercise deletes a workout exercise if its parent workout belongs to the profile.
func (uc *workoutExerciseUseCase) Delete(ctx context.Context, profileID, workoutExerciseID string) error {
	workoutExercise, err := uc.GetByID(ctx, profileID, workoutExerciseID)
	if err != nil {
		return err
	}
	return uc.repo.Delete(ctx, workoutExercise.ID)
}

// Reorder implements WorkoutExerciseUseCase.
func (uc *workoutExerciseUseCase) Reorder(ctx context.Context, profileID, workoutExerciseID string, request openapi.ReorderWorkoutExerciseRequest) error {
	workoutExercise, err := uc.GetByID(ctx, profileID, workoutExerciseID)
	if err != nil {
		return err
	}
	return uc.repo.Reorder(ctx, workoutExercise.ID, int(request.Position))
}
