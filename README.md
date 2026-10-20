# Weather Intelligence API

An API-first B2B weather intelligence service for turning forecasts into operational risk insights. It combines Visual Crossing weather data with an OpenAI-compatible model endpoint and provides bearer-token authentication for provisioned tenant organizations.

## What it does

- Returns cached weather forecasts for a location.
- Produces structured, AI-generated business impact assessments and recommendations from a forecast and a customer's business context.
- Authenticates tenant API requests with per-tenant bearer keys.
- Keeps the weather provider and AI provider credentials on the server.
- Preserves the original `GET /weather` endpoint for existing clients. New integrations should use the authenticated `/v1` endpoints.

This repository currently implements the API service and environment-based tenant provisioning. It does not include a customer dashboard, self-service signup, persistent account management, billing, or durable usage metering.

## Requirements

- Go 1.25.4 or later
- Redis at `localhost:6379`
- A Visual Crossing API key
- An API key for an OpenAI-compatible chat-completions provider

## Configuration

Copy `.env.example` to `.env` and fill in the provider credentials. Set `TENANT_API_KEYS` to one or more comma-separated `tenant-name=secret` pairs. Each secret must be at least 32 characters; generate a unique high-entropy value for every tenant and provide it to that tenant as its API key.

```env
API_KEY=your_visual_crossing_api_key
AI_API_KEY=your_ai_provider_api_key
AI_BASE_URL=https://api.openai.com/v1
AI_MODEL=gpt-4o-mini
TENANT_API_KEYS=acme=replace_with_a_unique_secret_of_at_least_32_characters
PORT=8080
```

`AI_BASE_URL` and `AI_MODEL` are optional; the values above are the defaults. The AI base URL must be an absolute HTTP(S) URL. `.env` is excluded from Git.

## Run

Start Redis, configure the environment, then run:

```powershell
go run ./cmd/server
```

The server listens on `http://localhost:8080` by default; set `PORT` to override the port. `GET /healthz` is available for health checks.

## API

All `/v1` routes require `Authorization: Bearer <tenant-api-key>`.

### Get a forecast

```http
GET /v1/weather?city=Berlin
Authorization: Bearer <tenant-api-key>
```

The weather cache is shared by location and retained for 15 minutes. Example response:

```json
{
  "address": "Berlin, Germany",
  "days": [
    {
      "datetime": "2026-10-12",
      "temp": 13.4,
      "conditions": "Partially cloudy"
    }
  ]
}
```

### Generate business insights

```http
POST /v1/insights
Authorization: Bearer <tenant-api-key>
Content-Type: application/json
```

Request:

```json
{
  "city": "Berlin",
  "business_context": "We operate regional delivery routes and need to protect drivers and avoid late shipments."
}
```

`business_context` is optional and limited to 2,000 characters. The JSON request body is limited to 8 KB. Example response:

```json
{
  "risk_level": "high",
  "summary": "Heavy rain may disrupt some deliveries.",
  "impacts": [
    "Longer travel times on exposed routes"
  ],
  "recommendations": [
    "Review dispatch timing and notify affected customers."
  ]
}
```

Risk levels are `low`, `moderate`, `high`, or `critical`. AI output is advisory and should be validated against operational requirements.

### Legacy forecast endpoint

`GET /weather?city=Berlin` remains available without tenant authentication for backward compatibility. It is not intended as the B2B integration surface; migrate clients to `/v1/weather`.

## Architecture

```text
HTTP API
  ├── legacy /weather
  └── authenticated /v1 routes
        ├── tenant bearer-key authentication
        ├── service layer and Redis forecast cache
        ├── Visual Crossing forecast client
        └── OpenAI-compatible AI insight client
```

Weather provider and AI provider calls use request contexts and bounded HTTP timeouts. Provider errors are logged server-side and returned as generic API errors; provider credentials and response bodies are not returned to clients.
