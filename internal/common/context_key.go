package common

import "context"

type ProfileID string

// Export the key so it can be used in other packages
const ProfileIDKey ProfileID = "ProfileID"

func ExtractProfileID(ctx context.Context) (string, error) {
	profileID, ok := ctx.Value(ProfileIDKey).(string)
	if !ok || profileID == "" {
		return "", ErrUnauthorized
	}
	return profileID, nil
}
