package cmd

import (
	"testing"

	"github.com/AbsolutOD/esp/internal/app"
)

// TestEspName pins formatParamName's rule:
//
//	if HasPrefix(n, OrgPrefix) -> n unchanged
//	else                       -> OrgPrefix + "_" + ReplaceAll(ToUpper(n), "-", "_")
func TestEspName(t *testing.T) {
	cfg := &app.Config{OrgPrefix: "ACME"}

	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "already prefixed (canonical) is unchanged", in: "ACME_FOO", want: "ACME_FOO"},
		{name: "lowercase input is uppercased and prefixed", in: "foo", want: "ACME_FOO"},
		{name: "hyphenated input is uppercased and hyphens become underscores", in: "foo-bar", want: "ACME_FOO_BAR"},
		{name: "multiple hyphens all convert to underscores", in: "foo-bar-baz", want: "ACME_FOO_BAR_BAZ"},
		{name: "mixed case input is fully uppercased", in: "fooBar", want: "ACME_FOOBAR"},
		{name: "any string starting with the prefix passes through unchanged (HasPrefix, not exact match)", in: "ACMEISH", want: "ACMEISH"},
		{name: "uppercase non-prefixed input still gets prefix", in: "FOO", want: "ACME_FOO"},
		{name: "already-prefixed input keeps its hyphens (HasPrefix branch returns verbatim)", in: "ACME-LEGACY-NAME", want: "ACME-LEGACY-NAME"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := formatParamName(cfg, tc.in)
			if got != tc.want {
				t.Errorf("formatParamName(cfg, %q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}
