package app

import "testing"

func TestGetAppPath(t *testing.T) {
	tests := []struct {
		name string
		cfg  Config
		want string
	}{
		{
			name: "all fields populated",
			cfg:  Config{OrgName: "acme", Env: "dev", AppName: "billing"},
			want: "/acme/dev/billing/",
		},
		{
			name: "single-letter fields",
			cfg:  Config{OrgName: "a", Env: "b", AppName: "c"},
			want: "/a/b/c/",
		},
		{
			name: "empty fields produce empty segments",
			cfg:  Config{},
			want: "////",
		},
		{
			name: "only OrgName set",
			cfg:  Config{OrgName: "acme"},
			want: "/acme///",
		},
		{
			name: "leading slashes in inputs are not stripped",
			cfg:  Config{OrgName: "/acme", Env: "/dev", AppName: "/billing"},
			want: "//acme//dev//billing/",
		},
		{
			name: "names with hyphens and underscores pass through",
			cfg:  Config{OrgName: "acme-corp", Env: "prod", AppName: "billing_svc"},
			want: "/acme-corp/prod/billing_svc/",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.cfg.GetAppPath()
			if got != tc.want {
				t.Errorf("GetAppPath() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestGetAppParamPath(t *testing.T) {
	base := Config{OrgName: "acme", Env: "dev", AppName: "billing"}

	tests := []struct {
		name  string
		cfg   Config
		param string
		want  string
	}{
		{
			name:  "simple param",
			cfg:   base,
			param: "DB_URL",
			want:  "/acme/dev/billing/DB_URL",
		},
		{
			name:  "lowercase param",
			cfg:   base,
			param: "secret",
			want:  "/acme/dev/billing/secret",
		},
		{
			name:  "empty param leaves trailing slash",
			cfg:   base,
			param: "",
			want:  "/acme/dev/billing/",
		},
		{
			name:  "leading slash on param produces double slash",
			cfg:   base,
			param: "/DB_URL",
			want:  "/acme/dev/billing//DB_URL",
		},
		{
			name:  "param containing slashes pass through",
			cfg:   base,
			param: "nested/key/name",
			want:  "/acme/dev/billing/nested/key/name",
		},
		{
			name:  "all-empty config with param",
			cfg:   Config{},
			param: "X",
			want:  "////X",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.cfg.GetAppParamPath(tc.param)
			if got != tc.want {
				t.Errorf("GetAppParamPath(%q) = %q, want %q", tc.param, got, tc.want)
			}
		})
	}
}
