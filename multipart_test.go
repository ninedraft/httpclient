package httpclient_test

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/ninedraft/httpclient"
)

var methodsMultiPart = map[string]methodMultipart{
	http.MethodPost:  (*httpclient.Client).PostMultipart,
	http.MethodPut:   (*httpclient.Client).PutMultipart,
	http.MethodPatch: (*httpclient.Client).PatchMultipart,
	MethodQuery:      (*httpclient.Client).QueryMultipart,
}

type methodMultipart = func(client *httpclient.Client, ctx context.Context, addr string, writeMultipart httpclient.WriteMultipart) (*http.Response, error)

func TestClient_MultipartFile(t *testing.T) {
	const field, filename, value = "field", "filename", "value"

	tc := func(method string, call methodMultipart) {
		t.Run(method, func(t *testing.T) {
			t.Parallel()

			server := testServer(t, func(w http.ResponseWriter, r *http.Request) {
				errParse := r.ParseMultipartForm(1 << 20)
				assertEqual(t, nil, errParse, "parse multipart form error")

				file, errFile := r.MultipartForm.File[field][0].Open()
				assertEqual(t, nil, errFile, "file open error")
				defer file.Close()

				got := readString(t, file)

				assertEqual(t, value, got, "field value")

				w.WriteHeader(http.StatusOK)
			})
			defer server.Assert(t)

			client := httpclient.NewFrom(server.Client())
			ctx := context.Background()

			resp, errCall := call(client, ctx, server.URL,
				httpclient.MultipartFile(field, filename, strings.NewReader(value)))

			assertEqual(t, nil, errCall, "call error")
			assertEqual(t, http.StatusOK, resp.StatusCode, "status code")
		})
	}

	for method, call := range methodsMultiPart {
		tc(method, call)
	}
}

func TestClient_MultipartFields(t *testing.T) {
	var values = url.Values{
		"field": {"value"},
	}

	tc := func(method string, call methodMultipart) {
		t.Run(method, func(t *testing.T) {
			t.Parallel()

			server := testServer(t, func(w http.ResponseWriter, r *http.Request) {
				errParse := r.ParseMultipartForm(1 << 20)
				assertEqual(t, nil, errParse, "parse multipart form error")

				form := r.MultipartForm.Value
				for field, value := range values {
					assertEqualSlices(t, value, form[field], "field value")
				}

				w.WriteHeader(http.StatusOK)
			})
			defer server.Assert(t)

			client := httpclient.NewFrom(server.Client())
			ctx := context.Background()

			resp, errCall := call(client, ctx, server.URL,
				httpclient.MultipartFields(values))

			assertEqual(t, nil, errCall, "call error")
			assertEqual(t, http.StatusOK, resp.StatusCode, "status code")
		})
	}

	for method, call := range methodsMultiPart {
		tc(method, call)
	}
}
