package usecase

import (
	"context"
	"time"

	"github.com/VladimirKholomyanskyy/gym-api/internal/auth"
	customerrors "github.com/VladimirKholomyanskyy/gym-api/internal/customErrors"
	"github.com/VladimirKholomyanskyy/gym-api/internal/training/model"
	"github.com/VladimirKholomyanskyy/gym-api/internal/training/repository"
)

type ScheduledWorkoutUseCase interface {
	Create(ctx context.Context, input model.CreateScheduledWorkoutInput) (*model.ScheduledWorkout, error)
	GetByID(ctx context.Context, profileId, scheduledWorkoutId string) (*model.ScheduledWorkout, error)
	List(ctx context.Context, profileID string, startDate, endDate time.Time, page, pageSize int) ([]model.ScheduledWorkout, int64, error)
	Update(ctx context.Context, input model.UpdateScheduledWorkoutInput) (*model.ScheduledWorkout, error)
	Delete(ctx context.Context, profileId, scheduledWorkoutId string) error
	GetUpcommingScheduledWorkout(ctx context.Context, profileID string) (*model.ScheduledWorkout, error)
}

type scheduledWorkoutUseCase struct {
	repo          repository.ScheduledWorkoutRepository
	authorization *auth.Authorization
}

// NewExerciseUseCase creates a new instance of ExerciseUseCase
func NewScheduledWorkoutUseCase(repo repository.ScheduledWorkoutRepository, authorization *auth.Authorization) ScheduledWorkoutUseCase {
	return &scheduledWorkoutUseCase{
		repo:          repo,
		authorization: authorization,
	}
}

func (uc *scheduledWorkoutUseCase) Create(ctx context.Context, input model.CreateScheduledWorkoutInput) (*model.ScheduledWorkout, error) {
	if err := uc.authorization.CanModifyWorkout(ctx, input.ProfileID, input.WorkoutID); err != nil {
		return nil, customerrors.ErrAccessForbidden
	}
	scheduledWorkout := &model.ScheduledWorkout{
		WorkoutID: input.WorkoutID,
		ProfileID: input.ProfileID,
		Date:      input.Date,
		Notes:     input.Notes,
	}
	if err := uc.repo.Create(ctx, scheduledWorkout); err != nil {
		return nil, err
	}
	return scheduledWorkout, nil
}

// GetScheduledWorkouts retrieves a user's scheduled workouts in a date range.
func (uc *scheduledWorkoutUseCase) List(ctx context.Context, profileID string, startDate, endDate time.Time, page, pageSize int) ([]model.ScheduledWorkout, int64, error) {
	return uc.repo.GetAllByProfileIDAndRange(ctx, profileID, startDate, endDate, page, pageSize)
}

// GetScheduledWorkoutByID retrieves a scheduled workout ensuring it belongs to the profile.
func (uc *scheduledWorkoutUseCase) GetByID(ctx context.Context, profileID, scheduledWorkoutID string) (*model.ScheduledWorkout, error) {
	scheduledWorkout, err := uc.repo.GetByID(ctx, scheduledWorkoutID)
	if err != nil {
		return nil, err
	}
	if scheduledWorkout.ProfileID != profileID {
		return nil, customerrors.ErrAccessForbidden
	}
	return scheduledWorkout, nil
}

// UpdateScheduledWorkout updates a scheduled workout.
func (uc *scheduledWorkoutUseCase) Update(ctx context.Context, input model.UpdateScheduledWorkoutInput) (*model.ScheduledWorkout, error) {
	scheduledWorkout, err := uc.GetByID(ctx, input.ProfileID, input.ScheduledWorkoutID)
	if err != nil {
		return nil, err
	}
	updates := make(map[string]any)
	if input.Date != nil {
		updates["date"] = input.Date
	}
	if input.Notes != nil {
		updates["notes"] = *input.Notes
	}
	if len(updates) == 0 {
		return scheduledWorkout, nil
	}
	return uc.repo.UpdatePartial(ctx, scheduledWorkout.ID, updates)
}

// DeleteScheduledWorkout deletes a scheduled workout if it belongs to the user.
func (uc *scheduledWorkoutUseCase) Delete(ctx context.Context, profileID, scheduledWorkoutID string) error {
	_, err := uc.GetByID(ctx, profileID, scheduledWorkoutID)
	if err != nil {
		return err
	}
	return uc.repo.Delete(ctx, scheduledWorkoutID)
}

// GetUpcommingScheduledWorkout implements ScheduledWorkoutUseCase.
func (uc *scheduledWorkoutUseCase) GetUpcommingScheduledWorkout(ctx context.Context, profileID string) (*model.ScheduledWorkout, error) {
	return uc.repo.GetUpcomming(ctx, profileID)
}
