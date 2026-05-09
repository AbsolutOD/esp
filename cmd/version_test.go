package cmd

import "testing"

func TestRunVersion(t *testing.T) {
	if err := runVersion(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
