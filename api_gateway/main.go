package main

import (
	"api_gateway/infrastructure/server"
	"github.com/MarcoFontana48/AUSL-Romagna-CCE-Microservices-Project-Proposal/utils/log"
	"log/slog"
)

func main() {
	log.InitAsJson()
	slog.Debug("api_gateway module started", "module", "api_gateway")

	server.StartServer()
}
