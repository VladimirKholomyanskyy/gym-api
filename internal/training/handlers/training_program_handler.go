package handlers

import (
	"context"
	"errors"
	"net/http"

	openapi "github.com/VladimirKholomyanskyy/gym-api/internal/api/go"
	"github.com/VladimirKholomyanskyy/gym-api/internal/common"

	"github.com/VladimirKholomyanskyy/gym-api/internal/training/usecase"
)

type TrainingProgramHandler struct {
	useCase usecase.TrainingProgramUseCase
}

func NewTrainingProgramHandler(useCase usecase.TrainingProgramUseCase) openapi.TrainingProgramsAPIServicer {
	return &TrainingProgramHandler{useCase: useCase}
}

func (h *TrainingProgramHandler) ListTrainingPrograms(ctx context.Context, page, pageSize int32) (openapi.ImplResponse, error) {
	profileId, err := common.ExtractProfileID(ctx)
	if err != nil {
		return common.ErrorResponse(http.StatusUnauthorized, openapi.FORBIDDEN, err.Error())
	}
	if !isPageValid(page) {
		return common.ErrorResponse(http.StatusBadRequest, openapi.INVALID_PAGE_NUMBER, "Page must be greater than 0")
	}

	if !isPageSizeValid(pageSize) {
		return common.ErrorResponse(http.StatusBadRequest, openapi.INVALID_PAGE_SIZE, "PageSize must be between 1 and 100")
	}
	userPrograms, totalCount, err := h.useCase.List(ctx, profileId, int(page), int(pageSize))
	if err != nil {
		return common.ErrorResponse(http.StatusInternalServerError, openapi.INTERNAL_SERVER_ERROR, err.Error())
	}
	return openapi.Response(
		http.StatusOK,
		openapi.ListTrainingPrograms200Response{
			TotalItems:  int32(totalCount),
			CurrentPage: page,
			PageSize:    pageSize,
			TotalPages:  common.CalculateTotalPages(totalCount, pageSize),
			Items:       common.ConvertTrainingPrograms(userPrograms)}), nil

}

func (h *TrainingProgramHandler) CreateTrainingProgram(ctx context.Context, request openapi.CreateTrainingProgramRequest) (openapi.ImplResponse, error) {
	profileId, err := common.ExtractProfileID(ctx)
	if err != nil {
		return common.ErrorResponse(http.StatusUnauthorized, openapi.FORBIDDEN, err.Error())
	}
	if common.HasText(&request.Name) {
		return common.ErrorResponse(http.StatusBadRequest, openapi.INVALID_REQUEST, "Training program name cannot be empty")
	}
	program, err := h.useCase.Create(ctx, profileId, request)
	if err != nil {
		return common.ErrorResponse(http.StatusInternalServerError, openapi.INTERNAL_SERVER_ERROR, err.Error())
	}
	return openapi.Response(http.StatusCreated, common.ConvertTrainingProgram(program)), nil
}

func (h *TrainingProgramHandler) GetTrainingProgramById(ctx context.Context, programId string) (openapi.ImplResponse, error) {
	profileId, err := common.ExtractProfileID(ctx)
	if err != nil {
		return common.ErrorResponse(http.StatusUnauthorized, openapi.FORBIDDEN, err.Error())
	}
	if isUUIDValid(programId) {
		return common.ErrorResponse(http.StatusBadRequest, openapi.INVALID_ID, "Program ID is not a valid UUID")
	}
	userProgram, err := h.useCase.GetByID(ctx, profileId, programId)
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
	return openapi.Response(http.StatusOK, common.ConvertTrainingProgram(userProgram)), nil
}

func (h *TrainingProgramHandler) DeleteTrainingProgram(ctx context.Context, programId string) (openapi.ImplResponse, error) {
	profileId, err := common.ExtractProfileID(ctx)
	if err != nil {
		return common.ErrorResponse(http.StatusUnauthorized, openapi.FORBIDDEN, err.Error())
	}
	if isUUIDValid(programId) {
		return common.ErrorResponse(http.StatusBadRequest, openapi.INVALID_ID, "Program ID is not a valid UUID")
	}
	err = h.useCase.Delete(ctx, profileId, programId)
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

func (h *TrainingProgramHandler) UpdateTrainingProgram(ctx context.Context, programId string, request openapi.PatchTrainingProgramRequest) (openapi.ImplResponse, error) {
	profileId, err := common.ExtractProfileID(ctx)
	if err != nil {
		return common.ErrorResponse(http.StatusUnauthorized, openapi.FORBIDDEN, err.Error())
	}
	if isUUIDValid(programId) {
		return common.ErrorResponse(http.StatusBadRequest, openapi.INVALID_ID, "Program ID is not a valid UUID")
	}
	userProgram, err := h.useCase.Update(ctx, profileId, programId, request)
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
	return openapi.Response(http.StatusOK, common.ConvertTrainingProgram(userProgram)), nil
}
