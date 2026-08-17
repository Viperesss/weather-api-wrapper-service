# Weather API Wrapper Service

A simple Go service that provides a local HTTP API for retrieving weather data from the Visual Crossing Weather API.

The project also uses Redis to cache weather responses and reduce unnecessary requests to the external API.

## Features

* HTTP endpoint for weather requests
* Integration with Visual Crossing Weather API
* Redis caching with a 15-minute TTL
* Cache-aside strategy
* JSON response
* Environment-based API key configuration
* Separate packages for weather logic, Redis cache, and error handling

## Architecture

```text
Client
  │
  ▼
Weather API Wrapper
  │
  ├── Redis
  │    ├── Cache hit → return cached data
  │    └── Cache miss
  │
  ▼
Visual Crossing API
  │
  ▼
Store response in Redis
  │
  ▼
Return weather data
```

Redis is used as an optimization layer. If Redis is unavailable, the service can still request weather data directly from the external API.

## Requirements

* Go 1.XX+
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

### Get weather

```http
GET /weather?city=Berlin
```

Example:

```bash
curl "http://localhost:8080/weather?city=Berlin"
```

The first request for a city retrieves data from Visual Crossing and stores the response in Redis.

Subsequent requests for the same city use the cached response until the 15-minute TTL expires.

## Project Structure

```text
weather-api-wrapper-service/
├── cmd/
│   └── server/
│       └── main.go
├── weather/
│   ├── types.go
│   └── weather.go
├── redis/
│   └── cache.go
├── lib/
│   └── e/
│       └── e.go
├── .env
├── .gitignore
├── go.mod
└── go.sum
```

## Technologies

* Go
* `net/http`
* JSON
* Redis
* `go-redis`
* Visual Crossing Weather API
* `godotenv`

## Purpose

This project was created as a practical Go project to work with:

* HTTP clients and servers
* REST APIs
* JSON encoding and decoding
* Redis and caching
* Contexts and timeouts
* Package separation
* Error handling
* Basic service architecture

https://roadmap.sh/projects/weather-api-wrapper-service
