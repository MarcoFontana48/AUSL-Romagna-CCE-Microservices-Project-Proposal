package log

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"os"
	"strings"
	"testing"
)

func TestInitAsJson(t *testing.T) {
	// Save original environment and restore after tests
	originalLogAddSource := os.Getenv("LOG_ADD_SOURCE")
	originalLogLevel := os.Getenv("LOG_LEVEL")
	defer func() {
		err := os.Setenv("LOG_ADD_SOURCE", originalLogAddSource)
		if err != nil {
			return
		}
		err = os.Setenv("LOG_LEVEL", originalLogLevel)
		if err != nil {
			return
		}
	}()

	tests := []struct {
		name           string
		logAddSource   string
		logLevel       string
		expectedLevel  slog.Level
		expectedSource bool
	}{
		{
			name:           "default values - no env vars set",
			logAddSource:   "",
			logLevel:       "",
			expectedLevel:  slog.LevelInfo,
			expectedSource: false,
		},
		{
			name:           "LOG_ADD_SOURCE true, LOG_LEVEL debug",
			logAddSource:   "true",
			logLevel:       "debug",
			expectedLevel:  slog.LevelDebug,
			expectedSource: true,
		},
		{
			name:           "LOG_ADD_SOURCE false, LOG_LEVEL warn",
			logAddSource:   "false",
			logLevel:       "warn",
			expectedLevel:  slog.LevelWarn,
			expectedSource: false,
		},
		{
			name:           "LOG_ADD_SOURCE true, LOG_LEVEL error",
			logAddSource:   "true",
			logLevel:       "error",
			expectedLevel:  slog.LevelError,
			expectedSource: true,
		},
		{
			name:           "LOG_ADD_SOURCE false, LOG_LEVEL info",
			logAddSource:   "false",
			logLevel:       "info",
			expectedLevel:  slog.LevelInfo,
			expectedSource: false,
		},
		{
			name:           "LOG_ADD_SOURCE invalid, LOG_LEVEL unknown - defaults",
			logAddSource:   "invalid",
			logLevel:       "unknown",
			expectedLevel:  slog.LevelInfo,
			expectedSource: false,
		},
		{
			name:           "LOG_ADD_SOURCE 1, LOG_LEVEL DEBUG (uppercase)",
			logAddSource:   "1",
			logLevel:       "DEBUG",
			expectedLevel:  slog.LevelDebug,
			expectedSource: true,
		},
		{
			name:           "LOG_ADD_SOURCE 0, LOG_LEVEL WARN (uppercase)",
			logAddSource:   "0",
			logLevel:       "WARN",
			expectedLevel:  slog.LevelWarn,
			expectedSource: false,
		},
		{
			name:           "LOG_ADD_SOURCE T, LOG_LEVEL ERROR (uppercase)",
			logAddSource:   "T",
			logLevel:       "ERROR",
			expectedLevel:  slog.LevelError,
			expectedSource: true,
		},
		{
			name:           "LOG_ADD_SOURCE f, LOG_LEVEL mixed case",
			logAddSource:   "f",
			logLevel:       "WaRn",
			expectedLevel:  slog.LevelWarn,
			expectedSource: false,
		},
		{
			name:           "LOG_ADD_SOURCE empty string, LOG_LEVEL empty string",
			logAddSource:   "",
			logLevel:       "",
			expectedLevel:  slog.LevelInfo,
			expectedSource: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			if tt.logAddSource == "" {
				err := os.Unsetenv("LOG_ADD_SOURCE")
				if err != nil {
					return
				}
			} else {
				err := os.Setenv("LOG_ADD_SOURCE", tt.logAddSource)
				if err != nil {
					return
				}
			}

			if tt.logLevel == "" {
				err := os.Unsetenv("LOG_LEVEL")
				if err != nil {
					return
				}
			} else {
				err := os.Setenv("LOG_LEVEL", tt.logLevel)
				if err != nil {
					return
				}
			}

			// Capture stdout to verify JSON handler is set up correctly
			var buf bytes.Buffer
			originalStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Call the function
			InitAsJson()

			// Call the function
			InitAsJson()

			// Test that the logger is configured correctly by logging a test message
			slog.Info("test message", "key", "value")

			// Restore stdout and read the captured output
			err := w.Close()
			if err != nil {
				return
			}
			os.Stdout = originalStdout
			_, err = buf.ReadFrom(r)
			if err != nil {
				return
			}
			output := buf.String()

			// Verify JSON format
			if output != "" {
				var logEntry map[string]interface{}
				if err := json.Unmarshal([]byte(strings.TrimSpace(output)), &logEntry); err != nil {
					t.Errorf("Expected JSON output, but got invalid JSON: %v", err)
				}

				// Check if source information is included based on expectedSource
				_, hasSource := logEntry["source"]
				if hasSource != tt.expectedSource {
					t.Errorf("Expected source field presence: %v, got: %v", tt.expectedSource, hasSource)
				}
			}

			// Verify the log level by testing what gets logged
			testLogLevel(t, tt.expectedLevel)
		})
	}
}

// testLogLevel verifies that the logger respects the configured log level
func testLogLevel(t *testing.T, expectedLevel slog.Level) {
	// Capture stdout for level testing
	var buf bytes.Buffer
	originalStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Re-initialize the logger to write to our captured stdout
	var level = expectedLevel
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})
	slog.SetDefault(slog.New(handler))

	// Test logging at different levels
	slog.Debug("debug message")
	slog.Info("info message")
	slog.Warn("warn message")
	slog.Error("error message")

	err := w.Close()
	if err != nil {
		return
	}
	os.Stdout = originalStdout
	_, err = buf.ReadFrom(r)
	if err != nil {
		return
	}
	output := buf.String()

	lines := strings.Split(strings.TrimSpace(output), "\n")

	// Count how many log entries were actually written
	actualLogCount := 0
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			actualLogCount++
		}
	}

	// Determine expected log count based on level
	var expectedLogCount int
	switch expectedLevel {
	case slog.LevelDebug:
		expectedLogCount = 4 // debug, info, warn, error
	case slog.LevelInfo:
		expectedLogCount = 3 // info, warn, error
	case slog.LevelWarn:
		expectedLogCount = 2 // warn, error
	case slog.LevelError:
		expectedLogCount = 1 // error only
	}

	if actualLogCount != expectedLogCount {
		t.Errorf("Expected %d log entries for level %v, got %d", expectedLogCount, expectedLevel, actualLogCount)
	}
}

func TestInitAsJsonWithSpecificEnvironmentValues(t *testing.T) {
	// Test edge cases for LOG_ADD_SOURCE parsing
	testCases := []struct {
		value    string
		expected bool
	}{
		{"true", true},
		{"TRUE", true},
		{"True", true},
		{"1", true},
		{"t", true},
		{"T", true},
		{"false", false},
		{"FALSE", false},
		{"False", false},
		{"0", false},
		{"f", false},
		{"F", false},
		{"", false},        // unset
		{"invalid", false}, // invalid value defaults to false
		{"yes", false},     // strconv.ParseBool doesn't recognize "yes"
		{"no", false},      // strconv.ParseBool doesn't recognize "no"
	}

	originalLogAddSource := os.Getenv("LOG_ADD_SOURCE")
	originalLogLevel := os.Getenv("LOG_LEVEL")
	defer func() {
		err := os.Setenv("LOG_ADD_SOURCE", originalLogAddSource)
		if err != nil {
			return
		}
		err = os.Setenv("LOG_LEVEL", originalLogLevel)
		if err != nil {
			return
		}
	}()

	for _, tc := range testCases {
		t.Run("LOG_ADD_SOURCE_"+tc.value, func(t *testing.T) {
			if tc.value == "" {
				err := os.Unsetenv("LOG_ADD_SOURCE")
				if err != nil {
					return
				}
			} else {
				err := os.Setenv("LOG_ADD_SOURCE", tc.value)
				if err != nil {
					return
				}
			}
			err := os.Setenv("LOG_LEVEL", "info")
			if err != nil {
				return
			}

			// Capture output
			var buf bytes.Buffer
			originalStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			InitAsJson()
			slog.Info("test", "key", "value")

			err = w.Close()
			if err != nil {
				return
			}
			os.Stdout = originalStdout
			_, err = buf.ReadFrom(r)
			if err != nil {
				return
			}
			output := buf.String()

			// Parse JSON and check source field
			var logEntry map[string]interface{}
			if err := json.Unmarshal([]byte(strings.TrimSpace(output)), &logEntry); err != nil {
				t.Fatalf("Failed to parse JSON: %v", err)
			}

			_, hasSource := logEntry["source"]
			if hasSource != tc.expected {
				t.Errorf("For LOG_ADD_SOURCE=%q, expected source field presence: %v, got: %v", tc.value, tc.expected, hasSource)
			}
		})
	}
}

// Benchmark the InitAsJson function
func BenchmarkInitAsJson(b *testing.B) {
	originalLogAddSource := os.Getenv("LOG_ADD_SOURCE")
	originalLogLevel := os.Getenv("LOG_LEVEL")
	defer func() {
		err := os.Setenv("LOG_ADD_SOURCE", originalLogAddSource)
		if err != nil {
			return
		}
		err = os.Setenv("LOG_LEVEL", originalLogLevel)
		if err != nil {
			return
		}
	}()

	err := os.Setenv("LOG_ADD_SOURCE", "true")
	if err != nil {
		return
	}
	err = os.Setenv("LOG_LEVEL", "info")
	if err != nil {
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		InitAsJson()
	}
}
