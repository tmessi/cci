// Package circleci provides a Client for making requests to CircleCI.
package circleci

import (
	"context"
	"fmt"
	"net/http"
)

// Project represents a CircleCI project.
type Project struct {
	Name         string
	Organization string
	VCSType      string
}

// Client is used to make HTTP requests to CircleCI.
type Client struct {
	client *http.Client

	rootURL string
	project *Project
	token   string
}

// New creates a Client.
func New(client *http.Client, rootURL string, project *Project, token string) *Client {
	return &Client{
		client:  client,
		rootURL: rootURL,
		project: project,
		token:   token,
	}
}

func (c *Client) baseURL() string {
	return fmt.Sprintf(
		"%s/api/v1.1/project/%s/%s/%s",
		c.rootURL,
		c.project.VCSType,
		c.project.Organization,
		c.project.Name,
	)
}

// do sends an HTTP request and returns an HTTP response.
// It is a thin wrapper around http.Client.Do that will
// set the auth using the Client's token, and set headers.
func (c *Client) do(ctx context.Context, req *http.Request) (*http.Response, error) {
	req = req.WithContext(ctx)
	req.SetBasicAuth(c.token, "")
	req.Header.Set("Accept", "application/json")

	return c.client.Do(req)
}
