package main

import (
	"log"
	"os"
	"weather-api-wrapper-service/handler"
	cache "weather-api-wrapper-service/redis"
	"weather-api-wrapper-service/server"
	"weather-api-wrapper-service/service"
	"weather-api-wrapper-service/weatherAPI"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("failed load .env")
	}

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatal("empty API key")
	}

	cacheClient := cache.NewClient()
	weatherClient := weatherAPI.NewClient(apiKey)

	service := service.NewService(weatherClient, cacheClient)

	handler := handler.NewHandler(service)

	err := server.StartServe(handler)
	if err != nil {
		log.Fatal(err)
	}
}
