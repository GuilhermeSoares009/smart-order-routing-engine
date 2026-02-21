package routing

import "errors"

const minAvailability = 0.5

var ErrNoTargets = errors.New("no targets provided")

type Target struct {
    ID           string
    Name         string
    LatencyMs    int64
    Availability float64
    Priority     int
}

type Decision struct {
    Target   Target
    Score    float64
    Fallback bool
    Reason   string
}

func SelectTarget(targets []Target) (Decision, error) {
    if len(targets) == 0 {
        return Decision{}, ErrNoTargets
    }

    eligible := make([]Target, 0, len(targets))
    for _, target := range targets {
        if target.Availability >= minAvailability {
            eligible = append(eligible, target)
        }
    }

    if len(eligible) == 0 {
        fallback := pickBest(targets)
        return Decision{
            Target:   fallback,
            Score:    float64(fallback.LatencyMs),
            Fallback: true,
            Reason:   "fallback-no-healthy-targets",
        }, nil
    }

    best := pickBest(eligible)
    return Decision{
        Target:   best,
        Score:    float64(best.LatencyMs),
        Fallback: false,
        Reason:   "best-latency",
    }, nil
}

func pickBest(targets []Target) Target {
    best := targets[0]
    for _, candidate := range targets[1:] {
        if candidate.LatencyMs < best.LatencyMs {
            best = candidate
            continue
        }
        if candidate.LatencyMs == best.LatencyMs {
            if candidate.Availability > best.Availability {
                best = candidate
                continue
            }
            if candidate.Availability == best.Availability && candidate.Priority < best.Priority {
                best = candidate
            }
        }
    }
    return best
}
