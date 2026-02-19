# Plan: Smart Order Routing Engine

## Architecture
- Go service exposing /api/v1 routes.
- Internal gRPC for auxiliary services.
- Postgres for rules and audits.
- Redis for health and latency cache.

## Routing Logic
- Score destinations by latency, error rate, and SLA.
- Use circuit breaker and fallback list.

## Observability
- JSON logs with traceId, routeId, destination.
- OpenTelemetry traces for routing decisions.

## Security
- Input validation on all endpoints.
- Rate limiting for /routes.
- Audit trail for rule changes.

## Feature Flags
- routing_algo_v2
- regional_cache_enabled
- fallback_strict_mode

## Local Dev and CI
- Docker Compose for Postgres and Redis.
- CI runs unit, integration, and latency tests.

## ADRs
- gRPC internal contracts
- Regional cache strategy
