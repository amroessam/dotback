package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/amroessam/dotback/internal/common/logger"
	"github.com/amroessam/dotback/internal/common/storage"
	"github.com/amroessam/dotback/internal/common/types"
)

const (
	configFileName = "config.json"
)

// GetConfigDirFunc is a function type for getting the config directory
type GetConfigDirFunc func() (string, error)

// DefaultGetConfigDir returns the default config directory
func DefaultGetConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error getting home directory: %w", err)
	}

	return filepath.Join(homeDir, ".config", "dotback"), nil
}

// GetConfigDir is the current implementation of getting the config directory
var GetConfigDir GetConfigDirFunc = DefaultGetConfigDir

// Manager handles configuration storage and retrieval
type Manager struct {
	configPath string
	storage    *storage.KeyringStorage
	config     *types.Config
}

// NewManager creates a new configuration manager
func NewManager() (*Manager, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(configDir, configFileName)
	storage := storage.NewKeyringStorage()

	return &Manager{
		configPath: configPath,
		storage:    storage,
		config:     &types.Config{},
	}, nil
}

// Load loads the configuration from disk
func (m *Manager) Load() (*types.Config, error) {
	logger.Debug("Loading configuration from %s", m.configPath)

	// Read config file
	data, err := os.ReadFile(m.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Debug("Configuration file not found, using defaults")
			return m.config, nil
		}
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	// Parse config
	if err := json.Unmarshal(data, m.config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	return m.config, nil
}

// Save saves the configuration to disk
func (m *Manager) Save(config *types.Config) error {
	logger.Debug("Saving configuration to %s", m.configPath)

	// Create config directory if it doesn't exist
	configDir := filepath.Dir(m.configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("error creating config directory: %w", err)
	}

	// Marshal config
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("error encoding config: %w", err)
	}

	// Write config file
	if err := os.WriteFile(m.configPath, data, 0644); err != nil {
		return fmt.Errorf("error writing config file: %w", err)
	}

	m.config = config
	return nil
}

// GetToken retrieves the GitHub token from secure storage
func (m *Manager) GetToken() (string, error) {
	return m.storage.GetToken()
}

// SetToken stores the GitHub token in secure storage
func (m *Manager) SetToken(token string) error {
	if token == "" {
		return m.storage.DeleteToken()
	}
	return m.storage.StoreToken(token)
}

// GetMachine gets the current machine configuration
func (m *Manager) GetMachine() (*types.Machine, error) {
	if m.config == nil {
		if _, err := m.Load(); err != nil {
			return nil, err
		}
	}
	return &m.config.Machine, nil
}

// SetMachine updates the current machine configuration
func (m *Manager) SetMachine(machine *types.Machine) error {
	if m.config == nil {
		if _, err := m.Load(); err != nil {
			return err
		}
	}
	m.config.Machine = *machine
	return m.Save(m.config)
}
