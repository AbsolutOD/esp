package cmd

import (
	"errors"
	"testing"

	"github.com/AbsolutOD/esp/internal/app"
	"github.com/AbsolutOD/esp/internal/common"
	"github.com/spf13/cobra"
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

func TestRunList_NoArgsUsesAppPath(t *testing.T) {
	fake := &fakeBackend{manyOut: []common.EspParam{{Name: "/acme/dev/billing/X", Value: "v"}}}
	c := newTestEspClient(fake)
	cfg := testConfig()
	cmd := newCmdWithFlags(t, func(c *cobra.Command) {
		c.Flags().BoolP("decrypt", "d", false, "")
		c.Flags().BoolP("path", "p", false, "")
	}, []string{})

	if err := runList(cmd, nil, c, cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fake.manyIn.Path != "/acme/dev/billing/" {
		t.Errorf("Path = %q, want /acme/dev/billing/", fake.manyIn.Path)
	}
	if !fake.manyIn.Recursive {
		t.Error("Recursive = false, want true")
	}
}

func TestRunList_WithArgUsesArg(t *testing.T) {
	fake := &fakeBackend{}
	c := newTestEspClient(fake)
	cmd := newCmdWithFlags(t, func(c *cobra.Command) {
		c.Flags().BoolP("decrypt", "d", false, "")
		c.Flags().BoolP("path", "p", false, "")
	}, []string{})

	if err := runList(cmd, []string{"/elsewhere/"}, c, testConfig()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fake.manyIn.Path != "/elsewhere/" {
		t.Errorf("Path = %q, want /elsewhere/", fake.manyIn.Path)
	}
}

func TestRunList_DecryptFlagPropagates(t *testing.T) {
	fake := &fakeBackend{}
	c := newTestEspClient(fake)
	cmd := newCmdWithFlags(t, func(c *cobra.Command) {
		c.Flags().BoolP("decrypt", "d", false, "")
		c.Flags().BoolP("path", "p", false, "")
	}, []string{"--decrypt"})

	if err := runList(cmd, []string{"/x/"}, c, testConfig()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !fake.manyIn.Decrypt {
		t.Error("Decrypt = false, want true")
	}
}

func TestRunList_ErrorSurfaces(t *testing.T) {
	fake := &fakeBackend{manyErr: errors.New("api err")}
	c := newTestEspClient(fake)
	cmd := newCmdWithFlags(t, func(c *cobra.Command) {
		c.Flags().BoolP("decrypt", "d", false, "")
		c.Flags().BoolP("path", "p", false, "")
	}, []string{})

	if err := runList(cmd, []string{"/x/"}, c, testConfig()); err == nil {
		t.Fatal("expected error, got nil")
	}
}
