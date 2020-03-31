package ip_test

import (
	"context"
	"cloudflare-ddns/pkg/ip"
	"io"
	"net/http"
	"testing"
)

type (
	body          string
	testTransport func(r *http.Request) *http.Response
)

func (b *body) Read(p []byte) (n int, err error) {
	return copy(p, *b), io.EOF
}

func (b *body) Close() error {
	return nil
}

func (t testTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	return t(r), nil
}

func newTestClient(transport testTransport) *http.Client {
	return &http.Client{
		Transport: transport,
	}
}

func Test_VersionsGet(t *testing.T) {
	tests := []struct {
		name     string
		version  ip.Version
		want     string
		wantURL  string
		hasError bool
	}{
		{
			name:     "ipv4 should return ipv4 address",
			version:  ip.V4,
			want:     "192.168.0.1",
			wantURL:  "https://api.ipify.org",
			hasError: false,
		},
		{
			name:     "ipv6 should return ipv6 address",
			version:  ip.V6,
			want:     "fe80::bf51:5d53:8f20:9d4",
			wantURL:  "https://api6.ipify.org",
			hasError: false,
		},
		{
			name:     "ipv4 should fail and error is expected",
			version:  ip.V4,
			want:     "",
			wantURL:  "https://api.ipify.org",
			hasError: true,
		},
		{
			name:     "ipv6 should fail and error is expected",
			version:  ip.V6,
			want:     "",
			wantURL:  "https://api6.ipify.org",
			hasError: true,
		},
		{
			name:     "non-standard IP should fail and error is expected",
			version:  "badversion",
			want:     "",
			wantURL:  "badversion",
			hasError: true,
		},
	}

	t.Parallel()
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			client := ip.NewClient(ip.Client(newTestClient(func(r *http.Request) *http.Response {
				if r.URL.String() != tt.wantURL {
					t.Errorf("wrong URL, want %q, got %q", tt.wantURL, r.URL.String())
				}
				var want = body(tt.want)
				code := 200
				if tt.hasError {
					code = 500
				}
				return &http.Response{
					StatusCode: code,
					Header:     map[string][]string{"Content-Type": {"text/plain"}},
					Body:       &want,
				}
			})))

			got, err := client.Get(context.Background(), tt.version)
			if tt.hasError {
				if err == nil {
					t.Errorf("expected error, did not receive one")
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %s", err)
				}
			}

			if got != tt.want {
				t.Errorf("want %s, got %s", tt.want, got)
			}
		})
	}
}
