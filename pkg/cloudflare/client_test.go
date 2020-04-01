package cloudflare_test

import (
	"cloudflare-ddns/pkg/cloudflare"
	"cloudflare-ddns/pkg/test"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func jsonFixture(t *testing.T, file string) (body *test.Body) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		t.Fatalf("could not load fixture %q: %s", file, err)
	}

	b := test.FromBytes(data)
	return b
}

func Test_ClientGetRecord(t *testing.T) {
	tests := []struct {
		name     string
		fixture  string
		want     cloudflare.Record
		wantURL  string
		zoneName string
		err      string
		recName  string
		recType  cloudflare.Type
	}{
		{
			name:    "test single page get domain record",
			fixture: "testdata/single_page.json",
			wantURL: "https://api.cloudflare.com/client/v4/zones/zone12345/dns_records?page=1&per_page=50",
			want: cloudflare.Record{
				ID:      "a3bfb90b2b9d",
				Type:    "A",
				Name:    "nenad.dev",
				Content: "192.168.0.2",
			},
			recName:  "nenad.dev",
			recType:  cloudflare.A,
			zoneName: "nenad.dev",
		},
		{
			name:    "test single page get subdomain record",
			fixture: "testdata/single_page.json",
			wantURL: "https://api.cloudflare.com/client/v4/zones/zone12345/dns_records?page=1&per_page=50",
			want: cloudflare.Record{
				ID:      "6e04a018b79b9302244a56d",
				Type:    "AAAA",
				Name:    "home.nenad.dev",
				Content: "2a03:8103:98c1:7b0b:188c:243f:c6fc:86d2",
			},
			recName:  "home.nenad.dev",
			recType:  cloudflare.AAAA,
			zoneName: "nenad.dev",
		},
		{
			name:     "duplicate record will return error",
			fixture:  "testdata/duplicate_entry.json",
			wantURL:  "https://api.cloudflare.com/client/v4/zones/zone12345/dns_records?page=1&per_page=50",
			want:     cloudflare.Record{},
			recName:  "home.nenad.dev",
			recType:  cloudflare.A,
			err:      "duplicate",
			zoneName: "nenad.dev",
		},
		{
			name:     "single page no record found should return error",
			fixture:  "testdata/single_page.json",
			wantURL:  "https://api.cloudflare.com/client/v4/zones/zone12345/dns_records?page=1&per_page=50",
			want:     cloudflare.Record{},
			recName:  "not-found.hello.dev",
			recType:  cloudflare.A,
			err:      "no record",
			zoneName: "hello.dev",
		},
	}

	t.Parallel()
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			client, _ := cloudflare.NewClient("token", cloudflare.Client(test.NewTestClient(func(r *http.Request) *http.Response {
				if r.URL.String() == fmt.Sprintf("https://api.cloudflare.com/client/v4/zones?name=%s", tt.zoneName) {
					return &http.Response{
						StatusCode: 200,
						Header: http.Header{
							"Content-Type": {"application/json"},
						},
						Body: jsonFixture(t, "testdata/zone.json"),
					}
				}

				if r.URL.String() != tt.wantURL {
					t.Errorf("URL mismatch, want %q, got %q", tt.wantURL, r.URL.String())
				}
				return &http.Response{
					StatusCode: 200,
					Header: http.Header{
						"Content-Type": {"application/json"},
					},
					Body: jsonFixture(t, tt.fixture),
				}
			})))
			got, err := client.GetRecord(context.Background(), tt.recName, tt.recType)
			if tt.err != "" {
				if err == nil {
					t.Fatalf("expected an error")
				}

				if !strings.Contains(err.Error(), tt.err) {
					t.Fatalf("expected error to contain %q, got %q", tt.err, err)
				}
			} else {
				if err != nil {
					t.Fatalf("did not expect an error: %s", err)
				}
			}

			if !reflect.DeepEqual(tt.want, got) {
				t.Fatalf("result mismatch; want %#v, got %#v", tt.want, got)
			}
		})
	}

}
