package main

import (
	"cloudflare-ddns/pkg/cache"
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

	retriever := ip.Factory(cfg.App.Interface)
	myIP, err := retriever.Get(cfg.CloudFlare.IPVersion)
	if err != nil {
		log.Fatalf("could not get IP: %s", err)
	}

	cacher := cache.Factory(cfg.App.CacheEnabled)
	cached, err := cacher.GetRecord(cfg.CloudFlare.Domain, cfg.CloudFlare.Type)
	if err != nil {
		fmt.Printf("error while getting cache: %s\n", err)
	}

	if myIP == cached.Content {
		fmt.Println("no changes in IP, skipping update")
		return
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
	if err := cacher.SaveRecord(rec); err != nil {
		fmt.Printf("could not save cached record: %s\n", err)
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
