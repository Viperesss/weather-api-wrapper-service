package main

import (
	"errors"
	"log"
	"os"
	"weather-api-wrapper-service/ai"
	"weather-api-wrapper-service/auth"
	"weather-api-wrapper-service/handler"
	cache "weather-api-wrapper-service/redis"
	"weather-api-wrapper-service/server"
	"weather-api-wrapper-service/service"
	"weather-api-wrapper-service/weatherAPI"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Printf("could not load .env: %v", err)
	}

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatal("empty API key")
	}
	tenantAuth, err := auth.NewTenantAuth(os.Getenv("TENANT_API_KEYS"))
	if err != nil {
		log.Fatal(err)
	}
	aiClient, err := ai.NewClient(os.Getenv("AI_API_KEY"), os.Getenv("AI_BASE_URL"), os.Getenv("AI_MODEL"))
	if err != nil {
		log.Fatal(err)
	}

	cacheClient := cache.NewClient()
	weatherClient := weatherAPI.NewClient(apiKey)

	service := service.NewServiceWithInsights(weatherClient, cacheClient, aiClient)

	handler := handler.NewHandler(service)

	err = server.StartServe(handler, tenantAuth)
	if err != nil {
		log.Fatal(err)
	}
}
