package httpapi

import (
    "context"
    "crypto/rand"
    "encoding/hex"
    "encoding/json"
    "errors"
    "io"
    "net"
    "net/http"
    "os"
    "strings"
    "sync"
    "time"

    "github.com/GuilhermeSoares009/smart-order-routing-engine/internal/audit"
    "github.com/GuilhermeSoares009/smart-order-routing-engine/internal/routing"
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/metric"
    "go.opentelemetry.io/otel/propagation"
    "go.opentelemetry.io/otel/trace"
)

const (
    maxBodySize      = 1 << 20
    latencyBudgetMs  = 50
    serviceTraceName = "httpapi"
)

var (
    durationOnce      sync.Once
    durationHistogram metric.Int64Histogram
)

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
    ctx, span := startSpan(r.Context(), r)
    defer span.End()

    writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
    logRequest(ctx, logEntry{
        Message:     "health check",
        RouteID:     newID(),
        Destination: "",
        Status:      http.StatusOK,
        Path:        r.URL.Path,
        Method:      r.Method,
    })
}

func (s *Server) handleRoutes(w http.ResponseWriter, r *http.Request) {
    start := time.Now()
    ctx, span := startSpan(r.Context(), r)
    defer span.End()

    if r.Method != http.MethodPost {
        writeJSON(w, http.StatusMethodNotAllowed, errorResponse{Error: "method not allowed"})
        logRequest(ctx, logEntry{
            Message:     "method not allowed",
            RouteID:     newID(),
            Destination: "",
            Status:      http.StatusMethodNotAllowed,
            Path:        r.URL.Path,
            Method:      r.Method,
        })
        return
    }

    var payload routeRequest
    if err := readJSON(r, &payload); err != nil {
        writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
        logRequest(ctx, logEntry{
            Message:     "invalid request",
            RouteID:     newID(),
            Destination: "",
            Status:      http.StatusBadRequest,
            Path:        r.URL.Path,
            Method:      r.Method,
        })
        return
    }

    if err := payload.Validate(); err != nil {
        writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
        logRequest(ctx, logEntry{
            Message:     "validation failed",
            RouteID:     newID(),
            Destination: "",
            Status:      http.StatusBadRequest,
            Path:        r.URL.Path,
            Method:      r.Method,
        })
        return
    }

    routeID := payload.RouteID
    if strings.TrimSpace(routeID) == "" {
        routeID = newID()
    }

    targets := payload.TargetsToRouting()
    if s.metricCache != nil {
        targets = s.metricCache.Merge(targets, time.Now().UTC())
    }

    decision, err := routing.SelectTarget(targets)
    if err != nil {
        status := http.StatusInternalServerError
        message := "routing decision failed"
        if errors.Is(err, routing.ErrNoTargets) {
            status = http.StatusBadRequest
            message = "no targets provided"
        }
        writeJSON(w, status, errorResponse{Error: message})
        logRequest(ctx, logEntry{
            Message:     message,
            RouteID:     routeID,
            Destination: "",
            Status:      status,
            Path:        r.URL.Path,
            Method:      r.Method,
        })
        return
    }

    span.SetAttributes(
        attribute.String("route.id", routeID),
        attribute.String("routing.target", decision.Target.ID),
        attribute.Bool("routing.fallback", decision.Fallback),
    )

    response := routeResponse{
        RouteID: routeID,
        TraceID: span.SpanContext().TraceID().String(),
        Decision: decisionPayload{
            TargetID: decision.Target.ID,
            Reason:   decision.Reason,
            Fallback: decision.Fallback,
            Score:    decision.Score,
        },
    }

    writeJSON(w, http.StatusOK, response)

    durationMs := time.Since(start).Milliseconds()
    recordMetrics(ctx, r.URL.Path, durationMs, decision.Fallback)

    logRequest(ctx, logEntry{
        Message:        "routing decision",
        RouteID:        routeID,
        Destination:    decision.Target.ID,
        Status:         http.StatusOK,
        Path:           r.URL.Path,
        Method:         r.Method,
        DurationMs:     durationMs,
        BudgetExceeded: durationMs > latencyBudgetMs,
        Fallback:       decision.Fallback,
    })

    if s.auditStore != nil {
        s.auditStore.Add(audit.Entry{
            Timestamp:   time.Now().UTC(),
            RouteID:     routeID,
            OrderID:     payload.Order.ID,
            TargetID:    decision.Target.ID,
            Reason:      decision.Reason,
            Fallback:    decision.Fallback,
            Score:       decision.Score,
            TargetCount: len(targets),
        })
    }
}

func (s *Server) handleAudit(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        writeJSON(w, http.StatusMethodNotAllowed, errorResponse{Error: "method not allowed"})
        return
    }
    if s.auditStore == nil {
        writeJSON(w, http.StatusOK, auditResponse{Entries: []audit.Entry{}})
        return
    }
    limit := parseLimit(r.URL.Query().Get("limit"), 50)
    writeJSON(w, http.StatusOK, auditResponse{Entries: s.auditStore.List(limit)})
}

func readJSON(r *http.Request, dst any) error {
    decoder := json.NewDecoder(io.LimitReader(r.Body, maxBodySize))
    decoder.DisallowUnknownFields()
    if err := decoder.Decode(dst); err != nil {
        return err
    }
    if err := decoder.Decode(&struct{}{}); err != io.EOF {
        return errors.New("unexpected data after json body")
    }
    return nil
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    _ = json.NewEncoder(w).Encode(payload)
}

func clientIP(r *http.Request) string {
    if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
        parts := strings.Split(forwarded, ",")
        return strings.TrimSpace(parts[0])
    }
    if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
        return strings.TrimSpace(realIP)
    }
    host, _, err := net.SplitHostPort(r.RemoteAddr)
    if err == nil {
        return host
    }
    return r.RemoteAddr
}

func newID() string {
    bytes := make([]byte, 16)
    if _, err := rand.Read(bytes); err != nil {
        return hex.EncodeToString([]byte(time.Now().Format("20060102150405.000000")))
    }
    return hex.EncodeToString(bytes)
}

func startSpan(ctx context.Context, r *http.Request) (context.Context, trace.Span) {
    propagator := otel.GetTextMapPropagator()
    ctx = propagator.Extract(ctx, propagation.HeaderCarrier(r.Header))
    tracer := otel.Tracer(serviceTraceName)
    ctx, span := tracer.Start(ctx, r.Method+" "+r.URL.Path)
    span.SetAttributes(
        attribute.String("http.method", r.Method),
        attribute.String("http.route", r.URL.Path),
    )
    return ctx, span
}

type logEntry struct {
    Message        string `json:"message"`
    TraceID        string `json:"traceId"`
    RouteID        string `json:"routeId"`
    Destination    string `json:"destination"`
    Status         int    `json:"status"`
    Path           string `json:"path"`
    Method         string `json:"method"`
    DurationMs     int64  `json:"durationMs,omitempty"`
    BudgetExceeded bool   `json:"budgetExceeded,omitempty"`
    Fallback       bool   `json:"fallback,omitempty"`
}

func logRequest(ctx context.Context, entry logEntry) {
    span := trace.SpanFromContext(ctx)
    entry.TraceID = span.SpanContext().TraceID().String()
    data, err := json.Marshal(entry)
    if err != nil {
        return
    }
    _, _ = os.Stdout.Write(append(data, '\n'))
}

func recordMetrics(ctx context.Context, path string, durationMs int64, fallback bool) {
    duration := getDurationHistogram()
    duration.Record(ctx, durationMs,
        metric.WithAttributes(
            attribute.String("http.route", path),
            attribute.Bool("routing.fallback", fallback),
        ),
    )
}

func getDurationHistogram() metric.Int64Histogram {
    durationOnce.Do(func() {
        meter := otel.Meter(serviceTraceName)
        histogram, _ := meter.Int64Histogram("http.server.duration", metric.WithUnit("ms"))
        durationHistogram = histogram
    })
    return durationHistogram
}
