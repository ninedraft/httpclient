package httpclient_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/ninedraft/httpclient"
)

type exampleRoundTripper struct{}

func (exampleRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	_, _ = io.Copy(io.Discard, req.Body)

	resp := httptest.NewRecorder()

	resp.WriteHeader(http.StatusOK)
	resp.WriteString("OK")

	return resp.Result(), nil
}

func ExampleClient_Get() {
	ctx := context.Background()
	client := httpclient.NewFrom(&http.Client{
		Transport: exampleRoundTripper{},
	})

	// GET request
	resp, err := client.Get(ctx, "https://example.com")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

func ExampleClient_Post() {
	ctx := context.Background()
	client := httpclient.NewFrom(&http.Client{
		Transport: exampleRoundTripper{},
	})

	// POST request
	resp, err := client.Post(ctx, "https://example.com", "application/json",
		strings.NewReader(`{"foo": "bar"}`))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

func ExampleClient_PostJSON() {
	ctx := context.Background()
	client := httpclient.NewFrom(&http.Client{
		Transport: exampleRoundTripper{},
	})

	// POST request
	resp, err := client.PostJSON(ctx, "https://example.com",
		map[string]string{
			"foo": "bar",
		})
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

func ExampleClient_PostForm() {
	ctx := context.Background()
	client := httpclient.NewFrom(&http.Client{
		Transport: exampleRoundTripper{},
	})

	// POST request
	resp, err := client.PostForm(ctx, "https://example.com",
		url.Values{
			"foo": []string{"bar"},
		})
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// POST request
	resp, err = client.PostMultipart(ctx, "https://example.com",
		httpclient.MultipartFile("file", "file.txt", strings.NewReader("file content")))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// POST request
	resp, err = client.PostMultipart(ctx, "https://example.com",
		httpclient.MultipartFile("file", "file.txt", strings.NewReader("file content")))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// POST request
	resp, err = client.PostMultipart(ctx, "https://example.com",
		httpclient.WriteMultiparts(
			httpclient.MultipartFields(url.Values{
				"foo": []string{"bar"},
			}),
			httpclient.MultipartFile("file", "file.txt", strings.NewReader("file content")),
		))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}
