package progress

import (
	"gorm.io/gorm"
)

type ExerciseLogRepository struct {
	DB *gorm.DB
}

func NewExerciseLogRepository(db *gorm.DB) *ExerciseLogRepository {
	return &ExerciseLogRepository{DB: db}
}

// Create a new exercise log
func (r *ExerciseLogRepository) Create(exerciseLog *ExerciseLog) error {
	return r.DB.Create(exerciseLog).Error
}

// Get an exercise log by ID
func (r *ExerciseLogRepository) GetByID(id uint) (*ExerciseLog, error) {
	var exerciseLog ExerciseLog
	err := r.DB.First(&exerciseLog, id).Error
	return &exerciseLog, err
}

// Get all exercise logs for a workout log
func (r *ExerciseLogRepository) GetAllByWorkoutLogID(workoutLogID uint) ([]ExerciseLog, error) {
	var exerciseLogs []ExerciseLog
	err := r.DB.Where("session_id = ?", workoutLogID).Find(&exerciseLogs).Error
	return exerciseLogs, err
}

// Update an exercise log
func (r *ExerciseLogRepository) Update(exerciseLog *ExerciseLog) error {
	return r.DB.Save(exerciseLog).Error
}

// Delete an exercise log
func (r *ExerciseLogRepository) Delete(id uint) error {
	return r.DB.Delete(&ExerciseLog{}, id).Error
}
