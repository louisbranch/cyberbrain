package test

import (
	"reflect"
	"testing"
)

func Equal(t *testing.T, desc string, expected, got interface{}) {
	t.Helper()
	if !reflect.DeepEqual(expected, got) {
		t.Fatalf("%s expected %v, got %v", desc, expected, got)
	}
}

func OK(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("expected no error got %q", err)
	}
}

func Error(t *testing.T, err error, msg ...string) {
	t.Helper()
	if err == nil {
		t.Fatal("expected error not to be nil")
	}

	emsg := err.Error()

	if len(msg) > 0 && emsg != msg[0] {
		t.Fatalf("expected error to be %q, got %q", emsg, msg[0])
	}
}
