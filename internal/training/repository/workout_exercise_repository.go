package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/VladimirKholomyanskyy/gym-api/internal/training/model"

	"gorm.io/gorm"
)

var (
	// Custom error types for more specific error handling
	ErrWorkoutExerciseNotFound = errors.New("workout exercise not found")
)

type WorkoutExerciseRepository interface {
	Create(ctx context.Context, exercise *model.WorkoutExercise) error
	FindByID(ctx context.Context, id string) (*model.WorkoutExercise, error)
	FindByWorkoutID(ctx context.Context, workoutID string, page, pageSize int) ([]model.WorkoutExercise, int64, error)
	FindByExerciseID(ctx context.Context, exerciseID string) ([]model.WorkoutExercise, error)
	Update(ctx context.Context, workoutExercise *model.WorkoutExercise) error
	Delete(ctx context.Context, id string) error
	Reorder(ctx context.Context, workoutExerciseID string, newPosition int) error
}

type workoutExerciseRepository struct {
	db *gorm.DB
}

// NewWorkoutExerciseRepository creates a new instance of WorkoutExerciseRepository
func NewWorkoutExerciseRepository(db *gorm.DB) WorkoutExerciseRepository {
	return &workoutExerciseRepository{db: db}
}

// Create inserts a new workout exercise into the database
func (r *workoutExerciseRepository) Create(ctx context.Context, exercise *model.WorkoutExercise) error {
	// Validate input
	if exercise == nil {
		return errors.New("workout exercise cannot be nil")
	}

	// Use a transaction to ensure atomic last position calculation
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Calculate the last position within the specific workout
		var lastPosition int
		err := tx.Model(&model.WorkoutExercise{}).
			Where("workout_id = ?", exercise.WorkoutID).
			Select("COALESCE(MAX(position), 0)").
			Scan(&lastPosition).Error
		if err != nil {
			return fmt.Errorf("failed to calculate last position: %w", err)
		}

		// Set the new position
		exercise.Position = lastPosition + 1

		// Create the exercise
		return tx.Create(exercise).Error
	})
}

// FindByID retrieves a workout exercise by its ID
func (r *workoutExerciseRepository) FindByID(ctx context.Context, id string) (*model.WorkoutExercise, error) {
	var exercise model.WorkoutExercise
	result := r.db.WithContext(ctx).First(&exercise, "id = ?", id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrWorkoutExerciseNotFound
		}
		return nil, result.Error
	}

	return &exercise, nil
}

// FindByWorkoutID retrieves all exercises associated with a given workout
// Now returns the total count for pagination purposes
func (r *workoutExerciseRepository) FindByWorkoutID(
	ctx context.Context,
	workoutID string,
	page,
	pageSize int,
) ([]model.WorkoutExercise, int64, error) {
	var exercises []model.WorkoutExercise
	var totalCount int64

	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10 // Default page size
	}

	offset := (page - 1) * pageSize

	// Count total exercises for the workout
	countResult := r.db.WithContext(ctx).
		Model(&model.WorkoutExercise{}).
		Where("workout_id = ?", workoutID).
		Count(&totalCount)
	if countResult.Error != nil {
		return nil, 0, countResult.Error
	}

	// Fetch paginated exercises
	err := r.db.WithContext(ctx).
		Where("workout_id = ?", workoutID).
		Order("position ASC").
		Limit(pageSize).
		Offset(offset).
		Find(&exercises).Error

	if err != nil {
		return nil, 0, err
	}

	return exercises, totalCount, nil
}

// FindByExerciseID retrieves all instances of a specific exercise in any workout
func (r *workoutExerciseRepository) FindByExerciseID(ctx context.Context, exerciseID string) ([]model.WorkoutExercise, error) {
	var exercises []model.WorkoutExercise

	result := r.db.WithContext(ctx).
		Where("exercise_id = ?", exerciseID).
		Find(&exercises)

	if result.Error != nil {
		return nil, result.Error
	}

	return exercises, nil
}

// Update modifies an existing workout exercise
func (r *workoutExerciseRepository) Update(ctx context.Context, workoutExercise *model.WorkoutExercise) error {
	// Validate input
	if workoutExercise == nil {
		return errors.New("workout exercise cannot be nil")
	}

	return r.db.WithContext(ctx).Save(workoutExercise).Error
}

// Delete removes a workout exercise from the database
func (r *workoutExerciseRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// First, get the workout exercise to ensure it exists and get its workout ID
		var workoutExercise model.WorkoutExercise
		if err := tx.First(&workoutExercise, "id = ?", id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrWorkoutExerciseNotFound
			}
			return err
		}

		// Delete the exercise
		result := tx.Delete(&model.WorkoutExercise{}, "id = ?", id)
		if result.Error != nil {
			return result.Error
		}

		// Reorder remaining exercises in the workout
		return tx.Model(&model.WorkoutExercise{}).
			Where("workout_id = ? AND position > ?", workoutExercise.WorkoutID, workoutExercise.Position).
			Update("position", gorm.Expr("position - 1")).Error
	})
}

// Reorder handles repositioning a workout exercise within its workout
func (r *workoutExerciseRepository) Reorder(ctx context.Context, workoutExerciseID string, newPosition int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Find the workout exercise
		var workoutExercise model.WorkoutExercise
		if err := tx.First(&workoutExercise, "id = ?", workoutExerciseID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrWorkoutExerciseNotFound
			}
			return err
		}

		// Validate new position
		var totalExercises int64
		if err := tx.Model(&model.WorkoutExercise{}).
			Where("workout_id = ?", workoutExercise.WorkoutID).
			Count(&totalExercises).Error; err != nil {
			return err
		}

		if newPosition < 1 || newPosition > int(totalExercises) {
			return fmt.Errorf("%w: position must be between 1 and %d", ErrInvalidPosition, totalExercises)
		}

		// No change needed
		if workoutExercise.Position == newPosition {
			return nil
		}

		// Shift positions
		if workoutExercise.Position < newPosition {
			// Move exercises down
			if err := tx.Model(&model.WorkoutExercise{}).
				Where("workout_id = ? AND position > ? AND position <= ?",
					workoutExercise.WorkoutID, workoutExercise.Position, newPosition).
				Update("position", gorm.Expr("position - 1")).Error; err != nil {
				return err
			}
		} else {
			// Move exercises up
			if err := tx.Model(&model.WorkoutExercise{}).
				Where("workout_id = ? AND position < ? AND position >= ?",
					workoutExercise.WorkoutID, workoutExercise.Position, newPosition).
				Update("position", gorm.Expr("position + 1")).Error; err != nil {
				return err
			}
		}

		// Update the exercise's position
		return tx.Model(&model.WorkoutExercise{}).
			Where("id = ?", workoutExerciseID).
			Update("position", newPosition).Error
	})
}
