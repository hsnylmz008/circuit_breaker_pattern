package circuitbreaker

import (
	"github.com/sony/gobreaker"
	"time"
)

type CircuitBreaker struct {
	breaker *gobreaker.CircuitBreaker
}

type Config struct {
	Name            string
	MaxRequests     uint32
	Interval        time.Duration
	Timeout         time.Duration
	FailureRatio    float64
	MinimumRequests int64
}

func NewCircuitBreaker(cfg Config) *CircuitBreaker {
	settings := gobreaker.Settings{
		Name:        cfg.Name,
		MaxRequests: cfg.MaxRequests,
		Interval:    cfg.Interval,
		Timeout:     cfg.Timeout,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= cfg.MinimumRequests && failureRatio >= cfg.FailureRatio
		},
	}

	return &CircuitBreaker{
		breaker: gobreaker.NewCircuitBreaker(settings),
	}
}

func (cb *CircuitBreaker) Execute(req func() (interface{}, error)) (interface{}, error) {
	return cb.breaker.Execute(req)
}

func (cb *CircuitBreaker) State() string {
	return cb.breaker.State().String()
}
