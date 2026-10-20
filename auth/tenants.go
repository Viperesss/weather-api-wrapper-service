package auth

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type contextKey struct{}

type TenantAuth struct {
	keys map[string][sha256.Size]byte
}

func NewTenantAuth(config string) (*TenantAuth, error) {
	if strings.TrimSpace(config) == "" {
		return nil, errors.New("TENANT_API_KEYS is required")
	}

	keys := make(map[string][sha256.Size]byte)
	for _, entry := range strings.Split(config, ",") {
		tenant, key, ok := strings.Cut(strings.TrimSpace(entry), "=")
		tenant = strings.TrimSpace(tenant)
		key = strings.TrimSpace(key)
		if !ok || tenant == "" || len(key) < 32 {
			return nil, fmt.Errorf("invalid tenant API key entry %q: expected tenant=key with a key of at least 32 characters", entry)
		}
		if _, exists := keys[tenant]; exists {
			return nil, fmt.Errorf("duplicate tenant %q", tenant)
		}
		keys[tenant] = sha256.Sum256([]byte(key))
	}

	return &TenantAuth{keys: keys}, nil
}

func (a *TenantAuth) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
		if !ok || token == "" {
			http.Error(w, "missing or invalid bearer token", http.StatusUnauthorized)
			return
		}

		tokenHash := sha256.Sum256([]byte(token))
		var tenant string
		for candidate, keyHash := range a.keys {
			if subtle.ConstantTimeCompare(tokenHash[:], keyHash[:]) == 1 {
				tenant = candidate
			}
		}
		if tenant == "" {
			http.Error(w, "invalid API key", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), contextKey{}, tenant)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func TenantFromContext(ctx context.Context) (string, bool) {
	tenant, ok := ctx.Value(contextKey{}).(string)
	return tenant, ok
}
