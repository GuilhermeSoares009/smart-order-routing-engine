# Smart Order Routing Engine

Go service that decides which target should receive an order, given latency, availability and priority for each candidate, and exposes both the decision and the reasoning behind it through a REST API.

## How the decision works
- Candidates below a minimum availability threshold (0.5) are discarded.
Among the healthy ones the lowest latency wins. Ties are broken first by higher availability, then by lower priority number.
When no candidate is healthy the engine still returns the best of the unhealthy set and flags the decision as a fallback, so the caller can react instead of receiving an error.
Every decision carries a machine-readable reason: best-latency or fallback-no-healthy-targets.

## Metric smoothing

Reported latency and availability are merged with the last values seen for the same target inside a TTL window, so a single noisy sample does not flip the routing decision. The cache is in-memory and concurrency-safe.

## Audit trail

Every routing decision is recorded with timestamp, route id, order id, chosen target, reason, fallback flag, score and how many candidates were considered. Entries live in a bounded in-memory store that keeps the last 1000 decisions and can be read back through the API.

## API
- GET /api/v1/health - liveness check
POST /api/v1/routes - submit candidate targets and receive the routing decision
GET /api/v1/audit/routes - list recent routing decisions

Requests are rate limited per client IP, configurable through RATE_LIMIT_PER_MIN.

## Observability

Structured request logging, OpenTelemetry spans and per-endpoint duration histograms.

## Running locally

```bash
docker compose up -d --build
curl http://localhost:8082/api/v1/health
```

The compose stack also provisions PostgreSQL on port 5434 and Redis on port 6381 for local development. The current implementation keeps routing metrics and the audit trail in memory.

## Tests

```bash
go test ./...
```

## Stack

Go, net/http, OpenTelemetry, Docker Compose and GitHub Actions.
