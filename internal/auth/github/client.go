package github

import (
	"context"
	"fmt"

	"github.com/amroessam/dotback/internal/common/types"
	"github.com/google/go-github/v60/github"
	"golang.org/x/oauth2"
)

type Client struct {
	client *github.Client
	ctx    context.Context
}

// NewClient creates a new GitHub client
func NewClient(token string) *Client {
	ctx := context.Background()
	var ghClient *github.Client

	if token != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc := oauth2.NewClient(ctx, ts)
		ghClient = github.NewClient(tc)
	} else {
		ghClient = github.NewClient(nil)
	}

	return &Client{
		client: ghClient,
		ctx:    ctx,
	}
}

// ValidateToken validates the GitHub token
func (c *Client) ValidateToken(token string) error {
	// For testing, if we already have a client set up, use it
	var client *github.Client
	if c.client != nil && c.client.BaseURL != nil && c.client.BaseURL.String() != "https://api.github.com/" {
		client = c.client
	} else {
		// For production, create a new client with the token
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc := oauth2.NewClient(c.ctx, ts)
		client = github.NewClient(tc)
	}

	_, _, err := client.Users.Get(c.ctx, "")
	if err != nil {
		return fmt.Errorf("invalid token: %w", err)
	}
	return nil
}

// GetUser gets the authenticated user's information
func (c *Client) GetUser() (string, error) {
	user, _, err := c.client.Users.Get(c.ctx, "")
	if err != nil {
		return "", fmt.Errorf("error getting user: %w", err)
	}
	return user.GetLogin(), nil
}

// ListRepositories lists the authenticated user's repositories
func (c *Client) ListRepositories() ([]types.Repository, error) {
	opt := &github.RepositoryListOptions{
		Type: "owner",
	}
	repos, _, err := c.client.Repositories.List(c.ctx, "", opt)
	if err != nil {
		return nil, fmt.Errorf("error listing repositories: %w", err)
	}

	var result []types.Repository
	for _, repo := range repos {
		result = append(result, types.Repository{
			Owner:       repo.GetOwner().GetLogin(),
			Name:        repo.GetName(),
			Description: repo.GetDescription(),
			Private:     repo.GetPrivate(),
		})
	}
	return result, nil
}

// CreateRepository creates a new repository
func (c *Client) CreateRepository(name, description string, private bool) error {
	repo := &github.Repository{
		Name:        github.String(name),
		Description: github.String(description),
		Private:     github.Bool(private),
	}
	_, _, err := c.client.Repositories.Create(c.ctx, "", repo)
	if err != nil {
		return fmt.Errorf("error creating repository: %w", err)
	}
	return nil
}

// DeleteRepository deletes a repository
func (c *Client) DeleteRepository(name string) error {
	user, _, err := c.client.Users.Get(c.ctx, "")
	if err != nil {
		return fmt.Errorf("error getting user: %w", err)
	}

	_, err = c.client.Repositories.Delete(c.ctx, user.GetLogin(), name)
	if err != nil {
		return fmt.Errorf("error deleting repository: %w", err)
	}
	return nil
}

// UploadFile uploads a file to a repository
func (c *Client) UploadFile(repo, path string, content []byte, message string) error {
	user, _, err := c.client.Users.Get(c.ctx, "")
	if err != nil {
		return fmt.Errorf("error getting user: %w", err)
	}

	// Check if file exists to get the SHA
	var sha *string
	fileContent, _, _, err := c.client.Repositories.GetContents(c.ctx, user.GetLogin(), repo, path, &github.RepositoryContentGetOptions{})
	if err == nil && fileContent != nil {
		sha = fileContent.SHA
	}

	// Create or update file
	opts := &github.RepositoryContentFileOptions{
		Message: github.String(message),
		Content: content,
		SHA:     sha,
	}

	_, _, err = c.client.Repositories.CreateFile(c.ctx, user.GetLogin(), repo, path, opts)
	if err != nil {
		return fmt.Errorf("error uploading file: %w", err)
	}

	return nil
}

// DownloadFile downloads a file from a repository
func (c *Client) DownloadFile(repo, path string) ([]byte, error) {
	user, _, err := c.client.Users.Get(c.ctx, "")
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	fileContent, _, _, err := c.client.Repositories.GetContents(c.ctx, user.GetLogin(), repo, path, &github.RepositoryContentGetOptions{})
	if err != nil {
		return nil, fmt.Errorf("error downloading file: %w", err)
	}

	content, err := fileContent.GetContent()
	if err != nil {
		return nil, fmt.Errorf("error decoding content: %w", err)
	}

	return []byte(content), nil
}

// ListFiles lists files in a repository path
func (c *Client) ListFiles(repo, path string) ([]string, error) {
	user, _, err := c.client.Users.Get(c.ctx, "")
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	_, contents, _, err := c.client.Repositories.GetContents(c.ctx, user.GetLogin(), repo, path, &github.RepositoryContentGetOptions{})
	if err != nil {
		return nil, fmt.Errorf("error listing files: %w", err)
	}

	var files []string
	for _, content := range contents {
		files = append(files, content.GetName())
	}

	return files, nil
}
