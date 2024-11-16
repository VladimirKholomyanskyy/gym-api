package repository

import (
	"github.com/VladimirKholomyanskyy/gym-api/internal/models"
	"gorm.io/gorm"
)

type ExerciseRepository struct {
	DB *gorm.DB
}

func (r *ExerciseRepository) GetAllExercises() ([]models.Exercise, error) {
	var exercises []models.Exercise
	err := r.DB.Find(&exercises).Error
	return exercises, err
}

func (r *ExerciseRepository) GetExercisesByPrimaryMuscle(primaryMuscle string) ([]models.Exercise, error) {
	var exercises []models.Exercise
	if err := r.DB.Where("primary_muscle = ?", primaryMuscle).Find(&exercises).Error; err != nil {
		return nil, err
	}
	return exercises, nil
}
