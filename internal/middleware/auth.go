package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

type AuthMiddleware struct {
	BaseURL       string
	WebKeySetsUrl string
	Realm         string
	KeySet        jwk.Set
	Issuer        string
}

func GetUserDocIDFromContext(ctx context.Context) (string, error) {
	userId, ok := ctx.Value("userDocId").(string)
	if !ok {
		return "", fmt.Errorf("user identifier not found in request context")
	}
	return userId, nil

}

func GetUserUUIDFromContext(ctx context.Context) (string, error) {
	userUUID, ok := ctx.Value("userUUID").(string)
	if !ok {
		return "", fmt.Errorf("user UUID not found in request context")
	}
	return userUUID, nil

}

func NewAuthMiddleware(baseURL string, realm string, hostName string) *AuthMiddleware {
	WebKeySetsUrl := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/certs", baseURL, realm)

	keySet, err := jwk.Fetch(context.Background(), WebKeySetsUrl)
	if err != nil {
		fmt.Println(err)
		fmt.Println("All authenticated requests will be rejected until the JWKS is available")
	}

	return &AuthMiddleware{
		BaseURL:       baseURL,
		WebKeySetsUrl: WebKeySetsUrl,
		Realm:         realm,
		KeySet:        keySet,
		Issuer:        fmt.Sprintf("%s/realms/%s", hostName, realm),
	}
}

func (kam *AuthMiddleware) RefreshKeySet(interval time.Duration) {
	var mu sync.Mutex

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			keySet, err := jwk.Fetch(context.Background(), kam.WebKeySetsUrl)
			if err != nil {
				fmt.Printf("Error fetching JWKS from %s: \n", kam.WebKeySetsUrl, err)
				continue
			}

			mu.Lock()
			kam.KeySet = keySet
			mu.Unlock()
			fmt.Printf("Refreshed JWKS from %s\n", kam.WebKeySetsUrl)
		}
	}()
}

func WriteError(w http.ResponseWriter, httpErr []map[string]string) {
	w.WriteHeader(http.StatusUnauthorized)
	res, err := json.Marshal(httpErr)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(res)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (kam *AuthMiddleware) Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if kam.KeySet == nil {
			httpError := []map[string]string{
				{
					"error":             "key_set_not_available",
					"error_description": "Key set is not available, please try again later",
				},
			}
			WriteError(w, httpError)
			return
		}
		ah := r.Header.Get("Authorization")
		if ah == "" {
			httpError := []map[string]string{
				{
					"error":             "missing_authorization_header",
					"error_description": "Missing Authorization header in request",
				},
			}
			WriteError(w, httpError)
			return
		}

		ts := strings.TrimPrefix(ah, "Bearer ")
		if ts == ah {
			httpError := []map[string]string{
				{
					"error":             "bad_formatted_authorization_header",
					"error_description": "Missing bearer prefix in authorization header",
				},
			}
			WriteError(w, httpError)
			return
		}

		token, err := jwt.Parse([]byte(ts), jwt.WithKeySet(kam.KeySet), jwt.WithValidate(true))
		if err != nil {
			httpError := []map[string]string{
				{
					"error":             "invalid_token",
					"error_description": "Invalid token or signature",
				},
			}
			WriteError(w, httpError)
			return
		}

		if token.Issuer() != kam.Issuer {
			httpError := []map[string]string{
				{
					"error":             "invalid_token_issuer",
					"error_description": "Invalid issuer in token",
				},
			}
			WriteError(w, httpError)
			return
		}

		userId, _ := token.PrivateClaims()["identifier"].(string)
		if userId == "" {
			httpError := []map[string]string{
				{
					"error":             "token_user_id_not_found",
					"error_description": "Token does not contain user identifier",
				},
			}
			WriteError(w, httpError)
			return
		}

		realmAccess, ok := token.PrivateClaims()["realm_access"].(map[string]interface{})
		if !ok {
			httpError := []map[string]string{
				{
					"error":             "token_realm_access_not_found",
					"error_description": "Token does not contain realm access information",
				},
			}
			WriteError(w, httpError)
			return
		}

		ri, ok := realmAccess["roles"].([]interface{})
		if !ok {
			httpError := []map[string]string{
				{
					"error":             "token_roles_not_found",
					"error_description": "Token does not contain roles information",
				},
			}
			WriteError(w, httpError)
			return
		}
		var roles []string
		for _, role := range ri {
			strRole, ok := role.(string)
			if !ok {
				httpError := []map[string]string{
					{
						"error":             "token_invalid_role",
						"error_description": "Token contains an invalid role format",
					},
				}
				WriteError(w, httpError)
				return
			}
			roles = append(roles, strRole)
		}

		userUUID := token.Subject()
		if userUUID == "" {
			httpError := []map[string]string{
				{
					"error":             "token_user_uuid_not_found",
					"error_description": "Token does not contain user UUID",
				},
			}
			WriteError(w, httpError)
			return
		}

		ctx := context.WithValue(r.Context(), "userDocId", userId)
		ctx = context.WithValue(ctx, "roles", roles)
		ctx = context.WithValue(ctx, "userUUID", userUUID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
