package test

import (
	"io"
	"net/http"
)

type (
	// Body is in-memory structure that implements io.ReadCloser.
	Body struct {
		offset      int
		contentSize int
		content     []byte
	}
	// Transport allows simulating request-response in any http.Client.
	Transport func(r *http.Request) *http.Response
)

// Creates a Body from given input.
func FromBytes(input []byte) *Body {
	b := &Body{
		offset:      0,
		contentSize: len(input),
	}
	b.content = make([]byte, len(input))
	copy(b.content, input)
	return b
}

// Reads bytes in p.
func (b *Body) Read(p []byte) (n int, err error) {
	n = copy(p, b.content[b.offset:])
	b.offset += n
	if b.offset >= b.contentSize {
		return n, io.EOF
	}

	return n, nil
}

// Closes the body.
func (b *Body) Close() error {
	b.offset = 0
	b.content = nil
	b.contentSize = 0
	return nil
}

// RoundTrip implements RoundTripper.
func (t Transport) RoundTrip(r *http.Request) (*http.Response, error) {
	return t(r), nil
}

// NewTestClient returns a http.Client that has test transport injected.
func NewTestClient(transport Transport) *http.Client {
	return &http.Client{
		Transport: transport,
	}
}
