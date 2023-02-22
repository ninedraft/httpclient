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

type Client struct {
	Header     http.Header
	Transport  Transport
	Middleware func(req *http.Request) (*http.Request, error)
}

func New() *Client {
	return NewFrom(newTransport())
}

func NewFrom(transport Transport) *Client {
	return &Client{
		Transport: transport,
		Header:    make(http.Header),
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

type Transport interface {
	Do(req *http.Request) (*http.Response, error)
}

func (client *Client) Get(ctx context.Context, addr string) (*http.Response, error) {
	req, err := client.newRequest(ctx, http.MethodGet, addr, nil)
	if err != nil {
		return nil, err
	}
	return client.Transport.Do(req)
}

func (client *Client) Post(ctx context.Context, addr, bodyType string, body io.Reader) (*http.Response, error) {
	return client.doBody(ctx, http.MethodPost, addr, bodyType, body)
}

func (client *Client) Put(ctx context.Context, addr, bodyType string, body io.Reader) (*http.Response, error) {
	return client.doBody(ctx, http.MethodPut, addr, bodyType, body)
}

func (client *Client) Patch(ctx context.Context, addr, bodyType string, body io.Reader) (*http.Response, error) {
	return client.doBody(ctx, http.MethodPatch, addr, bodyType, body)
}

func (client *Client) doBody(ctx context.Context, method, addr, bodyType string, body io.Reader) (*http.Response, error) {
	req, err := client.newRequest(ctx, method, addr, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set(headerContentType, bodyType)

	return client.Transport.Do(req)
}

func (client *Client) Delete(ctx context.Context, addr string) (*http.Response, error) {
	req, err := client.newRequest(ctx, http.MethodDelete, addr, nil)
	if err != nil {
		return nil, err
	}
	return client.Transport.Do(req)
}

func (client *Client) Head(ctx context.Context, addr string) (*http.Response, error) {
	req, err := client.newRequest(ctx, http.MethodHead, addr, nil)
	if err != nil {
		return nil, err
	}
	return client.Transport.Do(req)
}

func (client *Client) Options(ctx context.Context, addr string) (*http.Response, error) {
	req, err := client.newRequest(ctx, http.MethodOptions, addr, nil)
	if err != nil {
		return nil, err
	}
	return client.Transport.Do(req)
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
