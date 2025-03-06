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
	ErrWorkoutNotFound      = errors.New("workout not found")
	ErrInvalidPosition      = errors.New("invalid workout position")
	ErrNilWorkout           = errors.New("workout cannot be nil")
	ErrEmptyTrainingProgram = errors.New("training program ID cannot be empty")
)

type WorkoutRepository interface {
	Create(ctx context.Context, workout *model.Workout) error
	FindByID(ctx context.Context, id string) (*model.Workout, error)
	FindByTrainingProgramID(ctx context.Context, id string, page, pageSize int) ([]model.Workout, int64, error)
	Update(ctx context.Context, workout *model.Workout) error
	Delete(ctx context.Context, id string) error
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
	// Validate input
	if workout == nil {
		return ErrNilWorkout
	}
	if workout.TrainingProgramID == "" {
		return ErrEmptyTrainingProgram
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Find the max position in the training program
		var maxPosition int
		err := tx.Model(&model.Workout{}).
			Where("training_program_id = ?", workout.TrainingProgramID).
			Select("COALESCE(MAX(position), 0)").
			Scan(&maxPosition).Error
		if err != nil {
			return fmt.Errorf("failed to calculate max position: %w", err)
		}

		// Assign the next available position
		workout.Position = maxPosition + 1

		// Create workout
		return tx.Create(workout).Error
	})
}

// FindByID retrieves a workout by ID with exercises preloaded.
func (r *workoutRepository) FindByID(ctx context.Context, id string) (*model.Workout, error) {
	var workout model.Workout
	result := r.db.WithContext(ctx).
		Preload("Exercises.Exercise").
		First(&workout, "id = ?", id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrWorkoutNotFound
		}
		return nil, result.Error
	}

	return &workout, nil
}

// FindByTrainingProgramID retrieves workouts for a training program with pagination and total count.
func (r *workoutRepository) FindByTrainingProgramID(
	ctx context.Context,
	programID string,
	page,
	pageSize int,
) ([]model.Workout, int64, error) {
	// Validate and normalize pagination parameters
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10 // Default page size
	}
	offset := (page - 1) * pageSize

	var workouts []model.Workout
	var totalCount int64

	// Count total workouts
	countErr := r.db.WithContext(ctx).
		Model(&model.Workout{}).
		Where("training_program_id = ?", programID).
		Count(&totalCount).Error
	if countErr != nil {
		return nil, 0, countErr
	}

	// Fetch paginated workouts
	err := r.db.WithContext(ctx).
		Where("training_program_id = ?", programID).
		Order("position ASC").
		Preload("Exercises.Exercise").
		Limit(pageSize).
		Offset(offset).
		Find(&workouts).Error

	if err != nil {
		return nil, 0, err
	}

	return workouts, totalCount, nil
}

// Update modifies an existing workout.
func (r *workoutRepository) Update(ctx context.Context, workout *model.Workout) error {
	// Validate input
	if workout == nil {
		return ErrNilWorkout
	}
	if workout.TrainingProgramID == "" {
		return ErrEmptyTrainingProgram
	}

	return r.db.WithContext(ctx).Save(workout).Error
}

// Delete removes a workout and shifts positions of the remaining workouts.
func (r *workoutRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Find the workout to get the training program ID and position
		var workout model.Workout
		if err := tx.First(&workout, "id = ?", id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrWorkoutNotFound
			}
			return err
		}

		// Delete the workout
		if err := tx.Delete(&model.Workout{}, "id = ?", id).Error; err != nil {
			return err
		}

		// Shift positions of remaining workouts in the same program
		return tx.Model(&model.Workout{}).
			Where("training_program_id = ? AND position > ?", workout.TrainingProgramID, workout.Position).
			Update("position", gorm.Expr("position - 1")).Error
	})
}

// Reorder changes the position of a workout within a training program.
func (r *workoutRepository) Reorder(ctx context.Context, workoutID string, newPosition int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Find the workout to get current position and training program ID
		var workout model.Workout
		if err := tx.First(&workout, "id = ?", workoutID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrWorkoutNotFound
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
			return fmt.Errorf("%w: position must be between 1 and %d", ErrInvalidPosition, totalWorkouts)
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
