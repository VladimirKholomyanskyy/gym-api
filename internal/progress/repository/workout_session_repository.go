package repository

import (
	"context"
	"fmt"
	"time"

	customerrors "github.com/VladimirKholomyanskyy/gym-api/internal/customErrors"
	"github.com/VladimirKholomyanskyy/gym-api/internal/progress/model"
	"gorm.io/gorm"
)

// WorkoutSessionRepository defines the interface for workout session operations
type WorkoutSessionRepository interface {
	Create(ctx context.Context, workoutSession *model.WorkoutSession) error
	GetByID(ctx context.Context, id string) (*model.WorkoutSession, error)
	GetAllByProfileID(ctx context.Context, profileID string, page, pageSize int) ([]model.WorkoutSession, int64, error)
	GetAllByProfileIDAndDateRange(ctx context.Context, profileID string, startDate, endDate time.Time, page, pageSize int) ([]model.WorkoutSession, int64, error)
	UpdatePartial(ctx context.Context, id string, updates map[string]any) error
	Delete(ctx context.Context, id string) error
	PermanentDelete(ctx context.Context, id string) error
}

// workoutSessionRepository implements WorkoutSessionRepository
type workoutSessionRepository struct {
	db *gorm.DB
}

// NewWorkoutSessionRepository creates a new repository instance
func NewWorkoutSessionRepository(db *gorm.DB) WorkoutSessionRepository {
	return &workoutSessionRepository{db: db}
}

// Create a new workout session with context and transaction support
func (r *workoutSessionRepository) Create(ctx context.Context, workoutSession *model.WorkoutSession) error {
	if err := r.db.WithContext(ctx).Create(workoutSession).Error; err != nil {
		return fmt.Errorf("failed to create workout session: %w", err)
	}
	return nil
}

// GetByID retrieves a workout session by ID with eager loading and error handling
func (r *workoutSessionRepository) GetByID(ctx context.Context, id string) (*model.WorkoutSession, error) {
	var workoutSession model.WorkoutSession

	err := r.db.WithContext(ctx).
		First(&workoutSession, "id = ?", id).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, customerrors.ErrEntityNotFound
		}
		return nil, fmt.Errorf("failed to fetch workout session: %w", err)
	}

	return &workoutSession, nil
}

// GetAllByProfileID retrieves all non-deleted workout sessions for a specific profile
func (r *workoutSessionRepository) GetAllByProfileID(ctx context.Context, profileID string, page, pageSize int) ([]model.WorkoutSession, int64, error) {
	var workoutSessions []model.WorkoutSession
	var total int64
	offset := (page - 1) * pageSize

	// Count total records
	countQuery := r.db.WithContext(ctx).
		Model(&model.WorkoutSession{}).
		Where("profile_id = ?", profileID)
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count workout sessions: %w", err)
	}

	// Fetch the paginated results
	err := r.db.WithContext(ctx).
		Where("profile_id = ?", profileID).
		Preload("Exercises").
		Preload("Exercises.Sets").
		Order("started_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&workoutSessions).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch workout sessions: %w", err)
	}

	return workoutSessions, total, nil
}

// GetByDateRange retrieves workout sessions within a specific date range
func (r *workoutSessionRepository) GetAllByProfileIDAndDateRange(ctx context.Context, profileID string, startDate, endDate time.Time, page, pageSize int) ([]model.WorkoutSession, int64, error) {
	var workoutSessions []model.WorkoutSession
	var total int64
	offset := (page - 1) * pageSize

	// Count total records
	countQuery := r.db.WithContext(ctx).
		Model(&model.WorkoutSession{}).
		Where("profile_id = ? AND started_at BETWEEN ? AND ?", profileID, startDate, endDate)
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count workout sessions: %w", err)
	}

	// Fetch the paginated results
	err := r.db.WithContext(ctx).
		Where("profile_id = ? AND started_at BETWEEN ? AND ?", profileID, startDate, endDate).
		Preload("Exercises").
		Preload("Exercises.Sets").
		Order("started_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&workoutSessions).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch workout sessions: %w", err)
	}

	return workoutSessions, total, nil
}

// Update a workout session with optimistic locking and validation
func (r *workoutSessionRepository) UpdatePartial(ctx context.Context, id string, updates map[string]any) error {
	result := r.db.WithContext(ctx).Model(&model.WorkoutSession{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to update workout session: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return customerrors.ErrEntityNotFound
	}
	return nil
}

// SoftDelete marks a workout session as deleted without removing it from the database
func (r *workoutSessionRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&model.WorkoutSession{})

	if result.Error != nil {
		return fmt.Errorf("failed to delete workout session: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return customerrors.ErrEntityNotFound
	}

	return nil
}
func (r *workoutSessionRepository) PermanentDelete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Unscoped().Where("id = ?", id).Delete(&model.WorkoutSession{})
	if result.Error != nil {
		return fmt.Errorf("failed to permanently delete workout session: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return customerrors.ErrEntityNotFound
	}
	return nil
}
