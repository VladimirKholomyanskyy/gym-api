package handlers

import (
	"context"
	"errors"
	"net/http"

	openapi "github.com/VladimirKholomyanskyy/gym-api/internal/api/go"
	"github.com/VladimirKholomyanskyy/gym-api/internal/common"
	"github.com/VladimirKholomyanskyy/gym-api/internal/training/usecase"
)

type WorkoutHandler struct {
	useCase usecase.WorkoutUseCase
}

func NewWorkoutHandler(useCase usecase.WorkoutUseCase) openapi.WorkoutsAPIServicer {
	return &WorkoutHandler{useCase: useCase}
}

func (h *WorkoutHandler) ListWorkoutsForProgram(ctx context.Context, programId string, page, pageSize int32) (openapi.ImplResponse, error) {
	profileId, err := common.ExtractProfileID(ctx)
	if err != nil {
		return common.ErrorResponse(http.StatusUnauthorized, openapi.FORBIDDEN, err.Error())
	}
	if isUUIDValid(programId) {
		return common.ErrorResponse(http.StatusBadRequest, openapi.INVALID_ID, "Program ID is not a valid UUID")
	}
	if !isPageValid(page) {
		return common.ErrorResponse(http.StatusBadRequest, openapi.INVALID_PAGE_NUMBER, "Page must be greater than 0")
	}

	if !isPageSizeValid(pageSize) {
		return common.ErrorResponse(http.StatusBadRequest, openapi.INVALID_PAGE_SIZE, "PageSize must be between 1 and 100")
	}

	workouts, totalCount, err := h.useCase.List(ctx, profileId, programId, int(page), int(pageSize))
	if err != nil {
		switch {
		case errors.As(err, &common.NotFoundError{}):
			return common.ErrorResponse(http.StatusNotFound, openapi.RESOURCE_NOT_FOUND, err.Error())
		case errors.As(err, &common.ForbiddenError{}):
			return common.ErrorResponse(http.StatusForbidden, openapi.FORBIDDEN, err.Error())
		default:
			return common.ErrorResponse(http.StatusInternalServerError, openapi.INTERNAL_SERVER_ERROR, err.Error())
		}
	}
	return openapi.Response(
		http.StatusOK,
		openapi.ListWorkoutsForProgram200Response{
			TotalItems:  int32(totalCount),
			CurrentPage: page,
			PageSize:    pageSize,
			TotalPages:  common.CalculateTotalPages(totalCount, pageSize),
			Items:       common.ConvertWorkouts(workouts)}), nil
}

func (h *WorkoutHandler) AddWorkoutToProgram(ctx context.Context, programId string, request openapi.CreateWorkoutRequest) (openapi.ImplResponse, error) {
	profileId, err := common.ExtractProfileID(ctx)
	if err != nil {
		return common.ErrorResponse(http.StatusUnauthorized, openapi.FORBIDDEN, err.Error())
	}
	if isUUIDValid(programId) {
		return common.ErrorResponse(http.StatusBadRequest, openapi.INVALID_ID, "Program ID is not a valid UUID")
	}
	if common.HasText(&request.Name) {
		return common.ErrorResponse(http.StatusBadRequest, openapi.INVALID_REQUEST, "Workout name cannot be empty")
	}
	workout, err := h.useCase.Create(ctx, profileId, programId, request)
	if err != nil {
		switch {
		case errors.As(err, &common.NotFoundError{}):
			return common.ErrorResponse(http.StatusNotFound, openapi.RESOURCE_NOT_FOUND, err.Error())
		case errors.As(err, &common.ForbiddenError{}):
			return common.ErrorResponse(http.StatusForbidden, openapi.FORBIDDEN, err.Error())
		default:
			return common.ErrorResponse(http.StatusInternalServerError, openapi.INTERNAL_SERVER_ERROR, err.Error())
		}
	}
	return openapi.Response(http.StatusCreated, common.ConvertWorkout(workout)), nil
}

func (h *WorkoutHandler) GetWorkoutForProgram(ctx context.Context, programId string, workoutId string) (openapi.ImplResponse, error) {
	profileId, err := common.ExtractProfileID(ctx)
	if err != nil {
		return common.ErrorResponse(http.StatusUnauthorized, openapi.FORBIDDEN, err.Error())
	}
	if isUUIDValid(programId) {
		return common.ErrorResponse(http.StatusBadRequest, openapi.INVALID_ID, "Program ID is not a valid UUID")
	}
	if isUUIDValid(workoutId) {
		return common.ErrorResponse(http.StatusBadRequest, openapi.INVALID_ID, "Workout ID is not a valid UUID")
	}
	workout, err := h.useCase.GetByProgramIDAndWorkoutID(ctx, profileId, programId, workoutId)
	if err != nil {
		switch {
		case errors.As(err, &common.NotFoundError{}):
			return common.ErrorResponse(http.StatusNotFound, openapi.RESOURCE_NOT_FOUND, err.Error())
		case errors.As(err, &common.ForbiddenError{}):
			return common.ErrorResponse(http.StatusForbidden, openapi.FORBIDDEN, err.Error())
		default:
			return common.ErrorResponse(http.StatusInternalServerError, openapi.INTERNAL_SERVER_ERROR, err.Error())
		}
	}
	return openapi.Response(http.StatusOK, common.ConvertWorkout(workout)), nil
}

func (h *WorkoutHandler) UpdateWorkout(ctx context.Context, programId string, workoutId string, request openapi.PatchWorkoutRequest) (openapi.ImplResponse, error) {
	profileId, err := common.ExtractProfileID(ctx)
	if err != nil {
		return common.ErrorResponse(http.StatusUnauthorized, openapi.FORBIDDEN, err.Error())
	}
	if isUUIDValid(programId) {
		return common.ErrorResponse(http.StatusBadRequest, openapi.INVALID_ID, "Program ID is not a valid UUID")
	}
	if isUUIDValid(workoutId) {
		return common.ErrorResponse(http.StatusBadRequest, openapi.INVALID_ID, "Workout ID is not a valid UUID")
	}
	workout, err := h.useCase.Update(ctx, profileId, programId, workoutId, request)
	if err != nil {
		switch {
		case errors.As(err, &common.NotFoundError{}):
			return common.ErrorResponse(http.StatusNotFound, openapi.RESOURCE_NOT_FOUND, err.Error())
		case errors.As(err, &common.ForbiddenError{}):
			return common.ErrorResponse(http.StatusForbidden, openapi.FORBIDDEN, err.Error())
		default:
			return common.ErrorResponse(http.StatusInternalServerError, openapi.INTERNAL_SERVER_ERROR, err.Error())
		}
	}
	return openapi.Response(http.StatusOK, common.ConvertWorkout(workout)), nil
}

func (h *WorkoutHandler) DeleteWorkout(ctx context.Context, programId string, workoutId string) (openapi.ImplResponse, error) {
	profileId, err := common.ExtractProfileID(ctx)
	if err != nil {
		return common.ErrorResponse(http.StatusUnauthorized, openapi.FORBIDDEN, err.Error())
	}
	if isUUIDValid(programId) {
		return common.ErrorResponse(http.StatusBadRequest, openapi.INVALID_ID, "Program ID is not a valid UUID")
	}
	if isUUIDValid(workoutId) {
		return common.ErrorResponse(http.StatusBadRequest, openapi.INVALID_ID, "Workout ID is not a valid UUID")
	}
	err = h.useCase.Delete(ctx, profileId, programId, workoutId)
	if err != nil {
		switch {
		case errors.As(err, &common.NotFoundError{}):
			return common.ErrorResponse(http.StatusNotFound, openapi.RESOURCE_NOT_FOUND, err.Error())
		case errors.As(err, &common.ForbiddenError{}):
			return common.ErrorResponse(http.StatusForbidden, openapi.FORBIDDEN, err.Error())
		default:
			return common.ErrorResponse(http.StatusInternalServerError, openapi.INTERNAL_SERVER_ERROR, err.Error())
		}
	}

	return openapi.Response(http.StatusNoContent, nil), nil
}

func (h *WorkoutHandler) ReorderWorkout(ctx context.Context, programID string, workoutID string, request openapi.ReorderWorkoutRequest) (openapi.ImplResponse, error) {
	profileId, err := common.ExtractProfileID(ctx)
	if err != nil {
		return common.ErrorResponse(http.StatusUnauthorized, openapi.FORBIDDEN, err.Error())
	}
	if isUUIDValid(programID) {
		return common.ErrorResponse(http.StatusBadRequest, openapi.INVALID_ID, "Program ID is not a valid UUID")
	}
	if isUUIDValid(workoutID) {
		return common.ErrorResponse(http.StatusBadRequest, openapi.INVALID_ID, "Workout ID is not a valid UUID")
	}
	err = h.useCase.Reorder(ctx, profileId, programID, workoutID, request)
	if err != nil {
		return common.ErrorResponse(http.StatusInternalServerError, openapi.INTERNAL_SERVER_ERROR, err.Error())
	}
	return openapi.Response(http.StatusNoContent, nil), nil
}
