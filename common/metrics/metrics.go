package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sony/gobreaker/v2"
)

// Metrics holds all Prometheus metrics
type Metrics struct {
	httpRequestsTotal      *prometheus.CounterVec
	httpRequestDuration    *prometheus.HistogramVec
	httpRequestsInFlight   prometheus.Gauge
	circuitBreakerState    *prometheus.GaugeVec
	circuitBreakerRequests *prometheus.CounterVec
	healthCheckRequests    *prometheus.CounterVec
	healthCheckDuration    prometheus.Histogram
	routesRequests         *prometheus.CounterVec
	routesDuration         prometheus.Histogram
}

// New creates a new Metrics instance with all Prometheus metrics initialized
func New() *Metrics {
	return &Metrics{
		httpRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests by method, endpoint, and status code",
			},
			[]string{"method", "endpoint", "status_code"},
		),
		httpRequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "Duration of HTTP requests in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "endpoint", "status_code"},
		),
		httpRequestsInFlight: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "http_requests_in_flight",
				Help: "Current number of HTTP requests being processed",
			},
		),
		circuitBreakerState: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "circuit_breaker_state",
				Help: "Circuit breaker state (0=closed, 1=half-open, 2=open)",
			},
			[]string{"name"},
		),
		circuitBreakerRequests: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "circuit_breaker_requests_total",
				Help: "Total number of requests through circuit breaker",
			},
			[]string{"name", "result"},
		),
		healthCheckRequests: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "health_check_requests_total",
				Help: "Total number of health check requests",
			},
			[]string{"status"},
		),
		healthCheckDuration: promauto.NewHistogram(
			prometheus.HistogramOpts{
				Name:    "health_check_duration_seconds",
				Help:    "Duration of health check requests in seconds",
				Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
			},
		),
	}
}

// HTTPMiddleware wraps HTTP handlers to collect metrics
func (m *Metrics) HTTPMiddleware(endpoint string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// If endpoint is empty, try to get it from the request
			endpointLabel := endpoint
			if endpointLabel == "" {
				endpointLabel = r.URL.Path
			}

			m.httpRequestsInFlight.Inc()
			defer m.httpRequestsInFlight.Dec()

			ww := &responseWriterWrapper{ResponseWriter: w, statusCode: http.StatusOK}
			next.ServeHTTP(ww, r)

			duration := time.Since(start).Seconds()
			statusCode := strconv.Itoa(ww.statusCode)

			m.httpRequestsTotal.WithLabelValues(r.Method, endpointLabel, statusCode).Inc()
			m.httpRequestDuration.WithLabelValues(r.Method, endpointLabel, statusCode).Observe(duration)
		})
	}
}

// Middleware creates a middleware that automatically extracts route patterns
func (m *Metrics) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			m.httpRequestsInFlight.Inc()
			defer m.httpRequestsInFlight.Dec()

			ww := &responseWriterWrapper{ResponseWriter: w, statusCode: http.StatusOK}
			next.ServeHTTP(ww, r)

			// Extract route pattern from mux.Vars (if available)
			endpointLabel := r.URL.Path
			if route := mux.CurrentRoute(r); route != nil {
				if template, err := route.GetPathTemplate(); err == nil {
					endpointLabel = template
				}
			}

			duration := time.Since(start).Seconds()
			statusCode := strconv.Itoa(ww.statusCode)

			m.httpRequestsTotal.WithLabelValues(r.Method, endpointLabel, statusCode).Inc()
			m.httpRequestDuration.WithLabelValues(r.Method, endpointLabel, statusCode).Observe(duration)
		})
	}
}

// RecordCircuitBreakerStateChange updates circuit breaker state metric
func (m *Metrics) RecordCircuitBreakerStateChange(name string, state gobreaker.State) {
	var stateValue float64
	switch state {
	case gobreaker.StateClosed:
		stateValue = 0
	case gobreaker.StateHalfOpen:
		stateValue = 1
	case gobreaker.StateOpen:
		stateValue = 2
	}
	m.circuitBreakerState.WithLabelValues(name).Set(stateValue)
}

// RecordCircuitBreakerRequest records circuit breaker request result
func (m *Metrics) RecordCircuitBreakerRequest(name, result string) {
	m.circuitBreakerRequests.WithLabelValues(name, result).Inc()
}

// RecordHealthCheck records health check metrics
func (m *Metrics) RecordHealthCheck(status string, duration time.Duration) {
	m.healthCheckRequests.WithLabelValues(status).Inc()
	m.healthCheckDuration.Observe(duration.Seconds())
}

func (m *Metrics) RecordRoutesRequest(status string, duration time.Duration) {
	m.routesRequests.WithLabelValues(status).Inc()
	m.routesDuration.Observe(duration.Seconds())
}

// Handler returns the Prometheus metrics HTTP handler
func (m *Metrics) Handler() http.Handler {
	return promhttp.Handler()
}

// responseWriterWrapper wraps http.ResponseWriter to capture status code
type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriterWrapper) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
