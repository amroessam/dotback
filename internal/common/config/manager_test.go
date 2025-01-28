package config

import (
	"os"
	"testing"
	"time"

	"github.com/amroessam/dotback/internal/common/types"
)

func TestManager(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "dotback-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Override GetConfigDir for testing
	oldGetConfigDir := GetConfigDir
	defer func() { GetConfigDir = oldGetConfigDir }()
	GetConfigDir = func() (string, error) {
		return tempDir, nil
	}

	// Create manager
	manager, err := NewManager()
	if err != nil {
		t.Fatalf("NewManager() error = %v", err)
	}

	// Test loading empty config
	t.Run("Load empty config", func(t *testing.T) {
		config, err := manager.Load()
		if err != nil {
			t.Errorf("Load() error = %v", err)
		}
		if config == nil {
			t.Error("Load() returned nil config")
		}
	})

	// Test saving and loading config
	t.Run("Save and load config", func(t *testing.T) {
		testConfig := &types.Config{
			LastBackup: time.Now(),
			Machine: types.Machine{
				Hostname:    "test-host",
				LastSync:    time.Now(),
				Description: "Test machine",
				Labels: map[string]string{
					"env": "test",
				},
			},
		}

		// Save config
		if err := manager.Save(testConfig); err != nil {
			t.Errorf("Save() error = %v", err)
		}

		// Load config
		loadedConfig, err := manager.Load()
		if err != nil {
			t.Errorf("Load() error = %v", err)
		}

		// Compare configs
		if loadedConfig.Machine.Hostname != testConfig.Machine.Hostname {
			t.Errorf("Loaded config hostname = %v, want %v", loadedConfig.Machine.Hostname, testConfig.Machine.Hostname)
		}
		if loadedConfig.Machine.Description != testConfig.Machine.Description {
			t.Errorf("Loaded config description = %v, want %v", loadedConfig.Machine.Description, testConfig.Machine.Description)
		}
		if loadedConfig.Machine.Labels["env"] != testConfig.Machine.Labels["env"] {
			t.Errorf("Loaded config label = %v, want %v", loadedConfig.Machine.Labels["env"], testConfig.Machine.Labels["env"])
		}
	})

	// Test machine operations
	t.Run("Machine operations", func(t *testing.T) {
		testMachine := &types.Machine{
			Hostname:    "new-host",
			LastSync:    time.Now(),
			Description: "New test machine",
			Labels: map[string]string{
				"env": "prod",
			},
		}

		// Set machine
		if err := manager.SetMachine(testMachine); err != nil {
			t.Errorf("SetMachine() error = %v", err)
		}

		// Get machine
		loadedMachine, err := manager.GetMachine()
		if err != nil {
			t.Errorf("GetMachine() error = %v", err)
		}

		// Compare machines
		if loadedMachine.Hostname != testMachine.Hostname {
			t.Errorf("Loaded machine hostname = %v, want %v", loadedMachine.Hostname, testMachine.Hostname)
		}
		if loadedMachine.Description != testMachine.Description {
			t.Errorf("Loaded machine description = %v, want %v", loadedMachine.Description, testMachine.Description)
		}
		if loadedMachine.Labels["env"] != testMachine.Labels["env"] {
			t.Errorf("Loaded machine label = %v, want %v", loadedMachine.Labels["env"], testMachine.Labels["env"])
		}
	})

	// Test token operations
	t.Run("Token operations", func(t *testing.T) {
		// Set token
		if err := manager.SetToken("test-token"); err != nil {
			t.Errorf("SetToken() error = %v", err)
		}

		// Get token
		token, err := manager.GetToken()
		if err != nil {
			t.Errorf("GetToken() error = %v", err)
		}
		if token != "test-token" {
			t.Errorf("GetToken() = %v, want %v", token, "test-token")
		}

		// Delete token
		if err := manager.SetToken(""); err != nil {
			t.Errorf("SetToken(\"\") error = %v", err)
		}

		// Verify token was deleted
		token, err = manager.GetToken()
		if err != nil {
			t.Errorf("GetToken() error = %v", err)
		}
		if token != "" {
			t.Errorf("Token still exists after deletion")
		}
	})
}
