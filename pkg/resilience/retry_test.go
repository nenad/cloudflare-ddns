package resilience

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

var server string

func TestMain(m *testing.M) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))

	server = testServer.URL
	os.Exit(m.Run())
}

type testRoundTrip struct {
	fails       int
	failCounter int
}

func (t *testRoundTrip) RoundTrip(r *http.Request) (*http.Response, error) {
	if t.failCounter < t.fails {
		t.failCounter++
		return nil, fmt.Errorf("failure when connecting")
	}

	return &http.Response{StatusCode: 200}, nil
}

func newRoundTrip(fails int) http.RoundTripper {
	return &testRoundTrip{fails: fails}
}

func TestRetry_RoundTripFailsAfterAttemptsExhausted(t *testing.T) {
	tests := []struct {
		name        string
		serverFails int
		attempts    int
		shouldError bool
		statusCode  int
	}{
		{
			name:        "should error after two attempts and possible three failures",
			serverFails: 3,
			attempts:    2,
			shouldError: true,
		},
		{
			name:        "should error after three attempts and three failures",
			serverFails: 3,
			attempts:    3,
			shouldError: true,
		},
		{
			name:        "should get successful response after 3 attempts and 2 server failures",
			serverFails: 2,
			attempts:    3,
			shouldError: false,
			statusCode:  200,
		},
	}

	t.Parallel()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			retry := &Retry{
				NextRoundTrip: newRoundTrip(tt.serverFails),
				Attempts:      tt.attempts,
				Wait:          time.Millisecond * 1,
			}

			client := http.Client{
				Transport: retry,
			}

			resp, err := client.Get(server)
			if tt.shouldError {
				if err == nil {
					t.Errorf("expected error, got %#v", resp)
				}
			} else {
				if err != nil {
					t.Errorf("did not expect error, got %s", err)
				}
				if tt.statusCode != resp.StatusCode {
					t.Errorf("did not get 200 status code, got %d", resp.StatusCode)
				}
			}
		})
	}

}
