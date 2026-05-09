package cmd

import (
	"errors"
	"testing"
)

func TestRunDelete_HappyPath(t *testing.T) {
	fake := &fakeBackend{delOut: "/x"}
	if err := runDelete([]string{"/x"}, newTestEspClient(fake)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fake.delIn.Name != "/x" {
		t.Errorf("Delete name = %q, want /x", fake.delIn.Name)
	}
}

func TestRunDelete_BackendError(t *testing.T) {
	fake := &fakeBackend{delErr: errors.New("nope")}
	if err := runDelete([]string{"/x"}, newTestEspClient(fake)); err == nil {
		t.Fatal("expected error, got nil")
	}
}
