package routing

import (
	"sync"
	"time"
)

type MetricCache struct {
	mu      sync.RWMutex
	ttl     time.Duration
	metrics map[string]metricEntry
}

type metricEntry struct {
	latencyMs    int64
	availability float64
	updatedAt    time.Time
}

func NewMetricCache(ttl time.Duration) *MetricCache {
	return &MetricCache{
		ttl:     ttl,
		metrics: make(map[string]metricEntry),
	}
}

func (c *MetricCache) Merge(targets []Target, now time.Time) []Target {
	c.mu.Lock()
	defer c.mu.Unlock()

	merged := make([]Target, 0, len(targets))
	for _, target := range targets {
		if cached, ok := c.metrics[target.ID]; ok && now.Sub(cached.updatedAt) <= c.ttl {
			target.LatencyMs = (target.LatencyMs + cached.latencyMs) / 2
			target.Availability = (target.Availability + cached.availability) / 2
		}

		entry := metricEntry{
			latencyMs:    target.LatencyMs,
			availability: target.Availability,
			updatedAt:    now,
		}
		c.metrics[target.ID] = entry
		merged = append(merged, target)
	}
	return merged
}
