package controller

import (
	"encoding/json"
	"github.com/MarcoFontana48/AUSL-Romagna-CCE-Microservices-Project-Proposal/common/circuitbreaker"
	"github.com/MarcoFontana48/AUSL-Romagna-CCE-Microservices-Project-Proposal/utils/http/endpoint"
	"github.com/MarcoFontana48/AUSL-Romagna-CCE-Microservices-Project-Proposal/utils/http/response"
	"github.com/sony/gobreaker/v2"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

var circuitBreakerSettings = gobreaker.Settings{
	Name:     "api-gateway",
	Timeout:  time.Second * 30,
	Interval: time.Second * 60,
	ReadyToTrip: func(counts gobreaker.Counts) bool {
		failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
		return counts.Requests >= 10 && failureRatio >= 0.8
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

func RoutesHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("requested routes", "from", r.RemoteAddr)

	msg, err := circuitBreaker.Execute(func() ([]byte, error) {
		msg := endpoint.All
		return json.Marshal(msg)
	})

	if err != nil {
		response.Error(w, err)
		slog.Error("requested routes failed, sent response", "content", err, "to", r.RemoteAddr)
		return
	}

	response.Ok(w, msg)
	slog.Info("successful requested routes, sent response", "content", msg, "to", r.RemoteAddr)
}

func RerouteHandler(service string, serviceProxy *httputil.ReverseProxy) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// remove prefix before forwarding
		removePrefix(r, service)
		slog.Debug("Forwarding request to "+service, "endpoint", r.URL.Path)
		serviceProxy.ServeHTTP(w, r)
	}
}

func removePrefix(r *http.Request, service string) {
	r.URL.Path = strings.TrimPrefix(r.URL.Path, service)
	if r.URL.Path == "" {
		r.URL.Path = endpoint.Root
	}
}
