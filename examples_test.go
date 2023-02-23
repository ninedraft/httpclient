package httpclient_test

import (
	"context"
	"net/url"
	"strings"

	"github.com/ninedraft/httpclient"
)

func ExampleClient() {
	client := httpclient.New()
	ctx := context.Background()

	// GET request
	resp, err := client.Get(ctx, "https://httpbin.org/get")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// POST request
	resp, err = client.Post(ctx, "https://httpbin.org/post", "application/json",
		strings.NewReader(`{"foo": "bar"}`))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

func ExampleClient_JSON() {
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

func ExampleClient_Form() {
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
}

func ExampleClient_MultipartFields() {
	client := httpclient.New()
	ctx := context.Background()

	// POST request
	resp, err := client.PostMultipart(ctx, "https://httpbin.org/post",
		httpclient.MultipartFields(url.Values{
			"foo": []string{"bar"},
		}))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

func ExampleClient_MultipartFile() {
	client := httpclient.New()
	ctx := context.Background()

	// POST request
	resp, err := client.PostMultipart(ctx, "https://httpbin.org/post",
		httpclient.MultipartFile("file", "file.txt", strings.NewReader("file content")))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

func ExampleClient_Multipart() {
	client := httpclient.New()
	ctx := context.Background()

	// POST request
	resp, err := client.PostMultipart(ctx, "https://httpbin.org/post",
		func(w httpclient.MultipartWriter) error {
			// write form fields
			err := w.WriteField("foo", "bar")
			if err != nil {
				return err
			}
			// write file
			fw, err := w.CreateFormFile("file", "file.txt")
			if err != nil {
				return err
			}
			_, err = fw.Write([]byte("file content"))
			if err != nil {
				return err
			}
			return nil
		})
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}
