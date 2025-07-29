package circuitbreaker

import "github.com/sony/gobreaker/v2"

var _ *gobreaker.CircuitBreaker[[]byte]

func init() {
	var st gobreaker.Settings
	st.ReadyToTrip = func(counts gobreaker.Counts) bool {
		failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
		return counts.Requests >= 3 && failureRatio >= 0.6
	}

	_ = gobreaker.NewCircuitBreaker[[]byte](st)
}
