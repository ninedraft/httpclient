package httpclient_test

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/ninedraft/httpclient"
	"golang.org/x/exp/slices"
)

// special cases: GET, QUERY
var methodsForms = map[string]methodForm{
	http.MethodPost:  (*httpclient.Client).PostForm,
	http.MethodPut:   (*httpclient.Client).PutForm,
	http.MethodPatch: (*httpclient.Client).PatchForm,
	http.MethodGet:   (*httpclient.Client).GetForm,
}

type methodForm = func(client *httpclient.Client, ctx context.Context, addr string, data url.Values) (*http.Response, error)

func TestClient_Form(t *testing.T) {
	const formKey = "foo"
	values := url.Values{
		formKey: {"bar"},
	}

	const body = "hello, world!"

	tc := func(method string, call methodForm) {
		t.Run(method, func(t *testing.T) {
			t.Parallel()

			server := testServer(t, func(w http.ResponseWriter, r *http.Request) {
				errForm := r.ParseForm()
				assertEqual(t, nil, errForm, "parse form error")

				for key, v := range values {
					assertEqualSlices(t, v, r.Form[key], "form[%s]", key)
				}

				w.WriteHeader(http.StatusOK)
				_, _ = io.WriteString(w, body)
			})
			defer server.Assert(t)

			client := httpclient.NewFrom(server.Client())
			ctx := context.Background()

			resp, errCall := call(client, ctx, server.URL, values)

			if resp != nil {
				defer resp.Body.Close()
			}

			requireEqual(t, nil, errCall, "call error")
			assertEqual(t, http.StatusOK, resp.StatusCode, "status code")
			assertEqual(t, body, readString(t, resp.Body), "body")
		})
	}

	for method, call := range methodsForms {
		tc(method, call)
	}
}

func TestClient_FormQuery(t *testing.T) {
	t.Parallel()

	const formKey = "foo"
	values := url.Values{
		formKey: {"bar"},
	}

	const body = "hello, world!"

	server := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		form, errForm := url.ParseQuery(readString(t, r.Body))
		assertEqual(t, nil, errForm, "parse form error")

		for key, v := range values {
			assertEqualSlices(t, v, form[key], "form[%s]", key)
		}

		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, body)
	})
	defer server.Assert(t)

	client := httpclient.NewFrom(server.Client())
	ctx := context.Background()

	resp, errCall := client.QueryForm(ctx, server.URL, values)

	if resp != nil {
		defer resp.Body.Close()
	}

	requireEqual(t, nil, errCall, "call error")
	assertEqual(t, http.StatusOK, resp.StatusCode, "status code")
	assertEqual(t, body, readString(t, resp.Body), "body")
}

func readString(t *testing.T, re io.Reader) string {
	t.Helper()

	body, errRead := io.ReadAll(re)
	requireEqual(t, nil, errRead, "read body error")

	return string(body)
}

func assertEqualSlices[E comparable](t *testing.T, want, got []E, msg string, args ...any) {
	t.Helper()

	if !slices.Equal(want, got) {
		t.Errorf(msg, args...)
		t.Errorf("expected: %+v", want)
		t.Errorf("got:      %+v", got)
	}
}
