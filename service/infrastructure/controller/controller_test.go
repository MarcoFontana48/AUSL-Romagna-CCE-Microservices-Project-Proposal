package controller

import (
	"github.com/MarcoFontana48/AUSL-Romagna-CCE-Microservices-Project-Proposal/common/metrics"
	"github.com/MarcoFontana48/AUSL-Romagna-CCE-Microservices-Project-Proposal/utils/http/endpoint"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheckHandler_Success(t *testing.T) {
	metricsInstance := metrics.New()
	ctrl := NewController(metricsInstance)

	req, err := http.NewRequest("GET", endpoint.Health, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	ctrl.HealthCheckHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"status":"OK","service":"service"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
