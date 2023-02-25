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

var exampleService = func() string {
	const status = http.StatusOK
	var response = []byte(http.StatusText(status))

	handler := func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.Copy(io.Discard, r.Body)
		w.WriteHeader(status)
		w.Write(response)
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	return server.URL
}()

func ExampleClient_Get() {
	ctx := context.Background()
	client := httpclient.New()

	// GET request
	resp, err := client.Get(ctx, exampleService)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

func ExampleClient_Post() {
	ctx := context.Background()
	client := httpclient.New()

	// POST request
	resp, err := client.Post(ctx, exampleService, "application/json",
		strings.NewReader(`{"foo": "bar"}`))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

func ExampleClient_PostJSON() {
	ctx := context.Background()
	client := httpclient.New()

	// POST request
	resp, err := client.PostJSON(ctx, exampleService,
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
	client := httpclient.New()

	// POST request
	resp, err := client.PostForm(ctx, exampleService,
		url.Values{
			"foo": []string{"bar"},
		})
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// POST request
	resp, err = client.PostMultipart(ctx, exampleService,
		httpclient.MultipartFile("file", "file.txt", strings.NewReader("file content")))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// POST request
	resp, err = client.PostMultipart(ctx, exampleService,
		httpclient.MultipartFile("file", "file.txt", strings.NewReader("file content")))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// POST request
	resp, err = client.PostMultipart(ctx, exampleService,
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
