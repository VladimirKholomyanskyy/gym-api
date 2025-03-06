package handlers

import (
	"context"
	"errors"
	"net/http"

	openapi "github.com/VladimirKholomyanskyy/gym-api/internal/api/go"
	"github.com/VladimirKholomyanskyy/gym-api/internal/common"
	"github.com/VladimirKholomyanskyy/gym-api/internal/training/usecase"
)

type workoutExerciseHandler struct {
	useCase usecase.WorkoutExerciseUseCase
}

func NewWorkoutExerciseHandler(useCase usecase.WorkoutExerciseUseCase) openapi.WorkoutExercisesAPIServicer {
	return &workoutExerciseHandler{useCase: useCase}
}

func (h *workoutExerciseHandler) ListWorkoutExercises(ctx context.Context, workoutId string, page, pageSize int32) (openapi.ImplResponse, error) {
	profileId, err := common.ExtractProfileID(ctx)
	if err != nil {
		return common.ErrorResponse(http.StatusUnauthorized, openapi.FORBIDDEN, err.Error())
	}
	if isUUIDValid(workoutId) {
		return common.ErrorResponse(http.StatusBadRequest, openapi.INVALID_ID, "Workout ID is not a valid UUID")
	}
	if !isPageValid(page) {
		return common.ErrorResponse(http.StatusBadRequest, openapi.INVALID_PAGE_NUMBER, "Page must be greater than 0")
	}

	if !isPageSizeValid(pageSize) {
		return common.ErrorResponse(http.StatusBadRequest, openapi.INVALID_PAGE_SIZE, "PageSize must be between 1 and 100")
	}
	workoutExercises, totalCount, err := h.useCase.List(ctx, profileId, workoutId, int(page), int(pageSize))
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
		openapi.ListWorkoutExercises200Response{
			TotalItems:  int32(totalCount),
			CurrentPage: page,
			PageSize:    pageSize,
			TotalPages:  common.CalculateTotalPages(totalCount, pageSize),
			Items:       common.ConvertWorkoutExercises(workoutExercises)}), nil

}

func (h *workoutExerciseHandler) AddWorkoutExercise(ctx context.Context, request openapi.CreateWorkoutExerciseRequest) (openapi.ImplResponse, error) {
	profileId, err := common.ExtractProfileID(ctx)
	if err != nil {
		return common.ErrorResponse(http.StatusUnauthorized, openapi.FORBIDDEN, err.Error())
	}
	if isUUIDValid(request.WorkoutId) {
		return common.ErrorResponse(http.StatusBadRequest, openapi.INVALID_ID, "Workout ID is not a valid UUID")
	}
	if isUUIDValid(request.ExerciseId) {
		return common.ErrorResponse(http.StatusBadRequest, openapi.INVALID_ID, "Exercise ID is not a valid UUID")
	}
	if request.Sets < 1 {
		return common.ErrorResponse(http.StatusBadRequest, openapi.INVALID_REQUEST, "Sets must be greater then 0")
	}
	if request.Reps < 1 {
		return common.ErrorResponse(http.StatusBadRequest, openapi.INVALID_REQUEST, "Reps must be greater then 0")
	}
	workoutExercise, err := h.useCase.Create(ctx, profileId, request)
	if err != nil {
		return common.ErrorResponse(http.StatusInternalServerError, openapi.INTERNAL_SERVER_ERROR, err.Error())
	}

	return openapi.Response(http.StatusCreated, common.ConvertWorkoutExercise(workoutExercise)), nil
}

func (h *workoutExerciseHandler) UpdateWorkoutExercise(ctx context.Context, workoutExerciseId string, workoutExerciseRequest openapi.PatchWorkoutExerciseRequest) (openapi.ImplResponse, error) {
	profileId, err := common.ExtractProfileID(ctx)
	if err != nil {
		return common.ErrorResponse(http.StatusUnauthorized, openapi.FORBIDDEN, err.Error())
	}
	if isUUIDValid(workoutExerciseId) {
		return common.ErrorResponse(http.StatusBadRequest, openapi.INVALID_ID, "Workout exercise ID is not a valid UUID")
	}
	workoutExercise, err := h.useCase.Update(ctx, profileId, workoutExerciseId, workoutExerciseRequest)
	if err != nil {
		return common.ErrorResponse(http.StatusInternalServerError, openapi.INTERNAL_SERVER_ERROR, err.Error())
	}
	return openapi.Response(http.StatusCreated, common.ConvertWorkoutExercise(workoutExercise)), nil
}

func (h *workoutExerciseHandler) DeleteWorkoutExercise(ctx context.Context, workoutExerciseId string) (openapi.ImplResponse, error) {
	profileId, err := common.ExtractProfileID(ctx)
	if err != nil {
		return common.ErrorResponse(http.StatusUnauthorized, openapi.FORBIDDEN, err.Error())
	}
	if isUUIDValid(workoutExerciseId) {
		return common.ErrorResponse(http.StatusBadRequest, openapi.INVALID_ID, "Workout exercise ID is not a valid UUID")
	}
	err = h.useCase.Delete(ctx, profileId, workoutExerciseId)
	if err != nil {
		return common.ErrorResponse(http.StatusInternalServerError, openapi.INTERNAL_SERVER_ERROR, err.Error())
	}
	return openapi.Response(http.StatusNoContent, nil), nil
}

func (h *workoutExerciseHandler) ReorderWorkoutExercise(ctx context.Context, workoutExerciseId string, request openapi.ReorderWorkoutExerciseRequest) (openapi.ImplResponse, error) {
	profileId, err := common.ExtractProfileID(ctx)
	if err != nil {
		return common.ErrorResponse(http.StatusUnauthorized, openapi.FORBIDDEN, err.Error())
	}
	err = h.useCase.Reorder(ctx, profileId, workoutExerciseId, request)
	if err != nil {
		return common.ErrorResponse(http.StatusInternalServerError, openapi.INTERNAL_SERVER_ERROR, err.Error())
	}

	return openapi.Response(http.StatusNoContent, nil), nil
}
