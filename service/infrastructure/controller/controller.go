package controller

import (
	"github.com/MarcoFontana48/AUSL-Romagna-CCE-Microservices-Project-Proposal/utils/http/response"
	"log/slog"
	"net/http"
)

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("requested health check", "from", r.RemoteAddr)

	msg := response.HealthCheck{Status: "OK", Service: "service"}
	response.SendResponse(w, r, msg)
}
