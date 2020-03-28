package resilience

import (
	"fmt"
	"net/http"
	"time"
)

// Retry retries the round trip until no-500 response is received, or the attempts are exhausted.
type Retry struct {
	NextRoundTrip http.RoundTripper
	Wait          time.Duration
	Attempts      int
}

func (r *Retry) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	if r.Attempts <= 0 {
		r.Attempts = 3
	}
	if r.Wait == 0 {
		r.Wait = time.Second
	}

	for i := 0; i < r.Attempts; i++ {
		resp, err = r.NextRoundTrip.RoundTrip(req)

		if err == nil && resp.StatusCode < 500 {
			return resp, err
		}

		select {
		case <-req.Context().Done():
			return resp, req.Context().Err()
		case <-time.After(r.Wait):
		}
	}
	return resp, fmt.Errorf("could not get a response after %d attempts: %w", r.Attempts, err)
}
