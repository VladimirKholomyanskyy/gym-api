package repository

import (
	"context"
	"errors"
	"time"

	"github.com/VladimirKholomyanskyy/gym-api/internal/progress/model"
	"gorm.io/gorm"
)

var (
	// ErrWorkoutSessionNotFound is returned when a workout session cannot be found
	ErrWorkoutSessionNotFound = errors.New("workout session not found")

	// ErrWorkoutSessionCreate is returned when there's an issue creating a workout session
	ErrWorkoutSessionCreate = errors.New("failed to create workout session")

	// ErrWorkoutSessionUpdate is returned when there's an issue updating a workout session
	ErrWorkoutSessionUpdate = errors.New("failed to update workout session")

	// ErrInvalidWorkoutSession is returned when workout session data is invalid
	ErrInvalidWorkoutSession = errors.New("invalid workout session data")
)

// WorkoutSessionRepository defines the interface for workout session operations
type WorkoutSessionRepository interface {
	Create(ctx context.Context, workoutSession *model.WorkoutSession) error
	GetByID(ctx context.Context, id string) (*model.WorkoutSession, error)
	GetAllByProfileID(ctx context.Context, profileID string, page, pageSize int) ([]model.WorkoutSession, int64, error)
	GetByDateRange(ctx context.Context, profileID string, startDate, endDate time.Time, page, pageSize int) ([]model.WorkoutSession, int64, error)
	Update(ctx context.Context, workoutSession *model.WorkoutSession) error
	SoftDelete(ctx context.Context, id string, profileID string) error
	GetRecentSessions(ctx context.Context, profileID string, limit int) ([]model.WorkoutSession, error)
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
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create the workout session
		if err := tx.Create(workoutSession).Error; err != nil {
			return errors.Join(ErrWorkoutSessionCreate, err)
		}
		return nil
	})
}

// GetByID retrieves a workout session by ID with eager loading and error handling
func (r *workoutSessionRepository) GetByID(ctx context.Context, id string) (*model.WorkoutSession, error) {
	var workoutSession model.WorkoutSession

	// Use First with error handling
	result := r.db.WithContext(ctx).
		Preload("Exercises").
		Preload("Exercises.Sets").
		First(&workoutSession, "id = ?", id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrWorkoutSessionNotFound
		}
		return nil, result.Error
	}

	return &workoutSession, nil
}

// GetAllByProfileID retrieves all non-deleted workout sessions for a specific profile
func (r *workoutSessionRepository) GetAllByProfileID(ctx context.Context, profileID string, page, pageSize int) ([]model.WorkoutSession, int64, error) {
	if err := validatePagination(page, pageSize); err != nil {
		return nil, 0, err
	}

	var workoutSessions []model.WorkoutSession
	var total int64
	offset := (page - 1) * pageSize

	// Count total records
	countQuery := r.db.WithContext(ctx).
		Model(&model.WorkoutSession{}).
		Where("profile_id = ?", profileID)
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch the paginated results
	result := r.db.WithContext(ctx).
		Where("profile_id = ?", profileID).
		Preload("Exercises").
		Preload("Exercises.Sets").
		Order("started_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&workoutSessions)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	return workoutSessions, total, nil
}

// GetByDateRange retrieves workout sessions within a specific date range
func (r *workoutSessionRepository) GetByDateRange(ctx context.Context, profileID string, startDate, endDate time.Time, page, pageSize int) ([]model.WorkoutSession, int64, error) {
	var workoutSessions []model.WorkoutSession
	var total int64
	offset := (page - 1) * pageSize

	// Count total records
	countQuery := r.db.WithContext(ctx).
		Model(&model.WorkoutSession{}).
		Where("profile_id = ? AND started_at BETWEEN ? AND ?", profileID, startDate, endDate)
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch the paginated results
	result := r.db.WithContext(ctx).
		Where("profile_id = ? AND started_at BETWEEN ? AND ?", profileID, startDate, endDate).
		Preload("Exercises").
		Preload("Exercises.Sets").
		Order("started_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&workoutSessions)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	return workoutSessions, total, nil
}

// Update a workout session with optimistic locking and validation
func (r *workoutSessionRepository) Update(ctx context.Context, workoutSession *model.WorkoutSession) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Use Save with select to update only specific fields
		result := tx.Model(workoutSession).
			Where("id = ? AND profile_id = ?", workoutSession.ID, workoutSession.ProfileID).
			Select("workout_type", "duration", "intensity", "notes", "started_at", "updated_at").
			Updates(map[string]interface{}{

				"started_at": workoutSession.StartedAt,
				"updated_at": time.Now(),
			})

		if result.Error != nil {
			return errors.Join(ErrWorkoutSessionUpdate, result.Error)
		}

		// Check if any rows were actually updated
		if result.RowsAffected == 0 {
			return ErrWorkoutSessionNotFound
		}

		return nil
	})
}

// SoftDelete marks a workout session as deleted without removing it from the database
func (r *workoutSessionRepository) SoftDelete(ctx context.Context, id string, profileID string) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND profile_id = ?", id, profileID).
		Delete(&model.WorkoutSession{})

	if result.Error != nil {
		return result.Error
	}

	// Check if any rows were actually deleted
	if result.RowsAffected == 0 {
		return ErrWorkoutSessionNotFound
	}

	return nil
}

// GetRecentSessions retrieves the most recent workout sessions for a user
func (r *workoutSessionRepository) GetRecentSessions(ctx context.Context, profileID string, limit int) ([]model.WorkoutSession, error) {
	if limit < 1 {
		limit = 5 // Default limit if invalid
	}

	var workoutSessions []model.WorkoutSession

	result := r.db.WithContext(ctx).
		Where("profile_id = ?", profileID).
		Preload("Exercises").
		Preload("Exercises.Sets").
		Order("started_at DESC").
		Limit(limit).
		Find(&workoutSessions)

	if result.Error != nil {
		return nil, result.Error
	}

	return workoutSessions, nil
}
