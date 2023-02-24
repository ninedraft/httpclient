package httpclient_test

import (
	"context"
	"net/url"
	"strings"

	"github.com/ninedraft/httpclient"
)

func ExampleClient_Get() {
	client := httpclient.New()
	ctx := context.Background()

	// GET request
	resp, err := client.Get(ctx, "https://httpbin.org/get")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

func ExampleClient_Post() {
	client := httpclient.New()
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
	client := httpclient.New()
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
	client := httpclient.New()
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
