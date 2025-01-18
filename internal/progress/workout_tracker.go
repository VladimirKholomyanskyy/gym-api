package progress

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	openapi "github.com/VladimirKholomyanskyy/gym-api/internal/api/go"
	"github.com/VladimirKholomyanskyy/gym-api/internal/common"
	"github.com/VladimirKholomyanskyy/gym-api/internal/training"
)

type WorkoutProgressManager struct {
	trainingManager   *training.TrainingManager
	sessionRepository *WorkoutSessionRepository
	logRepository     *ExerciseLogRepository
}

func NewWorkoutProgressManager(trainingManager *training.TrainingManager,
	sessionRepository *WorkoutSessionRepository,
	logRepository *ExerciseLogRepository) *WorkoutProgressManager {
	return &WorkoutProgressManager{trainingManager: trainingManager, sessionRepository: sessionRepository, logRepository: logRepository}
}

func (s *WorkoutProgressManager) StartWorkout(userId uint, input openapi.StartWorkoutSessionRequest) (*WorkoutSession, error) {
	workoutId, err := strconv.Atoi(input.WorkoutId)
	if err != nil {
		return nil, err
	}
	workout, err := s.trainingManager.GetWorkout(uint(workoutId))
	if err != nil {
		return nil, err
	}
	program, err := s.trainingManager.GetTrainingProgram(userId, workout.TrainingProgramID)
	if err != nil {
		return nil, err
	}
	if program == nil {
		return nil, common.ErrAccessForbidden
	}
	jsonData, err := json.Marshal(workout)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal workout: %w", err)
	}
	workoutSession := &WorkoutSession{UserID: userId, WorkoutID: workout.ID, Snapshot: jsonData}
	err = s.sessionRepository.Create(workoutSession)
	if err != nil {
		return nil, err
	}
	return workoutSession, nil
}

func (s *WorkoutProgressManager) FinishWorkoutSession(userID uint, sessionID uint) (*WorkoutSession, error) {
	session, err := s.sessionRepository.GetByID(sessionID)
	if err != nil {
		return nil, err
	}
	if session.UserID != userID {
		return nil, common.ErrAccessForbidden
	}
	currnetTime := time.Now()
	session.CompletedAt = &currnetTime
	err = s.sessionRepository.Update(session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (s *WorkoutProgressManager) GetAllWorkoutSessions(userID uint) ([]WorkoutSession, error) {
	sessions, err := s.sessionRepository.GetAllByUserID(userID)
	if err != nil {
		return nil, err
	}
	return sessions, err
}

func (s *WorkoutProgressManager) GetWorkoutSession(userID uint, sessionID uint) (*WorkoutSession, error) {
	session, err := s.sessionRepository.GetByID(sessionID)
	if err != nil {
		return nil, err
	}
	if session.UserID != userID {
		return nil, common.ErrAccessForbidden
	}
	return session, err
}

func (s *WorkoutProgressManager) LogExercise(userID uint, sessionID uint, input openapi.LogExerciseRequest) (*ExerciseLog, error) {
	session, err := s.sessionRepository.GetByID(sessionID)
	if err != nil {
		return nil, err
	}
	if session.UserID != userID {
		return nil, common.ErrAccessForbidden
	}
	exerciseId, err := strconv.Atoi(input.ExerciseId)
	if err != nil {
		return nil, err
	}
	exerciseLog := &ExerciseLog{SessionID: sessionID, ExerciseID: uint(exerciseId), SetNumber: int(input.SetNumber), Reps: int(input.RepsCompleted), Weight: float64(input.WeightUsed)}
	err = s.logRepository.Create(exerciseLog)
	if err != nil {
		return nil, err
	}
	return exerciseLog, nil
}

func (s *WorkoutProgressManager) GetExerciseLog(userID uint, sessionID uint, logID uint) (*ExerciseLog, error) {
	session, err := s.sessionRepository.GetByID(sessionID)
	if err != nil {
		return nil, err
	}
	if session.UserID != userID {
		return nil, common.ErrAccessForbidden
	}
	log, err := s.logRepository.GetByID(logID)
	if err != nil {
		return nil, err
	}
	return log, nil

}

func (s *WorkoutProgressManager) GetExerciseLogs(userID uint, sessionID uint) ([]ExerciseLog, error) {
	session, err := s.sessionRepository.GetByID(sessionID)
	if err != nil {
		return nil, err
	}
	if session.UserID != userID {
		return nil, common.ErrAccessForbidden
	}
	logs, err := s.logRepository.GetAllByWorkoutLogID(sessionID)
	if err != nil {
		return nil, err
	}
	return logs, nil
}
