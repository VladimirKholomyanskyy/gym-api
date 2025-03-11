package repository

import (
	"context"
	"fmt"

	customerrors "github.com/VladimirKholomyanskyy/gym-api/internal/customErrors"
	"github.com/VladimirKholomyanskyy/gym-api/internal/training/model"

	"gorm.io/gorm"
)

type WorkoutExerciseRepository interface {
	Create(ctx context.Context, exercise *model.WorkoutExercise) error
	GetByID(ctx context.Context, id string) (*model.WorkoutExercise, error)
	GetAllByWorkoutID(ctx context.Context, workoutID string, page, pageSize int) ([]model.WorkoutExercise, int64, error)
	UpdatePartial(ctx context.Context, id string, updates map[string]any) (*model.WorkoutExercise, error)
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
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var lastPosition int
		err := tx.Model(&model.WorkoutExercise{}).
			Where("workout_id = ?", exercise.WorkoutID).
			Select("COALESCE(MAX(position), 0)").
			Scan(&lastPosition).Error
		if err != nil {
			return fmt.Errorf("failed to calculate last position: %w", err)
		}
		exercise.Position = lastPosition + 1
		if err := tx.Create(exercise).Error; err != nil {
			return fmt.Errorf("failed to create workout exercise: %w", err)
		}
		return nil
	})
}

// FindByID retrieves a workout exercise by its ID
func (r *workoutExerciseRepository) GetByID(ctx context.Context, id string) (*model.WorkoutExercise, error) {
	var exercise model.WorkoutExercise
	err := r.db.WithContext(ctx).First(&exercise, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, customerrors.ErrEntityNotFound
		}
		return nil, fmt.Errorf("failed to fetch workout exercise: %w", err)
	}
	return &exercise, nil
}

// FindByWorkoutID retrieves all exercises associated with a given workout
// Now returns the total count for pagination purposes
func (r *workoutExerciseRepository) GetAllByWorkoutID(
	ctx context.Context,
	workoutID string,
	page,
	pageSize int,
) ([]model.WorkoutExercise, int64, error) {
	var exercises []model.WorkoutExercise
	var totalCount int64
	offset := (page - 1) * pageSize

	countErr := r.db.WithContext(ctx).
		Model(&model.WorkoutExercise{}).
		Where("workout_id = ?", workoutID).
		Count(&totalCount).Error
	if countErr != nil {
		return nil, 0, fmt.Errorf("failed to count workouts exercises: %w", countErr)
	}

	err := r.db.WithContext(ctx).
		Where("workout_id = ?", workoutID).
		Order("position ASC").
		Limit(pageSize).
		Offset(offset).
		Find(&exercises).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch workouts exercises: %w", err)
	}

	return exercises, totalCount, nil
}

// Update modifies an existing workout exercise
func (r *workoutExerciseRepository) UpdatePartial(ctx context.Context, id string, updates map[string]any) (*model.WorkoutExercise, error) {
	result := r.db.WithContext(ctx).Model(&model.WorkoutExercise{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to update workout exercise: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, customerrors.ErrEntityNotFound
	}
	var updatedWorkoutExercise model.WorkoutExercise
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&updatedWorkoutExercise).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch updated workout exercise: %w", err)
	}

	return &updatedWorkoutExercise, nil
}

// Delete removes a workout exercise from the database
func (r *workoutExerciseRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var workoutExercise model.WorkoutExercise
		if err := tx.First(&workoutExercise, "id = ?", id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return customerrors.ErrEntityNotFound
			}
			return fmt.Errorf("failed to fetch workout exercise: %w", err)
		}

		if err := tx.Delete(&model.WorkoutExercise{}, "id = ?", id).Error; err != nil {
			return fmt.Errorf("failed to delete workout exercise: %w", err)
		}

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
			if err == gorm.ErrRecordNotFound {
				return customerrors.ErrEntityNotFound
			}
			return fmt.Errorf("failed to fetch workout exercise: %w", err)
		}

		var totalExercises int64
		if err := tx.Model(&model.WorkoutExercise{}).
			Where("workout_id = ?", workoutExercise.WorkoutID).
			Count(&totalExercises).Error; err != nil {
			return err
		}

		if newPosition < 1 || newPosition > int(totalExercises) {
			return customerrors.NewErrInvalidPosition(newPosition, totalExercises)
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
