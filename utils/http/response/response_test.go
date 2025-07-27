package response

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// mock utils package functions for testing
type TestData struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func TestSendResponse(t *testing.T) {
	tests := []struct {
		name           string
		message        interface{}
		expectedStatus int
		expectedBody   string
		expectError    bool
	}{
		{
			name:           "valid struct response",
			message:        TestData{Message: "success", Code: 200},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"success","code":200}`,
			expectError:    false,
		},
		{
			name:           "string response",
			message:        "hello world",
			expectedStatus: http.StatusOK,
			expectedBody:   `"hello world"`,
			expectError:    false,
		},
		{
			name:           "nil response",
			message:        nil,
			expectedStatus: http.StatusOK,
			expectedBody:   "null",
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request and response recorder
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			w := httptest.NewRecorder()

			// Call the function
			err := SendOkResponse(w, req, tt.message)

			// Check error expectation
			if (err != nil) != tt.expectError {
				t.Errorf("SendOkResponse() error = %v, expectError %v", err, tt.expectError)
				return
			}

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("SendOkResponse() status = %v, expected %v", w.Code, tt.expectedStatus)
			}

			// Check content type
			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("SendOkResponse() Content-Type = %v, expected application/json", contentType)
			}

			// Check response body
			body := strings.TrimSpace(w.Body.String())
			if body != tt.expectedBody {
				t.Errorf("SendOkResponse() body = %v, expected %v", body, tt.expectedBody)
			}
		})
	}
}

func TestSendErrorResponse(t *testing.T) {
	tests := []struct {
		name           string
		inputError     error
		expectedStatus int
		expectError    bool
	}{
		{
			name:           "generic error",
			inputError:     errors.New("something went wrong"),
			expectedStatus: http.StatusInternalServerError,
			expectError:    false,
		},
		{
			name:           "custom error message",
			inputError:     errors.New("database connection failed"),
			expectedStatus: http.StatusInternalServerError,
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request and response recorder
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			w := httptest.NewRecorder()

			// Call the function
			err := SendErrorResponse(w, req, tt.inputError)

			// Check error expectation
			if (err != nil) != tt.expectError {
				t.Errorf("SendErrorResponse() error = %v, expectError %v", err, tt.expectError)
				return
			}

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("SendErrorResponse() status = %v, expected %v", w.Code, tt.expectedStatus)
			}

			// Check content type
			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("SendErrorResponse() Content-Type = %v, expected application/json", contentType)
			}

			// Check response body contains error message
			body := w.Body.String()
			if !strings.Contains(body, tt.inputError.Error()) {
				t.Errorf("SendErrorResponse() body should contain error message %v, got %v", tt.inputError.Error(), body)
			}

			// Check JSON structure
			expectedPrefix := `{"error":"`
			expectedSuffix := `"}`
			if !strings.HasPrefix(body, expectedPrefix) || !strings.HasSuffix(strings.TrimSpace(body), expectedSuffix) {
				t.Errorf("SendErrorResponse() body format incorrect, got %v", body)
			}
		})
	}
}

// Test headers are set correctly
func TestResponseHeaders(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	err := SendOkResponse(w, req, "test")
	if err != nil {
		t.Fatalf("SendOkResponse() failed: %v", err)
	}

	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type header to be application/json")
	}
}

func TestErrorResponseHeaders(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	err := SendErrorResponse(w, req, errors.New("test error"))
	if err != nil {
		t.Fatalf("SendErrorResponse() failed: %v", err)
	}

	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type header to be application/json")
	}
}

// Test with nil request (edge case)
func TestSendResponseWithNilRequest(t *testing.T) {
	w := httptest.NewRecorder()

	// This should not panic even with nil request
	err := SendOkResponse(w, nil, "test message")
	if err == nil {
		t.Errorf("SendOkResponse() with nil request failed: %v", err)
	}

	if w.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %v", w.Code)
	}
}

func TestSendErrorResponseWithNilRequest(t *testing.T) {
	w := httptest.NewRecorder()

	// This should not panic even with nil request
	err := SendErrorResponse(w, nil, errors.New("test error"))
	if err == nil {
		t.Errorf("SendErrorResponse() with nil request failed: %v", err)
	}
}
