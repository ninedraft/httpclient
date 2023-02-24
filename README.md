# httpclient

[![codecov](https://codecov.io/gh/ninedraft/httpclient/branch/master/graph/badge.svg?token=6DWVJCYU4P)](https://codecov.io/gh/ninedraft/httpclient)


httpclient is a thin wrapper around [http.Client](https://pkg.go.dev/net/http#Client) with some useful features.

- methods: .Post, .Get, .Put, .Delete, etc.
- simplified JSON, form and multipart requests

## Installation
```bash
go get -v github.com/ninedraft/httpclient@latest
```

## Usage


**Simple POST request**
```go
resp, err := client.Post(ctx, "https://httpbin.org/post", "application/json",
		strings.NewReader(`{"foo": "bar"}`))
```

**JSON request**
```go
resp, err := client.PostJSON(ctx, "https://httpbin.org/post", map[string]string{"foo": "bar"})
```

**Form request**
```go
resp, err := client.PostForm(ctx, "https://httpbin.org/post", 
    url.Values{
        "foo": {"bar"}",
    })
```

**Multipart request**
```go
resp, err = client.PostMultipart(ctx, "https://httpbin.org/post",
    httpclient.WriteMultiparts(
	    httpclient.MultipartFields(url.Values{
		    "foo": []string{"bar"},
		}),
        httpclient.MultipartFile("file", "file.txt", strings.NewReader("file content")),
    ))
```