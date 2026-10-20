package ai

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"weather-api-wrapper-service/weather"
)

func TestGenerateInsightUsesProviderAndValidatesResult(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/chat/completions" {
			t.Errorf("request path = %q", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer test-provider-key" {
			t.Errorf("authorization header was not set")
		}
		var request completionRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Errorf("decode request: %v", err)
		}
		if request.Model != "test-model" || len(request.Messages) != 2 {
			t.Errorf("unexpected completion request: %+v", request)
		}
		_ = json.NewEncoder(w).Encode(completionResponse{
			Choices: []struct {
				Message message `json:"message"`
			}{{Message: message{
				Content: `{"risk_level":"high","summary":"Heavy rain may delay deliveries.","impacts":["Delivery delays"],"recommendations":["Review dispatch timing"]}`,
			}}},
		})
	}))
	defer server.Close()

	client, err := NewClient("test-provider-key", server.URL+"/v1", "test-model")
	if err != nil {
		t.Fatal(err)
	}
	insight, err := client.GenerateInsight(context.Background(), weather.WeatherReport{
		Address: "Berlin",
		Days:    []weather.Day{{Datetime: "2026-10-12", Conditions: "Rain"}},
	}, "Regional delivery business")
	if err != nil {
		t.Fatal(err)
	}
	if insight.RiskLevel != "high" || len(insight.Recommendations) != 1 {
		t.Fatalf("unexpected insight: %+v", insight)
	}
}

func TestNewClientRequiresValidConfiguration(t *testing.T) {
	if _, err := NewClient("", "", ""); err == nil {
		t.Fatal("NewClient with empty API key succeeded, want error")
	}
	if _, err := NewClient("key", "not-a-url", "model"); err == nil {
		t.Fatal("NewClient with invalid base URL succeeded, want error")
	}
}
