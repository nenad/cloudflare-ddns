package config

import (
	"cloudflare-ddns/pkg/ip"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestFromEnvironment(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		want        Configuration
		errKeywords []string
	}{
		{
			name: "minimal configuration should not throw an error and have expected defaults",
			args: []string{
				"-domain", "nenad.dev",
				"-token", "token",
				"-type", "A",
			},
			want: Configuration{
				CloudFlare: CloudFlare{
					Domain:    "nenad.dev",
					Token:     "token",
					Type:      "A",
					Timeout:   time.Second * time.Duration(10),
					Proxied:   true,
					TTL:       1,
					IPVersion: ip.V4,
				},
			},
		},
		{
			name: "full configuration should not throw an error",
			args: []string{
				"-domain", "nenad.dev",
				"-token", "token",
				"-type", "AAAA",
				"-timeout", "200",
				"-proxied",
				"-ttl", "300",
				"-interface", "wlp3s0",
			},
			want: Configuration{
				CloudFlare: CloudFlare{
					Domain:    "nenad.dev",
					Token:     "token",
					Type:      "AAAA",
					Timeout:   time.Second * time.Duration(200),
					Proxied:   true,
					TTL:       300,
					IPVersion: ip.V6,
				},
				App: App{
					Interface: "wlp3s0",
				},
			},
		},
		{
			name: "type is only A or AAAA",
			args: []string{
				"-domain", "nenad.dev",
				"-token", "token",
				"-type", "B",
			},
			want:        Configuration{},
			errKeywords: []string{"-type", "AAAA"},
		},
		{
			name:        "empty command line should fail for domain and token",
			args:        []string{},
			want:        Configuration{},
			errKeywords: []string{"-domain", "-token"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args)

			if err != nil && len(tt.errKeywords) == 0 {
				t.Fatalf("got error when none was expected: %s", err)
			} else if err != nil {
				for _, kw := range tt.errKeywords {
					if !strings.Contains(err.Error(), kw) {
						t.Fatalf("expected error to contain keyword %q, got %q", kw, err)
					}
				}
			} else if len(tt.errKeywords) > 0 {
				t.Fatalf("no error expected, got %q", tt.errKeywords)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromEnvironment() got = %#v, want %#v", got, tt.want)
			}
		})
	}
}
