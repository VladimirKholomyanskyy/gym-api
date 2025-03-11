package account

import (
	"context"
	"net/http"

	openapi "github.com/VladimirKholomyanskyy/gym-api/internal/api/go"
	"github.com/VladimirKholomyanskyy/gym-api/internal/common"
	customerrors "github.com/VladimirKholomyanskyy/gym-api/internal/customErrors"
	"github.com/VladimirKholomyanskyy/gym-api/internal/utils"
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
		return utils.ErrorResponse(http.StatusUnauthorized, openapi.UNAUTHORIZED, err.Error())
	}
	settings, err := h.settingsRepo.GetByProfileID(ctx, profileID)
	if err != nil {
		if err == customerrors.ErrEntityNotFound {
			return utils.ErrorResponse(http.StatusNotFound, openapi.RESOURCE_NOT_FOUND, "User settings not found")
		}
		return utils.ErrorResponse(http.StatusInternalServerError, openapi.INTERNAL_SERVER_ERROR, "Failed to fetch user settings")
	}

	return openapi.Response(http.StatusOK, ConvertSettingToOpenAPI(settings)), nil
}
func (h *settingsHandler) UpdateSettings(ctx context.Context, request openapi.PatchSettingsRequest) (openapi.ImplResponse, error) {
	profileID, err := common.ExtractProfileID(ctx)
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, openapi.UNAUTHORIZED, err.Error())
	}
	settings, err := h.settingsRepo.GetByProfileID(ctx, profileID)
	if err != nil {
		if err == customerrors.ErrEntityNotFound {
			return utils.ErrorResponse(http.StatusNotFound, openapi.RESOURCE_NOT_FOUND, "User settings not found")
		}
		return utils.ErrorResponse(http.StatusInternalServerError, openapi.INTERNAL_SERVER_ERROR, "Failed to fetch user settings")
	}
	updates := make(map[string]any)
	if request.Language != nil {
		updates["language"] = *request.Language
	}
	if request.MeasurementUnits != nil {
		updates["measurement_units"] = *request.MeasurementUnits
	}
	if request.NotificationsEnabled != nil {
		updates["notifications_enabled"] = *request.NotificationsEnabled
	}
	if request.Timezone != nil {
		updates["timezone"] = *request.Timezone
	}

	if err := h.settingsRepo.UpdatePartial(ctx, settings.ID, updates); err != nil {
		return utils.ErrorResponse(http.StatusInternalServerError, openapi.INTERNAL_SERVER_ERROR, "Failed to update user settings")
	}

	return openapi.Response(http.StatusOK, ConvertSettingToOpenAPI(settings)), nil
}
