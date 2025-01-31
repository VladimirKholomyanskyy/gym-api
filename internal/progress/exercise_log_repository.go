package progress

import (
	"time"

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
func (r *ExerciseLogRepository) GetAllByUserIDAndSessionID(userId, sessionId uint) ([]ExerciseLog, error) {
	var exerciseLogs []ExerciseLog
	err := r.DB.Where("user_id = ? AND session_id = ?", userId, sessionId).Find(&exerciseLogs).Error
	return exerciseLogs, err
}

func (r *ExerciseLogRepository) GetAllByUserIDAndExerciseID(userId, exerciseId uint) ([]ExerciseLog, error) {
	var exerciseLogs []ExerciseLog
	err := r.DB.Where("user_id = ? AND exercise_id = ?", userId, exerciseId).Find(&exerciseLogs).Error
	return exerciseLogs, err
}

func (r *ExerciseLogRepository) GetAllByUserID(userId uint) ([]ExerciseLog, error) {
	var exerciseLogs []ExerciseLog
	err := r.DB.Where("user_id = ?", userId).Find(&exerciseLogs).Error
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

func (r *ExerciseLogRepository) GetWeightPerDay(userID, exerciseID uint, startDate, endDate *time.Time) ([]WeightPerDay, error) {
	var results []WeightPerDay

	query := r.DB.Table("exercise_logs").
		Select("DATE(logged_at) AS date, SUM(weight * reps) AS total_weight").
		Where("user_id = ? AND exercise_id = ?", userID, exerciseID).
		Group("DATE(logged_at)").
		Order("date ASC")

	// Add date range filter if provided
	if startDate != nil {
		query = query.Where("logged_at >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("logged_at <= ?", *endDate)
	}

	// Execute the query and populate results
	err := query.Scan(&results).Error
	return results, err
}
