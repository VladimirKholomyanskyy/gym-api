package handlers

import (
	"context"
	"errors"
	"net/http"

	openapi "github.com/VladimirKholomyanskyy/gym-api/internal/api/go"
	"github.com/VladimirKholomyanskyy/gym-api/internal/common"
	"github.com/VladimirKholomyanskyy/gym-api/internal/training/usecase"
)

// ExerciseHandler handles HTTP requests for exercise resources
type exerciseHandler struct {
	useCase usecase.ExerciseUseCase
}

// NewExerciseHandler creates a new instance of ExerciseHandler
func NewExerciseHandler(useCase usecase.ExerciseUseCase) openapi.ExercisesAPIServicer {
	return &exerciseHandler{useCase: useCase}
}

// ListExercises returns a paginated list of exercises
func (h *exerciseHandler) ListExercises(ctx context.Context, page, pageSize int32) (openapi.ImplResponse, error) {
	if !isPageValid(page) {
		return common.ErrorResponse(http.StatusBadRequest, openapi.INVALID_PAGE_NUMBER, "page must be greater than 0")
	}

	if !isPageSizeValid(pageSize) {
		return common.ErrorResponse(http.StatusBadRequest, openapi.INVALID_PAGE_SIZE, "pageSize must be between 1 and 100")
	}

	exercises, totalCount, err := h.useCase.List(ctx, int(page), int(pageSize))
	if err != nil {
		return common.ErrorResponse(http.StatusInternalServerError, openapi.INTERNAL_SERVER_ERROR, "failed to list exercises")
	}

	return openapi.Response(
		http.StatusOK,
		openapi.ListExercises200Response{
			TotalItems:  int32(totalCount),
			CurrentPage: page,
			PageSize:    pageSize,
			TotalPages:  common.CalculateTotalPages(totalCount, pageSize),
			Items:       common.ConvertExercises(exercises)}), nil
}

// GetExerciseById returns an exercise by its ID
func (h *exerciseHandler) GetExerciseById(ctx context.Context, exerciseId string) (openapi.ImplResponse, error) {
	if isUUIDValid(exerciseId) {
		return common.ErrorResponse(http.StatusBadRequest, openapi.INVALID_ID, "exercise ID is not a valid UUID")
	}

	exercise, err := h.useCase.GetByID(ctx, exerciseId)
	if err != nil {
		switch {
		case errors.Is(err, common.NotFoundError{}):
			return common.ErrorResponse(http.StatusNotFound, openapi.RESOURCE_NOT_FOUND, "exercise not found")
		default:
			return common.ErrorResponse(http.StatusInternalServerError, openapi.INTERNAL_SERVER_ERROR, "failed to retrieve exercise")
		}
	}

	return openapi.Response(http.StatusOK, common.ConvertExercise(exercise)), nil
}
