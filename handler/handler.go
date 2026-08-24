package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
	"weather-api-wrapper-service/lib/e"
	"weather-api-wrapper-service/service"
)

type Handler struct {
	Service service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{
		Service: *s,
	}
}

func (h *Handler) HandleWeather(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	if city == "" {
		http.Error(w, "empty field city", http.StatusBadRequest)
		return
	}

	ctx, ctxCancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer ctxCancel()

	body, err := h.Service.GetWeather(ctx, city)
	if err != nil {
		http.Error(w, "failed to get weather", http.StatusBadGateway)
		log.Println(err)
		return
	}

	weatherStruct, err := getData(body)
	if err != nil {
		http.Error(w, "failed to get data from request", http.StatusBadGateway)
		log.Println(err)
		return
	}

	dataOutput(w, weatherStruct)
}

func getData(body []byte) (WeatherReport, error) {
	var ww WeatherReport

	if err := json.Unmarshal(body, &ww); err != nil {
		return WeatherReport{}, e.Wrap("unmarshal body", err)
	}

	return ww, nil
}

func dataOutput(w http.ResponseWriter, body any) {
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Println("encode response:", err)
	}
}

// func handleWeather(w http.ResponseWriter, r *http.Request) {
// 	city := r.URL.Query().Get("city")
// 	if city == "" {
// 		http.Error(w, "empty field city", http.StatusBadRequest)
// 		return
// 	}

// 	//Service

// 	weatherStruct, err := getData(body)
// 	if err != nil {
// 		http.Error(w, "failed to get data from request", http.StatusBadGateway)
// 		log.Println(err)
// 		return
// 	}

// 	dataOutput(w, weatherStruct)
// }

// package weatherAPI

// type WeatherReport struct {
// 	Address string `json:"address"`
// 	Days    []Day  `json:"days"`
// }

// type Day struct {
// 	Datetime   string  `json:"datetime"`
// 	Temp       float64 `json:"temp"`
// 	Conditions string  `json:"conditions"`
// }
