package controller

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/MarcoFontana48/AUSL-Romagna-CCE-Microservices-Project-Proposal/common/circuitbreaker"
	"github.com/MarcoFontana48/AUSL-Romagna-CCE-Microservices-Project-Proposal/common/metrics"
	"github.com/MarcoFontana48/AUSL-Romagna-CCE-Microservices-Project-Proposal/utils/http/response"
	"github.com/sony/gobreaker/v2"
)

/* === StandardController definition and settings === */

// StandardController holds dependencies for handling requests
type StandardController struct {
	metrics        *metrics.Metrics
	circuitBreaker *circuitbreaker.CircuitBreaker
}

// NewController creates a new controller with injected dependencies
func NewController(m *metrics.Metrics) *StandardController {
	c := &StandardController{
		metrics: m,
	}

	circuitBreakerSettings := getCircuitBreakerDefaultSettings(c)
	c.circuitBreaker = circuitbreaker.NewCircuitBreaker(circuitBreakerSettings)
	return c
}

/* === Handlers === */

// HealthCheckHandler handles health check requests
func (c *StandardController) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("requested health check", "from", r.RemoteAddr)

	// start measuring time for metrics
	startTime := time.Now()

	msg, err := c.generateHealthCheckMessageResponse()
	if err != nil {
		// record health check failure metrics
		c.recordHealthCheckData(startTime, "failure")

		// send error response
		response.Error(w, err)

		// log the error
		slog.Error("health check failed, sent response", "content", err, "to", r.RemoteAddr)
		return
	}
	// record health check success metrics
	c.recordHealthCheckData(startTime, "failure")

	// send success response
	response.Ok(w, msg)

	// log the successful response
	slog.Info("successful health check, sent response", "content", msg, "to", r.RemoteAddr)
}

// MetricsHandler returns the Prometheus metrics HTTP handler function
func (c *StandardController) MetricsHandler(w http.ResponseWriter, r *http.Request) {
	c.metrics.Handler().ServeHTTP(w, r)
}

/* === Helper Methods === */

func (c *StandardController) generateHealthCheckMessageResponse() ([]byte, error) {
	msg, err := c.circuitBreaker.Execute(func() ([]byte, error) {
		msg := response.HealthCheck{Status: "OK", Service: "service"}
		return json.Marshal(msg)
	})
	return msg, err
}

func (c *StandardController) recordHealthCheckData(startTime time.Time, status string) {
	elapsedTimeSinceStart := time.Since(startTime)
	c.metrics.RecordCircuitBreakerRequest("healthcheck", status)
	c.metrics.RecordHealthCheck(status, elapsedTimeSinceStart)
}

/* === Getters === */

// GetMetricsMiddleware returns the metrics middleware
func (c *StandardController) GetMetricsMiddleware() func(http.Handler) http.Handler {
	return c.metrics.Middleware()
}

// GetCircuitBreakerMetrics returns current circuit breaker statistics
func (c *StandardController) GetCircuitBreakerMetrics() map[string]interface{} {
	counts := c.circuitBreaker.Counts()
	return map[string]interface{}{
		"requests":              counts.Requests,
		"total_successes":       counts.TotalSuccesses,
		"total_failures":        counts.TotalFailures,
		"consecutive_successes": counts.ConsecutiveSuccesses,
		"consecutive_failures":  counts.ConsecutiveFailures,
		"state":                 c.circuitBreaker.State().String(),
	}
}

func getCircuitBreakerDefaultSettings(c *StandardController) gobreaker.Settings {
	circuitBreakerSettings := gobreaker.Settings{
		Name:     "service",
		Timeout:  time.Second * 30,
		Interval: time.Second * 60,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 5 && failureRatio >= 0.8
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			c.metrics.RecordCircuitBreakerStateChange(name, to)
			slog.Info("circuit breaker state changed", "name", name, "from", from, "to", to)
		},
	}
	return circuitBreakerSettings
}
