package account

import (
	"context"
	"net/http"

	openapi "github.com/VladimirKholomyanskyy/gym-api/internal/api/go"
	"github.com/VladimirKholomyanskyy/gym-api/internal/common"
)

type settingsHandler struct {
	settingsRepo SettingRepository
}

// NewSettingsHandler creates a new instance of SettingsHandler.
func NewSettingsHandler(settingsRepo SettingRepository) openapi.SettingsAPIServicer {
	return &settingsHandler{settingsRepo: settingsRepo}
}

func (h *settingsHandler) GetSettings(ctx context.Context) (openapi.ImplResponse, error) {
	profileID, err := common.ExtractProfileID(ctx)
	if err != nil {
		return common.ErrorResponse(http.StatusUnauthorized, openapi.UNAUTHORIZED, err.Error())
	}
	settings, err := h.settingsRepo.GetByUserID(ctx, profileID)
	if err != nil {
		return common.ErrorResponse(http.StatusInternalServerError, openapi.INTERNAL_SERVER_ERROR, "Failed to fetch user settings")
	}
	if settings == nil {
		return common.ErrorResponse(http.StatusNotFound, openapi.RESOURCE_NOT_FOUND, "User settings not found")
	}

	return openapi.Response(http.StatusOK, ConvertSettingToOpenAPI(settings)), nil
}
func (h *settingsHandler) UpdateSettings(ctx context.Context, request openapi.PatchSettingsRequest) (openapi.ImplResponse, error) {
	profileID, err := common.ExtractProfileID(ctx)
	if err != nil {
		return common.ErrorResponse(http.StatusUnauthorized, openapi.UNAUTHORIZED, err.Error())
	}
	settings, err := h.settingsRepo.GetByUserID(ctx, profileID)
	if err != nil {
		return common.ErrorResponse(http.StatusInternalServerError, openapi.INTERNAL_SERVER_ERROR, "Failed to fetch user settings")
	}
	if settings == nil {
		return common.ErrorResponse(http.StatusNotFound, openapi.RESOURCE_NOT_FOUND, "User settings not found")
	}

	if request.Language != nil {
		settings.Language = *request.Language
	}
	if request.MeasurementUnits != nil {
		settings.MeasurementUnits = *request.MeasurementUnits
	}
	if request.NotificationsEnabled != nil {
		settings.NotificationsEnabled = *request.NotificationsEnabled
	}
	if request.Timezone != nil {
		settings.Timezone = *request.Timezone
	}

	if err := h.settingsRepo.Update(ctx, settings); err != nil {
		return common.ErrorResponse(http.StatusInternalServerError, openapi.INTERNAL_SERVER_ERROR, "Failed to update user settings")
	}

	return openapi.Response(http.StatusOK, ConvertSettingToOpenAPI(settings)), nil
}
