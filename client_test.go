package httpclient_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/ninedraft/httpclient"
)

var methodsBody = map[string]callFn{
	http.MethodGet:     (*httpclient.Client).Get,
	http.MethodDelete:  (*httpclient.Client).Delete,
	http.MethodOptions: (*httpclient.Client).Options,
	http.MethodPatch: func(cl *httpclient.Client, ctx context.Context, addr string) (*http.Response, error) {
		return cl.Patch(ctx, addr, "text/plain", bytes.NewBufferString("hello, world"))
	},
	http.MethodPost: func(cl *httpclient.Client, ctx context.Context, addr string) (*http.Response, error) {
		return cl.Post(ctx, addr, "text/plain", bytes.NewBufferString("hello, world"))
	},
	http.MethodPut: func(cl *httpclient.Client, ctx context.Context, addr string) (*http.Response, error) {
		return cl.Put(ctx, addr, "text/plain", bytes.NewBufferString("hello, world"))
	},
}

var methodsNoBody = map[string]callFn{
	http.MethodHead: (*httpclient.Client).Head,
}

func TestClientNew(t *testing.T) {
	t.Parallel()

	cl := httpclient.New()

	assertNotEqual(t, nil, cl, "got nil client")

	server := testServer(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	defer server.Assert(t)

	resp, err := cl.Get(context.Background(), server.URL)

	assertEqual(t, nil, err, "response error")
	assertEqual(t, http.StatusOK, resp.StatusCode, "status code")
}

func TestClient_Simple(t *testing.T) {
	t.Parallel()

	tc := func(method, body string, call callFn) {
		t.Run(method, func(t *testing.T) {
			t.Parallel()

			testMethod{
				Method: method,
				Body:   body,
				Status: http.StatusOK,
				Call:   call,
			}.Run(t)
		})
	}

	for method, call := range methodsBody {
		tc(method, "hello, world", call)
	}

	for method, call := range methodsNoBody {
		tc(method, "", call)
	}
}

type callFn = func(cl *httpclient.Client, ctx context.Context, addr string) (*http.Response, error)

type testMethod struct {
	Method string
	Body   string
	Call   callFn
	Status int
}

func (tm testMethod) Run(t *testing.T) {
	const headerKey, headerValue = "X-Test", "test"
	const path = "/test"

	handler := func(w http.ResponseWriter, r *http.Request) {
		assertEqual(t, tm.Method, r.Method, "http method")
		assertEqual(t, headerValue, r.Header.Get(headerKey), "header value")
		assertEqual(t, path, r.URL.Path, "path")

		w.WriteHeader(tm.Status)
		io.WriteString(w, tm.Body)
	}

	server := testServer(t, handler)
	defer server.Assert(t)

	client := httpclient.NewFrom(server.Client())
	client.Header.Set(headerKey, headerValue)

	ctx := context.Background()

	resp, errCall := tm.Call(client, ctx, server.URL+path)

	requireEqual(t, nil, errCall, "call error")
	defer resp.Body.Close()

	gotBody, errRead := io.ReadAll(resp.Body)

	assertEqual(t, nil, errRead, "read body error")
	assertEqual(t, tm.Body, string(gotBody), "response body")
	assertEqual(t, tm.Status, resp.StatusCode, "status code")
}

func assertEqual[E comparable](t *testing.T, expected, actual E, msg string, args ...any) {
	t.Helper()

	if expected != actual {
		t.Errorf(msg, args...)
		t.Errorf("expected: %v", expected)
		t.Errorf("got:      %v", actual)
	}
}

func assertNotEqual[E comparable](t *testing.T, expected, actual E, msg string, args ...any) {
	t.Helper()

	if expected == actual {
		t.Errorf(msg, args...)
		t.Errorf("expected to be not equal: %v", expected)
		t.Errorf("got:                      %v", actual)
	}
}

func requireEqual[E comparable](t *testing.T, expected, actual E, msg string, args ...any) {
	t.Helper()

	if expected != actual {
		t.Errorf(msg, args...)
		t.Errorf("expected: %v", expected)
		t.Errorf("got:      %v", actual)
		t.FailNow()
	}
}

func asssertErrorIs(t *testing.T, expected, actual error, msg string, args ...any) bool {
	t.Helper()

	if !errors.Is(actual, expected) {
		t.Errorf(msg, args...)
		t.Errorf("expected: %v", expected)
		t.Errorf("got:      %v", actual)

		return false
	}
	return true
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
