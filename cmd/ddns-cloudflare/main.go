package main

import (
	"cloudflare-ddns/pkg/cloudflare"
	"cloudflare-ddns/pkg/ip"
	"context"
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {
	token := os.Getenv("TOKEN")
	if token == "" {
		fail("Environment variable TOKEN is not set, exiting.")
	}

	domain := os.Getenv("DOMAIN")
	if domain == "" {
		fail("Environment variable DOMAIN is not set, exiting.")
	}

	version := os.Getenv("VERSION")
	if version != "A" && version != "AAAA" {
		fail("Environment variable VERSION must be 'A' for IPv4 or 'AAAA' for IPv6")
	}
	// TODO Abstract IP version in a separate package
	// TODO Make ipify URL customizable
	ipifyVersion := ip.V4
	if version == "AAAA" {
		ipifyVersion = ip.V6
	}

	secTimeout, _ := strconv.Atoi(os.Getenv("TIMEOUT"))
	if secTimeout == 0 {
		secTimeout = 10
	}
	timeout := time.Second * time.Duration(secTimeout)
	proxied := os.Getenv("PROXIED") != "false"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cf, err := cloudflare.NewClient(token, cloudflare.Timeout(timeout), cloudflare.Retry(3))
	if err != nil {
		fail(err)
	}

	rec, err := cf.GetRecord(ctx, domain, cloudflare.Type(version))
	if err != nil {
		fail(err)
	}

	myIP, err := (ip.NewClient(ip.Retry(3))).Get(ctx, ipifyVersion)
	if err != nil {
		fail(err)
	}

	if myIP != rec.Content {
		err := cf.UpdateRecord(ctx, rec.ID, cloudflare.DNSUpdateRequest{
			Name:    rec.Name,
			Type:    rec.Type,
			Content: myIP,
			Proxied: proxied,
		})
		if err != nil {
			fail("could not update record, got " + err.Error())
		}
		fmt.Printf("Updated %q to point from %s to %s\n", domain, rec.Content, myIP)
	} else {
		fmt.Println("No updates")
	}
}

func fail(arg interface{}) {
	fmt.Println(arg)
	os.Exit(1)
}
