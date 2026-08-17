package main

import (
	"log"
	"os"
	cache "weather-api-wrapper-service/redis"
	"weather-api-wrapper-service/weather"

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
	weatherClient := weather.NewClient(cacheClient, apiKey)

	err := weatherClient.StartServe()
	if err != nil {
		log.Fatal(err)
	}
}
