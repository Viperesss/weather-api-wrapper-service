package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
	"weather-api-wrapper-service/service"
)

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

func dataOutput(w http.ResponseWriter, body any) {
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Println("encode response:", err)
	}
}
