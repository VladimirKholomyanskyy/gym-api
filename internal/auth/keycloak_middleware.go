package auth

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/VladimirKholomyanskyy/gym-api/internal/account"
	openapi "github.com/VladimirKholomyanskyy/gym-api/internal/api/go"
	"github.com/VladimirKholomyanskyy/gym-api/internal/common"
	"github.com/coreos/go-oidc"
)

type KeycloakMiddleware struct {
	verifier    *oidc.IDTokenVerifier
	clientID    string
	profileRepo account.ProfileRepository
}

func NewKeycloakMiddleware(profileRepo account.ProfileRepository, issuer, clientID string) (*KeycloakMiddleware, error) {
	provider, err := oidc.NewProvider(context.Background(), issuer)
	if err != nil {
		return nil, err
	}

	verifier := provider.Verifier(&oidc.Config{ClientID: clientID})
	return &KeycloakMiddleware{verifier: verifier, clientID: clientID, profileRepo: profileRepo}, nil
}

// Helper function to write error responses in JSON format
func writeErrorResponse(w http.ResponseWriter, statusCode int, errorCode string, message string, details []openapi.ErrorResponseDetailsInner) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errorResponse := openapi.ErrorResponse{
		ErrorCode: openapi.ErrorCodes(errorCode),
		Message:   message,
		Details:   details,
	}

	json.NewEncoder(w).Encode(errorResponse)
}

func (km *KeycloakMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			writeErrorResponse(w, http.StatusUnauthorized, "AUTHORIZATION_ERROR", "Missing or invalid Authorization header", nil)
			return
		}

		token := authHeader[7:]
		idToken, err := km.verifier.Verify(r.Context(), token)
		if err != nil {
			writeErrorResponse(w, http.StatusUnauthorized, "INVALID_TOKEN", "Invalid token", nil)
			return
		}

		var claims struct {
			Sub string `json:"sub"`
		}
		if err := idToken.Claims(&claims); err != nil {
			writeErrorResponse(w, http.StatusInternalServerError, "TOKEN_CLAIMS_ERROR", "Failed to parse claims", nil)
			return
		}

		profile, err := km.profileRepo.GetByExternalID(r.Context(), claims.Sub)
		if err != nil {
			log.Println("User profile doesn't exist, creating a new one.")
			profile = &account.Profile{ExternalID: claims.Sub}
			err = km.profileRepo.Create(r.Context(), profile)
			if err != nil {
				writeErrorResponse(w, http.StatusInternalServerError, "PROFILE_CREATION_ERROR", "Failed to create user profile", nil)
				return
			}
		}

		ctx := context.WithValue(r.Context(), common.ProfileIDKey, profile.ID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
