package response

import (
	"errors"
	"github.com/MarcoFontana48/AUSL-Romagna-CCE-Microservices-Project-Proposal/utils"
	"github.com/sony/gobreaker/v2"
	"io"
	"log/slog"
	"net/http"
)

func Ok(w http.ResponseWriter, jsonByteMsg []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(jsonByteMsg)
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

func Error(w http.ResponseWriter, err error) {
	if errors.Is(err, gobreaker.ErrOpenState) {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}
	http.Error(w, "Health check failed", http.StatusInternalServerError)
	return
}

// Deprecated: Use Ok to handle http ok responses and Error for error responses.
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

// Deprecated: Use Ok to handle http ok responses
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

// Deprecated: Use Error to handle http error responses
func SendErrorResponse(w http.ResponseWriter, r *http.Request, err error) error {
	err2 := checkSendResponseArguments(w, r)
	if err2 != nil {
		return err2
	}

	w.Header().Set("Content-Type", "application/json")

	var status int
	switch err.(type) {
	default:
		status = http.StatusInternalServerError
	}
	w.WriteHeader(status)

	errStruct := ErrorMsg{Error: err.Error()}
	jsonString, err3 := utils.ToJsonString(errStruct)
	if err3 != nil {
		slog.Error("Error marshaling error response", "error", err3)
		return err3
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
