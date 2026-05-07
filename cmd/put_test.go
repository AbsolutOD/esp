package cmd

import (
	"testing"

	"github.com/pinpt/esp/internal/app"
)

// withOrgPrefix swaps in a test OrgPrefix on the package-level esp
// config and restores the previous value when the test finishes.
// formatParamName depends on this global; resetting prevents test
// pollution across the cmd package.
func withOrgPrefix(t *testing.T, prefix string) {
	t.Helper()
	if esp == nil {
		esp = app.New(false)
	}
	prev := esp.OrgPrefix
	esp.OrgPrefix = prefix
	t.Cleanup(func() { esp.OrgPrefix = prev })
}

// TestEspName pins formatParamName's rule:
//
//	if HasPrefix(n, OrgPrefix) -> n unchanged
//	else                       -> OrgPrefix + "_" + ReplaceAll(ToUpper(n), "-", "_")
//
// Hyphens are normalized to underscores so the resulting name is a
// valid environment variable identifier.
func TestEspName(t *testing.T) {
	withOrgPrefix(t, "ACME")

	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "already prefixed (canonical) is unchanged",
			in:   "ACME_FOO",
			want: "ACME_FOO",
		},
		{
			name: "lowercase input is uppercased and prefixed",
			in:   "foo",
			want: "ACME_FOO",
		},
		{
			name: "hyphenated input is uppercased and hyphens become underscores",
			in:   "foo-bar",
			want: "ACME_FOO_BAR",
		},
		{
			name: "multiple hyphens all convert to underscores",
			in:   "foo-bar-baz",
			want: "ACME_FOO_BAR_BAZ",
		},
		{
			name: "mixed case input is fully uppercased",
			in:   "fooBar",
			want: "ACME_FOOBAR",
		},
		{
			name: "any string starting with the prefix passes through unchanged (HasPrefix, not exact match)",
			in:   "ACMEISH",
			want: "ACMEISH",
		},
		{
			name: "uppercase non-prefixed input still gets prefix",
			in:   "FOO",
			want: "ACME_FOO",
		},
		{
			name: "already-prefixed input keeps its hyphens (HasPrefix branch returns verbatim)",
			in:   "ACME-LEGACY-NAME",
			want: "ACME-LEGACY-NAME",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := formatParamName(tc.in)
			if got != tc.want {
				t.Errorf("formatParamName(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}
