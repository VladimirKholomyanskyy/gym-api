package progress

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

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
func (h *WorkoutProgressHandler) CompleteWorkoutSession(ctx context.Context, workoutSessionId string) (openapi.ImplResponse, error) {
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
func (h *WorkoutProgressHandler) ListExerciseLogs(ctx context.Context, workoutSessionIdParam, exerciseIdParam string) (openapi.ImplResponse, error) {
	userID := ctx.Value(common.UserIDKey).(uint)
	var (
		exerciseLogs  []ExerciseLog
		convertedLogs []openapi.LogExerciseResponse
		err           error
	)

	switch {
	case workoutSessionIdParam != "":
		// Parse workout session ID
		sessionID, parseErr := strconv.Atoi(workoutSessionIdParam)
		if parseErr != nil {
			return openapi.Response(http.StatusBadRequest, "Invalid workout session ID"), nil
		}
		// Fetch logs by session ID
		exerciseLogs, err = h.s.logRepository.GetAllByUserIDAndSessionID(userID, uint(sessionID))

	case exerciseIdParam != "":
		// Parse exercise ID
		exerciseID, parseErr := strconv.Atoi(exerciseIdParam)
		if parseErr != nil {
			return openapi.Response(http.StatusBadRequest, "Invalid exercise ID"), nil
		}
		// Fetch logs by exercise ID
		exerciseLogs, err = h.s.logRepository.GetAllByUserIDAndExerciseID(userID, uint(exerciseID))

	default:
		// Fetch all logs for the user
		exerciseLogs, err = h.s.logRepository.GetAllByUserID(userID)
	}

	// Handle repository error
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, "Failed to fetch exercise logs"), nil
	}

	// Convert logs to response format
	for _, log := range exerciseLogs {
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
func (h *WorkoutProgressHandler) LogExercise(ctx context.Context, logExerciseRequest openapi.LogExerciseRequest) (openapi.ImplResponse, error) {
	userID := ctx.Value(common.UserIDKey).(uint)
	sessionID, err := strconv.Atoi(logExerciseRequest.WorkoutSessionId)
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
func (h *WorkoutProgressHandler) GetExerciseLog(ctx context.Context, exerciseLogId string) (openapi.ImplResponse, error) {
	userID := ctx.Value(common.UserIDKey).(uint)
	logID, err := strconv.Atoi(exerciseLogId)
	if err != nil {
		return openapi.Response(http.StatusBadRequest, nil), nil
	}
	log, err := h.s.GetExerciseLog(userID, uint(logID))
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

func (h *WorkoutProgressHandler) GetWeightPerDay(ctx context.Context, exerciseId string, startDate string, endDate string) (openapi.ImplResponse, error) {
	userID := ctx.Value(common.UserIDKey).(uint)

	// Parse exerciseId
	exerciseID, err := strconv.Atoi(exerciseId)
	if err != nil {
		return openapi.Response(http.StatusBadRequest, "Invalid exerciseId"), nil
	}

	// Define date layout and initialize pointers
	layout := "2006-01-02"
	var startDatePtr, endDatePtr *time.Time

	// Parse startDate
	if startDate != "" {
		parsedStartDate, err := time.Parse(layout, startDate)
		if err != nil {
			return openapi.Response(http.StatusBadRequest, "Invalid startDate format"), nil
		}
		startDatePtr = &parsedStartDate
	}

	// Parse endDate
	if endDate != "" {
		parsedEndDate, err := time.Parse(layout, endDate)
		if err != nil {
			return openapi.Response(http.StatusBadRequest, "Invalid endDate format"), nil
		}
		endDatePtr = &parsedEndDate
	}

	// Fetch weight per day data from service
	weightPerDayList, err := h.s.GetWeightPerDay(userID, uint(exerciseID), startDatePtr, endDatePtr)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, "Failed to fetch weight per day"), nil
	}

	// Convert response to OpenAPI format
	var response []openapi.WeightPerDayResponseTotalWeightPerDayInner
	for _, v := range weightPerDayList {
		response = append(response, openapi.WeightPerDayResponseTotalWeightPerDayInner{
			Date:        v.Date.Format(layout),
			TotalWeight: float32(v.TotalWeight),
		})
	}

	// Return response
	return openapi.Response(http.StatusOK, openapi.WeightPerDayResponse{
		ExerciseId:        exerciseId,
		TotalWeightPerDay: response,
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
