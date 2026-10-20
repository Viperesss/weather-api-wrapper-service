package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

const tenantKey = "local-test-tenant-key-1234567890"

func TestNewTenantAuthRejectsInvalidConfiguration(t *testing.T) {
	for _, config := range []string{"", "tenant=short", "tenant,key"} {
		if _, err := NewTenantAuth(config); err == nil {
			t.Errorf("NewTenantAuth(%q) succeeded, want error", config)
		}
	}
}

func TestMiddlewareAuthenticatesTenant(t *testing.T) {
	authenticator, err := NewTenantAuth("acme=" + tenantKey)
	if err != nil {
		t.Fatal(err)
	}

	handler := authenticator.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tenant, ok := TenantFromContext(r.Context())
		if !ok || tenant != "acme" {
			t.Errorf("tenant from context = %q, %v; want acme, true", tenant, ok)
		}
		w.WriteHeader(http.StatusNoContent)
	}))

	for _, test := range []struct {
		name       string
		token      string
		statusCode int
	}{
		{name: "missing token", statusCode: http.StatusUnauthorized},
		{name: "invalid token", token: "invalid-token-that-is-long-enough-for-the-test", statusCode: http.StatusUnauthorized},
		{name: "valid token", token: tenantKey, statusCode: http.StatusNoContent},
	} {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/v1/weather", nil)
			if test.token != "" {
				req.Header.Set("Authorization", "Bearer "+test.token)
			}
			recorder := httptest.NewRecorder()
			handler.ServeHTTP(recorder, req)
			if recorder.Code != test.statusCode {
				t.Fatalf("status = %d, want %d", recorder.Code, test.statusCode)
			}
		})
	}
}
