package service

import (
	"github.com/VladimirKholomyanskyy/gym-api/internal/models"
	"github.com/VladimirKholomyanskyy/gym-api/internal/repository"
)

type TrainingProgramService struct {
	trainingProgramRepository *repository.TrainingProgramRepository
	workoutRepository         *repository.WorkoutRepository
	workoutExerciseRepository *repository.WorkoutExerciseRepository
}

func NewTrainingProgramService(trainingProgramRepository *repository.TrainingProgramRepository,
	workoutRepository *repository.WorkoutRepository,
	workoutExerciseRepository *repository.WorkoutExerciseRepository,
) *TrainingProgramService {
	return &TrainingProgramService{
		trainingProgramRepository: trainingProgramRepository,
		workoutRepository:         workoutRepository,
		workoutExerciseRepository: workoutExerciseRepository}
}

func (s *TrainingProgramService) CreateTrainingProgram(input models.TrainingProgramInput) (*models.TrainingProgram, error) {
	program := &models.TrainingProgram{
		Name:        input.Name,
		Description: input.Description,
		UserID:      input.UserID,
	}
	if err := s.trainingProgramRepository.Create(program); err != nil {
		return nil, err
	}
	return program, nil
}

func (s *TrainingProgramService) GetAllTrainingPrograms(userId uint) ([]models.TrainingProgram, error) {
	programs, err := s.trainingProgramRepository.GetByUserID(userId)
	if err != nil {
		return nil, err
	}
	return programs, nil
}

func (s *TrainingProgramService) DeleteTrainingProgram(userId, program_id uint) error {
	return s.trainingProgramRepository.Delete(program_id, userId)
}

func (s *TrainingProgramService) UpdateTrainingProgram(input models.TrainingProgramInput) (*models.TrainingProgram, error) {
	program := &models.TrainingProgram{
		Name:        input.Name,
		Description: input.Description,
		UserID:      input.UserID,
	}
	err := s.trainingProgramRepository.Update(program)
	if err != nil {
		return nil, err
	}
	return program, nil
}

func (s *TrainingProgramService) AddWorkoutToProgram(input models.WorkoutInput) (*models.Workout, error) {
	_, err := s.trainingProgramRepository.GetByID(input.UserID, input.TrainingProgramID)
	if err != nil {
		return nil, err
	}
	workout := &models.Workout{Name: input.Name, TrainingProgramID: input.TrainingProgramID}

	err = s.workoutRepository.Create(workout)
	if err != nil {
		return nil, err
	}
	return workout, nil
}
