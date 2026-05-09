package cmd

import (
	"testing"

	"github.com/AbsolutOD/esp/internal/app"
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
