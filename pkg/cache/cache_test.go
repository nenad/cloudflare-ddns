package cache_test

import (
	"cloudflare-ddns/pkg/cache"
	"cloudflare-ddns/pkg/cloudflare"
	"os"
	"path"
	"reflect"
	"testing"
)

func TestCache_SaveAndGetRecord(t *testing.T) {
	c := cache.Cache{}
	want := cloudflare.Record{
		ID:      "f0764281-2853-4e5c-842d-bea41fcbccdf",
		Type:    "A",
		Name:    "test.nenad.dev",
		Content: "192.168.0.1",
	}

	if err := c.SaveRecord(want); err != nil {
		t.Fatalf("could not save record: %s", err)
	}

	got, err := c.GetRecord("test.nenad.dev", "A")
	if err != nil {
		t.Fatalf("could not get record: %s", err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("records don't match. want %#v, got %#v", want, got)
	}

	cacheDir, _ := os.UserCacheDir()

	if err := os.Remove(path.Join(cacheDir, "cloudflare-ddns", "test.nenad.dev-A.json")); err != nil {
		t.Fatalf("error while removing file: %s", err)
	}
}
