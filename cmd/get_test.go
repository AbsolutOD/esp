package cmd

import (
	"errors"
	"testing"

	"github.com/AbsolutOD/esp/internal/app"
	"github.com/AbsolutOD/esp/internal/common"
	"github.com/spf13/cobra"
)

// TestGetParamPath pins getParamPath's rule: a leading "/" means the
// caller already wrote a literal SSM path, so it passes through; any
// other shape is routed through cfg.GetAppParamPath.
func TestGetParamPath(t *testing.T) {
	cfg := &app.Config{
		OrgName: "acme",
		Env:     "dev",
		AppName: "billing",
	}

	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "leading slash → literal path returned unchanged",
			in:   "/acme/dev/billing/DB_URL",
			want: "/acme/dev/billing/DB_URL",
		},
		{
			name: "absolute path with different segments → unchanged",
			in:   "/some/other/place",
			want: "/some/other/place",
		},
		{
			name: "short name → resolved through GetAppParamPath",
			in:   "DB_URL",
			want: "/acme/dev/billing/DB_URL",
		},
		{
			name: "lowercase short name passes through verbatim (no upcasing here)",
			in:   "secret",
			want: "/acme/dev/billing/secret",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := getParamPath(cfg, tc.in)
			if got != tc.want {
				t.Errorf("getParamPath(cfg, %q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestRunGet_LiteralPath(t *testing.T) {
	fake := &fakeBackend{getOut: common.EspParam{Name: "/x/y", Value: "v"}}
	c := newTestEspClient(fake)
	cfg := testConfig()
	cmd := newCmdWithFlags(t, func(c *cobra.Command) {
		c.Flags().BoolP("decrypt", "d", false, "")
		c.Flags().BoolP("details", "t", false, "")
	}, []string{})

	if err := runGet(cmd, []string{"/x/y"}, c, cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fake.getIn.Name != "/x/y" {
		t.Errorf("GetParam called with %q, want /x/y", fake.getIn.Name)
	}
}

func TestRunGet_RelativePathResolved(t *testing.T) {
	fake := &fakeBackend{getOut: common.EspParam{Name: "/acme/dev/billing/DB"}}
	c := newTestEspClient(fake)
	cfg := testConfig()
	cmd := newCmdWithFlags(t, func(c *cobra.Command) {
		c.Flags().BoolP("decrypt", "d", false, "")
		c.Flags().BoolP("details", "t", false, "")
	}, []string{})

	if err := runGet(cmd, []string{"DB"}, c, cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fake.getIn.Name != "/acme/dev/billing/DB" {
		t.Errorf("got %q, want /acme/dev/billing/DB", fake.getIn.Name)
	}
}

func TestRunGet_DecryptFlagPropagates(t *testing.T) {
	fake := &fakeBackend{getOut: common.EspParam{Name: "/x"}}
	c := newTestEspClient(fake)
	cmd := newCmdWithFlags(t, func(c *cobra.Command) {
		c.Flags().BoolP("decrypt", "d", false, "")
		c.Flags().BoolP("details", "t", false, "")
	}, []string{"--decrypt"})

	if err := runGet(cmd, []string{"/x"}, c, testConfig()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !fake.getIn.Decrypt {
		t.Error("Decrypt = false, want true")
	}
}

func TestRunGet_GetParamErrorSurfaces(t *testing.T) {
	fake := &fakeBackend{getErr: errors.New("nope")}
	c := newTestEspClient(fake)
	cmd := newCmdWithFlags(t, func(c *cobra.Command) {
		c.Flags().BoolP("decrypt", "d", false, "")
		c.Flags().BoolP("details", "t", false, "")
	}, []string{})

	err := runGet(cmd, []string{"/x"}, c, testConfig())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
