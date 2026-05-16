package cmd

import "testing"

func TestRunVersion(t *testing.T) {
	if err := runVersion(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestVersionString_Defaults(t *testing.T) {
	got := versionString()
	want := "esp dev (commit none, built unknown)"
	if got != want {
		t.Fatalf("versionString() = %q, want %q", got, want)
	}
}
