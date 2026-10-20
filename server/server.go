package server

import (
	"net/http"
	"os"
	"time"
	"weather-api-wrapper-service/auth"
	"weather-api-wrapper-service/handler"
	"weather-api-wrapper-service/lib/e"
)

func StartServe(h *handler.Handler, tenantAuth *auth.TenantAuth) error {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /weather", h.HandleWeather)
	mux.Handle("GET /v1/weather", tenantAuth.Middleware(http.HandlerFunc(h.HandleV1Weather)))
	mux.Handle("POST /v1/insights", tenantAuth.Middleware(http.HandlerFunc(h.HandleInsights)))
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	server := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	err := server.ListenAndServe()
	if err != nil {
		return e.Wrap("listen and serve: ", err)
	}
	return nil
}
