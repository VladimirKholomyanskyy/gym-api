package common

import (
	"context"

	customerrors "github.com/VladimirKholomyanskyy/gym-api/internal/customErrors"
)

type ProfileID string

// Export the key so it can be used in other packages
const ProfileIDKey ProfileID = "ProfileID"

func ExtractProfileID(ctx context.Context) (string, error) {
	profileID, ok := ctx.Value(ProfileIDKey).(string)
	if !ok || profileID == "" {
		return "", customerrors.ErrUnauthorized
	}
	return profileID, nil
}
