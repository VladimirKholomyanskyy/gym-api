package training

import (
	"errors"
	"fmt"
	"strconv"

	openapi "github.com/VladimirKholomyanskyy/gym-api/internal/api/go"
	"github.com/VladimirKholomyanskyy/gym-api/internal/common"
)

type TrainingManager struct {
	trainingProgramRepository *TrainingProgramRepository
	workoutRepository         *WorkoutRepository
	workoutExerciseRepository *WorkoutExerciseRepository
	exerciseRepository        *ExerciseRepository
}

func NewTrainingManager(
	trainingProgramRepository *TrainingProgramRepository,
	workoutRepository *WorkoutRepository,
	workoutExerciseRepository *WorkoutExerciseRepository,
	exerciseRepository *ExerciseRepository,
) *TrainingManager {
	return &TrainingManager{
		trainingProgramRepository: trainingProgramRepository,
		workoutRepository:         workoutRepository,
		workoutExerciseRepository: workoutExerciseRepository,
		exerciseRepository:        exerciseRepository,
	}
}
func (s *TrainingManager) CreateTrainingProgram(input openapi.CreateTrainingProgramRequest, userID uint) (*TrainingProgram, error) {
	program := &TrainingProgram{
		Name:        input.Name,
		Description: input.Description,
		UserID:      userID,
	}
	if err := s.trainingProgramRepository.Create(program); err != nil {
		return nil, err
	}
	return program, nil
}

func (s *TrainingManager) GetTrainingPrograms(userId uint) ([]TrainingProgram, error) {
	programs, err := s.trainingProgramRepository.FindByUserId(userId)
	if err != nil {
		return nil, err
	}
	return programs, nil
}

func (s *TrainingManager) GetTrainingProgram(userId, programId uint) (*TrainingProgram, error) {
	program, err := s.trainingProgramRepository.FindById(programId)
	if err != nil {
		return nil, err
	}
	return program, nil
}

func (s *TrainingManager) DeleteTrainingProgram(userId, program_id uint) error {
	return s.trainingProgramRepository.Delete(program_id, userId)
}

func (s *TrainingManager) UpdateTrainingProgram(input openapi.CreateTrainingProgramRequest, userID uint, programID uint) (*TrainingProgram, error) {
	program, err := s.trainingProgramRepository.FindById(programID)
	if err != nil {
		return nil, err
	}
	if program.UserID != userID {
		return nil, common.ErrAccessForbidden
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

func (s *TrainingManager) AddWorkoutToProgram(input openapi.WorkoutRequest, userID uint, programID uint) (*Workout, error) {
	if input.Name == "" {
		return nil, fmt.Errorf("workout name cannot be empty")
	}
	program, err := s.trainingProgramRepository.FindById(programID)
	if err != nil {
		return nil, fmt.Errorf("failed to find training program: %w", err)
	}
	if program == nil {
		return nil, common.ErrProgramNotFound
	}
	if program.UserID != userID {
		return nil, common.ErrAccessForbidden
	}
	workout := &Workout{
		Name:              input.Name,
		TrainingProgramID: programID,
	}
	if err := s.workoutRepository.Create(workout); err != nil {
		return nil, fmt.Errorf("failed to create workout: %w", err)
	}
	return workout, nil
}

func (s *TrainingManager) UpdateWorkout(input openapi.WorkoutRequest, userID uint, programID uint, workoutId uint) (*Workout, error) {
	program, err := s.trainingProgramRepository.FindById(programID)
	if err != nil {
		return nil, err
	}
	if program == nil {
		return nil, nil
	}
	if program.UserID != userID {
		return nil, common.ErrAccessForbidden
	}
	workout, err := s.workoutRepository.FindById(workoutId)
	if err != nil {
		return nil, err
	}
	if workout == nil {
		return nil, nil
	}
	workout.Name = input.Name
	err = s.workoutRepository.Update(workout)
	if err != nil {
		return nil, err
	}
	return workout, nil
}

func (s *TrainingManager) GetAllWorkouts(userID, programID uint) ([]Workout, error) {
	workouts, err := s.workoutRepository.FindByTrainingProgramId(programID)
	if err != nil {
		return nil, err
	}
	return workouts, err
}

func (s *TrainingManager) DeleteWorkout(userId, program_id, workoutID uint) error {
	return s.workoutRepository.Delete(workoutID)
}

func (s *TrainingManager) GetWorkout(workoutID uint) (*Workout, error) {
	return s.workoutRepository.FindById(workoutID)
}

func (s *TrainingManager) AddExerciseToWorkout(input openapi.WorkoutExerciseRequest, userID uint) (*WorkoutExercise, error) {
	workoutID, err := strconv.Atoi(input.WorkoutId)
	if err != nil {
		return nil, err
	}
	exerciseID, err := strconv.Atoi(input.ExerciseId)
	if err != nil {
		return nil, err
	}
	workout, err := s.workoutRepository.FindById(uint(workoutID))
	if err != nil || workout == nil {
		return nil, errors.New("workout not found")
	}
	program, err := s.trainingProgramRepository.FindById(workout.TrainingProgramID)
	if err != nil || program == nil {
		return nil, common.ErrAccessForbidden
	}

	workoutExercise := &WorkoutExercise{
		WorkoutID:  uint(workoutID),
		ExerciseID: uint(exerciseID),
		Sets:       int(input.Sets),
		Reps:       int(input.Reps),
	}
	err = s.workoutExerciseRepository.Create(workoutExercise)
	if err != nil {
		return nil, err
	}
	return workoutExercise, nil
}

func (s *TrainingManager) GetAllWorkoutExercisesByWorkout(userID uint, workoutID uint) ([]WorkoutExercise, error) {
	workout, err := s.workoutRepository.FindById(workoutID)
	if err != nil || workout == nil {
		return nil, errors.New("workout not found")
	}
	program, err := s.trainingProgramRepository.FindById(workout.TrainingProgramID)
	if err != nil || program == nil {
		return nil, errors.New("workout does not belong to user's training program")
	}
	workoutExercises, err := s.workoutExerciseRepository.FindByWorkoutId(workoutID)
	if err != nil {
		return nil, err
	}
	return workoutExercises, nil
}

func (s *TrainingManager) UpdateWorkoutExercise(input openapi.WorkoutExerciseRequest, userID uint, weID uint) (*WorkoutExercise, error) {
	workoutId, err := strconv.Atoi(input.WorkoutId)
	if err != nil {
		return nil, err
	}
	exerciseId, err := strconv.Atoi(input.ExerciseId)
	if err != nil {
		return nil, err
	}
	workout, err := s.workoutRepository.FindById(uint(workoutId))
	if err != nil || workout == nil {
		return nil, errors.New("workout not found")
	}
	program, err := s.trainingProgramRepository.FindById(workout.TrainingProgramID)
	if err != nil || program == nil {
		return nil, errors.New("workout does not belong to user's training program")
	}
	workoutExercise, err := s.workoutExerciseRepository.FindById(weID)
	if err != nil {
		return nil, err
	}
	if exerciseId > 0 {
		workoutExercise.ExerciseID = uint(exerciseId)
	}
	if input.Sets > 0 {
		workoutExercise.Sets = int(input.Sets)
	}
	if input.Reps > 0 {
		workoutExercise.Reps = int(input.Reps)
	}
	err = s.workoutExerciseRepository.Update(workoutExercise)
	if err != nil {
		return nil, err
	}
	return workoutExercise, nil
}
func (s *TrainingManager) DeleteWorkoutExercise(userID uint, weID uint) error {
	workoutExercise, err := s.workoutExerciseRepository.FindById(weID)
	if err != nil {
		return err
	}
	workout, err := s.workoutRepository.FindById(workoutExercise.WorkoutID)
	if err != nil {
		return err
	}
	_, err = s.trainingProgramRepository.FindById(workout.TrainingProgramID)
	if err != nil {
		return err
	}
	return s.workoutExerciseRepository.Delete(weID)

}
func (s *TrainingManager) GetExercisesByPrimaryMuscle(primaryMuscle string) ([]Exercise, error) {
	return s.exerciseRepository.FindByPrimaryMuscle(primaryMuscle)
}

func (s *TrainingManager) GetAllExercises() ([]Exercise, error) {
	return s.exerciseRepository.FindAll()
}

func (s *TrainingManager) GetExercise(id uint) (*Exercise, error) {
	return s.exerciseRepository.FindById(id)
}
