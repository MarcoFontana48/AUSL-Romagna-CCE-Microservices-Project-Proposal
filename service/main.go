package main

import (
	"github.com/MarcoFontana48/AUSL-Romagna-CCE-Microservices-Project-Proposal/utils/log"
	"log/slog"
)

func main() {
	log.InitAsJson()

	slog.Info("service module started", "module", "service")
	slog.Debug("service module started", "module", "service")
	slog.Warn("service module warning", "module", "service")
	slog.Error("service module error", "module", "service")
}
