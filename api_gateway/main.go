package main

import (
	"github.com/MarcoFontana48/AUSL-Romagna-CCE-Microservices-Project-Proposal/utils/log"
	"log/slog"
)

func main() {
	log.InitAsJson()

	slog.Info("api_gateway module started", "module", "api_gateway")
	slog.Debug("api_gateway module started", "module", "api_gateway")
	slog.Warn("api_gateway module warning", "module", "api_gateway")
	slog.Error("api_gateway module error", "module", "api_gateway")
}
