package httpclient

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/exp/maps"
)

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

func (client *Client) PostForm(ctx context.Context, addr string, data url.Values) (*http.Response, error) {
	body := strings.NewReader(data.Encode())

	return client.Post(ctx, addr, "application/x-www-form-urlencoded", body)
}

func (client *Client) QueryForm(ctx context.Context, addr string, data url.Values) (*http.Response, error) {
	body := strings.NewReader(data.Encode())

	return client.Query(ctx, addr, "application/x-www-form-urlencoded", body)
}
