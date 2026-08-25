package server

import (
	"net/http"
	"weather-api-wrapper-service/handler"
	"weather-api-wrapper-service/lib/e"
)

func StartServe(h *handler.Handler) error {
	http.HandleFunc("/weather", h.HandleWeather)
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		return e.Wrap("listen and serve: ", err)
	}
	return nil
}
