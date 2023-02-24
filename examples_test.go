package httpclient_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/ninedraft/httpclient"
)

var client = httpclient.NewFrom(&http.Client{
	Transport: &mockTransport{},
})

type mockTransport struct{}

func (t *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp := httptest.NewRecorder()

	resp.WriteHeader(http.StatusOK)
	resp.WriteString(http.StatusText(http.StatusOK))

	return resp.Result(), nil
}

func ExampleClient_Get() {
	ctx := context.Background()

	// GET request
	resp, err := client.Get(ctx, "https://httpbin.org/get")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

func ExampleClient_Post() {
	ctx := context.Background()

	// POST request
	resp, err := client.Post(ctx, "https://httpbin.org/post", "application/json",
		strings.NewReader(`{"foo": "bar"}`))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

func ExampleClient_PostJSON() {
	ctx := context.Background()

	// POST request
	resp, err := client.PostJSON(ctx, "https://httpbin.org/post",
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

	// POST request
	resp, err := client.PostForm(ctx, "https://httpbin.org/post",
		url.Values{
			"foo": []string{"bar"},
		})
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// POST request
	resp, err = client.PostMultipart(ctx, "https://httpbin.org/post",
		httpclient.MultipartFile("file", "file.txt", strings.NewReader("file content")))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// POST request
	resp, err = client.PostMultipart(ctx, "https://httpbin.org/post",
		httpclient.MultipartFile("file", "file.txt", strings.NewReader("file content")))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// POST request
	resp, err = client.PostMultipart(ctx, "https://httpbin.org/post",
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
