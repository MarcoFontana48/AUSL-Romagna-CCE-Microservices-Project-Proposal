package response

import (
	"errors"
	"github.com/MarcoFontana48/AUSL-Romagna-CCE-Microservices-Project-Proposal/utils"
	"io"
	"log/slog"
	"net/http"
)

func SendResponse(w http.ResponseWriter, r *http.Request, msg interface{}) {
	err := SendOkResponse(w, r, msg)
	if err != nil {
		err2 := SendErrorResponse(w, r, err)
		if err2 != nil {
			w.WriteHeader(http.StatusInternalServerError)
			slog.Error("Error sending error response", "error", err2, "to", r.RemoteAddr)
		}
	}
}

func SendOkResponse(w http.ResponseWriter, r *http.Request, msg interface{}) error {
	err2 := checkSendResponseArguments(w, r)
	if err2 != nil {
		return err2
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonString, err := utils.ToJsonString(msg)
	if err != nil {
		return err
	}
	_, err = io.WriteString(w, jsonString)
	if err != nil {
		slog.Error("Error writing response", "error", err)
		return err
	}
	slog.Info("Response sent", "response", jsonString, "to", r.RemoteAddr)
	return nil
}

func SendErrorResponse(w http.ResponseWriter, r *http.Request, err error) error {
	err2 := checkSendResponseArguments(w, r)
	if err2 != nil {
		return err2
	}

	w.Header().Set("Content-Type", "application/json")

	var status int
	switch err.(type) {
	//TODO: add specific error types for better granularity
	default:
		status = http.StatusInternalServerError
	}
	w.WriteHeader(status)

	errStruct := Error{Error: err.Error()}
	jsonString, errr := utils.ToJsonString(errStruct)
	if errr != nil {
		slog.Error("Error marshaling error response", "error", errr)
		return errr
	}
	writeString, err := io.WriteString(w, jsonString)
	if err != nil {
		return err
	}

	slog.Info("Error response sent", "response", writeString, "to", r.RemoteAddr, "status", status)
	return nil
}

func checkSendResponseArguments(w http.ResponseWriter, r *http.Request) error {
	if r == nil {
		slog.Error("Received nil request in SendOkResponse")
		return errors.New("nil request")
	}
	if w == nil {
		slog.Error("Received nil response writer in SendOkResponse")
		return errors.New("nil response writer")
	}
	return nil
}
