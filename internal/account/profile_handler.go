package account

import (
	"context"
	"net/http"

	openapi "github.com/VladimirKholomyanskyy/gym-api/internal/api/go"
	"github.com/VladimirKholomyanskyy/gym-api/internal/common"
	customerrors "github.com/VladimirKholomyanskyy/gym-api/internal/customErrors"
	"github.com/VladimirKholomyanskyy/gym-api/internal/utils"
)

type profileHandler struct {
	profileRepo ProfileRepository
}

func NewProfileHandler(profileRepo ProfileRepository) openapi.ProfileAPIServicer {
	return &profileHandler{profileRepo: profileRepo}
}

func (h *profileHandler) GetProfile(ctx context.Context) (openapi.ImplResponse, error) {
	profileID, err := common.ExtractProfileID(ctx)
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, openapi.UNAUTHORIZED, err.Error())
	}

	profile, err := h.profileRepo.GetByID(ctx, profileID)
	if err != nil {
		if err == customerrors.ErrEntityNotFound {
			return utils.ErrorResponse(http.StatusNotFound, openapi.RESOURCE_NOT_FOUND, "User profile not found")
		}
		return utils.ErrorResponse(http.StatusInternalServerError, openapi.INTERNAL_SERVER_ERROR, "Failed to fetch user profile")
	}

	return openapi.Response(http.StatusOK, ConvertProfileToOpenAPI(profile)), nil
}

func (h *profileHandler) UpdateProfile(ctx context.Context, request openapi.PatchProfileRequest) (openapi.ImplResponse, error) {
	profileID, err := common.ExtractProfileID(ctx)
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, openapi.UNAUTHORIZED, err.Error())
	}

	profile, err := h.profileRepo.GetByID(ctx, profileID)
	if err != nil {
		if err == customerrors.ErrEntityNotFound {
			return utils.ErrorResponse(http.StatusNotFound, openapi.RESOURCE_NOT_FOUND, "User profile not found")
		}
		return utils.ErrorResponse(http.StatusInternalServerError, openapi.INTERNAL_SERVER_ERROR, "Failed to fetch user profile")
	}
	updates := make(map[string]any)

	if request.AvatarUrl != nil {
		updates["avatar_url"] = request.AvatarUrl
	}

	if request.Birthday != nil {
		birthday, err := utils.ParseTime(*request.Birthday)
		if err != nil {
			return utils.ErrorResponse(http.StatusBadRequest, openapi.INVALID_DATE_FORMAT, "Invalid birthday date format")
		}
		updates["birthday"] = &birthday
	}

	if request.Height != nil {
		updates["height"] = request.Height
	}

	if request.Weight != nil {
		updates["weight"] = request.Weight
	}

	if request.Sex != nil {
		updates["sex"] = request.Sex
	}

	err = h.profileRepo.UpdatePartial(ctx, profile.ID, updates)
	if err != nil {
		return utils.ErrorResponse(http.StatusInternalServerError, openapi.INTERNAL_SERVER_ERROR, "Failed to update user profile")
	}

	return openapi.Response(http.StatusOK, ConvertProfileToOpenAPI(profile)), nil
}
