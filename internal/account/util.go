package account

import (
	openapi "github.com/VladimirKholomyanskyy/gym-api/internal/api/go"
	"github.com/VladimirKholomyanskyy/gym-api/internal/utils"
)

// ConvertProfileToOpenAPI converts a GORM Profile model to an OpenAPI Profile model
func ConvertProfileToOpenAPI(profile *Profile) openapi.Profile {
	converted := utils.FormatTime(profile.Birthday)
	return openapi.Profile{
		Id:        profile.ID,
		Sex:       profile.Sex,
		Weight:    profile.Weight,
		Height:    profile.Height,
		Birthday:  &converted,
		AvatarUrl: profile.AvatarURL,
	}
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
