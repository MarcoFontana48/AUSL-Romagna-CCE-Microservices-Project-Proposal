package controller

import (
	"encoding/json"
	"github.com/MarcoFontana48/AUSL-Romagna-CCE-Microservices-Project-Proposal/utils/http/common"
	"github.com/MarcoFontana48/AUSL-Romagna-CCE-Microservices-Project-Proposal/utils/http/response"
	"log/slog"
	"net/http"
)

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("requested health check", "from", r.RemoteAddr)

	msg, err := .Execute(func() ([]byte, error) {
		msg := response.HealthCheck{Status: "OK", Service: "service"}
		return json.Marshal(msg)
	})

	if err != nil {
		response.Error(w, err)
		return
	}

	response.Ok(w, msg)
	defer slog.Info("sent health check response", "content", msg, "to", r.RemoteAddr)
}

