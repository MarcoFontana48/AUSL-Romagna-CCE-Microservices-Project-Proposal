package main

import (
	"api_gateway/infrastructure/controller"
	"api_gateway/infrastructure/server"
	"github.com/MarcoFontana48/AUSL-Romagna-CCE-Microservices-Project-Proposal/common/metrics"
	"github.com/MarcoFontana48/AUSL-Romagna-CCE-Microservices-Project-Proposal/utils/log"
	"log/slog"
)

func main() {
	log.InitAsJson()
	slog.Debug("api_gateway module started", "module", "api_gateway")

	metricsInstance := metrics.New()

	ctrl := controller.NewController(metricsInstance)

	server.StartServer(ctrl)
}
