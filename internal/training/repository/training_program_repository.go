package repository

import (
	"context"
	"errors"
	"strings"

	"github.com/VladimirKholomyanskyy/gym-api/internal/training/model"
	"gorm.io/gorm"
)

var (
	// ErrTrainingProgramNotFound is returned when a training program cannot be found
	ErrTrainingProgramNotFound = errors.New("training program not found")

	// ErrTrainingProgramCreate is returned when there's an issue creating a training program
	ErrTrainingProgramCreate = errors.New("failed to create training program")

	// ErrTrainingProgramUpdate is returned when there's an issue updating a training program
	ErrTrainingProgramUpdate = errors.New("failed to update training program")
)

// TrainingProgramRepository defines the interface for training program operations
type TrainingProgramRepository interface {
	Create(ctx context.Context, trainingProgram *model.TrainingProgram) error
	FindByID(ctx context.Context, id string) (*model.TrainingProgram, error)
	FindByIDAndProfileID(ctx context.Context, programID, profileID string) (*model.TrainingProgram, error)
	FindByProfileID(ctx context.Context, profileID string, page, pageSize int) ([]model.TrainingProgram, int64, error)
	Update(ctx context.Context, trainingProgram *model.TrainingProgram) error
	SoftDelete(ctx context.Context, programID, profileID string) error
}

// trainingProgramRepository implements TrainingProgramRepository
type trainingProgramRepository struct {
	db *gorm.DB
}

// NewTrainingProgramRepository creates a new instance of the repository
func NewTrainingProgramRepository(db *gorm.DB) TrainingProgramRepository {
	return &trainingProgramRepository{db: db}
}

// validateTrainingProgram performs validation checks on the training program
func validateTrainingProgram(trainingProgram *model.TrainingProgram) error {
	// Trim and validate name
	trainingProgram.Name = strings.TrimSpace(trainingProgram.Name)
	if trainingProgram.Name == "" {
		return errors.New("training program name cannot be empty")
	}

	// Validate profile ID
	if strings.TrimSpace(trainingProgram.ProfileID) == "" {
		return errors.New("profile ID cannot be empty")
	}

	return nil
}

// Create inserts a new training program into the database
func (r *trainingProgramRepository) Create(ctx context.Context, trainingProgram *model.TrainingProgram) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Validate training program
		if err := validateTrainingProgram(trainingProgram); err != nil {
			return err
		}

		// Create the training program
		if err := tx.Create(trainingProgram).Error; err != nil {
			return errors.Join(ErrTrainingProgramCreate, err)
		}
		return nil
	})
}

// FindByID retrieves a training program by its ID
func (r *trainingProgramRepository) FindByID(ctx context.Context, id string) (*model.TrainingProgram, error) {
	var trainingProgram model.TrainingProgram
	result := r.db.WithContext(ctx).First(&trainingProgram, "id = ?", id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrTrainingProgramNotFound
		}
		return nil, result.Error
	}

	return &trainingProgram, nil
}

// FindByProfileID retrieves paginated training programs belonging to a user
func (r *trainingProgramRepository) FindByProfileID(ctx context.Context, profileID string, page, pageSize int) ([]model.TrainingProgram, int64, error) {
	if err := validatePagination(page, pageSize); err != nil {
		return nil, 0, err
	}

	var trainingPrograms []model.TrainingProgram
	var total int64
	offset := (page - 1) * pageSize

	// Count total records
	countQuery := r.db.WithContext(ctx).
		Model(&model.TrainingProgram{}).
		Where("profile_id = ?", profileID)
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch paginated results
	result := r.db.WithContext(ctx).
		Where("profile_id = ?", profileID).
		Limit(pageSize).
		Offset(offset).
		Order("created_at DESC").
		Find(&trainingPrograms)

	return trainingPrograms, total, result.Error
}

// FindByIDAndProfileID retrieves a training program by its ID and profile ID
func (r *trainingProgramRepository) FindByIDAndProfileID(ctx context.Context, programID, profileID string) (*model.TrainingProgram, error) {
	var trainingProgram model.TrainingProgram
	result := r.db.WithContext(ctx).
		Where("id = ? AND profile_id = ?", programID, profileID).
		First(&trainingProgram)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrTrainingProgramNotFound
		}
		return nil, result.Error
	}

	return &trainingProgram, nil
}

// Update modifies an existing training program
func (r *trainingProgramRepository) Update(ctx context.Context, trainingProgram *model.TrainingProgram) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Validate training program
		if err := validateTrainingProgram(trainingProgram); err != nil {
			return err
		}

		// Selective update of fields
		result := tx.Model(trainingProgram).
			Select("name", "description", "duration", "difficulty_level").
			Save(trainingProgram)

		if result.Error != nil {
			return errors.Join(ErrTrainingProgramUpdate, result.Error)
		}

		// Check if any rows were actually updated
		if result.RowsAffected == 0 {
			return ErrTrainingProgramNotFound
		}

		return nil
	})
}

// SoftDelete removes a training program, ensuring it belongs to the user
func (r *trainingProgramRepository) SoftDelete(ctx context.Context, programID, profileID string) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND profile_id = ?", programID, profileID).
		Delete(&model.TrainingProgram{})

	if result.Error != nil {
		return result.Error
	}

	// Check if any rows were actually deleted
	if result.RowsAffected == 0 {
		return ErrTrainingProgramNotFound
	}

	return nil
}
