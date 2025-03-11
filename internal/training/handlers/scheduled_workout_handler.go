package handlers

import (
	"context"
	"net/http"

	openapi "github.com/VladimirKholomyanskyy/gym-api/internal/api/go"
	"github.com/VladimirKholomyanskyy/gym-api/internal/common"
	customerrors "github.com/VladimirKholomyanskyy/gym-api/internal/customErrors"
	"github.com/VladimirKholomyanskyy/gym-api/internal/training/model"
	"github.com/VladimirKholomyanskyy/gym-api/internal/training/usecase"
	"github.com/VladimirKholomyanskyy/gym-api/internal/utils"
)

type scheduledWorkoutsHandler struct {
	useCase usecase.ScheduledWorkoutUseCase
}

func NewScheduledWorkoutsHandler(useCase usecase.ScheduledWorkoutUseCase) openapi.ScheduledWorkoutsAPIServicer {
	return &scheduledWorkoutsHandler{useCase: useCase}
}

func (h *scheduledWorkoutsHandler) ScheduleWorkout(ctx context.Context, request openapi.CreateScheduledWorkoutRequest) (openapi.ImplResponse, error) {
	profileId, err := common.ExtractProfileID(ctx)
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, openapi.FORBIDDEN, err.Error())
	}
	if common.IsUUIDValid(request.WorkoutId) {
		return utils.ErrorResponse(http.StatusBadRequest, openapi.INVALID_ID, "Workout ID is not a valid UUID")
	}
	input := model.CreateScheduledWorkoutInput{ProfileID: profileId, WorkoutID: request.WorkoutId}
	date, err := utils.ParseTime(request.Date)
	if err != nil {
		return utils.ErrorResponse(http.StatusBadRequest, openapi.INVALID_DATE_FORMAT, "Invalid date format")
	}
	input.Date = date
	if request.Notes != nil {
		input.Notes = utils.TrimPointer(request.Notes)
	}
	scheduledWorkout, err := h.useCase.Create(ctx, input)
	if err != nil {
		return utils.ErrorResponse(http.StatusInternalServerError, openapi.INTERNAL_SERVER_ERROR, "Failed to schedule workout")
	}
	return openapi.Response(http.StatusCreated, utils.ConvertScheduledWorkout(scheduledWorkout)), nil
}

func (h *scheduledWorkoutsHandler) GetScheduledWorkouts(ctx context.Context, startDate string, endDate string, page, pageSize int32) (openapi.ImplResponse, error) {
	profileId, err := common.ExtractProfileID(ctx)
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, openapi.FORBIDDEN, err.Error())
	}
	if !common.IsPageValid(page) {
		return utils.ErrorResponse(http.StatusBadRequest, openapi.INVALID_PAGE_NUMBER, "page must be greater than 0")
	}
	if !common.IsPageSizeValid(pageSize) {
		return utils.ErrorResponse(http.StatusBadRequest, openapi.INVALID_PAGE_SIZE, "pageSize must be between 1 and 100")
	}
	startDateTime, err := utils.ParseTime(startDate)
	if err != nil {
		return utils.ErrorResponse(http.StatusBadRequest, openapi.INVALID_DATE_FORMAT, "Invalid startDate format")
	}
	endDateTime, err := utils.ParseTime(endDate)
	if err != nil {
		return utils.ErrorResponse(http.StatusBadRequest, openapi.INVALID_DATE_FORMAT, "Invalid endDate format")
	}
	if !startDateTime.Before(endDateTime) {
		return utils.ErrorResponse(http.StatusBadRequest, openapi.INVALID_DATE_RANGE, "Invalid date range")
	}
	scheduledWorkouts, totalCount, err := h.useCase.List(ctx, profileId, startDateTime, endDateTime, int(page), int(pageSize))
	if err != nil {
		return utils.ErrorResponse(http.StatusInternalServerError, openapi.INTERNAL_SERVER_ERROR, "Failed to fetch scheduled workouts")
	}
	return openapi.Response(
		http.StatusOK,
		openapi.GetScheduledWorkouts200Response{
			TotalItems:  int32(totalCount),
			CurrentPage: page,
			PageSize:    pageSize,
			TotalPages:  utils.CalculateTotalPages(totalCount, pageSize),
			Items:       utils.ConvertScheduledWorkouts(scheduledWorkouts)}), nil
}

func (h *scheduledWorkoutsHandler) GetScheduledWorkout(ctx context.Context, id string) (openapi.ImplResponse, error) {
	profileID, err := common.ExtractProfileID(ctx)
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, openapi.FORBIDDEN, err.Error())
	}
	scheduledWorkout, err := h.useCase.GetByID(ctx, profileID, id)
	if err != nil {
		return utils.ErrorResponse(http.StatusInternalServerError, openapi.INTERNAL_SERVER_ERROR, "Failed to fetch scheduled workout")
	}
	return openapi.Response(http.StatusCreated, utils.ConvertScheduledWorkout(scheduledWorkout)), nil
}

func (h *scheduledWorkoutsHandler) UpdateScheduledWorkout(ctx context.Context, id string, request openapi.PatchScheduledWorkoutRequest) (openapi.ImplResponse, error) {
	profileId, err := common.ExtractProfileID(ctx)
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, openapi.FORBIDDEN, err.Error())
	}
	if !common.IsUUIDValid(id) {
		return utils.ErrorResponse(http.StatusBadRequest, openapi.INVALID_ID, "Invalid scheduled workout id")
	}
	input := model.UpdateScheduledWorkoutInput{ProfileID: profileId, ScheduledWorkoutID: id}
	if request.Date != nil {
		date, err := utils.ParseTime(*request.Date)
		if err != nil {
			return utils.ErrorResponse(http.StatusBadRequest, openapi.INVALID_DATE_FORMAT, "Invalid date format")
		}
		input.Date = &date
	}
	if request.Notes != nil {
		trimmedNotes := utils.TrimPointer(request.Notes)
		input.Notes = &trimmedNotes
	}
	scheduledWorkout, err := h.useCase.Update(ctx, input)
	if err != nil {
		if err == customerrors.ErrAccessForbidden {
			return utils.ErrorResponse(http.StatusForbidden, openapi.FORBIDDEN, err.Error())
		}
		if err == customerrors.ErrEntityNotFound {
			return utils.ErrorResponse(http.StatusNotFound, openapi.RESOURCE_NOT_FOUND, "Scheduled workout not found")
		}
		return utils.ErrorResponse(http.StatusInternalServerError, openapi.INTERNAL_SERVER_ERROR, "Failed to update scheduled workout")
	}
	return openapi.Response(http.StatusCreated, utils.ConvertScheduledWorkout(scheduledWorkout)), nil

}

func (h *scheduledWorkoutsHandler) DeleteScheduledWorkout(ctx context.Context, id string) (openapi.ImplResponse, error) {
	profileId, err := common.ExtractProfileID(ctx)
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, openapi.FORBIDDEN, err.Error())
	}
	if !common.IsUUIDValid(id) {
		return utils.ErrorResponse(http.StatusBadRequest, openapi.INVALID_ID, "Invalid scheduled workout id")
	}
	err = h.useCase.Delete(ctx, profileId, id)
	if err != nil {
		return utils.ErrorResponse(http.StatusInternalServerError, openapi.INTERNAL_SERVER_ERROR, "Failed to delete scheduled workout")
	}
	return openapi.Response(http.StatusNoContent, nil), nil
}

// GetNextScheduledWorkout implements openapi.ScheduledWorkoutsAPIServicer.
func (h *scheduledWorkoutsHandler) GetNextScheduledWorkout(ctx context.Context) (openapi.ImplResponse, error) {
	profileId, err := common.ExtractProfileID(ctx)
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, openapi.FORBIDDEN, err.Error())
	}
	scheduledWorkout, err := h.useCase.GetUpcommingScheduledWorkout(ctx, profileId)
	if err != nil {
		if err == customerrors.ErrEntityNotFound {
			return utils.ErrorResponse(http.StatusNotFound, openapi.RESOURCE_NOT_FOUND, "Scheduled workout not found")
		}
		return utils.ErrorResponse(http.StatusInternalServerError, openapi.INTERNAL_SERVER_ERROR, err.Error())
	}
	return openapi.Response(http.StatusOK, utils.ConvertScheduledWorkout(scheduledWorkout)), nil
}
