package repository

import (
	"context"
	"fmt"
	"time"

	customerrors "github.com/VladimirKholomyanskyy/gym-api/internal/customErrors"
	"github.com/VladimirKholomyanskyy/gym-api/internal/progress/model"
	"gorm.io/gorm"
)

// ExerciseLogRepository defines the interface for exercise log operations
type ExerciseLogRepository interface {
	Create(ctx context.Context, exerciseLog *model.ExerciseLog) error
	GetByID(ctx context.Context, id string) (*model.ExerciseLog, error)
	GetAllByProfileIDAndSessionID(ctx context.Context, profileID, sessionID string, page, pageSize int) ([]model.ExerciseLog, int64, error)
	GetAllByProfileIDAndExerciseID(ctx context.Context, profileID, exerciseID string, page, pageSize int) ([]model.ExerciseLog, int64, error)
	GetAllByProfileID(ctx context.Context, profileID string, page, pageSize int) ([]model.ExerciseLog, int64, error)
	UpdatePartial(ctx context.Context, id string, updates map[string]any) (*model.ExerciseLog, error)
	Delete(ctx context.Context, id string) error
	PermanentDelete(ctx context.Context, id string) error
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

// Create a new exercise log with context and transaction support
func (r *exerciseLogRepository) Create(ctx context.Context, exerciseLog *model.ExerciseLog) error {
	if err := r.db.WithContext(ctx).Create(exerciseLog).Error; err != nil {
		return fmt.Errorf("failed to create exercise log: %w", err)
	}
	return nil
}

// GetByID retrieves an exercise log by ID with error handling
func (r *exerciseLogRepository) GetByID(ctx context.Context, id string) (*model.ExerciseLog, error) {
	var exerciseLog model.ExerciseLog

	err := r.db.WithContext(ctx).First(&exerciseLog, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, customerrors.ErrEntityNotFound
		}
		return nil, fmt.Errorf("failed to fetch exercise log by id: %w", err)
	}

	return &exerciseLog, nil
}

// GetAllByUserIDAndSessionID retrieves paginated exercise logs for a specific user and session
func (r *exerciseLogRepository) GetAllByProfileIDAndSessionID(ctx context.Context, profileID, sessionID string, page, pageSize int) ([]model.ExerciseLog, int64, error) {
	var exerciseLogs []model.ExerciseLog
	var total int64
	offset := (page - 1) * pageSize

	// Count total records
	countQuery := r.db.WithContext(ctx).
		Model(&model.ExerciseLog{}).
		Where("profile_id = ? AND session_id = ?", profileID, sessionID)
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count exercise logs: %w", err)
	}

	// Fetch paginated results
	err := r.db.WithContext(ctx).
		Limit(pageSize).
		Offset(offset).
		Where("profile_id = ? AND session_id = ?", profileID, sessionID).
		Order("logged_at DESC").
		Find(&exerciseLogs).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch exercise logs by profile and session ids: %w", err)
	}
	return exerciseLogs, total, nil
}

// GetAllByUserIDAndExerciseID retrieves paginated exercise logs for a specific user and exercise
func (r *exerciseLogRepository) GetAllByProfileIDAndExerciseID(ctx context.Context, profileID, exerciseID string, page, pageSize int) ([]model.ExerciseLog, int64, error) {
	var exerciseLogs []model.ExerciseLog
	var total int64
	offset := (page - 1) * pageSize

	// Count total records
	countQuery := r.db.WithContext(ctx).
		Model(&model.ExerciseLog{}).
		Where("profile_id = ? AND exercise_id = ?", profileID, exerciseID)
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count exercise logs: %w", err)
	}

	// Fetch paginated results
	err := r.db.WithContext(ctx).
		Limit(pageSize).
		Offset(offset).
		Where("profile_id = ? AND exercise_id = ?", profileID, exerciseID).
		Order("logged_at DESC").
		Find(&exerciseLogs).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch exercise logs by profile and exercise ids: %w", err)
	}
	return exerciseLogs, total, nil
}

// GetAllByUserID retrieves paginated exercise logs for a specific user
func (r *exerciseLogRepository) GetAllByProfileID(ctx context.Context, profileID string, page, pageSize int) ([]model.ExerciseLog, int64, error) {
	var exerciseLogs []model.ExerciseLog
	var total int64
	offset := (page - 1) * pageSize

	// Count total records
	countQuery := r.db.WithContext(ctx).
		Model(&model.ExerciseLog{}).
		Where("profile_id = ?", profileID)
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count exercise logs: %w", err)
	}

	// Fetch paginated results
	err := r.db.WithContext(ctx).
		Limit(pageSize).
		Offset(offset).
		Where("profile_id = ?", profileID).
		Order("logged_at DESC").
		Find(&exerciseLogs).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch exercise logs by profile and exercise ids: %w", err)
	}
	return exerciseLogs, total, nil
}

// Update an exercise log with optimistic locking and validation
func (r *exerciseLogRepository) UpdatePartial(ctx context.Context, id string, updates map[string]any) (*model.ExerciseLog, error) {
	result := r.db.WithContext(ctx).Model(&model.ExerciseLog{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to update exercise log: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, customerrors.ErrEntityNotFound
	}
	var updatedExerciseLog model.ExerciseLog
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&updatedExerciseLog).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch updated exercise log: %w", err)
	}

	return &updatedExerciseLog, nil
}

// SoftDelete marks an exercise log as deleted without removing it from the database
func (r *exerciseLogRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).
		Delete(&model.ExerciseLog{}, id)

	if result.Error != nil {
		return result.Error
	}

	// Check if any rows were actually deleted
	if result.RowsAffected == 0 {
		return customerrors.ErrEntityNotFound
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
func (r *exerciseLogRepository) PermanentDelete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Unscoped().Where("id = ?", id).Delete(&model.ExerciseLog{})
	if result.Error != nil {
		return fmt.Errorf("failed to permanent delete exercise log: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return customerrors.ErrEntityNotFound
	}
	return nil
}
