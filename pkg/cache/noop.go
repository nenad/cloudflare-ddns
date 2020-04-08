package cache

import (
	"cloudflare-ddns/pkg/cloudflare"
)

type NoopCache struct{}

func (c *NoopCache) GetRecord(domain, recordType string) (cloudflare.Record, error) {
	return cloudflare.Record{}, nil
}
func (c *NoopCache) SaveRecord(record cloudflare.Record) error {
	return nil
}
