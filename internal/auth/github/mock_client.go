package github

import (
	"fmt"

	"github.com/amroessam/dotback/internal/common/types"
)

// MockClient implements GitHubClient interface for testing
type MockClient struct {
	token        string
	shouldFail   bool
	mockUsername string
}

// NewMockClient creates a new mock client
func NewMockClient(token string, shouldFail bool, mockUsername string) *MockClient {
	return &MockClient{
		token:        token,
		shouldFail:   shouldFail,
		mockUsername: mockUsername,
	}
}

func (c *MockClient) ValidateToken(token string) error {
	if c.shouldFail {
		return fmt.Errorf("mock validation failed")
	}
	return nil
}

func (c *MockClient) GetUser() (string, error) {
	if c.shouldFail {
		return "", fmt.Errorf("mock user fetch failed")
	}
	return c.mockUsername, nil
}

func (c *MockClient) ListRepositories() ([]types.Repository, error) {
	if c.shouldFail {
		return nil, fmt.Errorf("mock list repositories failed")
	}
	return []types.Repository{}, nil
}

func (c *MockClient) CreateRepository(name, description string, private bool) error {
	if c.shouldFail {
		return fmt.Errorf("mock create repository failed")
	}
	return nil
}

func (c *MockClient) DeleteRepository(name string) error {
	if c.shouldFail {
		return fmt.Errorf("mock delete repository failed")
	}
	return nil
}

func (c *MockClient) UploadFile(repo, path string, content []byte, message string) error {
	if c.shouldFail {
		return fmt.Errorf("mock upload file failed")
	}
	return nil
}

func (c *MockClient) DownloadFile(repo, path string) ([]byte, error) {
	if c.shouldFail {
		return nil, fmt.Errorf("mock download file failed")
	}
	return []byte("mock content"), nil
}

func (c *MockClient) ListFiles(repo, path string) ([]string, error) {
	if c.shouldFail {
		return nil, fmt.Errorf("mock list files failed")
	}
	return []string{}, nil
} 