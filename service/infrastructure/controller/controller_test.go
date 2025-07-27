package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	mockSendResponseCalled bool
	mockMessage            interface{}
)

// mock SendResponse function
func mockSendResponse(w http.ResponseWriter, r *http.Request, msg interface{}) error {
	mockSendResponseCalled = true
	mockMessage = msg
	w.WriteHeader(http.StatusOK)
	return nil
}

func TestHealthCheckHandler(t *testing.T) {
	// reset mock state
	mockSendResponseCalled = false
	mockMessage = nil

	// mock the response.SendResponse function
	originalSendResponse := sendResponse
	sendResponse = mockSendResponse
	defer func() { sendResponse = originalSendResponse }()

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	HealthCheckHandler(w, req)

	// verify response status
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

// variable to hold the actual SendResponse function for mocking
var sendResponse = func(w http.ResponseWriter, r *http.Request, msg interface{}) error {
	return nil
}
