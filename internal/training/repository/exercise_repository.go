package repository

import (
	"context"
	"errors"
	"strings"

	"github.com/VladimirKholomyanskyy/gym-api/internal/training/model"
	"gorm.io/gorm"
)

var (
	// ErrExerciseNotFound is returned when an exercise cannot be found
	ErrExerciseNotFound = errors.New("exercise not found")

	// ErrExerciseCreate is returned when there's an issue creating an exercise
	ErrExerciseCreate = errors.New("failed to create exercise")

	// ErrInvalidMuscleGroup is returned when an invalid muscle group is provided
	ErrInvalidMuscleGroup = errors.New("invalid muscle group")
)

// ExerciseRepository defines the interface for exercise operations
type ExerciseRepository interface {
	Create(ctx context.Context, exercise *model.Exercise) error
	FindAll(ctx context.Context, page, pageSize int) ([]model.Exercise, int64, error)
	FindByID(ctx context.Context, id string) (*model.Exercise, error)
	FindByPrimaryMuscle(ctx context.Context, primaryMuscle string, page, pageSize int) ([]model.Exercise, int64, error)
	Update(ctx context.Context, exercise *model.Exercise) error
	SoftDelete(ctx context.Context, id string) error
}

// exerciseRepository implements ExerciseRepository
type exerciseRepository struct {
	db *gorm.DB
}

// NewExerciseRepository creates a new instance of ExerciseRepository
func NewExerciseRepository(db *gorm.DB) ExerciseRepository {
	return &exerciseRepository{db: db}
}

// Create inserts a new exercise into the database
func (r *exerciseRepository) Create(ctx context.Context, exercise *model.Exercise) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Validate exercise before creation
		if err := validateExercise(exercise); err != nil {
			return err
		}

		// Create the exercise
		if err := tx.Create(exercise).Error; err != nil {
			return errors.Join(ErrExerciseCreate, err)
		}
		return nil
	})
}

// validateExercise performs validation checks on the exercise
func validateExercise(exercise *model.Exercise) error {
	// Trim and validate name
	exercise.Name = strings.TrimSpace(exercise.Name)
	if exercise.Name == "" {
		return errors.New("exercise name cannot be empty")
	}

	// Validate primary muscle
	exercise.PrimaryMuscle = strings.TrimSpace(exercise.PrimaryMuscle)
	if exercise.PrimaryMuscle == "" {
		return ErrInvalidMuscleGroup
	}

	return nil
}

// FindAll retrieves paginated exercises from the database
func (r *exerciseRepository) FindAll(ctx context.Context, page, pageSize int) ([]model.Exercise, int64, error) {
	var exercises []model.Exercise
	var total int64
	offset := (page - 1) * pageSize

	// Count total records
	countQuery := r.db.WithContext(ctx).Model(&model.Exercise{})
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch paginated results
	result := r.db.WithContext(ctx).
		Limit(pageSize).
		Offset(offset).
		Order("name ASC").
		Find(&exercises)

	return exercises, total, result.Error
}

// FindByID retrieves an exercise by its ID
func (r *exerciseRepository) FindByID(ctx context.Context, id string) (*model.Exercise, error) {
	var exercise model.Exercise
	result := r.db.WithContext(ctx).First(&exercise, "id = ?", id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrExerciseNotFound
		}
		return nil, result.Error
	}

	return &exercise, nil
}

// FindByPrimaryMuscle retrieves paginated exercises based on the primary muscle group
func (r *exerciseRepository) FindByPrimaryMuscle(ctx context.Context, primaryMuscle string, page, pageSize int) ([]model.Exercise, int64, error) {
	if err := validatePagination(page, pageSize); err != nil {
		return nil, 0, err
	}

	// Normalize muscle group
	primaryMuscle = strings.TrimSpace(primaryMuscle)
	if primaryMuscle == "" {
		return nil, 0, ErrInvalidMuscleGroup
	}

	var exercises []model.Exercise
	var total int64
	offset := (page - 1) * pageSize

	// Count total records
	countQuery := r.db.WithContext(ctx).
		Model(&model.Exercise{}).
		Where("primary_muscle = ?", primaryMuscle)
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch paginated results
	result := r.db.WithContext(ctx).
		Where("primary_muscle = ?", primaryMuscle).
		Limit(pageSize).
		Offset(offset).
		Order("name ASC").
		Find(&exercises)

	return exercises, total, result.Error
}

// Update updates an existing exercise
func (r *exerciseRepository) Update(ctx context.Context, exercise *model.Exercise) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Validate exercise before update
		if err := validateExercise(exercise); err != nil {
			return err
		}

		// Selective update of fields
		result := tx.Model(exercise).
			Select("name", "description", "primary_muscle", "secondary_muscles", "equipment").
			Save(exercise)

		if result.Error != nil {
			return result.Error
		}

		// Check if any rows were actually updated
		if result.RowsAffected == 0 {
			return ErrExerciseNotFound
		}

		return nil
	})
}

// SoftDelete marks an exercise as deleted without removing it from the database
func (r *exerciseRepository) SoftDelete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).
		Delete(&model.Exercise{}, id)

	if result.Error != nil {
		return result.Error
	}

	// Check if any rows were actually deleted
	if result.RowsAffected == 0 {
		return ErrExerciseNotFound
	}

	return nil
}
