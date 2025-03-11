package repository

import (
	"context"
	"errors"
	"fmt"

	customerrors "github.com/VladimirKholomyanskyy/gym-api/internal/customErrors"
	"github.com/VladimirKholomyanskyy/gym-api/internal/training/model"
	"gorm.io/gorm"
)

// TrainingProgramRepository defines the interface for training program operations
type TrainingProgramRepository interface {
	Create(ctx context.Context, trainingProgram *model.TrainingProgram) error
	FindByID(ctx context.Context, id string) (*model.TrainingProgram, error)
	FindByIDAndProfileID(ctx context.Context, programID, profileID string) (*model.TrainingProgram, error)
	FindByProfileID(ctx context.Context, profileID string, page, pageSize int) ([]model.TrainingProgram, int64, error)
	UpdatePartial(ctx context.Context, id string, updates map[string]any) (*model.TrainingProgram, error)
	Delete(ctx context.Context, id string) error
	PermanentDelete(ctx context.Context, id string) error
}

// trainingProgramRepository implements TrainingProgramRepository
type trainingProgramRepository struct {
	db *gorm.DB
}

// NewTrainingProgramRepository creates a new instance of the repository
func NewTrainingProgramRepository(db *gorm.DB) TrainingProgramRepository {
	return &trainingProgramRepository{db: db}
}

// Create inserts a new training program into the database
func (r *trainingProgramRepository) Create(ctx context.Context, trainingProgram *model.TrainingProgram) error {
	if err := r.db.WithContext(ctx).Create(trainingProgram).Error; err != nil {
		return fmt.Errorf("failed to create training program: %w", err)
	}
	return nil
}

// FindByID retrieves a training program by its ID
func (r *trainingProgramRepository) FindByID(ctx context.Context, id string) (*model.TrainingProgram, error) {
	var trainingProgram model.TrainingProgram
	err := r.db.WithContext(ctx).First(&trainingProgram, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, customerrors.ErrEntityNotFound
		}
		return nil, fmt.Errorf("failed to fetch training program: %w", err)
	}
	return &trainingProgram, nil
}

// FindByProfileID retrieves paginated training programs belonging to a user
func (r *trainingProgramRepository) FindByProfileID(ctx context.Context, profileID string, page, pageSize int) ([]model.TrainingProgram, int64, error) {
	var trainingPrograms []model.TrainingProgram
	var total int64
	offset := (page - 1) * pageSize

	// Count total records
	countQuery := r.db.WithContext(ctx).
		Model(&model.TrainingProgram{}).
		Where("profile_id = ?", profileID)
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count training programs: %w", err)
	}

	// Fetch paginated results
	err := r.db.WithContext(ctx).
		Where("profile_id = ?", profileID).
		Limit(pageSize).
		Offset(offset).
		Order("created_at DESC").
		Find(&trainingPrograms).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch training programs: %w", err)
	}

	return trainingPrograms, total, nil
}

// FindByIDAndProfileID retrieves a training program by its ID and profile ID
func (r *trainingProgramRepository) FindByIDAndProfileID(ctx context.Context, programID, profileID string) (*model.TrainingProgram, error) {
	var trainingProgram model.TrainingProgram
	result := r.db.WithContext(ctx).
		Where("id = ? AND profile_id = ?", programID, profileID).
		First(&trainingProgram)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, customerrors.ErrEntityNotFound
		}
		return nil, fmt.Errorf("failed to fetch training program: %w", result.Error)
	}

	return &trainingProgram, nil
}

func (r *trainingProgramRepository) UpdatePartial(ctx context.Context, id string, updates map[string]any) (*model.TrainingProgram, error) {
	result := r.db.WithContext(ctx).Model(&model.TrainingProgram{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to update training program: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, customerrors.ErrEntityNotFound
	}
	var updatedProgram model.TrainingProgram
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&updatedProgram).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch updated training program: %w", err)
	}

	return &updatedProgram, nil
}

// SoftDelete removes a training program, ensuring it belongs to the user
func (r *trainingProgramRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&model.TrainingProgram{})

	if result.Error != nil {
		return fmt.Errorf("failed to delete profile: %w", result.Error)
	}

	// Check if any rows were actually deleted
	if result.RowsAffected == 0 {
		return customerrors.ErrEntityNotFound
	}

	return nil
}

func (r *trainingProgramRepository) PermanentDelete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Unscoped().Where("id = ?", id).Delete(&model.TrainingProgram{})
	if result.Error != nil {
		return fmt.Errorf("failed to permanently delete profile: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return customerrors.ErrEntityNotFound
	}
	return nil
}
