package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	openapi "github.com/VladimirKholomyanskyy/gym-api/internal/api/go"
	"github.com/VladimirKholomyanskyy/gym-api/internal/auth"
	"github.com/VladimirKholomyanskyy/gym-api/internal/common"
	"github.com/VladimirKholomyanskyy/gym-api/internal/training/model"
	"github.com/VladimirKholomyanskyy/gym-api/internal/training/repository"
)

type ScheduledWorkoutUseCase interface {
	Create(ctx context.Context, profileId string, request openapi.CreateScheduledWorkoutRequest) (*model.ScheduledWorkout, error)
	GetByID(ctx context.Context, profileId, scheduledWorkoutId string) (*model.ScheduledWorkout, error)
	List(ctx context.Context, profileID string, startDate, endDate time.Time, page, pageSize int) ([]model.ScheduledWorkout, int64, error)
	Update(ctx context.Context, profileId, scheduledWorkoutId string, input openapi.PatchScheduledWorkoutRequest) (*model.ScheduledWorkout, error)
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

func (uc *scheduledWorkoutUseCase) Create(ctx context.Context, profileID string, input openapi.CreateScheduledWorkoutRequest) (*model.ScheduledWorkout, error) {
	if err := common.ValidateUUIDs(profileID, input.WorkoutId); err != nil {
		return nil, err
	}
	if err := uc.authorization.CanModifyWorkout(ctx, profileID, input.WorkoutId); err != nil {
		return nil, err
	}
	startDate, err := common.ParseTime(input.Date)
	if err != nil {
		return nil, common.NewValidationError("invalid date format, expected YYYY-MM-DD")
	}

	scheduledWorkout := &model.ScheduledWorkout{
		WorkoutID: input.WorkoutId,
		ProfileID: profileID,
		Date:      startDate,
		Notes:     common.TrimPointer(input.Notes),
	}
	if err := uc.repo.ScheduleWorkout(ctx, scheduledWorkout); err != nil {
		return nil, fmt.Errorf("failed to schedule workout: %w", err)
	}
	return scheduledWorkout, nil
}

// GetScheduledWorkouts retrieves a user's scheduled workouts in a date range.
func (uc *scheduledWorkoutUseCase) List(ctx context.Context, profileID string, startDate, endDate time.Time, page, pageSize int) ([]model.ScheduledWorkout, int64, error) {
	workouts, totalCount, err := uc.repo.GetUserScheduledWorkoutsByRange(ctx, profileID, startDate, endDate, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve scheduled workouts: %w", err)
	}
	return workouts, totalCount, nil
}

// GetScheduledWorkoutByID retrieves a scheduled workout ensuring it belongs to the profile.
func (uc *scheduledWorkoutUseCase) GetByID(ctx context.Context, profileID, scheduledWorkoutId string) (*model.ScheduledWorkout, error) {
	if err := common.ValidateUUIDs(profileID); err != nil {
		return nil, err
	}
	scheduledWorkout, err := uc.repo.GetScheduledWorkout(ctx, scheduledWorkoutId)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve scheduled workout: %w", err)
	}
	if scheduledWorkout == nil {
		return nil, common.NewNotFoundError("scheduled workout not found")
	}
	if scheduledWorkout.ProfileID != profileID {
		return nil, common.NewForbiddenError("you do not have permission to access this scheduled workout")
	}
	return scheduledWorkout, nil
}

// UpdateScheduledWorkout updates a scheduled workout.
func (uc *scheduledWorkoutUseCase) Update(ctx context.Context, profileId, scheduledWorkoutId string, input openapi.PatchScheduledWorkoutRequest) (*model.ScheduledWorkout, error) {
	scheduledWorkout, err := uc.GetByID(ctx, profileId, scheduledWorkoutId)
	if err != nil {
		return nil, err
	}
	var date time.Time
	updated := false
	if input.Date != nil {
		date, err = common.ParseTime(*input.Date)
		if err != nil {
			return nil, common.NewValidationError("invalid date format, expected YYYY-MM-DD")
		}
		scheduledWorkout.Date = date
		updated = true
	}
	if common.HasText(input.Notes) {
		scheduledWorkout.Notes = strings.TrimSpace(*input.Notes)
		updated = true
	}
	if !updated {
		return scheduledWorkout, nil
	}
	if err := uc.repo.UpdateScheduledWorkout(ctx, scheduledWorkout); err != nil {
		return nil, fmt.Errorf("failed to update scheduled workout: %w", err)
	}
	return scheduledWorkout, nil
}

// DeleteScheduledWorkout deletes a scheduled workout if it belongs to the user.
func (uc *scheduledWorkoutUseCase) Delete(ctx context.Context, profileId, scheduledWorkoutId string) error {
	_, err := uc.GetByID(ctx, profileId, scheduledWorkoutId)
	if err != nil {
		return err
	}
	if err = uc.repo.SoftDeleteScheduledWorkout(ctx, scheduledWorkoutId, profileId); err != nil {
		return fmt.Errorf("failed to delete scheduled workout: %w", err)
	}
	return nil
}

// GetUpcommingScheduledWorkout implements ScheduledWorkoutUseCase.
func (uc *scheduledWorkoutUseCase) GetUpcommingScheduledWorkout(ctx context.Context, profileID string) (*model.ScheduledWorkout, error) {
	scheduledWorkout, err := uc.repo.GetUpcommingScheduledWorkout(ctx, profileID)
	if err != nil {
		return nil, err
	}
	return scheduledWorkout, nil
}
