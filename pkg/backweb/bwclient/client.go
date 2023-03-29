package bwclient

import (
	"net/http"
	"net/url"
	"time"
)

type client struct {
	httpClient *http.Client
	url        string
}

// NewClient creates Client, and returns the pointer to it.
func NewClient(backURL string, timeout time.Duration) (*client, error) {

	//Trying to parse url
	_, err := url.Parse(backURL)
	if err != nil {
		return nil, err
	}

	// Fix timeout in case of negative value
	if timeout < 0 {
		timeout = 0
	}

	return &client{
		httpClient: &http.Client{
			Timeout: timeout,
		},
		url: backURL,
	}, nil
}
