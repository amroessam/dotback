package storage

import (
	"fmt"

	"github.com/amroessam/dotback/internal/common/logger"
	"github.com/zalando/go-keyring"
)

const (
	serviceName = "dotback"
	tokenKey    = "github_token"
)

// KeyringStorage implements secure storage using the system keyring
type KeyringStorage struct{}

// NewKeyringStorage creates a new keyring storage
func NewKeyringStorage() *KeyringStorage {
	return &KeyringStorage{}
}

// StoreToken stores the GitHub token securely
func (k *KeyringStorage) StoreToken(token string) error {
	logger.Debug("Storing GitHub token in keyring")
	err := keyring.Set(serviceName, tokenKey, token)
	if err != nil {
		return fmt.Errorf("failed to store token: %w", err)
	}
	return nil
}

// GetToken retrieves the GitHub token from secure storage
func (k *KeyringStorage) GetToken() (string, error) {
	logger.Debug("Retrieving GitHub token from keyring")
	token, err := keyring.Get(serviceName, tokenKey)
	if err != nil {
		if err == keyring.ErrNotFound {
			return "", nil
		}
		return "", fmt.Errorf("failed to retrieve token: %w", err)
	}
	return token, nil
}

// DeleteToken removes the GitHub token from secure storage
func (k *KeyringStorage) DeleteToken() error {
	logger.Debug("Deleting GitHub token from keyring")
	err := keyring.Delete(serviceName, tokenKey)
	if err != nil && err != keyring.ErrNotFound {
		return fmt.Errorf("failed to delete token: %w", err)
	}
	return nil
}
