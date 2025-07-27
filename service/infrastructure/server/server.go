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

func StartServer() {
	r := mux.NewRouter()

	// health check endpoint
	r.HandleFunc(endpoint.Health, controller.HealthCheckHandler).Methods("GET")

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
