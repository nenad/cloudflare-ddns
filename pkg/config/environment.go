package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type (
	// CloudFlare domain settings
	CloudFlare struct {
		Domain  string
		Token   string
		Type    string
		Timeout time.Duration
		Proxied bool
		TTL     int
	}
)

// FromEnvironment generates configuration from the running environment.
func FromEnvironment() (CloudFlare, error) {
	var errs []string
	token := os.Getenv("TOKEN")
	if token == "" {
		errs = append(errs, "environment variable TOKEN is not set")
	}

	domain := os.Getenv("DOMAIN")
	if domain == "" {
		errs = append(errs, "environment variable DOMAIN is not set")
	}

	domainType := os.Getenv("TYPE")
	if domainType != "A" && domainType != "AAAA" {
		errs = append(errs, "environment variable TYPE must be 'A' for IPv4 or 'AAAA' for IPv6")
	}

	secTimeout, _ := strconv.Atoi(os.Getenv("TIMEOUT"))
	if secTimeout <= 0 {
		secTimeout = 10
	}

	ttl, _ := strconv.Atoi(os.Getenv("TTL"))
	if ttl <= 0 {
		ttl = 1
	}

	proxied, err := strconv.ParseBool(os.Getenv("PROXIED"))
	if err != nil {
		proxied = true
	}

	if len(errs) > 0 {
		return CloudFlare{}, fmt.Errorf(strings.Join(errs, ";"))
	}

	return CloudFlare{
		Domain:  domain,
		Token:   token,
		Type:    domainType,
		Timeout: time.Second * time.Duration(secTimeout),
		Proxied: proxied,
		TTL:     ttl,
	}, nil
}
