# Spec: Smart Order Routing Engine (MVP)

## Product Vision
- Route orders to the best destination with low latency and resilience.
- Use health and latency signals to choose destinations dynamically.
- Provide fallbacks and clear routing rationale.

## User Scenarios and Testing

### User Story 1 - Route an order quickly (Priority: P1)
As a trading system, I want a route decision under tight latency budgets so orders are not delayed.

**Why this priority**: Low latency is the core value.

**Independent Test**: Submit a route request and verify p95 latency under 50ms in local tests.

**Acceptance Scenarios**:
1. **Given** two destinations with different latency, **When** I submit a route request, **Then** the lower latency destination is selected.
2. **Given** a valid order, **When** I submit /routes, **Then** I receive a routeId and destination.

---

### User Story 2 - Fallback on failure (Priority: P2)
As a system operator, I want automatic fallback when the best destination fails.

**Why this priority**: Resilience avoids failed orders.

**Independent Test**: Mark primary destination as failing and verify fallback selection.

**Acceptance Scenarios**:
1. **Given** destination A is failing, **When** I route an order, **Then** destination B is selected.

---

### User Story 3 - Update destination metrics (Priority: P3)
As a monitoring service, I want to publish health and latency metrics for routing decisions.

**Why this priority**: Routing relies on accurate signals.

**Independent Test**: Post metrics and confirm they influence subsequent route selection.

**Acceptance Scenarios**:
1. **Given** updated metrics, **When** routing occurs, **Then** decisions reflect the new metrics.

### Edge Cases
- Metrics missing or stale.
- Invalid route request payload.
- All destinations failing.

## Functional Requirements
- FR-01: Route an order based on rules and metrics.
- FR-02: Update destination health and latency signals.
- FR-03: Fallback when preferred destination fails.
- FR-04: Persist routing decisions for audit.
- FR-05: Expose /api/v1/health.

## Non-Functional Requirements
- NFR-01: p95 /api/v1/routes < 50ms in local env.
- NFR-02: 1000 requests/sec in local env.
- NFR-03: Structured JSON logs with traceId and routeId.
- NFR-04: OpenTelemetry traces for routing decisions.
- NFR-05: 100% Docker local environment.
- NFR-06: API versioned under /api/v1.

## Success Criteria
- SC-01: p95 routing latency under 50ms in local tests.
- SC-02: Fallback triggers when a destination fails.
- SC-03: Route status can be queried with a routeId.

## API Contracts
- OpenAPI: .specify/specs/001-smart-routing/contracts/openapi.yaml

## Roadmap
- Milestone 1: Basic routing algorithm.
- Milestone 2: Circuit breaker and fallbacks.
- Milestone 3: Observability and performance tuning.
- Milestone 4: Advanced heuristics.

## Trade-offs
- gRPC vs REST.
- Cache global vs per region.
- Strong vs eventual consistency.
- Simple vs sophisticated algorithm.
- Horizontal scaling vs local optimization.
