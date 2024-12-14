package repository

import (
	"github.com/VladimirKholomyanskyy/gym-api/internal/models"
	"gorm.io/gorm"
)

type ExerciseRepository struct {
	db *gorm.DB
}

func NewExerciseRepository(db *gorm.DB) *ExerciseRepository {
	return &ExerciseRepository{db: db}
}

func (r *ExerciseRepository) Create(exercise *models.Exercise) error {
	return r.db.Create(exercise).Error
}
func (r *ExerciseRepository) FindAll() ([]models.Exercise, error) {
	var exercises []models.Exercise
	err := r.db.Find(&exercises).Error
	return exercises, err
}

func (r *ExerciseRepository) FindByID(id uint) (*models.Exercise, error) {
	var exercise models.Exercise
	err := r.db.First(&exercise, id).Error
	return &exercise, err
}

func (r *ExerciseRepository) FindByPrimaryMuscle(primaryMuscle string) ([]models.Exercise, error) {
	var exercises []models.Exercise
	if err := r.db.Where("primary_muscle = ?", primaryMuscle).Find(&exercises).Error; err != nil {
		return nil, err
	}
	return exercises, nil
}
