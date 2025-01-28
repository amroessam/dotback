package logger

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

func TestLogger(t *testing.T) {
	// Capture stdout and stderr
	var stdout, stderr bytes.Buffer
	debugLogger = log.New(&stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	infoLogger = log.New(&stdout, "INFO: ", log.Ldate|log.Ltime)
	errorLogger = log.New(&stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	tests := []struct {
		name     string
		logFunc  func(string, ...interface{})
		message  string
		debug    bool
		output   *bytes.Buffer
		expected string
	}{
		{
			name:     "Info message",
			logFunc:  Info,
			message:  "test info message",
			output:   &stdout,
			expected: "INFO: ",
		},
		{
			name:     "Error message",
			logFunc:  Error,
			message:  "test error message",
			output:   &stderr,
			expected: "ERROR: ",
		},
		{
			name:     "Debug message with debug mode off",
			logFunc:  Debug,
			message:  "test debug message",
			debug:    false,
			output:   &stdout,
			expected: "",
		},
		{
			name:     "Debug message with debug mode on",
			logFunc:  Debug,
			message:  "test debug message",
			debug:    true,
			output:   &stdout,
			expected: "DEBUG: ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear buffers
			stdout.Reset()
			stderr.Reset()

			// Set debug mode
			SetDebugMode(tt.debug)

			// Execute log function
			tt.logFunc(tt.message)

			// Check output
			output := tt.output.String()
			if tt.expected == "" {
				if output != "" {
					t.Errorf("expected no output, got %q", output)
				}
			} else if !strings.Contains(output, tt.expected) {
				t.Errorf("expected output containing %q, got %q", tt.expected, output)
			}
			if tt.expected != "" && !strings.Contains(output, tt.message) {
				t.Errorf("expected output containing message %q, got %q", tt.message, output)
			}
		})
	}
}
