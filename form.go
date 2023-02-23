package httpclient

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/exp/maps"
)

const encodingURL = "application/x-www-form-urlencoded"

// GetForm makes a GET request to the given address with the given form data.
// The data is encoded as URL query parameters.
func (client *Client) GetForm(ctx context.Context, addr string, data url.Values) (*http.Response, error) {
	u, errParse := url.Parse(addr)
	if errParse != nil {
		return nil, errParse
	}

	q := u.Query()
	maps.Copy(q, data)
	u.RawQuery = q.Encode()

	return client.Get(ctx, u.String())
}

// PostForm makes a POST request to the given address with the given form data.
// The data is encoded as "application/x-www-form-urlencoded".
func (client *Client) PostForm(ctx context.Context, addr string, data url.Values) (*http.Response, error) {
	body := strings.NewReader(data.Encode())

	return client.Post(ctx, addr, encodingURL, body)
}

// PutForm makes a PUT request to the given address with the given form data.
// The data is encoded as "application/x-www-form-urlencoded".
func (client *Client) PutForm(ctx context.Context, addr string, data url.Values) (*http.Response, error) {
	body := strings.NewReader(data.Encode())

	return client.Put(ctx, addr, encodingURL, body)
}

// PatchForm makes a PATCH request to the given address with the given form data.
// The data is encoded as "application/x-www-form-urlencoded".
func (client *Client) PatchForm(ctx context.Context, addr string, data url.Values) (*http.Response, error) {
	body := strings.NewReader(data.Encode())

	return client.Patch(ctx, addr, encodingURL, body)
}

// QueryForm makes a QUERY request to the given address with the given form data.
// The data is encoded as "application/x-www-form-urlencoded".
func (client *Client) QueryForm(ctx context.Context, addr string, data url.Values) (*http.Response, error) {
	body := strings.NewReader(data.Encode())

	return client.Query(ctx, addr, encodingURL, body)
}
