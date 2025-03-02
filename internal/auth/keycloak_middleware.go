package auth

import (
	"context"
	"log"
	"net/http"

	"github.com/VladimirKholomyanskyy/gym-api/internal/account"
	"github.com/VladimirKholomyanskyy/gym-api/internal/common"
	"github.com/coreos/go-oidc"
)

type KeycloakMiddleware struct {
	verifier    *oidc.IDTokenVerifier
	clientID    string
	userService *account.UserService
}

func NewKeycloakMiddleware(userService *account.UserService, issuer, clientID string) (*KeycloakMiddleware, error) {
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

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" || len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}
		token := authHeader[7:]

		idToken, err := km.verifier.Verify(r.Context(), token)
		if err != nil {
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
			log.Panicln("User profile doesn't exist, creating a new one.")
			km.userService.CreateUser(claims.Sub)
		}

		ctx := context.WithValue(r.Context(), common.UserIDKey, user.ID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
