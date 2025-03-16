package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/VladimirKholomyanskyy/gym-api/internal/account"
	"github.com/VladimirKholomyanskyy/gym-api/internal/common"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
)

// CognitoMiddleware handles authentication using AWS Cognito tokens
type CognitoMiddleware struct {
	userPoolID  string
	region      string
	clientID    string
	profileRepo account.ProfileRepository
	jwkCache    jwk.Set
	cacheTTL    time.Time
}

// NewCognitoMiddleware creates a new instance of CognitoMiddleware
func NewCognitoMiddleware(profileRepo account.ProfileRepository, userPoolID, region, clientID string) (*CognitoMiddleware, error) {
	// Fetch the JWKs from Cognito
	keySet, err := fetchJWKS(userPoolID, region)
	if err != nil {
		return nil, err
	}

	return &CognitoMiddleware{
		userPoolID:  userPoolID,
		region:      region,
		clientID:    clientID,
		profileRepo: profileRepo,
		jwkCache:    keySet,
		cacheTTL:    time.Now().Add(24 * time.Hour), // Cache JWKs for 24 hours
	}, nil
}

// fetchJWKS fetches the JSON Web Key Set from Cognito
func fetchJWKS(userPoolID, region string) (jwk.Set, error) {
	jwksURL := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json", region, userPoolID)
	keySet, err := jwk.Fetch(context.Background(), jwksURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JWKs: %w", err)
	}
	return keySet, nil
}

// Authenticate middleware function to authenticate requests
func (cm *CognitoMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if JWKS cache needs refreshing
		if time.Now().After(cm.cacheTTL) {
			keySet, err := fetchJWKS(cm.userPoolID, cm.region)
			if err != nil {
				log.Printf("Failed to refresh JWKS: %v", err)
			} else {
				cm.jwkCache = keySet
				cm.cacheTTL = time.Now().Add(24 * time.Hour)
			}
		}

		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || len(authHeader) < 7 || !strings.HasPrefix(authHeader, "Bearer ") {
			writeErrorResponse(w, http.StatusUnauthorized, "AUTHORIZATION_ERROR", "Missing or invalid Authorization header", nil)
			return
		}

		tokenString := authHeader[7:]

		// Parse and verify the JWT using the JWKs
		token, err := jwt.Parse(
			[]byte(tokenString),
			jwt.WithKeySet(cm.jwkCache),
			jwt.WithValidate(true),
			jwt.WithIssuer(fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s", cm.region, cm.userPoolID)),
			jwt.WithAudience(cm.clientID),
		)
		if err != nil {
			writeErrorResponse(w, http.StatusUnauthorized, "INVALID_TOKEN", "Invalid token", nil)
			return
		}

		// Extract the subject (user ID) from the token
		sub, ok := token.Get("sub")
		if !ok {
			writeErrorResponse(w, http.StatusUnauthorized, "TOKEN_CLAIMS_ERROR", "Missing sub claim in token", nil)
			return
		}

		subStr, ok := sub.(string)
		if !ok {
			writeErrorResponse(w, http.StatusUnauthorized, "TOKEN_CLAIMS_ERROR", "Invalid sub claim in token", nil)
			return
		}

		// Get or create user profile
		profile, err := cm.profileRepo.GetByExternalID(r.Context(), subStr)
		if err != nil {
			log.Println("User profile doesn't exist, creating a new one.")
			profile = &account.Profile{ExternalID: subStr}
			err = cm.profileRepo.Create(r.Context(), profile)
			if err != nil {
				writeErrorResponse(w, http.StatusInternalServerError, "PROFILE_CREATION_ERROR", "Failed to create user profile", nil)
				return
			}

		}

		// Add profile ID to the request context
		ctx := context.WithValue(r.Context(), common.ProfileIDKey, profile.ID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
