package server

import (
	"github.com/MarcoFontana48/AUSL-Romagna-CCE-Microservices-Project-Proposal/utils/http/endpoint"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

var (
	healthHandlerCalled bool
)

func mockHealthCheckHandler(w http.ResponseWriter, _ *http.Request) {
	healthHandlerCalled = true
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(`{"status":"OK","service":"service"}`))
	if err != nil {
		return
	}
}

func createTestRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc(endpoint.Health, mockHealthCheckHandler).Methods("GET")
	return r
}

func TestHealthEndpoint(t *testing.T) {
	healthHandlerCalled = false
	router := createTestRouter()

	req := httptest.NewRequest("GET", endpoint.Health, nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// verify handler was called
	if !healthHandlerCalled {
		t.Error("HealthCheckHandler was not called")
	}

	// Verify response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestInvalidMethod(t *testing.T) {
	healthHandlerCalled = false
	router := createTestRouter()

	req := httptest.NewRequest("POST", endpoint.Health, nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// should not call handler
	if healthHandlerCalled {
		t.Error("Handler should not be called for POST method")
	}

	// should return method not allowed
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
}

func TestNotFound(t *testing.T) {
	healthHandlerCalled = false
	router := createTestRouter()

	req := httptest.NewRequest("GET", "/invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// should not call handler
	if healthHandlerCalled {
		t.Error("Handler should not be called for invalid path")
	}

	// should return not found
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}
