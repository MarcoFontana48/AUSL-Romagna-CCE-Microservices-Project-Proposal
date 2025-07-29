package controller

import (
	"encoding/json"
	"github.com/MarcoFontana48/AUSL-Romagna-CCE-Microservices-Project-Proposal/common/circuitbreaker"
	"github.com/MarcoFontana48/AUSL-Romagna-CCE-Microservices-Project-Proposal/utils/http/response"
	"github.com/sony/gobreaker/v2"
	"log/slog"
	"net/http"
	"time"
)

var circuitBreakerSettings = gobreaker.Settings{
	Name:     "healthcheck",
	Timeout:  time.Second * 30,
	Interval: time.Second * 60,
	ReadyToTrip: func(counts gobreaker.Counts) bool {
		failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
		return counts.Requests >= 5 && failureRatio >= 0.8
	},
}

var circuitBreaker = circuitbreaker.NewCircuitBreaker(circuitBreakerSettings)

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("requested health check", "from", r.RemoteAddr)

	msg, err := circuitBreaker.Execute(func() ([]byte, error) {
		msg := response.HealthCheck{Status: "OK", Service: "service"}
		return json.Marshal(msg)
	})

	if err != nil {
		response.Error(w, err)
		slog.Error("health check failed, sent response", "content", err, "to", r.RemoteAddr)
		return
	}

	response.Ok(w, msg)
	slog.Info("successful health check, sent response", "content", msg, "to", r.RemoteAddr)
}
