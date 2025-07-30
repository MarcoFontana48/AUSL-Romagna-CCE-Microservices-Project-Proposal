package server

import (
	"github.com/MarcoFontana48/AUSL-Romagna-CCE-Microservices-Project-Proposal/utils/http/endpoint"
	"github.com/MarcoFontana48/AUSL-Romagna-CCE-Microservices-Project-Proposal/utils/http/port"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
	"service/infrastructure/controller"
	"strconv"
)

func StartServer(controller *controller.StandardController) {
	r := mux.NewRouter()

	// apply metrics middleware to all routes
	r.Use(controller.GetMetricsMiddleware())

	// health check endpoint
	r.HandleFunc(endpoint.Health, controller.HealthCheckHandler).Methods("GET")

	// metrics endpoint
	r.HandleFunc(endpoint.Metrics, controller.MetricsHandler).Methods("GET")

	startServing(r)
}

func startServing(r *mux.Router) {
	portString := ":" + strconv.Itoa(port.Http)
	slog.Info("Service listening on " + portString)
	err := http.ListenAndServe(portString, r)
	if err != nil {
		return
	}
}
