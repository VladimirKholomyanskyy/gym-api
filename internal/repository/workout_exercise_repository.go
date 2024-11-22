package repository

import (
	"github.com/VladimirKholomyanskyy/gym-api/internal/models"
	"gorm.io/gorm"
)

type WorkoutExerciseRepository struct {
	db *gorm.DB
}

// NewWorkoutExerciseRepository creates a new instance of WorkoutExerciseRepository
func NewWorkoutExerciseRepository(db *gorm.DB) *WorkoutExerciseRepository {
	return &WorkoutExerciseRepository{db: db}
}

// Create adds a new WorkoutExercise to the database
func (r *WorkoutExerciseRepository) Create(exercise *models.WorkoutExercise) error {
	return r.db.Create(exercise).Error
}

// Get retrieves a WorkoutExercise by ID
func (r *WorkoutExerciseRepository) Get(id uint) (*models.WorkoutExercise, error) {
	var exercise models.WorkoutExercise
	err := r.db.First(&exercise, id).Error
	return &exercise, err
}

// GetByWorkout retrieves all WorkoutExercises for a specific Workout
func (r *WorkoutExerciseRepository) GetByWorkout(workoutID uint) ([]models.WorkoutExercise, error) {
	var exercises []models.WorkoutExercise
	err := r.db.Where("workout_id = ?", workoutID).Find(&exercises).Error
	return exercises, err
}

// Update updates an existing WorkoutExercise
func (r *WorkoutExerciseRepository) Update(exercise *models.WorkoutExercise) error {
	return r.db.Save(exercise).Error
}

// Delete removes a WorkoutExercise by ID
func (r *WorkoutExerciseRepository) Delete(id uint) error {
	return r.db.Delete(&models.WorkoutExercise{}, id).Error
}
