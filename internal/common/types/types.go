package types

import (
	"time"
)

// Config represents the application configuration
type Config struct {
	GitHubToken string    `json:"github_token"`
	LastBackup  time.Time `json:"last_backup"`
	Machine     Machine   `json:"machine"`
}

// Machine represents a machine configuration
type Machine struct {
	Hostname    string            `json:"hostname"`
	LastSync    time.Time         `json:"last_sync"`
	DotFiles    []DotFile         `json:"dot_files"`
	Apps        []App             `json:"apps"`
	Labels      map[string]string `json:"labels"`
	Description string            `json:"description"`
}

// DotFile represents a dotfile configuration
type DotFile struct {
	Path         string    `json:"path"`
	LastModified time.Time `json:"last_modified"`
	Hash         string    `json:"hash"`
	IsSymlink    bool      `json:"is_symlink"`
}

// App represents an application configuration
type App struct {
	Name         string            `json:"name"`
	Version      string            `json:"version"`
	ConfigFiles  []DotFile         `json:"config_files"`
	Dependencies map[string]string `json:"dependencies"`
}

// Repository represents a GitHub repository
type Repository struct {
	Owner       string `json:"owner"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Private     bool   `json:"private"`
}

// GitHubClient interface defines the methods needed for GitHub operations
type GitHubClient interface {
	// Authentication
	ValidateToken(token string) error
	GetUser() (string, error)

	// Repository operations
	ListRepositories() ([]Repository, error)
	CreateRepository(name, description string, private bool) error
	DeleteRepository(name string) error

	// Content operations
	UploadFile(repo, path string, content []byte, message string) error
	DownloadFile(repo, path string) ([]byte, error)
	ListFiles(repo, path string) ([]string, error)
}

// FileSystem interface defines the methods needed for file operations
type FileSystem interface {
	// File operations
	ReadFile(path string) ([]byte, error)
	WriteFile(path string, content []byte) error
	CreateSymlink(source, target string) error

	// Path operations
	Exists(path string) bool
	IsSymlink(path string) bool
	GetAbsolutePath(path string) string

	// Search operations
	FindDotFiles(paths []string) ([]DotFile, error)
	FindAppConfigs(apps []string) ([]App, error)
}

// ConfigManager interface defines the methods needed for configuration management
type ConfigManager interface {
	// Config operations
	Load() (*Config, error)
	Save(config *Config) error

	// Token operations
	GetToken() (string, error)
	SetToken(token string) error

	// Machine operations
	GetMachine() (*Machine, error)
	SetMachine(machine *Machine) error
}
