package repository

import (
	"github.com/VladimirKholomyanskyy/gym-api/internal/models"
	"gorm.io/gorm"
)

type WorkoutSessionRepository struct {
	DB *gorm.DB
}

func NewWorkoutLogRepository(db *gorm.DB) *WorkoutSessionRepository {
	return &WorkoutSessionRepository{DB: db}
}

// Create a new workout log
func (r *WorkoutSessionRepository) Create(workoutSession *models.WorkoutSession) error {
	return r.DB.Create(workoutSession).Error
}

// Get a workout log by ID
func (r *WorkoutSessionRepository) GetByID(id uint) (*models.WorkoutSession, error) {
	var workoutSession models.WorkoutSession
	err := r.DB.First(&workoutSession, id).Error
	return &workoutSession, err
}

// Get all workout logs for a user
func (r *WorkoutSessionRepository) GetAllByUserID(userID uint) ([]models.WorkoutSession, error) {
	var workoutSession []models.WorkoutSession
	err := r.DB.Where("user_id = ?", userID).Find(&workoutSession).Error
	return workoutSession, err
}

// Update a workout log
func (r *WorkoutSessionRepository) Update(workoutSession *models.WorkoutSession) error {
	return r.DB.Save(workoutSession).Error
}

// Delete a workout log
func (r *WorkoutSessionRepository) Delete(id uint) error {
	return r.DB.Delete(&models.WorkoutSession{}, id).Error
}
