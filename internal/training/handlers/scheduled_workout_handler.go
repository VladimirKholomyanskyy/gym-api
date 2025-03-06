package handlers

import (
	"context"
	"net/http"

	openapi "github.com/VladimirKholomyanskyy/gym-api/internal/api/go"
	"github.com/VladimirKholomyanskyy/gym-api/internal/common"
	"github.com/VladimirKholomyanskyy/gym-api/internal/training/usecase"
)

type scheduledWorkoutsHandler struct {
	useCase usecase.ScheduledWorkoutUseCase
}

func NewScheduledWorkoutsHandler(useCase usecase.ScheduledWorkoutUseCase) openapi.ScheduledWorkoutsAPIServicer {
	return &scheduledWorkoutsHandler{useCase: useCase}
}

func (h *scheduledWorkoutsHandler) GetScheduledWorkouts(ctx context.Context, startDate string, endDate string, page, pageSize int32) (openapi.ImplResponse, error) {
	profileId, err := common.ExtractProfileID(ctx)
	if err != nil {
		return common.ErrorResponse(http.StatusUnauthorized, openapi.FORBIDDEN, err.Error())
	}
	if !isPageValid(page) {
		return common.ErrorResponse(http.StatusBadRequest, openapi.INVALID_PAGE_NUMBER, "page must be greater than 0")
	}
	if !isPageSizeValid(pageSize) {
		return common.ErrorResponse(http.StatusBadRequest, openapi.INVALID_PAGE_SIZE, "pageSize must be between 1 and 100")
	}
	startDateTime, err := common.ParseTime(startDate)
	if err != nil {
		return common.ErrorResponse(http.StatusBadRequest, openapi.INVALID_DATE_FORMAT, "Invalid startDate format")
	}
	endDateTime, err := common.ParseTime(endDate)
	if err != nil {
		return common.ErrorResponse(http.StatusBadRequest, openapi.INVALID_DATE_FORMAT, "Invalid endDate format")
	}

	scheduledWorkouts, totalCount, err := h.useCase.List(ctx, profileId, startDateTime, endDateTime, int(page), int(pageSize))
	if err != nil {
		return common.ErrorResponse(http.StatusInternalServerError, openapi.INTERNAL_SERVER_ERROR, "Failed to fetch scheduled workouts")
	}
	return openapi.Response(
		http.StatusOK,
		openapi.GetScheduledWorkouts200Response{
			TotalItems:  int32(totalCount),
			CurrentPage: page,
			PageSize:    pageSize,
			TotalPages:  common.CalculateTotalPages(totalCount, pageSize),
			Items:       common.ConvertScheduledWorkouts(scheduledWorkouts)}), nil
}

func (h *scheduledWorkoutsHandler) ScheduleWorkout(ctx context.Context, request openapi.CreateScheduledWorkoutRequest) (openapi.ImplResponse, error) {
	profileId, err := common.ExtractProfileID(ctx)
	if err != nil {
		return common.ErrorResponse(http.StatusUnauthorized, openapi.FORBIDDEN, err.Error())
	}
	if isUUIDValid(request.WorkoutId) {
		return common.ErrorResponse(http.StatusBadRequest, openapi.INVALID_ID, "workout ID is not a valid UUID")
	}
	scheduledWorkout, err := h.useCase.Create(ctx, profileId, request)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), nil
	}
	return openapi.Response(http.StatusCreated, common.ConvertScheduledWorkout(scheduledWorkout)), nil
}

func (h *scheduledWorkoutsHandler) GetScheduledWorkout(ctx context.Context, id string) (openapi.ImplResponse, error) {
	profileId, err := common.ExtractProfileID(ctx)
	if err != nil {
		return common.ErrorResponse(http.StatusUnauthorized, openapi.FORBIDDEN, err.Error())
	}
	scheduledWorkout, err := h.useCase.GetByID(ctx, profileId, id)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), nil
	}
	return openapi.Response(http.StatusCreated, common.ConvertScheduledWorkout(scheduledWorkout)), nil
}

func (h *scheduledWorkoutsHandler) UpdateScheduledWorkout(ctx context.Context, id string, request openapi.PatchScheduledWorkoutRequest) (openapi.ImplResponse, error) {
	profileId, err := common.ExtractProfileID(ctx)
	if err != nil {
		return common.ErrorResponse(http.StatusUnauthorized, openapi.FORBIDDEN, err.Error())
	}

	scheduledWorkout, err := h.useCase.Update(ctx, profileId, id, request)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), nil
	}
	return openapi.Response(http.StatusCreated, common.ConvertScheduledWorkout(scheduledWorkout)), nil

}

func (h *scheduledWorkoutsHandler) DeleteScheduledWorkout(ctx context.Context, id string) (openapi.ImplResponse, error) {
	profileId, err := common.ExtractProfileID(ctx)
	if err != nil {
		return common.ErrorResponse(http.StatusUnauthorized, openapi.FORBIDDEN, err.Error())
	}
	err = h.useCase.Delete(ctx, profileId, id)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), nil
	}
	return openapi.Response(http.StatusNoContent, nil), nil
}

// GetNextScheduledWorkout implements openapi.ScheduledWorkoutsAPIServicer.
func (h *scheduledWorkoutsHandler) GetNextScheduledWorkout(ctx context.Context) (openapi.ImplResponse, error) {
	profileId, err := common.ExtractProfileID(ctx)
	if err != nil {
		return common.ErrorResponse(http.StatusUnauthorized, openapi.FORBIDDEN, err.Error())
	}
	scheduledWorkout, err := h.useCase.GetUpcommingScheduledWorkout(ctx, profileId)
	if err != nil {
		return common.ErrorResponse(http.StatusInternalServerError, openapi.INTERNAL_SERVER_ERROR, err.Error())
	}
	if scheduledWorkout == nil {
		return common.ErrorResponse(http.StatusNotFound, openapi.RESOURCE_NOT_FOUND, "No upcomming workouts found")
	}
	return openapi.Response(http.StatusOK, common.ConvertScheduledWorkout(scheduledWorkout)), nil
}
