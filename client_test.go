package httpclient_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/ninedraft/httpclient"
)

func TestClient_Simple(t *testing.T) {
	t.Parallel()

	tc := func(method, body string, call callFn) {
		t.Run(method, func(t *testing.T) {
			t.Parallel()

			testMethod{
				Method: method,
				Body:   body,
				Call:   call,
				Status: http.StatusOK,
				Err:    nil,
			}.Run(t)
		})
	}

	tc(http.MethodHead, "", (*httpclient.Client).Head)

	const body = "hello, world"

	tc(http.MethodGet, body, (*httpclient.Client).Get)
	tc(http.MethodOptions, body, (*httpclient.Client).Options)
	tc(http.MethodDelete, body, (*httpclient.Client).Delete)

	tc(http.MethodPatch, body,
		func(cl *httpclient.Client, ctx context.Context, addr string) (*http.Response, error) {
			return cl.Patch(ctx, addr, "text/plain", bytes.NewBufferString("hello, world"))
		})

	tc(http.MethodPut, body,
		func(cl *httpclient.Client, ctx context.Context, addr string) (*http.Response, error) {
			return cl.Put(ctx, addr, "text/plain", bytes.NewBufferString("hello, world"))
		})

	tc(http.MethodPost, body,
		func(cl *httpclient.Client, ctx context.Context, addr string) (*http.Response, error) {
			return cl.Post(ctx, addr, "text/plain", bytes.NewBufferString("hello, world"))
		})
}

type callFn = func(cl *httpclient.Client, ctx context.Context, addr string) (*http.Response, error)

type testMethod struct {
	Method string
	Body   string
	Call   callFn
	Status int
	Err    error
}

func (tm testMethod) Run(t *testing.T) {
	const headerKey, headerValue = "X-Test", "test"
	const path = "/test"

	server := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		assertEqual(t, tm.Method, r.Method, "http method")
		assertEqual(t, headerValue, r.Header.Get(headerKey), "header value")
		assertEqual(t, path, r.URL.Path, "path")

		w.WriteHeader(tm.Status)
		io.WriteString(w, tm.Body)
	})
	defer server.Assert(t)

	client := httpclient.NewFrom(server.Client())
	client.Header.Set(headerKey, headerValue)

	ctx := context.Background()

	resp, errCall := tm.Call(client, ctx, server.URL+path)

	assertEqual(t, nil, errCall, "call error")
	defer resp.Body.Close()

	gotBody, errRead := io.ReadAll(resp.Body)

	assertEqual(t, tm.Err, errRead, "read error")
	if errRead == nil {
		assertEqual(t, tm.Body, string(gotBody), "response body")
		assertEqual(t, tm.Status, resp.StatusCode, "status code")
	}
}

func assertEqual[E comparable](t *testing.T, expected, actual E, msg string, args ...any) {
	t.Helper()

	if expected != actual {
		t.Errorf(msg, args...)
		t.Errorf("expected: %v", expected)
		t.Errorf("got:      %v", actual)
	}
}

type serverAssert struct {
	isCalled atomic.Bool
	*httptest.Server
}

func (sa *serverAssert) Assert(t *testing.T) {
	t.Helper()

	if !sa.isCalled.Load() {
		t.Error("server was not called")
	}
}

func testServer(t *testing.T, handler http.HandlerFunc) *serverAssert {
	assert := &serverAssert{}

	handle := func(w http.ResponseWriter, r *http.Request) {
		assert.isCalled.Store(true)
		handler(w, r)
	}

	assert.Server = httptest.NewServer(http.HandlerFunc(handle))
	t.Cleanup(assert.Close)

	return assert
}
