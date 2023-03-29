package aclient

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"gitlab.assistagro.com/back/back.auth.go/pkg/model"
)

type client struct {
	httpClient *http.Client
	url        string
}

// NewClient creates Client, and returns the pointer to it.
//
//	backURL - url of the backend service
//	timeout - timeout for http requests, if timeout <= 0, then no timeout
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

// New creates new session client by token, and returns the pointer to it.
//
//	token - user token
func (c *client) NewSession(token string) (*Session, error) {

	if (*c == client{}) {
		return nil, fmt.Errorf("client is not initialized")
	}

	if len(strings.TrimSpace(token)) == 0 {
		return nil, fmt.Errorf("token is empty")
	}

	return &Session{
		client: c,
		token:  token,
	}, nil
}

// Returns user by user token
//
// token - user token
func (c *client) GetUser(token string) (model.User, error) {
	ses, err := c.NewSession(token)
	if err != nil {
		return model.User{}, err
	}
	return ses.GetUser()
}
