package main

import (
	"cloudflare-ddns/pkg/cloudflare"
	"cloudflare-ddns/pkg/config"
	"cloudflare-ddns/pkg/resolver/external"
	"context"
	"fmt"
	"os"
)

func main() {
	cfg, err := config.Parse(os.Args[1:])
	if err != nil {
		fail(err)
	}

	// TODO Abstract IP version in a separate package
	// TODO Make ipify URL customizable
	ipifyVersion := external.V4
	if cfg.CloudFlare.Type == "AAAA" {
		ipifyVersion = external.V6
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cf, err := cloudflare.NewClient(cfg.CloudFlare.Token, cloudflare.Timeout(cfg.CloudFlare.Timeout), cloudflare.Retry(3))
	if err != nil {
		fail(err)
	}

	rec, err := cf.GetRecord(ctx, cfg.CloudFlare.Domain, cloudflare.Type(cfg.CloudFlare.Type))
	if err != nil {
		fail(err)
	}

	myIP, err := (external.NewClient(external.Retry(3))).Get(ctx, ipifyVersion)
	if err != nil {
		fail(err)
	}

	if myIP != rec.Content {
		err := cf.UpdateRecord(ctx, rec.ID, cloudflare.DNSUpdateRequest{
			Name:    rec.Name,
			Type:    rec.Type,
			Content: myIP,
			Proxied: cfg.CloudFlare.Proxied,
		})
		if err != nil {
			fail("could not update record, got " + err.Error())
		}
		fmt.Printf("Updated %q to point from %s to %s\n", cfg.CloudFlare.Domain, rec.Content, myIP)
	} else {
		fmt.Println("No updates")
	}
}

func fail(arg interface{}) {
	fmt.Println(arg)
	os.Exit(1)
}
