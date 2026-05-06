package cmd

import (
	"testing"

	"github.com/pinpt/esp/internal/app"
)

func TestGetPathWithFullPath(t *testing.T) {
	testPath := "/corpa/dev/foo_app/"
	path := getPath([]string{testPath})
	if path != testPath {
		t.Errorf("want: %s | got %s", testPath, path)
	}
}

func TestGetPathEnvVarName(t *testing.T) {
	envVar := "TEST_VAR"
	path := getPath([]string{envVar})
	if path != envVar {
		t.Errorf("want: %s | got %s", envVar, path)
	}
}

// TestGetPathRelative pins getPath's empty-args branch: no positional
// argument means "list the project's base path", which is built from
// esp.GetAppPath() and follows the /OrgName/Env/AppName/ convention.
func TestGetPathRelative(t *testing.T) {
	withAppConfig(t, app.Config{
		OrgName: "acme",
		Env:     "dev",
		AppName: "billing",
	})

	got := getPath(nil)
	want := "/acme/dev/billing/"
	if got != want {
		t.Errorf("getPath(nil) = %q, want %q", got, want)
	}
}
