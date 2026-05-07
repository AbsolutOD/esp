package cmd

import (
	"testing"

	"github.com/pinpt/esp/internal/app"
)

// withAppConfig swaps in a test app.Config on the package-level esp
// global and restores the previous value when the test finishes.
// getParamPath delegates to esp.GetAppParamPath, which depends on
// these fields.
func withAppConfig(t *testing.T, cfg app.Config) {
	t.Helper()
	prev := esp
	cfg.Filename = ".espFile"
	cfg.Path = ""
	scoped := cfg
	esp = &scoped
	t.Cleanup(func() { esp = prev })
}

// TestGetParamPath pins getParamPath's rule: a leading "/" means the
// caller already wrote a literal SSM path, so it passes through; any
// other shape is routed through esp.GetAppParamPath.
func TestGetParamPath(t *testing.T) {
	withAppConfig(t, app.Config{
		OrgName: "acme",
		Env:     "dev",
		AppName: "billing",
	})

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
			got := getParamPath(tc.in)
			if got != tc.want {
				t.Errorf("getParamPath(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}
