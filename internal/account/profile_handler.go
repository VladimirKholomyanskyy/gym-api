package account

import (
	"context"
	"net/http"

	openapi "github.com/VladimirKholomyanskyy/gym-api/internal/api/go"
	"github.com/VladimirKholomyanskyy/gym-api/internal/common"
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
		return common.ErrorResponse(http.StatusUnauthorized, openapi.UNAUTHORIZED, err.Error())
	}

	profile, err := h.profileRepo.GetByID(ctx, profileID)
	if err != nil {
		return common.ErrorResponse(http.StatusInternalServerError, openapi.INTERNAL_SERVER_ERROR, "Failed to fetch user profile")
	}
	if profile == nil {
		return common.ErrorResponse(http.StatusNotFound, openapi.RESOURCE_NOT_FOUND, "User profile not found")
	}

	return openapi.Response(http.StatusOK, ConvertProfileToOpenAPI(profile)), nil
}

func (h *profileHandler) UpdateProfile(ctx context.Context, request openapi.PatchProfileRequest) (openapi.ImplResponse, error) {
	profileID, err := common.ExtractProfileID(ctx)
	if err != nil {
		return common.ErrorResponse(http.StatusUnauthorized, openapi.UNAUTHORIZED, err.Error())
	}

	profile, err := h.profileRepo.GetByID(ctx, profileID)
	if err != nil {
		return common.ErrorResponse(http.StatusInternalServerError, openapi.INTERNAL_SERVER_ERROR, "Failed to fetch user profile")
	}
	if profile == nil {
		return common.ErrorResponse(http.StatusNotFound, openapi.RESOURCE_NOT_FOUND, "User profile not found")
	}

	if request.AvatarUrl != nil {
		profile.AvatarURL = request.AvatarUrl
	}

	if request.Birthday != nil {
		birthday, err := common.ParseTime(*request.Birthday)
		if err != nil {
			return common.ErrorResponse(http.StatusBadRequest, openapi.INVALID_DATE_FORMAT, "Invalid birthday date format")
		}
		profile.Birthday = &birthday
	}

	if request.Height != nil {
		profile.Height = request.Height
	}

	if request.Weight != nil {
		profile.Weight = request.Weight
	}

	if request.Sex != nil {
		profile.Sex = request.Sex
	}

	err = h.profileRepo.Update(ctx, profile)
	if err != nil {
		return common.ErrorResponse(http.StatusInternalServerError, openapi.INTERNAL_SERVER_ERROR, "Failed to update user profile")
	}

	return openapi.Response(http.StatusOK, ConvertProfileToOpenAPI(profile)), nil
}
