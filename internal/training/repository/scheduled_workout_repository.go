package repository

import (
	"context"
	"fmt"
	"time"

	customerrors "github.com/VladimirKholomyanskyy/gym-api/internal/customErrors"
	"github.com/VladimirKholomyanskyy/gym-api/internal/training/model"
	"gorm.io/gorm"
)

// ScheduledWorkoutRepository defines the interface for scheduled workout operations
type ScheduledWorkoutRepository interface {
	Create(ctx context.Context, scheduledWorkout *model.ScheduledWorkout) error
	GetByID(ctx context.Context, id string) (*model.ScheduledWorkout, error)
	GetAllByProfileID(ctx context.Context, profileID string, page, pageSize int) ([]model.ScheduledWorkout, int64, error)
	GetAllByProfileIDAndDate(ctx context.Context, profileID string, date time.Time, page, pageSize int) ([]model.ScheduledWorkout, int64, error)
	GetAllByProfileIDAndRange(ctx context.Context, profileID string, startDate, endDate time.Time, page, pageSize int) ([]model.ScheduledWorkout, int64, error)
	UpdatePartial(ctx context.Context, id string, updates map[string]any) (*model.ScheduledWorkout, error)
	Delete(ctx context.Context, id string) error
	PermanentDelete(ctx context.Context, id string) error
	GetUpcomming(ctx context.Context, profileID string) (*model.ScheduledWorkout, error)
}

// scheduledWorkoutRepository implements ScheduledWorkoutRepository
type scheduledWorkoutRepository struct {
	db *gorm.DB
}

// NewScheduledWorkoutRepository creates a new instance of ScheduledWorkoutRepository
func NewScheduledWorkoutRepository(db *gorm.DB) ScheduledWorkoutRepository {
	return &scheduledWorkoutRepository{db: db}
}

// ScheduleWorkout schedules a workout for a user on a given date
func (r *scheduledWorkoutRepository) Create(ctx context.Context, scheduledWorkout *model.ScheduledWorkout) error {
	if err := r.db.WithContext(ctx).Create(scheduledWorkout).Error; err != nil {
		return fmt.Errorf("failed to create scheduled workout: %w", err)
	}
	return nil
}

// GetScheduledWorkout retrieves a scheduled workout by ID
func (r *scheduledWorkoutRepository) GetByID(ctx context.Context, id string) (*model.ScheduledWorkout, error) {
	var scheduledWorkout model.ScheduledWorkout
	err := r.db.WithContext(ctx).
		Preload("Workout").
		First(&scheduledWorkout, "id = ?", id).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, customerrors.ErrEntityNotFound
		}
		return nil, fmt.Errorf("failed to fetch scheduled workout: %w", err)
	}

	return &scheduledWorkout, nil
}

// GetUserScheduledWorkouts retrieves paginated scheduled workouts for a specific user
func (r *scheduledWorkoutRepository) GetAllByProfileID(ctx context.Context, profileID string, page, pageSize int) ([]model.ScheduledWorkout, int64, error) {
	var workouts []model.ScheduledWorkout
	var total int64
	offset := (page - 1) * pageSize

	countQuery := r.db.WithContext(ctx).
		Model(&model.ScheduledWorkout{}).
		Where("profile_id = ?", profileID)
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count scheduled workouts: %w", err)
	}

	err := r.db.WithContext(ctx).
		Where("profile_id = ?", profileID).
		Preload("Workout").
		Order("date ASC").
		Limit(pageSize).
		Offset(offset).
		Find(&workouts).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch scheduled workouts: %w", err)
	}
	return workouts, total, nil
}

// GetUserScheduledWorkoutsByDate retrieves paginated workouts scheduled for a user on a specific date
func (r *scheduledWorkoutRepository) GetAllByProfileIDAndDate(ctx context.Context, profileID string, date time.Time, page, pageSize int) ([]model.ScheduledWorkout, int64, error) {
	var workouts []model.ScheduledWorkout
	var total int64
	offset := (page - 1) * pageSize

	countQuery := r.db.WithContext(ctx).
		Model(&model.ScheduledWorkout{}).
		Where("profile_id = ? AND date = ?", profileID, date)
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count scheduled workouts: %w", err)
	}

	err := r.db.WithContext(ctx).
		Where("profile_id = ? AND date = ?", profileID, date).
		Preload("Workout").
		Limit(pageSize).
		Offset(offset).
		Find(&workouts).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch scheduled workouts: %w", err)
	}
	return workouts, total, nil
}

// GetUserScheduledWorkoutsByRange retrieves paginated scheduled workouts for a user within a date range
func (r *scheduledWorkoutRepository) GetAllByProfileIDAndRange(ctx context.Context, profileID string, startDate, endDate time.Time, page, pageSize int) ([]model.ScheduledWorkout, int64, error) {
	var workouts []model.ScheduledWorkout
	var total int64
	offset := (page - 1) * pageSize

	countQuery := r.db.WithContext(ctx).
		Model(&model.ScheduledWorkout{}).
		Where("profile_id = ? AND date BETWEEN ? AND ?", profileID, startDate, endDate)
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count scheduled workouts: %w", err)
	}

	err := r.db.WithContext(ctx).
		Where("profile_id = ? AND date BETWEEN ? AND ?", profileID, startDate, endDate).
		Preload("Workout").
		Order("date ASC").
		Limit(pageSize).
		Offset(offset).
		Find(&workouts).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch scheduled workouts: %w", err)
	}
	return workouts, total, nil
}

// UpdateScheduledWorkout updates an existing scheduled workout
func (r *scheduledWorkoutRepository) UpdatePartial(ctx context.Context, id string, updates map[string]any) (*model.ScheduledWorkout, error) {
	result := r.db.WithContext(ctx).Model(&model.ScheduledWorkout{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to update scheduled workout: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, customerrors.ErrEntityNotFound
	}
	var updatedScheduledWorkout model.ScheduledWorkout
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&updatedScheduledWorkout).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch updated scheduled workout: %w", err)
	}

	return &updatedScheduledWorkout, nil
}

// SoftDeleteScheduledWorkout removes a scheduled workout by ID and user
func (r *scheduledWorkoutRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&model.ScheduledWorkout{})

	if result.Error != nil {
		return fmt.Errorf("failed to delete profile: %w", result.Error)
	}

	// Check if any rows were actually deleted
	if result.RowsAffected == 0 {
		return customerrors.ErrEntityNotFound
	}

	return nil
}

func (r *scheduledWorkoutRepository) PermanentDelete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Unscoped().Where("id = ?", id).Delete(&model.ScheduledWorkout{})
	if result.Error != nil {
		return fmt.Errorf("failed to permanent delete profile: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return customerrors.ErrEntityNotFound
	}
	return nil
}

// GetClosestScheduledWorkout retrieves the scheduled workout that is closest to today's date
func (r *scheduledWorkoutRepository) GetUpcomming(ctx context.Context, profileID string) (*model.ScheduledWorkout, error) {
	var scheduledWorkout model.ScheduledWorkout
	today := time.Now()

	result := r.db.WithContext(ctx).
		Where("profile_id = ? AND date >= ?", profileID, today).
		Preload("Workout").
		Order("date ASC").
		Limit(1).
		Find(&scheduledWorkout)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to permanent delete profile: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, customerrors.ErrEntityNotFound
	}

	return &scheduledWorkout, nil
}
