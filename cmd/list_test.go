package cmd

import (
	"testing"

	"github.com/AbsolutOD/esp/internal/app"
)

func TestGetPathWithFullPath(t *testing.T) {
	cfg := &app.Config{}
	testPath := "/corpa/dev/foo_app/"
	if got := getPath(cfg, []string{testPath}); got != testPath {
		t.Errorf("want: %s | got %s", testPath, got)
	}
}

func TestGetPathEnvVarName(t *testing.T) {
	cfg := &app.Config{}
	envVar := "TEST_VAR"
	if got := getPath(cfg, []string{envVar}); got != envVar {
		t.Errorf("want: %s | got %s", envVar, got)
	}
}

// TestGetPathRelative pins getPath's empty-args branch.
func TestGetPathRelative(t *testing.T) {
	cfg := &app.Config{
		OrgName: "acme",
		Env:     "dev",
		AppName: "billing",
	}
	got := getPath(cfg, nil)
	want := "/acme/dev/billing/"
	if got != want {
		t.Errorf("getPath(cfg, nil) = %q, want %q", got, want)
	}
}
