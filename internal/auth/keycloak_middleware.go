package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VladimirKholomyanskyy/gym-api/internal/handlers"
	"github.com/VladimirKholomyanskyy/gym-api/internal/service"
	"github.com/coreos/go-oidc"
)

type KeycloakMiddleware struct {
	verifier    *oidc.IDTokenVerifier
	clientID    string
	userService *service.UserService
}

func NewKeycloakMiddleware(userService *service.UserService, issuer, clientID string) (*KeycloakMiddleware, error) {
	// Create OIDC provider
	provider, err := oidc.NewProvider(context.Background(), issuer)
	if err != nil {
		return nil, err
	}

	// Configure the verifier
	verifier := provider.Verifier(&oidc.Config{ClientID: clientID})
	return &KeycloakMiddleware{verifier: verifier, clientID: clientID, userService: userService}, nil
}

func (km *KeycloakMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}
		token := authHeader[7:]

		// Verify token
		idToken, err := km.verifier.Verify(r.Context(), token)
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		var claims struct {
			Sub string `json:"sub"`
		}
		if err := idToken.Claims(&claims); err != nil {
			http.Error(w, "Failed to parse claims", http.StatusInternalServerError)
			return
		}
		user, err := km.userService.FindUserByExternalID(claims.Sub)
		if err != nil {
			km.userService.CreateUser(claims.Sub)
		}
		// Save user ID in the context
		ctx := context.WithValue(r.Context(), handlers.UserIDKey, user.ID)

		// Token is valid; proceed with the next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
