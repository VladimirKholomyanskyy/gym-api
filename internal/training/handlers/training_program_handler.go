package handlers

import (
	"context"
	"net/http"

	openapi "github.com/VladimirKholomyanskyy/gym-api/internal/api/go"
	"github.com/VladimirKholomyanskyy/gym-api/internal/common"
	customerrors "github.com/VladimirKholomyanskyy/gym-api/internal/customErrors"
	"github.com/VladimirKholomyanskyy/gym-api/internal/utils"

	"github.com/VladimirKholomyanskyy/gym-api/internal/training/usecase"
)

type TrainingProgramHandler struct {
	useCase usecase.TrainingProgramUseCase
}

func NewTrainingProgramHandler(useCase usecase.TrainingProgramUseCase) openapi.TrainingProgramsAPIServicer {
	return &TrainingProgramHandler{useCase: useCase}
}

func (h *TrainingProgramHandler) ListTrainingPrograms(ctx context.Context, page, pageSize int32) (openapi.ImplResponse, error) {
	profileID, err := common.ExtractProfileID(ctx)
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, openapi.FORBIDDEN, err.Error())
	}
	if !common.IsPageValid(page) {
		return utils.ErrorResponse(http.StatusBadRequest, openapi.INVALID_PAGE_NUMBER, "Page must be greater than 0")
	}

	if !common.IsPageSizeValid(pageSize) {
		return utils.ErrorResponse(http.StatusBadRequest, openapi.INVALID_PAGE_SIZE, "PageSize must be between 1 and 100")
	}
	userPrograms, totalCount, err := h.useCase.List(ctx, profileID, int(page), int(pageSize))
	if err != nil {
		return utils.ErrorResponse(http.StatusInternalServerError, openapi.INTERNAL_SERVER_ERROR, err.Error())
	}
	return openapi.Response(
		http.StatusOK,
		openapi.ListTrainingPrograms200Response{
			TotalItems:  int32(totalCount),
			CurrentPage: page,
			PageSize:    pageSize,
			TotalPages:  utils.CalculateTotalPages(totalCount, pageSize),
			Items:       utils.ConvertTrainingPrograms(userPrograms)}), nil

}

func (h *TrainingProgramHandler) CreateTrainingProgram(ctx context.Context, request openapi.CreateTrainingProgramRequest) (openapi.ImplResponse, error) {
	profileID, err := common.ExtractProfileID(ctx)
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, openapi.FORBIDDEN, err.Error())
	}
	if utils.HasText(&request.Name) {
		return utils.ErrorResponse(http.StatusBadRequest, openapi.INVALID_REQUEST, "Training program name cannot be empty")
	}
	program, err := h.useCase.Create(ctx, profileID, request)
	if err != nil {
		return utils.ErrorResponse(http.StatusInternalServerError, openapi.INTERNAL_SERVER_ERROR, err.Error())
	}
	return openapi.Response(http.StatusCreated, utils.ConvertTrainingProgram(program)), nil
}

func (h *TrainingProgramHandler) GetTrainingProgramById(ctx context.Context, programID string) (openapi.ImplResponse, error) {
	profileID, err := common.ExtractProfileID(ctx)
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, openapi.FORBIDDEN, err.Error())
	}
	if common.IsUUIDValid(programID) {
		return utils.ErrorResponse(http.StatusBadRequest, openapi.INVALID_ID, "Program ID is not a valid UUID")
	}
	userProgram, err := h.useCase.GetByID(ctx, profileID, programID)
	if err != nil {
		if err == customerrors.ErrAccessForbidden {
			return utils.ErrorResponse(http.StatusForbidden, openapi.FORBIDDEN, "Access denied to training program")
		}
		if err == customerrors.ErrEntityNotFound {
			return utils.ErrorResponse(http.StatusNotFound, openapi.RESOURCE_NOT_FOUND, "Scheduled workout not found")
		}
		return utils.ErrorResponse(http.StatusInternalServerError, openapi.INTERNAL_SERVER_ERROR, "Failed to update scheduled workout")
	}
	return openapi.Response(http.StatusOK, utils.ConvertTrainingProgram(userProgram)), nil
}

func (h *TrainingProgramHandler) DeleteTrainingProgram(ctx context.Context, programID string) (openapi.ImplResponse, error) {
	profileID, err := common.ExtractProfileID(ctx)
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, openapi.FORBIDDEN, err.Error())
	}
	if common.IsUUIDValid(programID) {
		return utils.ErrorResponse(http.StatusBadRequest, openapi.INVALID_ID, "Program ID is not a valid UUID")
	}
	err = h.useCase.Delete(ctx, profileID, programID)
	if err != nil {
		if err == customerrors.ErrAccessForbidden {
			return utils.ErrorResponse(http.StatusForbidden, openapi.FORBIDDEN, "Access denied to delete training program")
		}
		if err == customerrors.ErrEntityNotFound {
			return utils.ErrorResponse(http.StatusNotFound, openapi.RESOURCE_NOT_FOUND, "Training program not found")
		}
		return utils.ErrorResponse(http.StatusInternalServerError, openapi.INTERNAL_SERVER_ERROR, "Failed to delete training program")
	}

	return openapi.Response(http.StatusNoContent, nil), nil
}

func (h *TrainingProgramHandler) UpdateTrainingProgram(ctx context.Context, programID string, request openapi.PatchTrainingProgramRequest) (openapi.ImplResponse, error) {
	profileID, err := common.ExtractProfileID(ctx)
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, openapi.FORBIDDEN, err.Error())
	}
	if common.IsUUIDValid(programID) {
		return utils.ErrorResponse(http.StatusBadRequest, openapi.INVALID_ID, "Program ID is not a valid UUID")
	}
	userProgram, err := h.useCase.Update(ctx, profileID, programID, request)
	if err != nil {
		if err == customerrors.ErrAccessForbidden {
			return utils.ErrorResponse(http.StatusForbidden, openapi.FORBIDDEN, "Access denied to update training program")
		}
		if err == customerrors.ErrEntityNotFound {
			return utils.ErrorResponse(http.StatusNotFound, openapi.RESOURCE_NOT_FOUND, "Training program not found")
		}
		return utils.ErrorResponse(http.StatusInternalServerError, openapi.INTERNAL_SERVER_ERROR, "Failed to update training program")
	}
	return openapi.Response(http.StatusOK, utils.ConvertTrainingProgram(userProgram)), nil
}
