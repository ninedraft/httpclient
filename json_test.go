package httpclient_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/ninedraft/httpclient"
)

var methodsJSON = map[string]methodJSON{
	http.MethodPost:        (*httpclient.Client).PostJSON,
	http.MethodPut:         (*httpclient.Client).PutJSON,
	http.MethodPatch:       (*httpclient.Client).PatchJSON,
	httpclient.MethodQuery: (*httpclient.Client).QueryJSON,
}

type methodJSON = func(client *httpclient.Client, ctx context.Context, addr string, obj any) (*http.Response, error)

func TestClient_JSON(t *testing.T) {
	type request struct {
		Foo string `json:"foo"`
	}

	var requestBody = request{
		Foo: "bar",
	}

	tc := func(method string, call methodJSON) {
		t.Run(method, func(t *testing.T) {
			t.Parallel()

			server := testServer(t, func(w http.ResponseWriter, r *http.Request) {
				assertEqual(t, "application/json", r.Header.Get("Content-Type"), "content-type")

				body, errBody := io.ReadAll(r.Body)
				requireEqual(t, nil, errBody, "read body error")

				got := request{}
				assertEqual(t, nil, json.Unmarshal(body, &got), "request body")
				assertEqual(t, requestBody, got, "request body")

				w.WriteHeader(http.StatusOK)
			})
			defer server.Assert(t)

			client := httpclient.NewFrom(server.Client())
			ctx := context.Background()

			resp, errCall := call(client, ctx, server.URL, requestBody)

			if resp != nil {
				defer resp.Body.Close()
			}

			requireEqual(t, nil, errCall, "call error")
			assertEqual(t, http.StatusOK, resp.StatusCode, "status code")
		})
	}

	for method, call := range methodsJSON {
		tc(method, call)
	}
}
