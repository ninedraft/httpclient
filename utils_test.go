package httpclient_test

import (
	"io"
	"testing"

	"golang.org/x/exp/slices"
)

const MethodQuery = "QUERY"

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
