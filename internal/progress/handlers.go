package progress

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	openapi "github.com/VladimirKholomyanskyy/gym-api/internal/api/go"
	"github.com/VladimirKholomyanskyy/gym-api/internal/common"
	"github.com/VladimirKholomyanskyy/gym-api/internal/training"
	"gorm.io/datatypes"
)

type WorkoutProgressHandler struct {
	s *WorkoutProgressManager
}

// NewWorkoutSessionsAPIService creates a default api service
func NewWorkoutProgressHandler(s *WorkoutProgressManager) *WorkoutProgressHandler {
	return &WorkoutProgressHandler{s: s}
}

func (h *WorkoutProgressHandler) GetWorkoutSession(ctx context.Context, workoutSessionId string) (openapi.ImplResponse, error) {
	userID := ctx.Value(common.UserIDKey).(uint)
	sessionID, err := strconv.Atoi(workoutSessionId)
	if err != nil {
		return openapi.Response(http.StatusBadRequest, nil), nil
	}
	session, err := h.s.GetWorkoutSession(userID, uint(sessionID))
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), nil
	}
	snapshot, err := h.convert(&session.Snapshot)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), nil
	}
	if session.CompletedAt != nil {
		return openapi.Response(http.StatusOK, openapi.WorkoutSessionResponse{
			Id:              fmt.Sprintf("%d", session.ID),
			StartedAt:       session.StartedAt,
			CompletedAt:     session.CompletedAt,
			WorkoutSnapshot: *snapshot,
		}), nil
	}
	return openapi.Response(http.StatusOK, openapi.WorkoutSessionResponse{
		Id:              fmt.Sprintf("%d", session.ID),
		StartedAt:       session.StartedAt,
		WorkoutSnapshot: *snapshot,
	}), nil
}

func (h *WorkoutProgressHandler) ListWorkoutSessions(ctx context.Context) (openapi.ImplResponse, error) {
	userID := ctx.Value(common.UserIDKey).(uint)

	sessions, err := h.s.GetAllWorkoutSessions(userID)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), nil
	}
	var convertedSessions []openapi.WorkoutSessionResponse
	for _, session := range sessions {
		snapshot, err := h.convert(&session.Snapshot)
		if err != nil {
			return openapi.Response(http.StatusInternalServerError, nil), nil
		}
		if session.CompletedAt != nil {
			convertedSessions = append(convertedSessions, openapi.WorkoutSessionResponse{
				Id:              fmt.Sprintf("%d", session.ID),
				StartedAt:       session.StartedAt,
				CompletedAt:     session.CompletedAt,
				WorkoutSnapshot: *snapshot,
			})
		} else {
			convertedSessions = append(convertedSessions, openapi.WorkoutSessionResponse{
				Id:              fmt.Sprintf("%d", session.ID),
				StartedAt:       session.StartedAt,
				WorkoutSnapshot: *snapshot,
			})
		}

	}
	return openapi.Response(http.StatusOK, convertedSessions), nil
}

func (s *WorkoutProgressHandler) AddWorkoutSession(ctx context.Context, startWorkoutSessionRequest openapi.StartWorkoutSessionRequest) (openapi.ImplResponse, error) {
	userID := ctx.Value(common.UserIDKey).(uint)
	workoutSession, err := s.s.StartWorkout(userID, startWorkoutSessionRequest)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), nil
	}
	snapshot, err := s.convert(&workoutSession.Snapshot)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), nil
	}
	return openapi.Response(http.StatusCreated, openapi.WorkoutSessionResponse{
		Id:              fmt.Sprintf("%d", workoutSession.ID),
		StartedAt:       workoutSession.StartedAt,
		WorkoutSnapshot: *snapshot,
	}), nil
}

// FinishWorkoutSession - Mark a workout session as completed
func (h *WorkoutProgressHandler) FinishWorkoutSession(ctx context.Context, workoutSessionId string) (openapi.ImplResponse, error) {
	userID := ctx.Value(common.UserIDKey).(uint)
	sessionID, err := strconv.Atoi(workoutSessionId)
	if err != nil {
		return openapi.Response(http.StatusBadRequest, nil), nil
	}
	workoutSession, err := h.s.FinishWorkoutSession(userID, uint(sessionID))
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), nil
	}

	snapshot, err := h.convert(&workoutSession.Snapshot)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), nil
	}
	return openapi.Response(http.StatusCreated, openapi.WorkoutSessionResponse{
		Id:              fmt.Sprintf("%d", workoutSession.ID),
		StartedAt:       workoutSession.StartedAt,
		WorkoutSnapshot: *snapshot,
	}), nil
}

// ListExerciseLogs - Retrieve all logged exercises for a workout session
func (h *WorkoutProgressHandler) ListExerciseLogs(ctx context.Context, workoutSessionId string) (openapi.ImplResponse, error) {
	userID := ctx.Value(common.UserIDKey).(uint)
	sessionID, err := strconv.Atoi(workoutSessionId)
	if err != nil {
		return openapi.Response(http.StatusBadRequest, nil), nil
	}
	logs, err := h.s.GetExerciseLogs(userID, uint(sessionID))
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), nil
	}
	var convertedLogs []openapi.LogExerciseResponse
	for _, log := range logs {
		convertedLogs = append(convertedLogs, openapi.LogExerciseResponse{
			Id:            fmt.Sprintf("%d", log.ID),
			ExerciseId:    fmt.Sprintf("%d", log.ExerciseID),
			SetNumber:     int32(log.SetNumber),
			RepsCompleted: int32(log.Reps),
			WeightUsed:    int32(log.Weight),
			LoggedAt:      log.LoggedAt,
		})
	}

	return openapi.Response(http.StatusOK, convertedLogs), nil
}

// LogExercise - Log an exercise during a workout session
func (h *WorkoutProgressHandler) LogExercise(ctx context.Context, workoutSessionId string, logExerciseRequest openapi.LogExerciseRequest) (openapi.ImplResponse, error) {
	userID := ctx.Value(common.UserIDKey).(uint)
	sessionID, err := strconv.Atoi(workoutSessionId)
	if err != nil {
		return openapi.Response(http.StatusBadRequest, nil), nil
	}
	log, err := h.s.LogExercise(userID, uint(sessionID), logExerciseRequest)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), nil
	}
	return openapi.Response(http.StatusCreated, openapi.LogExerciseResponse{
		Id:            fmt.Sprintf("%d", log.ID),
		ExerciseId:    fmt.Sprintf("%d", log.ExerciseID),
		SetNumber:     int32(log.SetNumber),
		RepsCompleted: int32(log.Reps),
		WeightUsed:    int32(log.Weight),
		LoggedAt:      log.LoggedAt,
	}), nil
}

// GetExerciseLog - Retrieve details of a specific exercise log
func (h *WorkoutProgressHandler) GetExerciseLog(ctx context.Context, workoutSessionId string, exerciseLogId string) (openapi.ImplResponse, error) {
	userID := ctx.Value(common.UserIDKey).(uint)
	sessionID, err := strconv.Atoi(workoutSessionId)
	if err != nil {
		return openapi.Response(http.StatusBadRequest, nil), nil
	}
	logID, err := strconv.Atoi(exerciseLogId)
	if err != nil {
		return openapi.Response(http.StatusBadRequest, nil), nil
	}
	log, err := h.s.GetExerciseLog(userID, uint(sessionID), uint(logID))
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), nil
	}
	return openapi.Response(http.StatusCreated, openapi.LogExerciseResponse{
		Id:            fmt.Sprintf("%d", log.ID),
		ExerciseId:    fmt.Sprintf("%d", log.ExerciseID),
		SetNumber:     int32(log.SetNumber),
		RepsCompleted: int32(log.Reps),
		WeightUsed:    int32(log.Weight),
		LoggedAt:      log.LoggedAt,
	}), nil
}

func (s *WorkoutProgressHandler) convert(workoutJson *datatypes.JSON) (*openapi.WorkoutSnapshot, error) {
	var workout training.Workout
	err := json.Unmarshal(*workoutJson, &workout)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal workout JSON: %w", err)
	}
	var snapshotExercises []openapi.WorkoutExercise
	for _, we := range workout.Exercises {
		snapshotExercises = append(snapshotExercises, openapi.WorkoutExercise{
			Id:        fmt.Sprintf("%d", we.ID),
			Sets:      int32(we.Sets),
			Reps:      int32(we.Reps),
			CreatedAt: we.CreatedAt,
			DeletedAt: &we.DeletedAt.Time,
			UpdatedAt: we.UpdatedAt,
			Exercise: openapi.Exercise{
				Id:              fmt.Sprintf("%d", we.Exercise.ID),
				Name:            we.Exercise.Name,
				PrimaryMuscle:   we.Exercise.PrimaryMuscle,
				SecondaryMuscle: we.Exercise.SecondaryMuscle,
				Equipment:       we.Exercise.Equipment,
				Description:     we.Exercise.Description,
			},
		})
	}
	return &openapi.WorkoutSnapshot{
		Id:                fmt.Sprintf("%d", workout.ID),
		Name:              workout.Name,
		CreatedAt:         workout.CreatedAt,
		DeletedAt:         &workout.DeletedAt.Time,
		UpdatedAt:         workout.UpdatedAt,
		TrainingProgramId: fmt.Sprintf("%d", workout.TrainingProgramID),
		WorkoutExercises:  snapshotExercises,
	}, nil
}
