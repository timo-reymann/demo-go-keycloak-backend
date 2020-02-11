package main

import (
	"encoding/json"
	"fmt"
	"github.com/coreos/go-oidc"
	"log"
	"net/http"
	"os"
)

// Claims of keycloak token
type Claims struct {
	Username       string                         `json:"preferred_username"`
	Email          string                         `json:"email"`
	EmailVerified  bool                           `json:"email_verified"`
	Locale         string                         `json:"locale"`
	GivenName      string                         `json:"given_name"`
	LastName       string                         `json:"last_name"`
	Name           string                         `json:"name"`
	Scope          string                         `json:"scope"`
	ResourceAccess map[string]map[string][]string `json:"resource_access"`
}

func main() {
	// Create middleware passing basic config
	keycloakMiddleware := CreateMiddleware(os.Getenv("KEYCLOAK_URL"), os.Getenv("KEYCLOAK_CLIENT"))

	// sample routing
	http.HandleFunc("/hello", keycloakMiddleware(func(writer http.ResponseWriter, request *http.Request) {
		// Get token
		token := request.Context().Value("token").(*oidc.IDToken)

		// Get claims, convert them to json again and use it to output some info
		claims := request.Context().Value("claims").(*Claims)
		_, _ = fmt.Fprintf(writer, "Hello, %s brought to us by %s!\n\n", claims.Username, token.Issuer)
		claimsJson, _ := json.MarshalIndent(&claims, " ", "\t")
		_, _ = fmt.Fprintf(writer, "Token: %s", string(claimsJson))
	}))

	log.Fatal(http.ListenAndServe(":3000", nil))
}
