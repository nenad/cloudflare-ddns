package cloudflare

import (
	"bytes"
	"context"
	"ddns-cloudflare/pkg/resilience"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"time"
)

const (
	baseAPI = "https://api.cloudflare.com/client/v4"
	perPage = 50
)

// API is an HTTP client that can invoke CloudFlare's API functions.
type API struct {
	client *http.Client
	zone   string
	token  string
}

// Retry sets the attempt number when making API calls to CloudFlare.
func Retry(attempts int) func(*API) {
	return func(a *API) {
		a.client.Transport = &resilience.Retry{
			NextRoundTrip: a.client.Transport,
			Wait:          time.Second,
			Attempts:      attempts,
		}
	}
}

// Client sets the HTTP client used for making requests to CloudFlare.
func Client(client *http.Client) func(*API) {
	return func(a *API) {
		a.client = client
	}
}

// Timeout sets the HTTP request timeouts when making API calls to CloudFlare.
func Timeout(duration time.Duration) func(*API) {
	return func(a *API) {
		a.client.Timeout = duration
	}
}

// NewClient returns an HTTP client that can invoke CloudFlare's API
func NewClient(zone, token string, options ...func(*API)) (*API, error) {
	if zone == "" {
		return nil, fmt.Errorf("zone is empty")
	}
	if token == "" {
		return nil, fmt.Errorf("token is empty")
	}

	a := &API{
		client: &http.Client{
			Timeout:   time.Second * 10,
			Transport: http.DefaultTransport,
		},
		zone:  zone,
		token: token,
	}
	for _, o := range options {
		o(a)
	}

	return a, nil
}

// GetRecord will return the DNS record matching the given record and type.
// Returns an error if there are either no records or duplicate records found.
func (a *API) GetRecord(ctx context.Context, name string, recordType Type) (rec Record, err error) {
	page := 1

	for {
		dnsResp := DNSResponse{}
		err := a.send(ctx, "GET", api("/zones/%s/dns_records?page=%d&per_page=%d", a.zone, page, perPage), nil, &dnsResp)
		if err != nil {
			return rec, err
		}

		for _, r := range *dnsResp.Result {
			if r.Name == name && r.Type == recordType {
				if rec != (Record{}) {
					return Record{}, fmt.Errorf("found duplicate entry for %q and type %q", r.Name, r.Type)
				}

				rec = r
			}
		}

		if page >= dnsResp.ResultInfo.TotalPages {
			if rec == (Record{}) {
				return rec, fmt.Errorf("no record for %q of type %s found", name, recordType)
			}

			return rec, err
		}
		page++
	}
}

// UpdateRecord will update the record with the given value.
// Returns an error if there are any CloudFlare errors returned.
func (a *API) UpdateRecord(ctx context.Context, id string, request DNSUpdateRequest) error {
	if request.TTL == 0 {
		request.TTL = 1
	}

	return a.send(ctx, "PUT", api("/zones/%s/dns_records/%s", a.zone, id), &request, nil)
}

func (a *API) send(ctx context.Context, method, url string, send, recv interface{}) error {
	buf := &bytes.Buffer{}
	if send != nil {
		if err := json.NewEncoder(buf).Encode(send); err != nil {
			return fmt.Errorf("could not encode structure: %w", err)
		}
	}

	req, err := a.request(ctx, method, url, buf)
	if err != nil {
		return fmt.Errorf("could not create request: %w", err)
	}

	resp, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("could not get response: %w", err)
	}

	if err := json.NewDecoder(resp.Body).Decode(&recv); err != nil {
		return fmt.Errorf("could not decode response: %w", err)
	}

	if err := checkError(recv); err != nil {
		return fmt.Errorf("error from CloudFlare: %w", err)
	}

	return nil
}

func (a *API) request(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.token))
	req.Header.Add("Content-Type", "application/json")
	return req, err
}

func api(format string, args ...interface{}) string {
	return fmt.Sprintf("%s%s", baseAPI, fmt.Sprintf(format, args...))
}

func checkError(resp interface{}) error {
	v := reflect.ValueOf(resp)
	if v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return fmt.Errorf("received non-struct item: %s", v.Kind())
	}

	// If we received Success, no need to check for further errors
	successField := v.FieldByName("Success").Interface()
	if val, ok := successField.(bool); ok {
		if val {
			return nil
		}
	}

	field := v.FieldByName("Errors").Interface()

	errs, ok := field.([]Error)
	if !ok {
		return fmt.Errorf("field Errors is not of type []Error")
	}

	if len(errs) == 0 {
		return nil
	}

	var err string
	for _, e := range errs {
		err = fmt.Sprintf("%s: [%d] %s", err, e.Code, e.Message)
		for _, er := range e.ErrorChain {
			err = fmt.Sprintf("%s: [%d] %s", err, e.Code, er.Message)
		}
	}

	return fmt.Errorf("error in the response%s", err)
}
