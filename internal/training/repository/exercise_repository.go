package repository

import (
	"context"
	"errors"
	"fmt"

	customerrors "github.com/VladimirKholomyanskyy/gym-api/internal/customErrors"
	"github.com/VladimirKholomyanskyy/gym-api/internal/training/model"
	"gorm.io/gorm"
)

// ExerciseRepository defines the interface for exercise operations
type ExerciseRepository interface {
	Create(ctx context.Context, exercise *model.Exercise) error
	FindAll(ctx context.Context, page, pageSize int) ([]model.Exercise, int64, error)
	FindByID(ctx context.Context, id string) (*model.Exercise, error)
	FindByPrimaryMuscle(ctx context.Context, primaryMuscle string, page, pageSize int) ([]model.Exercise, int64, error)
	UpdatePartial(ctx context.Context, id string, updates map[string]any) error
	Delete(ctx context.Context, id string) error
	PermanentDelete(ctx context.Context, id string) error
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
	if err := r.db.WithContext(ctx).Create(exercise).Error; err != nil {
		return fmt.Errorf("failed to create exercise: %w", err)
	}
	return nil
}

// FindAll retrieves paginated exercises from the database
func (r *exerciseRepository) FindAll(ctx context.Context, page, pageSize int) ([]model.Exercise, int64, error) {
	var exercises []model.Exercise
	var total int64
	offset := (page - 1) * pageSize
	countQuery := r.db.WithContext(ctx).Model(&model.Exercise{})
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count exercises: %w", err)
	}
	result := r.db.WithContext(ctx).
		Limit(pageSize).
		Offset(offset).
		Order("name ASC").
		Find(&exercises)
	if result.Error != nil {
		return nil, 0, fmt.Errorf("failed to fetch exercises: %w", result.Error)
	}
	return exercises, total, nil
}

// FindByID retrieves an exercise by its ID
func (r *exerciseRepository) FindByID(ctx context.Context, id string) (*model.Exercise, error) {
	var exercise model.Exercise
	result := r.db.WithContext(ctx).First(&exercise, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, customerrors.ErrEntityNotFound
		}
		return nil, fmt.Errorf("failed to fetch exercise by id: %w", result.Error)
	}
	return &exercise, nil
}

// FindByPrimaryMuscle retrieves paginated exercises based on the primary muscle group
func (r *exerciseRepository) FindByPrimaryMuscle(ctx context.Context, primaryMuscle string, page, pageSize int) ([]model.Exercise, int64, error) {
	var exercises []model.Exercise
	var total int64
	offset := (page - 1) * pageSize
	countQuery := r.db.WithContext(ctx).
		Model(&model.Exercise{}).
		Where("primary_muscle = ?", primaryMuscle)
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count exercises: %w", err)
	}
	result := r.db.WithContext(ctx).
		Where("primary_muscle = ?", primaryMuscle).
		Limit(pageSize).
		Offset(offset).
		Order("name ASC").
		Find(&exercises)
	if result.Error != nil {
		return nil, 0, fmt.Errorf("failed to fetch exercises by primary muscle: %w", result.Error)
	}
	return exercises, total, nil
}

// Update updates an existing exercise
func (r *exerciseRepository) UpdatePartial(ctx context.Context, id string, updates map[string]any) error {
	result := r.db.WithContext(ctx).Model(&model.Exercise{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to update exercise: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return customerrors.ErrEntityNotFound
	}
	return nil
}

// SoftDelete marks an exercise as deleted without removing it from the database
func (r *exerciseRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&model.Exercise{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete profile: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return customerrors.ErrEntityNotFound
	}
	return nil
}

func (r *exerciseRepository) PermanentDelete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Unscoped().Where("id = ?", id).Delete(&model.Exercise{})
	if result.Error != nil {
		return fmt.Errorf("failed to permanent delete profile: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return customerrors.ErrEntityNotFound
	}
	return nil
}
