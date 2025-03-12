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
	if !common.IsPageValid(page) {
		return utils.ErrorResponse(http.StatusBadRequest, openapi.INVALID_PAGE_NUMBER, "page must be greater than 0")
	}

	if !common.IsPageSizeValid(pageSize) {
		return utils.ErrorResponse(http.StatusBadRequest, openapi.INVALID_PAGE_SIZE, "pageSize must be between 1 and 100")
	}

	exercises, totalCount, err := h.useCase.List(ctx, int(page), int(pageSize))
	if err != nil {
		return utils.ErrorResponse(http.StatusInternalServerError, openapi.INTERNAL_SERVER_ERROR, "failed to list exercises")
	}

	return openapi.Response(
		http.StatusOK,
		openapi.ListExercises200Response{
			TotalItems:  int32(totalCount),
			CurrentPage: page,
			PageSize:    pageSize,
			TotalPages:  utils.CalculateTotalPages(totalCount, pageSize),
			Items:       utils.ConvertExercises(exercises)}), nil
}

// GetExerciseById returns an exercise by its ID
func (h *exerciseHandler) GetExerciseById(ctx context.Context, exerciseId string) (openapi.ImplResponse, error) {
	if !common.IsUUIDValid(exerciseId) {
		return utils.ErrorResponse(http.StatusBadRequest, openapi.INVALID_ID, "exercise ID is not a valid UUID")
	}

	exercise, err := h.useCase.GetByID(ctx, exerciseId)
	if err != nil {
		if err == customerrors.ErrEntityNotFound {
			return utils.ErrorResponse(http.StatusNotFound, openapi.RESOURCE_NOT_FOUND, "Exercise not found")
		}
		return utils.ErrorResponse(http.StatusInternalServerError, openapi.INTERNAL_SERVER_ERROR, "Failed to retrieve exercise")
	}

	return openapi.Response(http.StatusOK, utils.ConvertExercise(exercise)), nil
}
