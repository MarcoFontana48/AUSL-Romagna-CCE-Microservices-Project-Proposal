package controller

import (
	"net/http"
)

type ServiceController interface {
	HealthCheckHandler(w http.ResponseWriter, r *http.Request)
	MetricsHandler(w http.ResponseWriter, r *http.Request)
}
