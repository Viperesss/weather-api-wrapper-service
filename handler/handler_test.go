package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"weather-api-wrapper-service/service"
)

func TestHandleInsightsRejectsInvalidRequests(t *testing.T) {
	h := NewHandler(service.NewService(nil, nil))
	tests := []struct {
		name       string
		body       string
		statusCode int
	}{
		{name: "invalid JSON", body: `{`, statusCode: http.StatusBadRequest},
		{name: "unknown field", body: `{"city":"Berlin","extra":true}`, statusCode: http.StatusBadRequest},
		{name: "missing city", body: `{"business_context":"logistics"}`, statusCode: http.StatusBadRequest},
		{name: "multiple JSON values", body: `{"city":"Berlin"} {}`, statusCode: http.StatusBadRequest},
		{name: "oversized business context", body: `{"city":"Berlin","business_context":"` + strings.Repeat("x", 2001) + `"}`, statusCode: http.StatusBadRequest},
		{name: "oversized request", body: strings.Repeat("x", 8*1024+1), statusCode: http.StatusRequestEntityTooLarge},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/v1/insights", strings.NewReader(test.body))
			recorder := httptest.NewRecorder()
			h.HandleInsights(recorder, req)
			if recorder.Code != test.statusCode {
				t.Fatalf("status = %d, want %d; body: %s", recorder.Code, test.statusCode, recorder.Body.String())
			}
		})
	}
}
