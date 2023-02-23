package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

// PostJSON makes a POST request to the given address with JSON-encoded body.
func (client *Client) PostJSON(ctx context.Context, addr string, obj any) (*http.Response, error) {
	return client.doJSON(ctx, http.MethodPost, addr, obj)
}

// PutJSON makes a PUT request to the given address with JSON-encoded body.
func (client *Client) PutJSON(ctx context.Context, addr string, obj any) (*http.Response, error) {
	return client.doJSON(ctx, http.MethodPut, addr, obj)
}

// PatchJSON makes a PATCH request to the given address with JSON-encoded body.
func (client *Client) PatchJSON(ctx context.Context, addr string, obj any) (*http.Response, error) {
	return client.doJSON(ctx, http.MethodPatch, addr, obj)
}

// QueryJSON makes a QUERY request to the given address with JSON-encoded body.
func (client *Client) QueryJSON(ctx context.Context, addr string, obj any) (*http.Response, error) {
	return client.doJSON(ctx, methodQuery, addr, obj)
}

func (client *Client) doJSON(ctx context.Context, method, addr string, obj any) (*http.Response, error) {
	body, errJSON := json.Marshal(obj)
	if errJSON != nil {
		return nil, errJSON
	}

	return client.doBody(ctx, method, addr, "application/json", bytes.NewReader(body))
}
