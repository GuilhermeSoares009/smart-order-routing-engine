package httpapi

import (
    "errors"
    "strconv"
    "strings"

    "github.com/GuilhermeSoares009/smart-order-routing-engine/internal/audit"
    "github.com/GuilhermeSoares009/smart-order-routing-engine/internal/routing"
)

type routeRequest struct {
    RouteID string        `json:"routeId"`
    Order   orderRequest  `json:"order"`
    Targets []targetInput `json:"targets"`
}

type orderRequest struct {
    ID       string `json:"id"`
    Symbol   string `json:"symbol"`
    Quantity int64  `json:"quantity"`
    Side     string `json:"side"`
}

type targetInput struct {
    ID           string  `json:"id"`
    Name         string  `json:"name"`
    LatencyMs    int64   `json:"latencyMs"`
    Availability float64 `json:"availability"`
    Priority     int     `json:"priority"`
}

type routeResponse struct {
    RouteID string           `json:"routeId"`
    TraceID string           `json:"traceId"`
    Decision decisionPayload `json:"decision"`
}

type decisionPayload struct {
    TargetID string  `json:"targetId"`
    Reason   string  `json:"reason"`
    Fallback bool    `json:"fallback"`
    Score    float64 `json:"score"`
}

type errorResponse struct {
    Error string `json:"error"`
}

type auditResponse struct {
    Entries []audit.Entry `json:"entries"`
}

func (req routeRequest) Validate() error {
    if strings.TrimSpace(req.Order.ID) == "" {
        return errors.New("order.id is required")
    }
    if strings.TrimSpace(req.Order.Symbol) == "" {
        return errors.New("order.symbol is required")
    }
    if req.Order.Quantity <= 0 {
        return errors.New("order.quantity must be greater than 0")
    }
    side := strings.ToLower(strings.TrimSpace(req.Order.Side))
    if side != "buy" && side != "sell" {
        return errors.New("order.side must be 'buy' or 'sell'")
    }
    if len(req.Targets) == 0 {
        return errors.New("targets must include at least one target")
    }
    for idx, target := range req.Targets {
        if strings.TrimSpace(target.ID) == "" {
            return errors.New("targets[" + strconv.Itoa(idx) + "].id is required")
        }
        if target.LatencyMs < 0 {
            return errors.New("targets[" + strconv.Itoa(idx) + "].latencyMs must be >= 0")
        }
        if target.Availability < 0 || target.Availability > 1 {
            return errors.New("targets[" + strconv.Itoa(idx) + "].availability must be between 0 and 1")
        }
    }
    return nil
}

func (req routeRequest) TargetsToRouting() []routing.Target {
    targets := make([]routing.Target, 0, len(req.Targets))
    for _, target := range req.Targets {
        targets = append(targets, routing.Target{
            ID:           target.ID,
            Name:         target.Name,
            LatencyMs:    target.LatencyMs,
            Availability: target.Availability,
            Priority:     target.Priority,
        })
    }
    return targets
}

func parseLimit(raw string, fallback int) int {
    if raw == "" {
        return fallback
    }
    value, err := strconv.Atoi(raw)
    if err != nil || value <= 0 {
        return fallback
    }
    return value
}
