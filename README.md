# Weather API Wrapper Service

A Go service that provides a local HTTP API for retrieving weather data from the Visual Crossing Weather API.

The service uses Redis as a cache to reduce repeated requests to the external API.

## Features

* HTTP endpoint for weather requests
* Integration with Visual Crossing Weather API
* Redis caching with a 15-minute TTL
* Cache-aside strategy
* JSON responses
* Environment-based API key configuration
* Contexts and request timeouts
* Separate layers for HTTP handling, business logic, external API and caching
* Interfaces for API and cache abstractions
* Wrapped errors for easier error tracing

## Architecture

```text
Client
  │
  ▼
HTTP Handler
  │
  ▼
Service
  │
  ├── Redis Cache
  │      │
  │      ├── Cache hit ──────► return cached data
  │      │
  │      └── Cache miss
  │
  ▼
Visual Crossing API
  │
  ▼
Store weather data in Redis
  │
  ▼
Return weather data
```

The service follows a cache-aside strategy:

1. The service checks Redis for the requested city.
2. If cached data exists, it is returned immediately.
3. If the cache does not contain the requested city, the service requests data from Visual Crossing.
4. The received weather data is stored in Redis with a 15-minute TTL.
5. The weather data is returned to the client.

## Project Structure

```text
weather-api-wrapper-service/
│
├── cmd/
│   └── server/
│       └── main.go
│
├── handler/
│   └── handler.go
│
├── server/
│   └── server.go
│
├── service/
│   ├── interfaces.go
│   └── service.go
│
├── redis/
│   └── cache.go
│
├── weather/
│   ├── types.go
│   └── weather.go
│
├── weatherAPI/
│   └── client.go
│
├── lib/
│   └── e/
│       └── e.go
│
├── .env
├── .gitignore
├── go.mod
└── go.sum
```

### Package Responsibilities

* `cmd/server` — application entry point and dependency initialization
* `handler` — HTTP request handling and response encoding
* `server` — HTTP server configuration and startup
* `service` — application/business logic and cache-aside flow
* `redis` — Redis client and cache serialization
* `weatherAPI` — HTTP client for the Visual Crossing API
* `weather` — weather data structures and JSON parsing
* `lib/e` — error wrapping helper

## Requirements

* Go
* Redis
* Visual Crossing API key

## Configuration

Create a `.env` file in the project root:

```env
API_KEY=your_visual_crossing_api_key
```

The `.env` file is excluded from Git.

Redis is expected to run on:

```text
localhost:6379
```

## Running

Start Redis first, then run the service:

```bash
go run ./cmd/server
```

The server will start on:

```text
http://localhost:8080
```

## API

### Get Weather

```http
GET /weather?city=Berlin
```

Example:

```bash
curl "http://localhost:8080/weather?city=Berlin"
```

On the first request for a city, the service retrieves the weather data from Visual Crossing and stores it in Redis.

Subsequent requests for the same city are served from Redis until the 15-minute TTL expires.

## Example Response

```json
{
  "address": "Berlin, Germany",
  "days": [
    {
      "datetime": "2026-08-25",
      "temp": 21.4,
      "conditions": "Partially cloudy"
    }
  ]
}
```

## Technologies

* Go
* `net/http`
* JSON
* Redis
* `github.com/redis/go-redis/v9`
* Visual Crossing Weather API
* `github.com/joho/godotenv`

## What I Practiced

This project was created as a practical Go project to work with:

* HTTP clients and servers
* REST APIs
* JSON encoding and decoding
* Redis and caching
* Cache-aside pattern
* Contexts and timeouts
* Interfaces and dependency abstraction
* Package separation
* Error wrapping
* Basic service-layer architecture

## Project Reference

This project is based on the [Weather API Wrapper Service](https://roadmap.sh/projects/weather-api-wrapper-service) project from roadmap.sh.
