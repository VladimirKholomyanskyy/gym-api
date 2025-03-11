package repository

import (
	"context"
	"errors"
	"fmt"

	customerrors "github.com/VladimirKholomyanskyy/gym-api/internal/customErrors"
	"github.com/VladimirKholomyanskyy/gym-api/internal/training/model"
	"gorm.io/gorm"
)

type WorkoutRepository interface {
	Create(ctx context.Context, workout *model.Workout) error
	GetByID(ctx context.Context, id string) (*model.Workout, error)
	GetAllByTrainingProgramID(ctx context.Context, id string, page, pageSize int) ([]model.Workout, int64, error)
	UpdatePartial(ctx context.Context, id string, updates map[string]any) (*model.Workout, error)
	Delete(ctx context.Context, id string) error
	PermanentDelete(ctx context.Context, id string) error
	Reorder(ctx context.Context, workoutID string, newPosition int) error
}

type workoutRepository struct {
	db *gorm.DB
}

func NewWorkoutRepository(db *gorm.DB) WorkoutRepository {
	return &workoutRepository{db: db}
}

// Create inserts a new workout, setting its position automatically.
func (r *workoutRepository) Create(ctx context.Context, workout *model.Workout) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var maxPosition int
		err := tx.Model(&model.Workout{}).
			Where("training_program_id = ?", workout.TrainingProgramID).
			Select("COALESCE(MAX(position), 0)").
			Scan(&maxPosition).Error
		if err != nil {
			return fmt.Errorf("failed to calculate max position: %w", err)
		}
		workout.Position = maxPosition + 1
		if err := tx.Create(workout).Error; err != nil {
			return fmt.Errorf("failed to create workout: %w", err)
		}
		return nil
	})
}

// FindByID retrieves a workout by ID with exercises preloaded.
func (r *workoutRepository) GetByID(ctx context.Context, id string) (*model.Workout, error) {
	var workout model.Workout
	err := r.db.WithContext(ctx).
		Preload("Exercises.Exercise").
		First(&workout, "id = ?", id).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, customerrors.ErrEntityNotFound
		}
		return nil, fmt.Errorf("failed to fetch workout: %w", err)
	}

	return &workout, nil
}

// FindByTrainingProgramID retrieves workouts for a training program with pagination and total count.
func (r *workoutRepository) GetAllByTrainingProgramID(
	ctx context.Context,
	programID string,
	page,
	pageSize int,
) ([]model.Workout, int64, error) {
	offset := (page - 1) * pageSize

	var workouts []model.Workout
	var totalCount int64

	countErr := r.db.WithContext(ctx).
		Model(&model.Workout{}).
		Where("training_program_id = ?", programID).
		Count(&totalCount).Error
	if countErr != nil {
		return nil, 0, fmt.Errorf("failed to count workouts: %w", countErr)
	}

	err := r.db.WithContext(ctx).
		Where("training_program_id = ?", programID).
		Order("position ASC").
		Preload("Exercises.Exercise").
		Limit(pageSize).
		Offset(offset).
		Find(&workouts).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch workouts: %w", err)
	}

	return workouts, totalCount, nil
}

// Update modifies an existing workout.
func (r *workoutRepository) UpdatePartial(ctx context.Context, id string, updates map[string]any) (*model.Workout, error) {
	result := r.db.WithContext(ctx).Model(&model.Workout{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to update workout: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, customerrors.ErrEntityNotFound
	}
	var updatedWorkout model.Workout
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&updatedWorkout).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch updated workout: %w", err)
	}

	return &updatedWorkout, nil
}

// Delete removes a workout and shifts positions of the remaining workouts.
func (r *workoutRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var workout model.Workout
		if err := tx.First(&workout, "id = ?", id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return customerrors.ErrEntityNotFound
			}
			return err
		}
		if err := tx.Delete(&model.Workout{}, "id = ?", id).Error; err != nil {
			return err
		}
		return tx.Model(&model.Workout{}).
			Where("training_program_id = ? AND position > ?", workout.TrainingProgramID, workout.Position).
			Update("position", gorm.Expr("position - 1")).Error
	})
}

func (r *workoutRepository) PermanentDelete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Unscoped().Where("id = ?", id).Delete(&model.Workout{})
	if result.Error != nil {
		return fmt.Errorf("failed to permanent delete profile: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return customerrors.ErrEntityNotFound
	}
	return nil
}

// Reorder changes the position of a workout within a training program.
func (r *workoutRepository) Reorder(ctx context.Context, workoutID string, newPosition int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Find the workout to get current position and training program ID
		var workout model.Workout
		if err := tx.First(&workout, "id = ?", workoutID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return customerrors.ErrEntityNotFound
			}
			return err
		}

		// Validate new position
		var totalWorkouts int64
		if err := tx.Model(&model.Workout{}).
			Where("training_program_id = ?", workout.TrainingProgramID).
			Count(&totalWorkouts).Error; err != nil {
			return err
		}

		// Validate position range
		if newPosition < 1 || newPosition > int(totalWorkouts) {
			return customerrors.NewErrInvalidPosition(newPosition, totalWorkouts)
		}

		// No change needed
		if workout.Position == newPosition {
			return nil
		}

		// Shift positions
		if workout.Position < newPosition {
			// Move workouts down
			if err := tx.Model(&model.Workout{}).
				Where("training_program_id = ? AND position > ? AND position <= ?",
					workout.TrainingProgramID, workout.Position, newPosition).
				Update("position", gorm.Expr("position - 1")).Error; err != nil {
				return err
			}
		} else {
			// Move workouts up
			if err := tx.Model(&model.Workout{}).
				Where("training_program_id = ? AND position < ? AND position >= ?",
					workout.TrainingProgramID, workout.Position, newPosition).
				Update("position", gorm.Expr("position + 1")).Error; err != nil {
				return err
			}
		}

		// Update the workout's position
		return tx.Model(&model.Workout{}).
			Where("id = ?", workoutID).
			Update("position", newPosition).Error
	})
}
