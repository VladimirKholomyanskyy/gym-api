package account

import (
	openapi "github.com/VladimirKholomyanskyy/gym-api/internal/api/go"
	"github.com/VladimirKholomyanskyy/gym-api/internal/common"
)

// ConvertProfileToOpenAPI converts a GORM Profile model to an OpenAPI Profile model
func ConvertProfileToOpenAPI(profile *Profile) *openapi.Profile {
	if profile == nil {
		return nil
	}

	apiProfile := &openapi.Profile{
		Id: profile.ID,
	}

	if profile.Sex != nil {
		converted := openapi.Sex(*profile.Sex)
		apiProfile.Sex = &converted
	}
	if profile.AvatarURL != nil {
		apiProfile.AvatarUrl = profile.AvatarURL
	}
	if profile.Birthday != nil {
		converted := common.FormatTime(profile.Birthday)
		apiProfile.Birthday = &converted
	}

	if profile.Weight != nil {
		apiProfile.Weight = profile.Weight
	}

	if profile.Height != nil {
		apiProfile.Height = profile.Height
	}

	return apiProfile
}
func ConvertSettingToOpenAPI(setting *Setting) openapi.Settings {
	return openapi.Settings{
		Id:                   setting.ID,
		Language:             setting.Language,
		Timezone:             setting.Timezone,
		MeasurementUnits:     (*openapi.MeasurementUnits)(&setting.MeasurementUnits),
		NotificationsEnabled: setting.NotificationsEnabled,
	}
}

// formatTime safely converts *time.Time to string (ISO 8601 format: YYYY-MM-DD)
