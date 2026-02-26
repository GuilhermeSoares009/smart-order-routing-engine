package httpapi

import (
    "net/http"
    "time"

    "github.com/GuilhermeSoares009/smart-order-routing-engine/internal/audit"
    "github.com/GuilhermeSoares009/smart-order-routing-engine/internal/ratelimit"
    "github.com/GuilhermeSoares009/smart-order-routing-engine/internal/routing"
)

type Server struct {
    limiter     *ratelimit.Limiter
    auditStore  *audit.Store
    metricCache *routing.MetricCache
    mux         *http.ServeMux
}

func NewServer(limiter *ratelimit.Limiter) *Server {
    server := &Server{
        limiter:     limiter,
        auditStore:  audit.NewStore(),
        metricCache: routing.NewMetricCache(30 * time.Second),
        mux:         http.NewServeMux(),
    }
    server.routes()
    return server
}

func (s *Server) routes() {
    s.mux.HandleFunc("/api/v1/health", s.handleHealth)
    s.mux.HandleFunc("/api/v1/routes", s.handleRoutes)
    s.mux.HandleFunc("/api/v1/audit/routes", s.handleAudit)
}

func (s *Server) Handler() http.Handler {
    return s.withRateLimit(s.mux)
}

func (s *Server) withRateLimit(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if !s.limiter.Allow(clientIP(r), time.Now()) {
            writeJSON(w, http.StatusTooManyRequests, errorResponse{Error: "rate limit exceeded"})
            return
        }
        next.ServeHTTP(w, r)
    })
}
