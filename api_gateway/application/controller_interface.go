package application

import (
	"net/http"
	"net/http/httputil"
)

type ApiGatewayController interface {
	HealthCheckHandler(w http.ResponseWriter, r *http.Request)
	RoutesHandler(w http.ResponseWriter, r *http.Request)
	RerouteHandler(service string, serviceProxy *httputil.ReverseProxy) func(w http.ResponseWriter, r *http.Request)
	MetricsHandler(w http.ResponseWriter, r *http.Request)
}
