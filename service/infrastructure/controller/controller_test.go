package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func mockSendResponse(w http.ResponseWriter, _ *http.Request, _ interface{}) error {
	w.WriteHeader(http.StatusOK)
	return nil
}

func TestHealthCheckHandler(t *testing.T) {
	originalSendResponse := sendResponse
	sendResponse = mockSendResponse
	defer func() { sendResponse = originalSendResponse }()

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	HealthCheckHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

var sendResponse = func(w http.ResponseWriter, r *http.Request, msg interface{}) error {
	return nil
}
