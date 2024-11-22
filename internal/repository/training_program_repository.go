package repository

import (
	"github.com/VladimirKholomyanskyy/gym-api/internal/models"
	"gorm.io/gorm"
)

type TrainingProgramRepository struct {
	db *gorm.DB
}

// NewTrainingProgramRepository creates a new instance of the repository
func NewTrainingProgramRepository(db *gorm.DB) *TrainingProgramRepository {
	return &TrainingProgramRepository{db: db}
}

// Create adds a new training program
func (r *TrainingProgramRepository) Create(trainingProgram *models.TrainingProgram) error {
	return r.db.Create(trainingProgram).Error
}

// GetByID retrieves a training program by its ID
func (r *TrainingProgramRepository) GetByID(userID, trainingProgramID uint) (*models.TrainingProgram, error) {
	var trainingProgram models.TrainingProgram
	err := r.db.Where("user_id = ?", userID).First(&trainingProgram, trainingProgramID).Error
	if err != nil {
		return nil, err
	}
	return &trainingProgram, nil
}

// GetAll retrieves all training programs
func (r *TrainingProgramRepository) GetAll() ([]models.TrainingProgram, error) {
	var trainingPrograms []models.TrainingProgram
	err := r.db.Find(&trainingPrograms).Error
	if err != nil {
		return nil, err
	}
	return trainingPrograms, nil
}

// GetByUserID retrieves all training programs for a specific user
func (r *TrainingProgramRepository) GetByUserID(userID uint) ([]models.TrainingProgram, error) {
	var trainingPrograms []models.TrainingProgram
	err := r.db.Where("user_id = ?", userID).Find(&trainingPrograms).Error
	if err != nil {
		return nil, err
	}
	return trainingPrograms, nil
}

// Update modifies an existing training program
func (r *TrainingProgramRepository) Update(trainingProgram *models.TrainingProgram) error {
	return r.db.Save(trainingProgram).Error
}

// Delete removes a training program by its ID
func (r *TrainingProgramRepository) Delete(program_id uint, user_id uint) error {
	return r.db.Where("user_id = ?", user_id).Delete(&models.TrainingProgram{}, program_id).Error
}
