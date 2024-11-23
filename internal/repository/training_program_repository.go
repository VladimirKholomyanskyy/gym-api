package repository

import (
	"github.com/VladimirKholomyanskyy/gym-api/internal/models"
	"gorm.io/gorm"
)

type TrainingProgramRepository struct {
	db *gorm.DB
}

func NewTrainingProgramRepository(db *gorm.DB) *TrainingProgramRepository {
	return &TrainingProgramRepository{db: db}
}

func (r *TrainingProgramRepository) Create(trainingProgram *models.TrainingProgram) error {
	return r.db.Create(trainingProgram).Error
}

func (r *TrainingProgramRepository) FindByIDAndUserID(trainingProgramID, userID uint) (*models.TrainingProgram, error) {
	var trainingProgram models.TrainingProgram
	err := r.db.Where("user_id = ?", userID).First(&trainingProgram, trainingProgramID).Error
	if err != nil {
		return nil, err
	}
	return &trainingProgram, nil
}

// FindByUserID retrieves all training programs for a specific user
func (r *TrainingProgramRepository) FindByUserID(userID uint) ([]models.TrainingProgram, error) {
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
