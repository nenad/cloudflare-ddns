package config

import (
	"flag"
	"fmt"
	"os"
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

// Parse generates configuration from the command arguments.
func Parse(args []string) (CloudFlare, error) {
	fs := flag.NewFlagSet("cf", flag.ExitOnError)
	domain := ""
	token := ""
	recordType := "A"
	timeout := 10
	ttl := 1
	proxied := true

	fs.Usage = func() {
		_, _ = fmt.Fprintf(fs.Output(), "USAGE:\n\t%s -token xxx -domain example.com\n\nCONFIGURATION:\n", os.Args[0])
		fs.PrintDefaults()
	}

	fs.StringVar(&token, "token", "", "A CloudFlare token with Zone.Zone (Read), Zone.DNS (Edit) permissions (Required)")
	fs.StringVar(&domain, "domain", "", "The domain you would like to update (Required)")
	fs.StringVar(&recordType, "type", "A", "The record type you would like to update, must be A or AAAA")
	fs.IntVar(&timeout, "timeout", 10, "API request timeout to CloudFlare and external IP service")
	fs.IntVar(&ttl, "ttl", 1, "TTL for the domain record")
	fs.BoolVar(&proxied, "proxied", true, "Is the request proxied through CloudFlare's servers")

	if err := fs.Parse(args); err != nil {
		return CloudFlare{}, fmt.Errorf("could not parse command parameters: %w", err)
	}

	var errs []string
	if token == "" {
		errs = append(errs, "-token is required and must not be empty")
	}

	if domain == "" {
		errs = append(errs, "-domain is required and must not be empty")
	}

	if recordType != "A" && recordType != "AAAA" {
		errs = append(errs, "-type must be 'A' for IPv4 or 'AAAA' for IPv6")
	}

	if timeout <= 0 {
		timeout = 10
	}

	if ttl <= 0 {
		ttl = 1
	}

	if len(errs) > 0 {
		return CloudFlare{}, fmt.Errorf(strings.Join(errs, "; "))
	}

	return CloudFlare{
		Domain:  domain,
		Token:   token,
		Type:    recordType,
		Timeout: time.Second * time.Duration(timeout),
		Proxied: proxied,
		TTL:     ttl,
	}, nil
}