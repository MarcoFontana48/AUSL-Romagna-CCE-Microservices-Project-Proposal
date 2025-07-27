package controller

import (
	"github.com/MarcoFontana48/AUSL-Romagna-CCE-Microservices-Project-Proposal/utils/http/endpoint"
	"github.com/MarcoFontana48/AUSL-Romagna-CCE-Microservices-Project-Proposal/utils/http/response"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"strings"
)

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("requested health check", "from", r.RemoteAddr)

	msg := response.HealthCheck{Status: "OK", Service: "api-gateway"}
	response.SendResponse(w, r, msg)
}

func RoutesHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("requested routes", "from", r.RemoteAddr)

	msg := endpoint.All
	response.SendResponse(w, r, msg)
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
