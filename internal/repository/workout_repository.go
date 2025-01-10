package repository

import (
	"github.com/VladimirKholomyanskyy/gym-api/internal/models"
	"gorm.io/gorm"
)

type WorkoutRepository struct {
	db *gorm.DB
}

func NewWorkoutRepository(db *gorm.DB) *WorkoutRepository {
	return &WorkoutRepository{db: db}
}

// Create adds a new Workout to the database
func (r *WorkoutRepository) Create(workout *models.Workout) error {
	return r.db.Create(workout).Error
}

// Get retrieves a Workout by ID
func (r *WorkoutRepository) FindByID(id uint) (*models.Workout, error) {
	var workout models.Workout
	err := r.db.Debug().Preload("Exercises.Exercise").First(&workout, id).Error
	return &workout, err
}

func (r *WorkoutRepository) FindByTrainingProgramID(trainingProgramID uint) ([]models.Workout, error) {
	var workouts []models.Workout
	err := r.db.Where("training_program_id = ?", trainingProgramID).
		Preload("Exercises").
		Find(&workouts).Error
	return workouts, err
}

// Update updates an existing Workout
func (r *WorkoutRepository) Update(workout *models.Workout) error {
	return r.db.Save(workout).Error
}

// Delete removes a Workout by ID
func (r *WorkoutRepository) Delete(id uint) error {
	return r.db.Delete(&models.Workout{}, id).Error
}
