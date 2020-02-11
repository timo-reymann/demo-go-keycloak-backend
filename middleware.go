package main

import (
	"context"
	"github.com/coreos/go-oidc"
	"net/http"
	"strings"
)

func CreateMiddleware(keycloakUrl string, clientId string) func(next http.HandlerFunc) http.HandlerFunc {
	ctx := context.Background()
	provider, _ := oidc.NewProvider(ctx, keycloakUrl)

	verifier := provider.Verifier(&oidc.Config{
		ClientID: clientId,
	})

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Get auth header
			rawAccessToken := r.Header.Get("Authorization")
			if rawAccessToken == "" {
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte("Authorization header is missing"))
				return
			}

			// Try to extract bearer part
			parts := strings.Split(rawAccessToken, " ")
			if len(parts) != 2 {
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte("Bearer token is missing"))
				return
			}

			// Verify token is signed and valid
			token, err := verifier.Verify(ctx, parts[1])
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte(err.Error()))
				return
			}

			// Extract claims
			claims := &Claims{}
			if err := token.Claims(claims); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(err.Error()))
				return
			}

			// Add meta data to request
			tokenContext := context.WithValue(r.Context(), "token", token)
			claimContext := context.WithValue(tokenContext, "claims", claims)
			next(w, r.WithContext(claimContext))
		}
	}
}