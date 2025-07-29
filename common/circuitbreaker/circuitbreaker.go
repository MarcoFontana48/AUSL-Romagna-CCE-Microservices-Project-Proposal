package circuitbreaker

import "github.com/sony/gobreaker/v2"

func NewCircuitBreaker(settings gobreaker.Settings) *gobreaker.CircuitBreaker[[]byte] {
	return gobreaker.NewCircuitBreaker[[]byte](settings)
}

func DefaultSettings() gobreaker.Settings {
	return gobreaker.Settings{
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 3 && failureRatio >= 0.6
		},
	}
}
