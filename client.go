package httpclient

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/textproto"
	"time"

	"golang.org/x/exp/slices"
)

const headerContentType = "Content-Type"

// Client is a thin wrapper around http.Client.
// It provides a convenient way to set default headers and middleware.
type Client struct {
	Header     http.Header
	Doer       Doer
	Middleware func(req *http.Request) (*http.Request, error)
}

// New returns a new Client with default settings.
func New() *Client {
	return NewFrom(newTransport())
}

// NewFrom returns a new Client with the given transport.
func NewFrom(doer Doer) *Client {
	return &Client{
		Doer:   doer,
		Header: make(http.Header),
	}
}

func newTransport() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
}

// Doer describes a HTTP request executor.
type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Get makes a GET request to the given address.
func (client *Client) Get(ctx context.Context, addr string) (*http.Response, error) {
	req, err := client.newRequest(ctx, http.MethodGet, addr, nil)
	if err != nil {
		return nil, err
	}
	return client.Doer.Do(req)
}

// Post makes a POST request to the given address.
func (client *Client) Post(ctx context.Context, addr, bodyType string, body io.Reader) (*http.Response, error) {
	return client.doBody(ctx, http.MethodPost, addr, bodyType, body)
}

// Put makes a PUT request to the given address.
func (client *Client) Put(ctx context.Context, addr, bodyType string, body io.Reader) (*http.Response, error) {
	return client.doBody(ctx, http.MethodPut, addr, bodyType, body)
}

// Patch makes a PATCH request to the given address.
func (client *Client) Patch(ctx context.Context, addr, bodyType string, body io.Reader) (*http.Response, error) {
	return client.doBody(ctx, http.MethodPatch, addr, bodyType, body)
}

const methodQuery = "QUERY"

// Query makes a QUERY request to the given address.
// QUERY methods // https://www.ietf.org/archive/id/draft-ietf-httpbis-safe-method-w-body-02.html#name-introduction.
func (client *Client) Query(ctx context.Context, addr, bodyType string, body io.Reader) (*http.Response, error) {
	return client.doBody(ctx, methodQuery, addr, bodyType, body)
}

func (client *Client) doBody(ctx context.Context, method, addr, bodyType string, body io.Reader) (*http.Response, error) {
	req, err := client.newRequest(ctx, method, addr, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set(headerContentType, bodyType)

	return client.Doer.Do(req)
}

// Delete makes a DELETE request to the given address.
func (client *Client) Delete(ctx context.Context, addr string) (*http.Response, error) {
	req, err := client.newRequest(ctx, http.MethodDelete, addr, nil)
	if err != nil {
		return nil, err
	}
	return client.Doer.Do(req)
}

// Head makes a HEAD request to the given address.
func (client *Client) Head(ctx context.Context, addr string) (*http.Response, error) {
	req, err := client.newRequest(ctx, http.MethodHead, addr, nil)
	if err != nil {
		return nil, err
	}
	return client.Doer.Do(req)
}

// Options makes a OPTIONS request to the given address.
func (client *Client) Options(ctx context.Context, addr string) (*http.Response, error) {
	req, err := client.newRequest(ctx, http.MethodOptions, addr, nil)
	if err != nil {
		return nil, err
	}
	return client.Doer.Do(req)
}

func (client *Client) newRequest(ctx context.Context, method, addr string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, addr, body)
	if err != nil {
		return nil, err
	}

	for k, v := range client.Header {
		key := textproto.CanonicalMIMEHeaderKey(k)
		req.Header[key] = slices.Clone(v)
	}

	req, errReq := client.wrapRequest(req)
	if errReq != nil {
		return nil, errReq
	}

	return req, nil
}

func (client *Client) wrapRequest(req *http.Request) (*http.Request, error) {
	if client.Middleware != nil {
		return client.Middleware(req)
	}
	return req, nil
}
