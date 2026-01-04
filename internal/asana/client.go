package asana

import (
	"net/http"
	"strings"
	"time"

	"github.com/michalvavra/asncli/internal/auth"
)

const defaultBaseURL = "https://app.asana.com/api/1.0"

type Client struct {
	baseURL     string
	httpClient  *http.Client
	tokenSource auth.TokenSource
}

func NewClient(tokenSource auth.TokenSource, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 15 * time.Second}
	}
	return &Client{
		baseURL:     defaultBaseURL,
		httpClient:  httpClient,
		tokenSource: tokenSource,
	}
}

func (c *Client) WithBaseURL(baseURL string) *Client {
	clone := *c
	clone.baseURL = strings.TrimRight(baseURL, "/")
	return &clone
}
