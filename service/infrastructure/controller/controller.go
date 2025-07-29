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

/* === Controller definition and settings === */

// Controller holds dependencies for handling requests
type Controller struct {
	metrics        *metrics.Metrics
	circuitBreaker *circuitbreaker.CircuitBreaker
}

// NewController creates a new controller with injected dependencies
func NewController(m *metrics.Metrics) *Controller {
	c := &Controller{
		metrics: m,
	}

	// initialize circuit breaker with metrics integration
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

	c.circuitBreaker = circuitbreaker.NewCircuitBreaker(circuitBreakerSettings)
	return c
}

/* === Handlers === */

// HealthCheckHandler handles health check requests
func (c *Controller) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("requested health check", "from", r.RemoteAddr)

	// start measuring time for metrics
	startTime := time.Now()

	msg, err := c.generateHealthCheckMessageResponse()
	if err != nil {
		c.sendHealthCheckErrorResponse(w, r, err, startTime)
		return
	}
	c.sendHealthCheckOkResponse(w, r, msg, startTime)
}

/* === Helper Methods === */

func (c *Controller) generateHealthCheckMessageResponse() ([]byte, error) {
	msg, err := c.circuitBreaker.Execute(func() ([]byte, error) {
		msg := response.HealthCheck{Status: "OK", Service: "service"}
		return json.Marshal(msg)
	})
	return msg, err
}

func (c *Controller) sendHealthCheckErrorResponse(w http.ResponseWriter, r *http.Request, err error, startTime time.Time) {
	// send error response
	response.Error(w, err)

	// record metrics for failure
	elapsedTimeSinceStart := time.Since(startTime)
	c.recordHealthCheckMetrics("failure", elapsedTimeSinceStart)

	// log the error
	slog.Error("health check failed, sent response", "content", err, "to", r.RemoteAddr)
	return
}

func (c *Controller) sendHealthCheckOkResponse(w http.ResponseWriter, r *http.Request, msg []byte, startTime time.Time) {
	// send success response
	response.Ok(w, msg)

	// record metrics for success
	elapsedTimeSinceStart := time.Since(startTime)
	c.recordHealthCheckMetrics("success", elapsedTimeSinceStart)

	// log the successful response
	slog.Info("successful health check, sent response", "content", msg, "to", r.RemoteAddr)
}

func (c *Controller) recordHealthCheckMetrics(result string, duration time.Duration) {
	c.metrics.RecordCircuitBreakerRequest("healthcheck", result)
	c.metrics.RecordHealthCheck(result, duration)
}

// GetMetricsMiddleware returns the metrics middleware
func (c *Controller) GetMetricsMiddleware() func(http.Handler) http.Handler {
	return c.metrics.Middleware()
}

// GetMetricsHandler returns the Prometheus metrics HTTP handler function
func (c *Controller) GetMetricsHandler(w http.ResponseWriter, r *http.Request) {
	c.metrics.Handler().ServeHTTP(w, r)
}

// GetCircuitBreakerMetrics returns current circuit breaker statistics
func (c *Controller) GetCircuitBreakerMetrics() map[string]interface{} {
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
