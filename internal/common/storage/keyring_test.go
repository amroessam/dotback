package storage

import (
	"testing"

	"github.com/zalando/go-keyring"
)

func TestKeyringStorage(t *testing.T) {
	storage := NewKeyringStorage()

	// Clean up any existing token
	storage.DeleteToken()

	// Test storing token
	t.Run("Store token", func(t *testing.T) {
		err := storage.StoreToken("test-token")
		if err != nil {
			t.Errorf("StoreToken() error = %v", err)
		}

		// Verify token was stored
		token, err := keyring.Get(serviceName, tokenKey)
		if err != nil {
			t.Errorf("Failed to verify token storage: %v", err)
		}
		if token != "test-token" {
			t.Errorf("Stored token = %v, want %v", token, "test-token")
		}
	})

	// Test retrieving token
	t.Run("Get token", func(t *testing.T) {
		token, err := storage.GetToken()
		if err != nil {
			t.Errorf("GetToken() error = %v", err)
		}
		if token != "test-token" {
			t.Errorf("GetToken() = %v, want %v", token, "test-token")
		}
	})

	// Test deleting token
	t.Run("Delete token", func(t *testing.T) {
		err := storage.DeleteToken()
		if err != nil {
			t.Errorf("DeleteToken() error = %v", err)
		}

		// Verify token was deleted
		token, err := storage.GetToken()
		if err != nil {
			t.Errorf("GetToken() error = %v", err)
		}
		if token != "" {
			t.Errorf("Token still exists after deletion")
		}
	})

	// Test getting non-existent token
	t.Run("Get non-existent token", func(t *testing.T) {
		token, err := storage.GetToken()
		if err != nil {
			t.Errorf("GetToken() error = %v", err)
		}
		if token != "" {
			t.Errorf("Expected empty token, got %v", token)
		}
	})

	// Test deleting non-existent token
	t.Run("Delete non-existent token", func(t *testing.T) {
		err := storage.DeleteToken()
		if err != nil {
			t.Errorf("DeleteToken() error = %v", err)
		}
	})
}
