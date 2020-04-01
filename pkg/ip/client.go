package ip

import (
	"cloudflare-ddns/pkg/resilience"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	V4 Version = "https://api.ipify.org"
	V6 Version = "https://api6.ipify.org"
)

type (
	// Version of the IP.
	Version string

	// API is an HTTP client that can invoke ipify's API.
	API struct {
		client *http.Client
	}
)

// Client sets the HTTP client used for sending the request on the wire.
func Client(client *http.Client) func(*API) {
	return func(a *API) {
		a.client = client
	}
}

// Timeout sets the timeout of HTTP requests to ipify.
func Timeout(duration time.Duration) func(*API) {
	return func(a *API) {
		a.client.Timeout = duration
	}
}

// Retry sets the amount of attempts when trying to access ipify's API.
func Retry(attempts int) func(*API) {
	return func(a *API) {
		a.client.Transport = &resilience.Retry{Attempts: attempts, NextRoundTrip: http.DefaultTransport}
	}
}

// NewClients returns an HTTP client tha can access ipify's API.
func NewClient(options ...func(*API)) *API {
	c := &API{
		client: &http.Client{},
	}

	for _, o := range options {
		o(c)
	}

	return c
}

// Get returns the queried IP version, or an error if there are issues getting it.
func (c *API) Get(ctx context.Context, version Version) (ip string, err error) {
	// TODO Verify the requested version and returned response
	req, err := http.NewRequest("GET", string(version), nil)
	if err != nil {
		return "", fmt.Errorf("could not construct request: %w", err)
	}
	req = req.WithContext(ctx)
	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("could not get response: %w", err)
	}
	defer func() {
		cErr := resp.Body.Close()
		if err == nil {
			err = cErr
		}
	}()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("could not get process response, got %d status code", resp.StatusCode)
	}

	ipRaw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("could not read response: %w", err)
	}
	return string(ipRaw), err
}
