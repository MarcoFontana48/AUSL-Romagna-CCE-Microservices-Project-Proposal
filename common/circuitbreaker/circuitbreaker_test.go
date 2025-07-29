package circuitbreaker

import (
	"errors"
	"github.com/sony/gobreaker/v2"
	"testing"
)

func TestNewCircuitBreaker(t *testing.T) {
	cb := NewCircuitBreaker(DefaultSettings())

	if cb == nil {
		t.Fatal("Expected circuit breaker to be created, got nil")
	}
}

func TestCircuitBreakerInitialState(t *testing.T) {
	cb := NewCircuitBreaker(DefaultSettings())

	if cb.State() != gobreaker.StateClosed {
		t.Errorf("Expected initial state to be Closed, got %v", cb.State())
	}
}

func TestCircuitBreakerTripsOnFailureRatio(t *testing.T) {
	cb := NewCircuitBreaker(DefaultSettings())

	for i := 0; i < 3; i++ {
		_, err := cb.Execute(func() ([]byte, error) {
			return nil, errors.New("test error")
		})
		if err == nil {
			t.Error("Expected error from failed execution")
		}
	}

	if cb.State() != gobreaker.StateOpen {
		t.Errorf("Expected circuit to be Open after failures, got %v", cb.State())
	}
}

func TestCircuitBreakerDoesNotTripWithLowFailureRatio(t *testing.T) {
	cb := NewCircuitBreaker(DefaultSettings())

	_, err := cb.Execute(func() ([]byte, error) {
		return nil, errors.New("test error")
	})
	if err == nil {
		t.Error("Expected error from failed execution")
	}

	for i := 0; i < 2; i++ {
		_, err := cb.Execute(func() ([]byte, error) {
			return []byte("success"), nil
		})
		if err != nil {
			t.Errorf("Expected success, got error: %v", err)
		}
	}

	if cb.State() != gobreaker.StateClosed {
		t.Errorf("Expected circuit to remain Closed, got %v", cb.State())
	}
}

func TestCircuitBreakerDoesNotTripWithFewRequests(t *testing.T) {
	cb := NewCircuitBreaker(DefaultSettings())

	for i := 0; i < 2; i++ {
		_, err := cb.Execute(func() ([]byte, error) {
			return nil, errors.New("test error")
		})
		if err == nil {
			t.Error("Expected error from failed execution")
		}
	}

	if cb.State() != gobreaker.StateClosed {
		t.Errorf("Expected circuit to remain Closed, got %v", cb.State())
	}
}

func TestCustomSettings(t *testing.T) {
	customSettings := gobreaker.Settings{
		Name: "test-breaker",
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.Requests >= 2 // Trip after 2 requests
		},
	}

	cb := NewCircuitBreaker(customSettings)

	for i := 0; i < 2; i++ {
		_, err := cb.Execute(func() ([]byte, error) {
			return nil, errors.New("test error")
		})
		if err == nil {
			t.Error("Expected error from failed execution")
		}
	}

	if cb.State() != gobreaker.StateOpen {
		t.Errorf("Expected circuit to be Open with custom settings, got %v", cb.State())
	}
}
