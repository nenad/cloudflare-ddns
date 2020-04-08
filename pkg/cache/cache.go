package cache

import (
	"cloudflare-ddns/pkg/cloudflare"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

type Cacher interface {
	GetRecord(domain, recordType string) (cloudflare.Record, error)
	SaveRecord(record cloudflare.Record) error
}

type Cache struct{}

func Factory(enabled bool) Cacher {
	if !enabled {
		return &NoopCache{}
	}

	return &Cache{}
}

func (c *Cache) GetRecord(domain, recordType string) (rec cloudflare.Record, err error) {
	filename, err := getFilename(domain, recordType)
	if err != nil {
		return rec, fmt.Errorf("could not get filename: %w", err)
	}
	cacheBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return rec, fmt.Errorf("could not read cache file: %w", err)
	}

	if err := json.Unmarshal(cacheBytes, &rec); err != nil {
		return rec, fmt.Errorf("could not unmarshal cached record: %w", err)
	}

	return rec, nil
}

func (c *Cache) SaveRecord(rec cloudflare.Record) error {
	filename, err := getFilename(rec.Name, string(rec.Type))
	if err != nil {
		return fmt.Errorf("could not get filename: %w", err)
	}

	cacheFile, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("could not open file for writing: %w", err)
	}

	if err := json.NewEncoder(cacheFile).Encode(rec); err != nil {
		return fmt.Errorf("could not marshal record to file: %w", err)
	}

	return nil
}

func getFilename(domain, recordType string) (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("could not get home directory: %w", err)
	}

	dir := "cloudflare-ddns"
	if err := os.MkdirAll(path.Join(cacheDir, dir), os.ModePerm); err != nil {
		return "", fmt.Errorf("could not create config directory: %w", err)
	}
	filename := fmt.Sprintf("%s-%s.json", domain, recordType)

	return path.Join(cacheDir, dir, filename), nil
}
