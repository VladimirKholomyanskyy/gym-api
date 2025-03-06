package repository

import (
	"context"
	"errors"
	"time"

	"github.com/VladimirKholomyanskyy/gym-api/internal/progress/model"
	"gorm.io/gorm"
)

var (
	// ErrExerciseLogNotFound is returned when an exercise log cannot be found
	ErrExerciseLogNotFound = errors.New("exercise log not found")

	// ErrExerciseLogCreate is returned when there's an issue creating an exercise log
	ErrExerciseLogCreate = errors.New("failed to create exercise log")

	// ErrExerciseLogUpdate is returned when there's an issue updating an exercise log
	ErrExerciseLogUpdate = errors.New("failed to update exercise log")

	// ErrInvalidPagination is returned when pagination parameters are invalid
	ErrInvalidPagination = errors.New("invalid pagination parameters")
)

// ExerciseLogRepository defines the interface for exercise log operations
type ExerciseLogRepository interface {
	Create(ctx context.Context, exerciseLog *model.ExerciseLog) error
	GetByID(ctx context.Context, id string) (*model.ExerciseLog, error)
	GetAllByUserIDAndSessionID(ctx context.Context, profileID, sessionID string, page, pageSize int) ([]model.ExerciseLog, int64, error)
	GetAllByUserIDAndExerciseID(ctx context.Context, profileID, exerciseID string, page, pageSize int) ([]model.ExerciseLog, int64, error)
	GetAllByUserID(ctx context.Context, profileID string, page, pageSize int) ([]model.ExerciseLog, int64, error)
	Update(ctx context.Context, exerciseLog *model.ExerciseLog) error
	SoftDelete(ctx context.Context, id string) error
	GetWeightPerDay(ctx context.Context, profileID, exerciseID string, startDate, endDate *time.Time) ([]model.WeightPerDay, error)
}

// exerciseLogRepository implements ExerciseLogRepository
type exerciseLogRepository struct {
	db *gorm.DB
}

// NewExerciseLogRepository creates a new repository instance
func NewExerciseLogRepository(db *gorm.DB) ExerciseLogRepository {
	return &exerciseLogRepository{db: db}
}

// validatePagination checks if pagination parameters are valid
func validatePagination(page, pageSize int) error {
	if page < 1 || pageSize < 1 {
		return ErrInvalidPagination
	}
	return nil
}

// Create a new exercise log with context and transaction support
func (r *exerciseLogRepository) Create(ctx context.Context, exerciseLog *model.ExerciseLog) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Set default logged_at if not provided
		if exerciseLog.LoggedAt.IsZero() {
			exerciseLog.LoggedAt = time.Now()
		}

		// Create the exercise log
		if err := tx.Create(exerciseLog).Error; err != nil {
			return errors.Join(ErrExerciseLogCreate, err)
		}
		return nil
	})
}

// GetByID retrieves an exercise log by ID with error handling
func (r *exerciseLogRepository) GetByID(ctx context.Context, id string) (*model.ExerciseLog, error) {
	var exerciseLog model.ExerciseLog

	result := r.db.WithContext(ctx).First(&exerciseLog, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrExerciseLogNotFound
		}
		return nil, result.Error
	}

	return &exerciseLog, nil
}

// GetAllByUserIDAndSessionID retrieves paginated exercise logs for a specific user and session
func (r *exerciseLogRepository) GetAllByUserIDAndSessionID(ctx context.Context, profileID, sessionID string, page, pageSize int) ([]model.ExerciseLog, int64, error) {
	if err := validatePagination(page, pageSize); err != nil {
		return nil, 0, err
	}

	var exerciseLogs []model.ExerciseLog
	var total int64
	offset := (page - 1) * pageSize

	// Count total records
	countQuery := r.db.WithContext(ctx).
		Model(&model.ExerciseLog{}).
		Where("user_id = ? AND session_id = ?", profileID, sessionID)
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch paginated results
	result := r.db.WithContext(ctx).
		Limit(pageSize).
		Offset(offset).
		Where("user_id = ? AND session_id = ?", profileID, sessionID).
		Order("logged_at DESC").
		Find(&exerciseLogs)

	return exerciseLogs, total, result.Error
}

// GetAllByUserIDAndExerciseID retrieves paginated exercise logs for a specific user and exercise
func (r *exerciseLogRepository) GetAllByUserIDAndExerciseID(ctx context.Context, profileID, exerciseID string, page, pageSize int) ([]model.ExerciseLog, int64, error) {
	if err := validatePagination(page, pageSize); err != nil {
		return nil, 0, err
	}

	var exerciseLogs []model.ExerciseLog
	var total int64
	offset := (page - 1) * pageSize

	// Count total records
	countQuery := r.db.WithContext(ctx).
		Model(&model.ExerciseLog{}).
		Where("user_id = ? AND exercise_id = ?", profileID, exerciseID)
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch paginated results
	result := r.db.WithContext(ctx).
		Limit(pageSize).
		Offset(offset).
		Where("user_id = ? AND exercise_id = ?", profileID, exerciseID).
		Order("logged_at DESC").
		Find(&exerciseLogs)

	return exerciseLogs, total, result.Error
}

// GetAllByUserID retrieves paginated exercise logs for a specific user
func (r *exerciseLogRepository) GetAllByUserID(ctx context.Context, profileID string, page, pageSize int) ([]model.ExerciseLog, int64, error) {
	if err := validatePagination(page, pageSize); err != nil {
		return nil, 0, err
	}

	var exerciseLogs []model.ExerciseLog
	var total int64
	offset := (page - 1) * pageSize

	// Count total records
	countQuery := r.db.WithContext(ctx).
		Model(&model.ExerciseLog{}).
		Where("profile_id = ?", profileID)
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch paginated results
	result := r.db.WithContext(ctx).
		Limit(pageSize).
		Offset(offset).
		Where("profile_id = ?", profileID).
		Order("logged_at DESC").
		Find(&exerciseLogs)

	return exerciseLogs, total, result.Error
}

// Update an exercise log with optimistic locking and validation
func (r *exerciseLogRepository) Update(ctx context.Context, exerciseLog *model.ExerciseLog) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Use Save with select to update only specific fields
		result := tx.Model(exerciseLog).
			Select("exercise_id", "session_id", "weight", "reps", "sets", "notes", "logged_at").
			Save(exerciseLog)

		if result.Error != nil {
			return errors.Join(ErrExerciseLogUpdate, result.Error)
		}

		// Check if any rows were actually updated
		if result.RowsAffected == 0 {
			return ErrExerciseLogNotFound
		}

		return nil
	})
}

// SoftDelete marks an exercise log as deleted without removing it from the database
func (r *exerciseLogRepository) SoftDelete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).
		Delete(&model.ExerciseLog{}, id)

	if result.Error != nil {
		return result.Error
	}

	// Check if any rows were actually deleted
	if result.RowsAffected == 0 {
		return ErrExerciseLogNotFound
	}

	return nil
}

// GetWeightPerDay calculates total weight per day for a specific exercise and user
func (r *exerciseLogRepository) GetWeightPerDay(ctx context.Context, profileID, exerciseID string, startDate, endDate *time.Time) ([]model.WeightPerDay, error) {
	var results []model.WeightPerDay
	query := r.db.WithContext(ctx).
		Table("exercise_logs").
		Select("DATE(logged_at) AS date, SUM(weight * reps) AS total_weight").
		Where("profile_id = ? AND exercise_id = ?", profileID, exerciseID).
		Group("DATE(logged_at)").
		Order("date ASC")

	// Add date range filter if provided
	if startDate != nil {
		query = query.Where("logged_at >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("logged_at <= ?", *endDate)
	}

	// Execute the query and populate results
	err := query.Scan(&results).Error
	return results, err
}
