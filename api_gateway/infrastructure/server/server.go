package server

import (
	"api_gateway/infrastructure/controller"
	"github.com/MarcoFontana48/AUSL-Romagna-CCE-Microservices-Project-Proposal/utils/http/dns"
	"github.com/MarcoFontana48/AUSL-Romagna-CCE-Microservices-Project-Proposal/utils/http/endpoint"
	"github.com/MarcoFontana48/AUSL-Romagna-CCE-Microservices-Project-Proposal/utils/http/port"
	"github.com/MarcoFontana48/AUSL-Romagna-CCE-Microservices-Project-Proposal/utils/http/prefix"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
)

func StartServer(controller *controller.Controller) {
	r := mux.NewRouter()

	r.Use(controller.GetMetricsMiddleware())

	/* API GATEWAY ENDPOINTS */
	// health
	r.HandleFunc(endpoint.Health, controller.HealthCheckHandler).Methods("GET")
	// route
	r.HandleFunc(endpoint.Route, controller.RoutesHandler).Methods("GET")
	// metrics endpoint
	r.HandleFunc(endpoint.Metrics, controller.GetMetricsHandler).Methods("GET")

	/* REROUTES */
	// service
	serviceURL, _ := url.Parse(prefix.HttpPrefix + dns.Service + ":" + strconv.Itoa(port.Http))
	serviceProxy := httputil.NewSingleHostReverseProxy(serviceURL)
	r.PathPrefix(endpoint.Service).HandlerFunc(controller.RerouteHandler(endpoint.Service, serviceProxy))

	startServing(r)
}

func startServing(r *mux.Router) {
	portString := ":" + strconv.Itoa(port.Http)
	slog.Info("API Gateway listening on " + portString)
	err := http.ListenAndServe(portString, r)
	if err != nil {
		return
	}
}
