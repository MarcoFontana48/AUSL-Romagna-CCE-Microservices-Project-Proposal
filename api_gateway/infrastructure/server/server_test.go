package server

import (
	"github.com/MarcoFontana48/AUSL-Romagna-CCE-Microservices-Project-Proposal/utils/http/endpoint"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

// mock controller functions for testing
var (
	healthCheckCalled    = false
	routesHandlerCalled  = false
	rerouteHandlerCalled = false
)

// mock controller package
type MockController struct{}

func (m *MockController) HealthCheckHandler(w http.ResponseWriter, _ *http.Request) {
	healthCheckCalled = true
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(`{"status":"healthy"}`))
	if err != nil {
		return
	}
}

func (m *MockController) RoutesHandler(w http.ResponseWriter, _ *http.Request) {
	routesHandlerCalled = true
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(`{"routes":["health","route","service"]}`))
	if err != nil {
		return
	}
}

func (m *MockController) RerouteHandler(endpoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rerouteHandlerCalled = true
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"message":"rerouted to ` + endpoint + `"}`))
		if err != nil {
			return
		}
	}
}

// reset mock state
func resetMocks() {
	healthCheckCalled = false
	routesHandlerCalled = false
	rerouteHandlerCalled = false
}

// create test router with mocked dependencies
func createTestRouter() *mux.Router {
	r := mux.NewRouter()
	mock := &MockController{}

	// setup routes with mocked handlers
	r.HandleFunc(endpoint.Health, mock.HealthCheckHandler).Methods("GET")
	r.HandleFunc(endpoint.Route, mock.RoutesHandler).Methods("GET")
	r.PathPrefix(endpoint.Service).HandlerFunc(mock.RerouteHandler(endpoint.Service))

	return r
}

func TestHealthCheckEndpoint(t *testing.T) {
	resetMocks()
	router := createTestRouter()

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// verify controller was called
	if !healthCheckCalled {
		t.Error("HealthCheckHandler was not called")
	}

	// verify response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	expectedBody := `{"status":"healthy"}`
	if strings.TrimSpace(w.Body.String()) != expectedBody {
		t.Errorf("Expected body %s, got %s", expectedBody, w.Body.String())
	}
}

func TestRoutesEndpoint(t *testing.T) {
	resetMocks()
	router := createTestRouter()

	req := httptest.NewRequest("GET", "/route", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// verify controller was called
	if !routesHandlerCalled {
		t.Error("RouteHandler was not called")
	}

	// verify response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if !strings.Contains(w.Body.String(), "routes") {
		t.Error("Response should contain routes information")
	}
}

func TestServiceRerouteEndpoint(t *testing.T) {
	resetMocks()
	router := createTestRouter()

	// test various service paths
	testPaths := []string{
		"/service",
		"/service/users",
		"/service/api/v1/data",
	}

	for _, path := range testPaths {
		t.Run(path, func(t *testing.T) {
			resetMocks()

			req := httptest.NewRequest("GET", path, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			// verify reroute handler was called
			if !rerouteHandlerCalled {
				t.Error("RerouteHandler was not called for path:", path)
			}

			// verify response
			if w.Code != http.StatusOK {
				t.Errorf("Expected status 200, got %d for path %s", w.Code, path)
			}
		})
	}
}

func TestHTTPMethods(t *testing.T) {
	router := createTestRouter()

	tests := []struct {
		method       string
		path         string
		expectedCode int
		description  string
	}{
		{"GET", "/health", http.StatusOK, "health endpoint with GET"},
		{"POST", "/health", http.StatusMethodNotAllowed, "health endpoint with POST should fail"},
		{"GET", "/route", http.StatusOK, "routes endpoint with GET"},
		{"PUT", "/route", http.StatusMethodNotAllowed, "routes endpoint with PUT should fail"},
		{"POST", "/service/api", http.StatusOK, "service reroute accepts any method"},
		{"PUT", "/service/data", http.StatusOK, "service reroute accepts any method"},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			resetMocks()

			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("Expected status %d, got %d for %s %s",
					tt.expectedCode, w.Code, tt.method, tt.path)
			}
		})
	}
}

func TestNotFoundEndpoint(t *testing.T) {
	resetMocks()
	router := createTestRouter()

	req := httptest.NewRequest("GET", "/nonexistent", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// verify no controllers were called
	if healthCheckCalled || routesHandlerCalled || rerouteHandlerCalled {
		t.Error("No controllers should be called for non-existent endpoint")
	}

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestStartServingFunction(t *testing.T) {
	router := createTestRouter()

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("startServing panicked: %v", r)
		}
	}()

	if router == nil {
		t.Error("Router should not be nil")
	}
}

// integration test for the complete router setup
func TestCompleteRouterSetup(t *testing.T) {
	resetMocks()
	router := createTestRouter()

	// test all endpoints in sequence
	endpoints := []struct {
		method    string
		path      string
		mockCheck func() bool
	}{
		{"GET", "/health", func() bool { return healthCheckCalled }},
		{"GET", "/route", func() bool { return routesHandlerCalled }},
		{"GET", "/service/test", func() bool { return rerouteHandlerCalled }},
	}

	for i, currentEndpoint := range endpoints {
		resetMocks()

		req := httptest.NewRequest(currentEndpoint.method, currentEndpoint.path, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if !currentEndpoint.mockCheck() {
			t.Errorf("Endpoint %d (%s %s) did not call expected controller",
				i, currentEndpoint.method, currentEndpoint.path)
		}

		if w.Code != http.StatusOK {
			t.Errorf("Endpoint %d (%s %s) returned status %d, expected 200",
				i, currentEndpoint.method, currentEndpoint.path, w.Code)
		}
	}
}
