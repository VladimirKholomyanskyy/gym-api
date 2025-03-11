package handlers

import (
	"context"
	"net/http"

	openapi "github.com/VladimirKholomyanskyy/gym-api/internal/api/go"
	"github.com/VladimirKholomyanskyy/gym-api/internal/common"
	customerrors "github.com/VladimirKholomyanskyy/gym-api/internal/customErrors"
	"github.com/VladimirKholomyanskyy/gym-api/internal/training/usecase"
	"github.com/VladimirKholomyanskyy/gym-api/internal/utils"
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
		return utils.ErrorResponse(http.StatusUnauthorized, openapi.FORBIDDEN, err.Error())
	}
	if common.IsUUIDValid(programId) {
		return utils.ErrorResponse(http.StatusBadRequest, openapi.INVALID_ID, "Program ID is not a valid UUID")
	}
	if !common.IsPageValid(page) {
		return utils.ErrorResponse(http.StatusBadRequest, openapi.INVALID_PAGE_NUMBER, "Page must be greater than 0")
	}

	if !common.IsPageSizeValid(pageSize) {
		return utils.ErrorResponse(http.StatusBadRequest, openapi.INVALID_PAGE_SIZE, "PageSize must be between 1 and 100")
	}

	workouts, totalCount, err := h.useCase.List(ctx, profileId, programId, int(page), int(pageSize))
	if err != nil {
		if err == customerrors.ErrAccessForbidden {
			return utils.ErrorResponse(http.StatusForbidden, openapi.FORBIDDEN, "Access denied to workout")
		}
		if err == customerrors.ErrEntityNotFound {
			return utils.ErrorResponse(http.StatusNotFound, openapi.RESOURCE_NOT_FOUND, "Workout not found")
		}
		return utils.ErrorResponse(http.StatusInternalServerError, openapi.INTERNAL_SERVER_ERROR, "Failed to fetch workouts")
	}
	return openapi.Response(
		http.StatusOK,
		openapi.ListWorkoutsForProgram200Response{
			TotalItems:  int32(totalCount),
			CurrentPage: page,
			PageSize:    pageSize,
			TotalPages:  utils.CalculateTotalPages(totalCount, pageSize),
			Items:       utils.ConvertWorkouts(workouts)}), nil
}

func (h *WorkoutHandler) AddWorkoutToProgram(ctx context.Context, programId string, request openapi.CreateWorkoutRequest) (openapi.ImplResponse, error) {
	profileId, err := common.ExtractProfileID(ctx)
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, openapi.FORBIDDEN, err.Error())
	}
	if common.IsUUIDValid(programId) {
		return utils.ErrorResponse(http.StatusBadRequest, openapi.INVALID_ID, "Program ID is not a valid UUID")
	}
	if utils.HasText(&request.Name) {
		return utils.ErrorResponse(http.StatusBadRequest, openapi.INVALID_REQUEST, "Workout name cannot be empty")
	}
	workout, err := h.useCase.Create(ctx, profileId, programId, request)
	if err != nil {
		if err == customerrors.ErrAccessForbidden {
			return utils.ErrorResponse(http.StatusForbidden, openapi.FORBIDDEN, "Access denied to workout")
		}
		if err == customerrors.ErrEntityNotFound {
			return utils.ErrorResponse(http.StatusNotFound, openapi.RESOURCE_NOT_FOUND, "Workout not found")
		}
		return utils.ErrorResponse(http.StatusInternalServerError, openapi.INTERNAL_SERVER_ERROR, "Failed to create workout")
	}
	return openapi.Response(http.StatusCreated, utils.ConvertWorkout(workout)), nil
}

func (h *WorkoutHandler) GetWorkoutForProgram(ctx context.Context, programId string, workoutId string) (openapi.ImplResponse, error) {
	profileId, err := common.ExtractProfileID(ctx)
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, openapi.FORBIDDEN, err.Error())
	}
	if common.IsUUIDValid(programId) {
		return utils.ErrorResponse(http.StatusBadRequest, openapi.INVALID_ID, "Program ID is not a valid UUID")
	}
	if common.IsUUIDValid(workoutId) {
		return utils.ErrorResponse(http.StatusBadRequest, openapi.INVALID_ID, "Workout ID is not a valid UUID")
	}
	workout, err := h.useCase.GetByProgramIDAndWorkoutID(ctx, profileId, programId, workoutId)
	if err != nil {
		if err == customerrors.ErrAccessForbidden {
			return utils.ErrorResponse(http.StatusForbidden, openapi.FORBIDDEN, "Access denied to workout")
		}
		if err == customerrors.ErrEntityNotFound {
			return utils.ErrorResponse(http.StatusNotFound, openapi.RESOURCE_NOT_FOUND, "Workout not found")
		}
		return utils.ErrorResponse(http.StatusInternalServerError, openapi.INTERNAL_SERVER_ERROR, "Failed to fetch workout")
	}
	return openapi.Response(http.StatusOK, utils.ConvertWorkout(workout)), nil
}

func (h *WorkoutHandler) UpdateWorkout(ctx context.Context, programId string, workoutId string, request openapi.PatchWorkoutRequest) (openapi.ImplResponse, error) {
	profileId, err := common.ExtractProfileID(ctx)
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, openapi.FORBIDDEN, err.Error())
	}
	if common.IsUUIDValid(programId) {
		return utils.ErrorResponse(http.StatusBadRequest, openapi.INVALID_ID, "Program ID is not a valid UUID")
	}
	if common.IsUUIDValid(workoutId) {
		return utils.ErrorResponse(http.StatusBadRequest, openapi.INVALID_ID, "Workout ID is not a valid UUID")
	}
	workout, err := h.useCase.Update(ctx, profileId, programId, workoutId, request)
	if err != nil {
		if err == customerrors.ErrAccessForbidden {
			return utils.ErrorResponse(http.StatusForbidden, openapi.FORBIDDEN, "Access denied to workout")
		}
		if err == customerrors.ErrEntityNotFound {
			return utils.ErrorResponse(http.StatusNotFound, openapi.RESOURCE_NOT_FOUND, "Workout not found")
		}
		return utils.ErrorResponse(http.StatusInternalServerError, openapi.INTERNAL_SERVER_ERROR, "Failed to update workout")
	}
	return openapi.Response(http.StatusOK, utils.ConvertWorkout(workout)), nil
}

func (h *WorkoutHandler) DeleteWorkout(ctx context.Context, programId string, workoutId string) (openapi.ImplResponse, error) {
	profileId, err := common.ExtractProfileID(ctx)
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, openapi.FORBIDDEN, err.Error())
	}
	if common.IsUUIDValid(programId) {
		return utils.ErrorResponse(http.StatusBadRequest, openapi.INVALID_ID, "Program ID is not a valid UUID")
	}
	if common.IsUUIDValid(workoutId) {
		return utils.ErrorResponse(http.StatusBadRequest, openapi.INVALID_ID, "Workout ID is not a valid UUID")
	}
	err = h.useCase.Delete(ctx, profileId, programId, workoutId)
	if err != nil {
		if err == customerrors.ErrAccessForbidden {
			return utils.ErrorResponse(http.StatusForbidden, openapi.FORBIDDEN, "Access denied to workout")
		}
		if err == customerrors.ErrEntityNotFound {
			return utils.ErrorResponse(http.StatusNotFound, openapi.RESOURCE_NOT_FOUND, "Workout not found")
		}
		return utils.ErrorResponse(http.StatusInternalServerError, openapi.INTERNAL_SERVER_ERROR, "Failed to delete workout")
	}

	return openapi.Response(http.StatusNoContent, nil), nil
}

func (h *WorkoutHandler) ReorderWorkout(ctx context.Context, programID string, workoutID string, request openapi.ReorderWorkoutRequest) (openapi.ImplResponse, error) {
	profileId, err := common.ExtractProfileID(ctx)
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, openapi.FORBIDDEN, err.Error())
	}
	if common.IsUUIDValid(programID) {
		return utils.ErrorResponse(http.StatusBadRequest, openapi.INVALID_ID, "Program ID is not a valid UUID")
	}
	if common.IsUUIDValid(workoutID) {
		return utils.ErrorResponse(http.StatusBadRequest, openapi.INVALID_ID, "Workout ID is not a valid UUID")
	}
	err = h.useCase.Reorder(ctx, profileId, programID, workoutID, request)
	if err != nil {
		return utils.ErrorResponse(http.StatusInternalServerError, openapi.INTERNAL_SERVER_ERROR, err.Error())
	}
	return openapi.Response(http.StatusNoContent, nil), nil
}
