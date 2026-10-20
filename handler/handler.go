package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
	"weather-api-wrapper-service/service"
)

const maxInsightRequestSize = 8 * 1024

type Handler struct {
	service *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{
		service: s,
	}
}

func (h *Handler) HandleWeather(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	if city == "" {
		http.Error(w, "empty field city", http.StatusBadRequest)
		return
	}

	ctx, ctxCancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer ctxCancel()

	weatherStruct, err := h.service.GetWeather(ctx, city)
	if err != nil {
		http.Error(w, "failed to get weather", http.StatusBadGateway)
		log.Println(err)
		return
	}

	log.Println("Исполнен запрос по городу:", city)

	dataOutput(w, weatherStruct)
}

func (h *Handler) HandleV1Weather(w http.ResponseWriter, r *http.Request) {
	city := strings.TrimSpace(r.URL.Query().Get("city"))
	if city == "" || len(city) > 200 {
		http.Error(w, "city is required and must be at most 200 characters", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	report, err := h.service.GetWeather(ctx, city)
	if err != nil {
		log.Println("get tenant weather:", err)
		http.Error(w, "failed to get weather", http.StatusBadGateway)
		return
	}
	dataOutput(w, report)
}

func (h *Handler) HandleInsights(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(io.LimitReader(r.Body, maxInsightRequestSize+1))
	if err != nil {
		http.Error(w, "failed to read request body", http.StatusBadRequest)
		return
	}
	if len(body) > maxInsightRequestSize {
		http.Error(w, "request body exceeds 8 KB", http.StatusRequestEntityTooLarge)
		return
	}

	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.DisallowUnknownFields()

	var input struct {
		City            string `json:"city"`
		BusinessContext string `json:"business_context"`
	}
	if err := decoder.Decode(&input); err != nil {
		http.Error(w, "invalid JSON request body", http.StatusBadRequest)
		return
	}
	if err := decoder.Decode(new(any)); !errors.Is(err, io.EOF) {
		http.Error(w, "request body must contain one JSON object", http.StatusBadRequest)
		return
	}

	input.City = strings.TrimSpace(input.City)
	input.BusinessContext = strings.TrimSpace(input.BusinessContext)
	if input.City == "" || len(input.City) > 200 {
		http.Error(w, "city is required and must be at most 200 characters", http.StatusBadRequest)
		return
	}
	if len(input.BusinessContext) > 2000 {
		http.Error(w, "business_context must be at most 2000 characters", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 45*time.Second)
	defer cancel()

	insight, err := h.service.GetInsight(ctx, input.City, input.BusinessContext)
	if err != nil {
		log.Println("generate tenant weather insight:", err)
		http.Error(w, "failed to generate weather insight", http.StatusBadGateway)
		return
	}
	dataOutput(w, insight)
}

func dataOutput(w http.ResponseWriter, body any) {
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Println("encode response:", err)
	}
}
