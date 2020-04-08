package main

import (
	"cloudflare-ddns/pkg/cloudflare"
	"cloudflare-ddns/pkg/config"
	"cloudflare-ddns/pkg/ip"
	"context"
	"fmt"
	"log"
	"os"
)

func main() {
	cfg, err := config.Parse(os.Args[1:])
	if err != nil {
		log.Fatalf("could not parse flags: %s", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resolver := ip.Factory(cfg.App.Interface)
	myIP, err := resolver.Get(cfg.CloudFlare.IPVersion)
	if err != nil {
		log.Fatalf("could not get IP: %s", err)
	}

	cf, err := cloudflare.NewClient(
		cfg.CloudFlare.Token,
		cloudflare.Timeout(cfg.CloudFlare.Timeout),
		cloudflare.Retry(3),
	)
	if err != nil {
		log.Fatalf("could not initialize CloudFlare client: %s", err)
	}

	rec, err := cf.GetRecord(ctx, cfg.CloudFlare.Domain, cloudflare.Type(cfg.CloudFlare.Type))
	if err != nil {
		log.Fatalf("could not get CloudFlare record: %s", err)
	}
	if err := cf.UpdateRecord(ctx, rec.ID, cloudflare.DNSUpdateRequest{
		Name:    rec.Name,
		Type:    rec.Type,
		Content: myIP,
		Proxied: cfg.CloudFlare.Proxied,
	}); err != nil {
		log.Fatalf("could not update record: %s ", err)
	}

	fmt.Printf("Updated %q to point from %s to %s\n", cfg.CloudFlare.Domain, rec.Content, myIP)
}
