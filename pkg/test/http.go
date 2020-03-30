package test

import (
	"io"
	"net/http"
)

type (
	Body struct {
		offset      int
		contentSize int
		content     []byte
	}
	testTransport func(r *http.Request) *http.Response
)

func FromBytes(input []byte) *Body {
	b := &Body{
		offset:      0,
		contentSize: len(input),
	}
	b.content = make([]byte, len(input))
	copy(b.content, input)
	return b
}

func (b *Body) Read(p []byte) (n int, err error) {
	n = copy(p, b.content[b.offset:])
	b.offset += n
	if b.offset >= b.contentSize {
		return n, io.EOF
	}

	return n, nil
}

func (b *Body) Close() error {
	b.offset = 0
	b.content = nil
	b.contentSize = 0
	return nil
}

func (t testTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	return t(r), nil
}

func NewTestClient(transport testTransport) *http.Client {
	return &http.Client{
		Transport: transport,
	}
}
