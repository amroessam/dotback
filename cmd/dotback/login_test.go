package main

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/amroessam/dotback/internal/auth/github"
	"github.com/amroessam/dotback/internal/common/config"
	"github.com/amroessam/dotback/internal/common/logger"
)

func TestRunLogin(t *testing.T) {
	// Save original stdin and stdout
	oldStdin := os.Stdin
	oldStdout := os.Stdout
	defer func() {
		os.Stdin = oldStdin
		os.Stdout = oldStdout
	}()

	// Create a temporary directory for config
	tempDir, err := os.MkdirTemp("", "dotback-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Override config directory for testing
	oldGetConfigDir := config.GetConfigDir
	defer func() { config.GetConfigDir = oldGetConfigDir }()
	config.GetConfigDir = func() (string, error) {
		return tempDir, nil
	}

	tests := []struct {
		name          string
		token         string
		input         string
		mockClient    *github.MockClient
		expectedUser  string
		expectedError bool
	}{
		{
			name:          "Valid token from environment",
			token:         "valid-token",
			input:         "",
			mockClient:    github.NewMockClient("valid-token", false, "testuser"),
			expectedUser:  "testuser",
			expectedError: false,
		},
		{
			name:          "Valid token from input",
			token:         "",
			input:         "valid-token\n",
			mockClient:    github.NewMockClient("valid-token", false, "testuser"),
			expectedUser:  "testuser",
			expectedError: false,
		},
		{
			name:          "Empty token",
			token:         "",
			input:         "\n",
			mockClient:    github.NewMockClient("", true, ""),
			expectedUser:  "",
			expectedError: true,
		},
		{
			name:          "Invalid token",
			token:         "invalid-token",
			input:         "",
			mockClient:    github.NewMockClient("invalid-token", true, ""),
			expectedUser:  "",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up any existing token
			configManager, err := config.NewManager()
			if err != nil {
				t.Fatalf("Failed to create config manager: %v", err)
			}
			if err := configManager.SetToken(""); err != nil {
				t.Fatalf("Failed to clean up token: %v", err)
			}

			// Set up environment
			if tt.token != "" {
				os.Setenv("GITHUB_TOKEN", tt.token)
				defer os.Unsetenv("GITHUB_TOKEN")
			} else {
				os.Unsetenv("GITHUB_TOKEN")
			}

			// Create a pipe for stdin simulation
			r, w, _ := os.Pipe()
			os.Stdin = r

			// Create a pipe for stdout capture
			stdout, wout, _ := os.Pipe()
			os.Stdout = wout

			// Write test input
			if tt.input != "" {
				w.Write([]byte(tt.input))
			}
			w.Close()

			// Enable debug logging for tests
			logger.SetDebugMode(true)

			// Run the command
			err = runLogin(nil, nil, tt.mockClient)

			// Close stdout pipe and read output
			wout.Close()
			var buf bytes.Buffer
			buf.ReadFrom(stdout)
			output := buf.String()

			// Check results
			if tt.expectedError {
				if err == nil {
					t.Error("Expected error but got success")
				}
			} else {
				if err != nil {
					t.Errorf("Expected success but got error: %v", err)
				}
				if !strings.Contains(output, "Successfully logged in") {
					t.Errorf("Expected success message, got: %s", output)
				}
				if tt.expectedUser != "" && !strings.Contains(output, tt.expectedUser) {
					t.Errorf("Expected username %s in output, got: %s", tt.expectedUser, output)
				}

				// Verify token was stored
				storedToken, err := configManager.GetToken()
				if err != nil {
					t.Errorf("Failed to get stored token: %v", err)
				}
				if storedToken == "" {
					t.Error("Token was not stored")
				}
			}
		})
	}
}

func TestRunLogout(t *testing.T) {
	// Save original stdout
	oldStdout := os.Stdout
	defer func() {
		os.Stdout = oldStdout
	}()

	// Create a temporary directory for config
	tempDir, err := os.MkdirTemp("", "dotback-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Override config directory for testing
	oldGetConfigDir := config.GetConfigDir
	defer func() { config.GetConfigDir = oldGetConfigDir }()
	config.GetConfigDir = func() (string, error) {
		return tempDir, nil
	}

	tests := []struct {
		name          string
		setupToken    string
		expectedError bool
	}{
		{
			name:          "Successful logout",
			setupToken:    "test-token",
			expectedError: false,
		},
		{
			name:          "Logout when not logged in",
			setupToken:    "",
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a pipe for stdout capture
			stdout, wout, _ := os.Pipe()
			os.Stdout = wout

			// Set up initial token if needed
			configManager, err := config.NewManager()
			if err != nil {
				t.Fatalf("Failed to create config manager: %v", err)
			}
			if tt.setupToken != "" {
				if err := configManager.SetToken(tt.setupToken); err != nil {
					t.Fatalf("Failed to set up token: %v", err)
				}
			}

			// Run the command
			logoutErr := runLogout(nil, nil)

			// Close stdout pipe and read output
			wout.Close()
			var buf bytes.Buffer
			buf.ReadFrom(stdout)
			output := buf.String()

			// Check results
			if tt.expectedError {
				if logoutErr == nil {
					t.Error("Expected error but got success")
				}
			} else {
				if logoutErr != nil {
					t.Errorf("Expected success but got error: %v", logoutErr)
				}
				if !strings.Contains(output, "Successfully logged out") {
					t.Errorf("Expected success message, got: %s", output)
				}

				// Verify token was removed
				storedToken, err := configManager.GetToken()
				if err != nil {
					t.Errorf("Failed to get stored token: %v", err)
				}
				if storedToken != "" {
					t.Error("Token was not removed")
				}
			}
		})
	}
}
