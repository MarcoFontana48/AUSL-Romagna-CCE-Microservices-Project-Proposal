package main

import (
	"github.com/MarcoFontana48/AUSL-Romagna-CCE-Microservices-Project-Proposal/common/metrics"
	"github.com/MarcoFontana48/AUSL-Romagna-CCE-Microservices-Project-Proposal/utils/log"
	"log/slog"
	"service/infrastructure/controller"
	"service/infrastructure/server"
)

func main() {
	log.InitAsJson()
	slog.Debug("service module started", "module", "service")

	metricsInstance := metrics.New()

	ctrl := controller.NewController(metricsInstance)

	server.StartServer(ctrl)
}
