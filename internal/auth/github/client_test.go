package github

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/amroessam/dotback/internal/common/types"
	"github.com/google/go-github/v60/github"
)

func setupTestServer(t *testing.T, handler http.HandlerFunc) (*httptest.Server, *Client) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add /api/v3 prefix to match GitHub API
		r.URL.Path = "/api/v3" + r.URL.Path
		handler(w, r)
	}))

	// Create a custom HTTP client with auth header
	httpClient := &http.Client{
		Transport: roundTripperFunc(func(r *http.Request) (*http.Response, error) {
			r.Header.Set("Authorization", "token test-token")
			return http.DefaultTransport.RoundTrip(r)
		}),
	}

	// Create a GitHub client with our custom HTTP client
	ghClient := github.NewClient(httpClient)
	baseURL, _ := url.Parse(server.URL + "/")
	ghClient.BaseURL = baseURL

	// Create our client with the custom GitHub client
	client := &Client{
		client: ghClient,
		ctx:    context.Background(),
	}

	return server, client
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func decodeRequest(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func TestValidateToken(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		response   string
		wantErr    bool
	}{
		{
			name:       "Valid token",
			statusCode: http.StatusOK,
			response:   `{"login": "testuser"}`,
			wantErr:    false,
		},
		{
			name:       "Invalid token",
			statusCode: http.StatusUnauthorized,
			response:   `{"message": "Bad credentials"}`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, client := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Header.Get("Authorization") != "token test-token" {
					t.Error("Expected Authorization header with token")
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			})
			defer server.Close()

			err := client.ValidateToken("test-token")
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		response   string
		want       string
		wantErr    bool
	}{
		{
			name:       "Valid user",
			statusCode: http.StatusOK,
			response:   `{"login": "testuser"}`,
			want:       "testuser",
			wantErr:    false,
		},
		{
			name:       "Invalid response",
			statusCode: http.StatusUnauthorized,
			response:   `{"message": "Bad credentials"}`,
			want:       "",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, client := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			})
			defer server.Close()

			got, err := client.GetUser()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestListRepositories(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		response   string
		want       []types.Repository
		wantErr    bool
	}{
		{
			name:       "Valid repositories",
			statusCode: http.StatusOK,
			response:   `[{"name": "testrepo", "owner": {"login": "testuser"}, "description": "Test repository", "private": true}]`,
			want: []types.Repository{
				{
					Owner:       "testuser",
					Name:        "testrepo",
					Description: "Test repository",
					Private:     true,
				},
			},
			wantErr: false,
		},
		{
			name:       "Invalid response",
			statusCode: http.StatusUnauthorized,
			response:   `{"message": "Bad credentials"}`,
			want:       nil,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, client := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			})
			defer server.Close()

			got, err := client.ListRepositories()
			if (err != nil) != tt.wantErr {
				t.Errorf("ListRepositories() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(got) != len(tt.want) {
					t.Errorf("ListRepositories() = %v, want %v", got, tt.want)
				}
				for i := range got {
					if got[i].Name != tt.want[i].Name ||
						got[i].Owner != tt.want[i].Owner ||
						got[i].Description != tt.want[i].Description ||
						got[i].Private != tt.want[i].Private {
						t.Errorf("ListRepositories() = %v, want %v", got[i], tt.want[i])
					}
				}
			}
		})
	}
}

func TestCreateRepository(t *testing.T) {
	tests := []struct {
		name        string
		statusCode  int
		repoName    string
		description string
		private     bool
		wantErr     bool
	}{
		{
			name:        "Create private repository",
			statusCode:  http.StatusCreated,
			repoName:    "testrepo",
			description: "Test repository",
			private:     true,
			wantErr:     false,
		},
		{
			name:        "Failed to create repository",
			statusCode:  http.StatusUnauthorized,
			repoName:    "testrepo",
			description: "Test repository",
			private:     true,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, client := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				var repo github.Repository
				if err := decodeRequest(r, &repo); err != nil {
					t.Error(err)
				}

				if *repo.Name != tt.repoName ||
					*repo.Description != tt.description ||
					*repo.Private != tt.private {
					t.Error("Incorrect repository creation payload")
				}

				w.WriteHeader(tt.statusCode)
			})
			defer server.Close()

			err := client.CreateRepository(tt.repoName, tt.description, tt.private)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateRepository() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
