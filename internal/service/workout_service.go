package service

import "github.com/VladimirKholomyanskyy/gym-api/internal/repository"

type WorkoutService struct {
	workoutRepository         *repository.WorkoutRepository
	workoutExerciseRepository *repository.WorkoutExerciseRepository
}

func NewWorkoutService(workoutRepository *repository.WorkoutRepository, workoutExerciseRepository *repository.WorkoutExerciseRepository) *WorkoutService {
	return &WorkoutService{workoutRepository: workoutRepository, workoutExerciseRepository: workoutExerciseRepository}
}
