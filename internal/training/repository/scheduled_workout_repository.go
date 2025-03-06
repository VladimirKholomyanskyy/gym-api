package repository

import (
	"context"
	"errors"
	"time"

	"github.com/VladimirKholomyanskyy/gym-api/internal/training/model"
	"gorm.io/gorm"
)

var (
	// ErrScheduledWorkoutNotFound is returned when a scheduled workout cannot be found
	ErrScheduledWorkoutNotFound = errors.New("scheduled workout not found")

	// ErrScheduledWorkoutCreate is returned when there's an issue creating a scheduled workout
	ErrScheduledWorkoutCreate = errors.New("failed to create scheduled workout")

	// ErrScheduledWorkoutUpdate is returned when there's an issue updating a scheduled workout
	ErrScheduledWorkoutUpdate = errors.New("failed to update scheduled workout")

	// ErrInvalidDateRange is returned when the date range is invalid
	ErrInvalidDateRange = errors.New("invalid date range")
)

// ScheduledWorkoutRepository defines the interface for scheduled workout operations
type ScheduledWorkoutRepository interface {
	ScheduleWorkout(ctx context.Context, scheduledWorkout *model.ScheduledWorkout) error
	GetScheduledWorkout(ctx context.Context, id string) (*model.ScheduledWorkout, error)
	GetUserScheduledWorkouts(ctx context.Context, profileID string, page, pageSize int) ([]model.ScheduledWorkout, int64, error)
	GetUserScheduledWorkoutsByDate(ctx context.Context, profileID string, date time.Time, page, pageSize int) ([]model.ScheduledWorkout, int64, error)
	GetUserScheduledWorkoutsByRange(ctx context.Context, profileID string, startDate, endDate time.Time, page, pageSize int) ([]model.ScheduledWorkout, int64, error)
	UpdateScheduledWorkout(ctx context.Context, scheduledWorkout *model.ScheduledWorkout) error
	SoftDeleteScheduledWorkout(ctx context.Context, id string, profileID string) error
	GetAllScheduledWorkouts(ctx context.Context, page, pageSize int) ([]model.ScheduledWorkout, int64, error)
	GetUpcommingScheduledWorkout(ctx context.Context, profileID string) (*model.ScheduledWorkout, error)
}

// scheduledWorkoutRepository implements ScheduledWorkoutRepository
type scheduledWorkoutRepository struct {
	db *gorm.DB
}

// NewScheduledWorkoutRepository creates a new instance of ScheduledWorkoutRepository
func NewScheduledWorkoutRepository(db *gorm.DB) ScheduledWorkoutRepository {
	return &scheduledWorkoutRepository{db: db}
}

// validateDateRange checks if the date range is valid
func validateDateRange(startDate, endDate time.Time) error {
	if startDate.IsZero() || endDate.IsZero() {
		return ErrInvalidDateRange
	}
	if startDate.After(endDate) {
		return ErrInvalidDateRange
	}
	return nil
}

// ScheduleWorkout schedules a workout for a user on a given date
func (r *scheduledWorkoutRepository) ScheduleWorkout(ctx context.Context, scheduledWorkout *model.ScheduledWorkout) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Validate input
		if scheduledWorkout.ProfileID == "" {
			return errors.New("profile ID cannot be empty")
		}
		if scheduledWorkout.WorkoutID == "" {
			return errors.New("workout ID cannot be empty")
		}
		if scheduledWorkout.Date.IsZero() {
			scheduledWorkout.Date = time.Now()
		}

		// Create the scheduled workout
		if err := tx.Create(scheduledWorkout).Error; err != nil {
			return errors.Join(ErrScheduledWorkoutCreate, err)
		}
		return nil
	})
}

// GetScheduledWorkout retrieves a scheduled workout by ID
func (r *scheduledWorkoutRepository) GetScheduledWorkout(ctx context.Context, id string) (*model.ScheduledWorkout, error) {
	var scheduledWorkout model.ScheduledWorkout
	result := r.db.WithContext(ctx).
		Preload("Workout").
		First(&scheduledWorkout, "id = ?", id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrScheduledWorkoutNotFound
		}
		return nil, result.Error
	}

	return &scheduledWorkout, nil
}

// GetUserScheduledWorkouts retrieves paginated scheduled workouts for a specific user
func (r *scheduledWorkoutRepository) GetUserScheduledWorkouts(ctx context.Context, profileID string, page, pageSize int) ([]model.ScheduledWorkout, int64, error) {
	if err := validatePagination(page, pageSize); err != nil {
		return nil, 0, err
	}

	var workouts []model.ScheduledWorkout
	var total int64
	offset := (page - 1) * pageSize

	// Count total records
	countQuery := r.db.WithContext(ctx).
		Model(&model.ScheduledWorkout{}).
		Where("profile_id = ?", profileID)
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch paginated results
	result := r.db.WithContext(ctx).
		Where("profile_id = ?", profileID).
		Preload("Workout").
		Order("date ASC").
		Limit(pageSize).
		Offset(offset).
		Find(&workouts)

	return workouts, total, result.Error
}

// GetUserScheduledWorkoutsByDate retrieves paginated workouts scheduled for a user on a specific date
func (r *scheduledWorkoutRepository) GetUserScheduledWorkoutsByDate(ctx context.Context, profileID string, date time.Time, page, pageSize int) ([]model.ScheduledWorkout, int64, error) {
	if err := validatePagination(page, pageSize); err != nil {
		return nil, 0, err
	}

	var workouts []model.ScheduledWorkout
	var total int64
	offset := (page - 1) * pageSize

	// Count total records
	countQuery := r.db.WithContext(ctx).
		Model(&model.ScheduledWorkout{}).
		Where("profile_id = ? AND date = ?", profileID, date)
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch paginated results
	result := r.db.WithContext(ctx).
		Where("profile_id = ? AND date = ?", profileID, date).
		Preload("Workout").
		Limit(pageSize).
		Offset(offset).
		Find(&workouts)

	return workouts, total, result.Error
}

// GetUserScheduledWorkoutsByRange retrieves paginated scheduled workouts for a user within a date range
func (r *scheduledWorkoutRepository) GetUserScheduledWorkoutsByRange(ctx context.Context, profileID string, startDate, endDate time.Time, page, pageSize int) ([]model.ScheduledWorkout, int64, error) {
	if err := validatePagination(page, pageSize); err != nil {
		return nil, 0, err
	}

	if err := validateDateRange(startDate, endDate); err != nil {
		return nil, 0, err
	}

	var workouts []model.ScheduledWorkout
	var total int64
	offset := (page - 1) * pageSize

	// Count total records
	countQuery := r.db.WithContext(ctx).
		Model(&model.ScheduledWorkout{}).
		Where("profile_id = ? AND date BETWEEN ? AND ?", profileID, startDate, endDate)
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch paginated results
	result := r.db.WithContext(ctx).
		Where("profile_id = ? AND date BETWEEN ? AND ?", profileID, startDate, endDate).
		Preload("Workout").
		Order("date ASC").
		Limit(pageSize).
		Offset(offset).
		Find(&workouts)

	return workouts, total, result.Error
}

// UpdateScheduledWorkout updates an existing scheduled workout
func (r *scheduledWorkoutRepository) UpdateScheduledWorkout(ctx context.Context, scheduledWorkout *model.ScheduledWorkout) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Validate input
		if scheduledWorkout.ProfileID == "" {
			return errors.New("profile ID cannot be empty")
		}
		if scheduledWorkout.WorkoutID == "" {
			return errors.New("workout ID cannot be empty")
		}
		if scheduledWorkout.Date.IsZero() {
			scheduledWorkout.Date = time.Now()
		}

		// Update the scheduled workout
		result := tx.Model(scheduledWorkout).
			Select("workout_id", "date", "notes").
			Save(scheduledWorkout)

		if result.Error != nil {
			return errors.Join(ErrScheduledWorkoutUpdate, result.Error)
		}

		// Check if any rows were actually updated
		if result.RowsAffected == 0 {
			return ErrScheduledWorkoutNotFound
		}

		return nil
	})
}

// SoftDeleteScheduledWorkout removes a scheduled workout by ID and user
func (r *scheduledWorkoutRepository) SoftDeleteScheduledWorkout(ctx context.Context, id string, profileID string) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND profile_id = ?", id, profileID).
		Delete(&model.ScheduledWorkout{})

	if result.Error != nil {
		return result.Error
	}

	// Check if any rows were actually deleted
	if result.RowsAffected == 0 {
		return ErrScheduledWorkoutNotFound
	}

	return nil
}

// GetAllScheduledWorkouts retrieves all scheduled workouts with pagination
func (r *scheduledWorkoutRepository) GetAllScheduledWorkouts(ctx context.Context, page, pageSize int) ([]model.ScheduledWorkout, int64, error) {
	if err := validatePagination(page, pageSize); err != nil {
		return nil, 0, err
	}

	var workouts []model.ScheduledWorkout
	var total int64
	offset := (page - 1) * pageSize

	// Count total records
	countQuery := r.db.WithContext(ctx).Model(&model.ScheduledWorkout{})
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch paginated results
	result := r.db.WithContext(ctx).
		Preload("Workout").
		Order("date ASC").
		Limit(pageSize).
		Offset(offset).
		Find(&workouts)

	return workouts, total, result.Error
}

// GetClosestScheduledWorkout retrieves the scheduled workout that is closest to today's date
func (r *scheduledWorkoutRepository) GetUpcommingScheduledWorkout(ctx context.Context, profileID string) (*model.ScheduledWorkout, error) {
	var scheduledWorkout model.ScheduledWorkout
	today := time.Now()

	// First try to find the closest future workout
	futureResult := r.db.WithContext(ctx).
		Where("profile_id = ? AND date >= ?", profileID, today).
		Preload("Workout").
		Order("date ASC").
		Limit(1).
		Find(&scheduledWorkout)

	// If we found a future workout, return it
	if futureResult.Error == nil && futureResult.RowsAffected > 0 {
		return &scheduledWorkout, nil
	}
	return &scheduledWorkout, nil
}
