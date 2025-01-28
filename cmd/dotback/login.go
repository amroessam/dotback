package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/amroessam/dotback/internal/auth/github"
	"github.com/amroessam/dotback/internal/common/config"
	"github.com/amroessam/dotback/internal/common/logger"
	"github.com/amroessam/dotback/internal/common/types"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to GitHub using a Personal Access Token",
	Long: `Login to GitHub using a Personal Access Token (PAT).
The token should have the following scopes:
- repo (Full control of private repositories)
- read:user (Read all user profile data)
- user:email (Access user email addresses)`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runLogin(cmd, args, nil); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout from GitHub",
	Long:  `Logout from GitHub by removing the stored Personal Access Token.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runLogout(cmd, args); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(logoutCmd)
}

func runLogin(cmd *cobra.Command, args []string, testClient types.GitHubClient) error {
	logger.Info("Starting GitHub authentication")

	// Initialize config manager
	configManager, err := config.NewManager()
	if err != nil {
		logger.Error("Failed to initialize config manager: %v", err)
		return fmt.Errorf("Error: Could not initialize configuration")
	}

	// Check if already logged in
	existingToken, err := configManager.GetToken()
	if err != nil {
		logger.Error("Failed to check existing token: %v", err)
		return fmt.Errorf("Error: Could not check existing login state")
	}

	if existingToken != "" {
		logger.Info("Found existing login, validating...")
		if err := validateAndShowUser(existingToken, testClient); err == nil {
			return fmt.Errorf("Already logged in. Use 'logout' command to log out first")
		}
		// If validation fails, continue with new login
		logger.Info("Existing token is invalid, proceeding with new login")
	}

	// Get token from environment variable or prompt
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		fmt.Print("Enter your GitHub Personal Access Token: ")
		fmt.Scanln(&token)
		token = strings.TrimSpace(token)
	}

	if token == "" {
		logger.Error("No token provided")
		return fmt.Errorf("Error: GitHub token is required")
	}

	// Validate token and get username
	if err := validateAndShowUser(token, testClient); err != nil {
		return err
	}

	// Store token securely
	if err := configManager.SetToken(token); err != nil {
		logger.Error("Failed to store token: %v", err)
		return fmt.Errorf("Error: Could not store token securely")
	}

	return nil
}

func validateAndShowUser(token string, testClient types.GitHubClient) error {
	// Create GitHub client
	var client types.GitHubClient
	if testClient != nil {
		client = testClient
	} else {
		client = github.NewClient(token)
	}

	// Validate token
	logger.Debug("Validating GitHub token")
	if err := client.ValidateToken(token); err != nil {
		logger.Error("Token validation failed: %v", err)
		return fmt.Errorf("Error: Invalid GitHub token")
	}

	// Get username
	username, err := client.GetUser()
	if err != nil {
		logger.Error("Failed to get user info: %v", err)
		return fmt.Errorf("Error: Could not get user information")
	}

	logger.Info("Successfully authenticated as %s", username)
	fmt.Printf("Successfully logged in as %s\n", username)
	return nil
}

func runLogout(cmd *cobra.Command, args []string) error {
	logger.Info("Logging out from GitHub")

	// Initialize config manager
	configManager, err := config.NewManager()
	if err != nil {
		logger.Error("Failed to initialize config manager: %v", err)
		return fmt.Errorf("Error: Could not initialize configuration")
	}

	// Delete token
	if err := configManager.SetToken(""); err != nil {
		logger.Error("Failed to delete token: %v", err)
		return fmt.Errorf("Error: Could not remove stored token")
	}

	logger.Info("Successfully logged out")
	fmt.Println("Successfully logged out")
	return nil
}
