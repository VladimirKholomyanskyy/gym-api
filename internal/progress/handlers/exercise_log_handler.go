package handlers

import (
	"context"
	"net/http"

	openapi "github.com/VladimirKholomyanskyy/gym-api/internal/api/go"
	"github.com/VladimirKholomyanskyy/gym-api/internal/common"
	"github.com/VladimirKholomyanskyy/gym-api/internal/progress/model"
	usecase "github.com/VladimirKholomyanskyy/gym-api/internal/progress/usecase"
)

type exerciseLogHandler struct {
	useCase usecase.LogExerciseUseCase
}

// NewWorkoutSessionsAPIService creates a default api service
func NewExerciseLogHandler(useCase usecase.LogExerciseUseCase) openapi.ExerciseLogsAPIServicer {
	return &exerciseLogHandler{useCase: useCase}
}

// LogExercise - Log an exercise during a workout session
func (h *exerciseLogHandler) LogExercise(ctx context.Context, logExerciseRequest openapi.CreateExerciseLogRequest) (openapi.ImplResponse, error) {
	profileId, err := common.ExtractProfileID(ctx)
	if err != nil {
		return common.ErrorResponse(http.StatusUnauthorized, openapi.FORBIDDEN, err.Error())
	}

	log, err := h.useCase.Create(ctx, profileId, logExerciseRequest)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), nil
	}
	return openapi.Response(http.StatusCreated, convertExerciseLog(log)), nil
}

// ListExerciseLogs - Retrieve all logged exercises for a workout session
func (h *exerciseLogHandler) ListExerciseLogs(ctx context.Context, workoutSessionIdParam, exerciseIdParam string, page, pageSize int32) (openapi.ImplResponse, error) {
	profileId, err := common.ExtractProfileID(ctx)
	if err != nil {
		return common.ErrorResponse(http.StatusUnauthorized, openapi.FORBIDDEN, err.Error())
	}
	var (
		exerciseLogs []model.ExerciseLog
		totalCount   int64
	)
	switch {
	case workoutSessionIdParam != "":
		// Fetch logs by session ID
		exerciseLogs, totalCount, err = h.useCase.GetExerciseLogsBySessionID(ctx, profileId, workoutSessionIdParam, int(page), int(pageSize))

	case exerciseIdParam != "":
		// Fetch logs by exercise ID
		exerciseLogs, totalCount, err = h.useCase.GetExerciseLogsByExerciseID(ctx, profileId, exerciseIdParam, int(page), int(pageSize))

	default:
		// Fetch all logs for the user
		exerciseLogs, totalCount, err = h.useCase.GetExerciseLogsByProfileID(ctx, profileId, int(page), int(pageSize))
	}

	// Handle repository error
	if err != nil {
		return common.ErrorResponse(http.StatusInternalServerError, openapi.INTERNAL_SERVER_ERROR, "failed to fetch ")
	}

	return openapi.Response(
		http.StatusOK,
		openapi.ListExerciseLogs200Response{
			TotalItems:  int32(totalCount),
			CurrentPage: page,
			PageSize:    pageSize,
			TotalPages:  common.CalculateTotalPages(totalCount, pageSize),
			Items:       convertExerciseLogs(exerciseLogs)}), nil
}

// GetExerciseLog - Retrieve details of a specific exercise log
func (h *exerciseLogHandler) GetExerciseLog(ctx context.Context, exerciseLogId string) (openapi.ImplResponse, error) {
	profileId, err := common.ExtractProfileID(ctx)
	if err != nil {
		return common.ErrorResponse(http.StatusUnauthorized, openapi.FORBIDDEN, err.Error())
	}
	log, err := h.useCase.GetExerciseLog(ctx, profileId, exerciseLogId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), nil
	}
	return openapi.Response(http.StatusCreated, convertExerciseLog(log)), nil
}

func (h *exerciseLogHandler) GetWeightPerDay(ctx context.Context, exerciseId string, startDate string, endDate string) (openapi.ImplResponse, error) {
	profileId, err := common.ExtractProfileID(ctx)
	if err != nil {
		return common.ErrorResponse(http.StatusUnauthorized, openapi.FORBIDDEN, err.Error())
	}

	startDateTime, err := common.ParseTime(startDate)
	if err != nil {
		return common.ErrorResponse(http.StatusBadRequest, openapi.INVALID_DATE_FORMAT, "Invalid startDate format")
	}
	endDateTime, err := common.ParseTime(endDate)
	if err != nil {
		return common.ErrorResponse(http.StatusBadRequest, openapi.INVALID_DATE_FORMAT, "Invalid endDate format")
	}

	// Fetch weight per day data from service
	weightPerDayList, err := h.useCase.GetWeightPerDay(ctx, profileId, exerciseId, &startDateTime, &endDateTime)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, "Failed to fetch weight per day"), nil
	}

	// Convert response to OpenAPI format
	var response []openapi.GetWeightPerDayTotalWeightPerDayInner
	for _, v := range weightPerDayList {
		response = append(response, openapi.GetWeightPerDayTotalWeightPerDayInner{
			Date:        common.FormatTime(&v.Date),
			TotalWeight: v.TotalWeight,
		})
	}

	// Return response
	return openapi.Response(http.StatusOK, openapi.GetWeightPerDay{
		ExerciseId:        exerciseId,
		TotalWeightPerDay: response,
	}), nil
}
