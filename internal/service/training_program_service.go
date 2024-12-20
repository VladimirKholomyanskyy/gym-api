package service

import (
	"errors"

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

func (s *TrainingProgramService) CreateTrainingProgram(input models.CreateTrainingProgramRequest, userID uint) (*models.TrainingProgram, error) {
	program := &models.TrainingProgram{
		Name:        input.Name,
		Description: input.Description,
		UserID:      userID,
	}
	if err := s.trainingProgramRepository.Create(program); err != nil {
		return nil, err
	}
	return program, nil
}

func (s *TrainingProgramService) GetAllTrainingPrograms(userId uint) ([]models.TrainingProgram, error) {
	programs, err := s.trainingProgramRepository.FindByUserID(userId)
	if err != nil {
		return nil, err
	}
	return programs, nil
}

func (s *TrainingProgramService) GetTrainingProgram(userId uint, programId uint) (*models.TrainingProgram, error) {
	program, err := s.trainingProgramRepository.FindByIDAndUserID(programId, userId)
	if err != nil {
		return nil, err
	}
	return program, nil
}

func (s *TrainingProgramService) DeleteTrainingProgram(userId, program_id uint) error {
	return s.trainingProgramRepository.Delete(program_id, userId)
}

func (s *TrainingProgramService) UpdateTrainingProgram(input models.CreateTrainingProgramRequest, userID uint, programID uint) (*models.TrainingProgram, error) {
	program, err := s.trainingProgramRepository.FindByIDAndUserID(programID, userID)
	if err != nil {
		return nil, err
	}
	if input.Name != "" {
		program.Name = input.Name
	}
	if input.Description != "" {
		program.Description = input.Description
	}

	err = s.trainingProgramRepository.Update(program)
	if err != nil {
		return nil, err
	}
	return program, nil
}

func (s *TrainingProgramService) AddWorkoutToProgram(input models.CreateWorkoutRequest, userID uint, programID uint) (*models.Workout, error) {
	_, err := s.trainingProgramRepository.FindByIDAndUserID(programID, userID)
	if err != nil {
		return nil, err
	}
	workout := &models.Workout{Name: input.Name, TrainingProgramID: programID}

	err = s.workoutRepository.Create(workout)
	if err != nil {
		return nil, err
	}
	return workout, nil
}

func (s *TrainingProgramService) UpdateWorkout(input models.CreateWorkoutRequest, userID uint, programID uint, workoutId uint) (*models.Workout, error) {
	_, err := s.trainingProgramRepository.FindByIDAndUserID(programID, userID)
	if err != nil {
		return nil, err
	}
	workout, err := s.workoutRepository.FindByID(workoutId)
	if err != nil {
		return nil, err
	}
	workout.Name = input.Name
	err = s.workoutRepository.Update(workout)
	if err != nil {
		return nil, err
	}
	return workout, nil
}

func (s *TrainingProgramService) GetAllWorkouts(userID, programID uint) ([]models.Workout, error) {
	_, err := s.trainingProgramRepository.FindByIDAndUserID(programID, userID)
	if err != nil {
		return nil, err
	}
	workouts, err := s.workoutRepository.FindByTrainingProgramID(programID)
	if err != nil {
		return nil, err
	}
	return workouts, err
}

func (s *TrainingProgramService) DeleteWorkout(userId, program_id, workoutID uint) error {
	return s.workoutRepository.Delete(workoutID)
}

func (s *TrainingProgramService) GetWorkout(userId, program_id, workoutID uint) (*models.Workout, error) {
	return s.workoutRepository.FindByID(workoutID)
}

func (s *TrainingProgramService) AddExerciseToWorkout(input models.CreateWorkoutExerciseRequest, userID uint) (*models.WorkoutExercise, error) {
	workout, err := s.workoutRepository.FindByID(input.WorkoutID)
	if err != nil || workout == nil {
		return nil, errors.New("workout not found")
	}
	program, err := s.trainingProgramRepository.FindByIDAndUserID(workout.TrainingProgramID, userID)
	if err != nil || program == nil {
		return nil, errors.New("workout does not belong to user's training program")
	}

	workoutExercise := &models.WorkoutExercise{
		WorkoutID:  input.WorkoutID,
		ExerciseID: input.ExerciseID,
		Sets:       input.Sets,
		Reps:       input.Reps,
		Weight:     input.Weight,
	}
	err = s.workoutExerciseRepository.Create(workoutExercise)
	if err != nil {
		return nil, err
	}
	return workoutExercise, nil
}

func (s *TrainingProgramService) GetAllWorkoutExercisesByWorkout(userID uint, workoutID uint) ([]models.WorkoutExercise, error) {
	workout, err := s.workoutRepository.FindByID(workoutID)
	if err != nil || workout == nil {
		return nil, errors.New("workout not found")
	}
	program, err := s.trainingProgramRepository.FindByIDAndUserID(workout.TrainingProgramID, userID)
	if err != nil || program == nil {
		return nil, errors.New("workout does not belong to user's training program")
	}
	workoutExercises, err := s.workoutExerciseRepository.FindByWorkoutID(workoutID)
	if err != nil {
		return nil, err
	}
	return workoutExercises, nil
}

func (s *TrainingProgramService) UpdateWorkoutExercise(input models.CreateWorkoutExerciseRequest, userID uint, weID uint) (*models.WorkoutExercise, error) {
	workout, err := s.workoutRepository.FindByID(input.WorkoutID)
	if err != nil || workout == nil {
		return nil, errors.New("workout not found")
	}
	program, err := s.trainingProgramRepository.FindByIDAndUserID(workout.TrainingProgramID, userID)
	if err != nil || program == nil {
		return nil, errors.New("workout does not belong to user's training program")
	}
	workoutExercise, err := s.workoutExerciseRepository.FindByID(weID)
	if err != nil {
		return nil, err
	}
	if input.ExerciseID > 0 {
		workoutExercise.ExerciseID = input.ExerciseID
	}
	if input.Sets > 0 {
		workoutExercise.Sets = input.Sets
	}
	if input.Reps > 0 {
		workoutExercise.Reps = input.Reps
	}
	err = s.workoutExerciseRepository.Update(workoutExercise)
	if err != nil {
		return nil, err
	}
	return workoutExercise, nil
}
func (s *TrainingProgramService) DeleteWorkoutExercise(userID uint, weID uint) error {
	workoutExercise, err := s.workoutExerciseRepository.FindByID(weID)
	if err != nil {
		return err
	}
	workout, err := s.workoutRepository.FindByID(workoutExercise.WorkoutID)
	if err != nil {
		return err
	}
	_, err = s.trainingProgramRepository.FindByIDAndUserID(workout.TrainingProgramID, userID)
	if err != nil {
		return err
	}
	return s.workoutExerciseRepository.Delete(weID)

}
func (s *TrainingProgramService) GetAllWorkoutExercisesByExercise(userID uint, exerciseID uint) {

}
