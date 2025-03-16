package auth

import (
	"context"
	"net/http"
	"os"

	openapi "github.com/VladimirKholomyanskyy/gym-api/internal/api/go"
)

// AuthHandler handles authentication related endpoints
type AuthHandler struct {
	// Add any dependencies here
}

// NewAuthHandler creates a new instance of AuthHandler
func NewAuthHandler() openapi.AuthAPIServicer {
	return &AuthHandler{}
}

// GetAuthConfig returns the configuration needed for clients to authenticate
func (h *AuthHandler) GetAuthConfig(ctx context.Context) (openapi.ImplResponse, error) {
	config := openapi.AuthConfig{
		UserPoolId:             GetEnvOrDefault("AWS_COGNITO_USER_POOL_ID", ""),
		UserPoolRegion:         GetEnvOrDefault("AWS_COGNITO_REGION", ""),
		UserPoolClientId:       GetEnvOrDefault("AWS_COGNITO_CLIENT_ID", ""),
		IdentityPoolId:         GetEnvOrDefault("AWS_COGNITO_IDENTITY_POOL_ID", ""),
		AuthenticationFlowType: "USER_SRP_AUTH",
	}

	return openapi.Response(http.StatusOK, config), nil
}

// GetEnvOrDefault gets an environment variable or returns a default value
func GetEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
