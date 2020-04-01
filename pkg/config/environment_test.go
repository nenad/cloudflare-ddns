package config

import (
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestFromEnvironment(t *testing.T) {
	tests := []struct {
		name        string
		env         map[string]string
		want        CloudFlare
		errKeywords []string
	}{
		{
			name: "minimal configuration should not throw an error and have expected defaults",
			env: map[string]string{
				"DOMAIN": "nenad.dev",
				"TOKEN":  "token",
				"TYPE":   "A",
			},
			want: CloudFlare{
				Domain:  "nenad.dev",
				Token:   "token",
				Type:    "A",
				Timeout: time.Second * time.Duration(10),
				Proxied: true,
				TTL:     1,
			},
		},
		{
			name: "full configuration should not throw an error",
			env: map[string]string{
				"DOMAIN":  "nenad.dev",
				"TOKEN":   "token",
				"TYPE":    "AAAA",
				"TIMEOUT": "200",
				"PROXIED": "false",
				"TTL":     "300",
			},
			want: CloudFlare{
				Domain:  "nenad.dev",
				Token:   "token",
				Type:    "AAAA",
				Timeout: time.Second * time.Duration(200),
				Proxied: false,
				TTL:     300,
			},
		},
		{
			name: "type is only A or AAAA",
			env: map[string]string{
				"DOMAIN": "nenad.dev",
				"TOKEN":  "token",
				"TYPE":   "B",
			},
			want:        CloudFlare{},
			errKeywords: []string{"TYPE", "AAAA"},
		},
		{
			name:        "empty environment should fail for DOMAIN, TOKEN and TYPE",
			env:         map[string]string{},
			want:        CloudFlare{},
			errKeywords: []string{"DOMAIN", "TOKEN", "TYPE"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()
			for k, v := range tt.env {
				_ = os.Setenv(k, v)
			}

			got, err := FromEnvironment()

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
