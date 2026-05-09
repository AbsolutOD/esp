package cmd

import (
	"errors"
	"testing"

	"github.com/AbsolutOD/esp/internal/app"
	"github.com/AbsolutOD/esp/internal/common"
	"github.com/spf13/cobra"
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

func TestRunPut_HappyPath(t *testing.T) {
	fake := &fakeBackend{
		saveOut: common.SaveOutput{Version: 1},
		getOut:  common.EspParam{Name: "/acme/dev/billing/ACME_FOO", Value: "v"},
	}
	c := newTestEspClient(fake)
	cmd := newCmdWithFlags(t, func(c *cobra.Command) {
		c.Flags().StringP("name", "n", "", "")
		c.Flags().StringP("value", "v", "", "")
		c.Flags().BoolP("secure", "s", false, "")
	}, []string{"--name", "foo", "--value", "v"})

	if err := runPut(cmd, c, testConfig()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fake.saveIn.Name != "/acme/dev/billing/ACME_FOO" {
		t.Errorf("Save name = %q, want /acme/dev/billing/ACME_FOO", fake.saveIn.Name)
	}
	if fake.saveCalls != 1 || fake.getCalls != 1 {
		t.Errorf("Save calls=%d Get calls=%d, want 1,1", fake.saveCalls, fake.getCalls)
	}
}

func TestRunPut_HyphenInName(t *testing.T) {
	fake := &fakeBackend{
		saveOut: common.SaveOutput{Version: 1},
		getOut:  common.EspParam{Name: "/acme/dev/billing/ACME_FOO_BAR"},
	}
	c := newTestEspClient(fake)
	cmd := newCmdWithFlags(t, func(c *cobra.Command) {
		c.Flags().StringP("name", "n", "", "")
		c.Flags().StringP("value", "v", "", "")
		c.Flags().BoolP("secure", "s", false, "")
	}, []string{"--name", "foo-bar", "--value", "v"})

	if err := runPut(cmd, c, testConfig()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fake.saveIn.Name != "/acme/dev/billing/ACME_FOO_BAR" {
		t.Errorf("Save name = %q, want /acme/dev/billing/ACME_FOO_BAR", fake.saveIn.Name)
	}
}

func TestRunPut_SaveFailsNoRefetch(t *testing.T) {
	fake := &fakeBackend{saveErr: errors.New("limit")}
	c := newTestEspClient(fake)
	cmd := newCmdWithFlags(t, func(c *cobra.Command) {
		c.Flags().StringP("name", "n", "", "")
		c.Flags().StringP("value", "v", "", "")
		c.Flags().BoolP("secure", "s", false, "")
	}, []string{"--name", "foo", "--value", "v"})

	err := runPut(cmd, c, testConfig())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if fake.getCalls != 0 {
		t.Errorf("re-fetch called %d times after Save failure; want 0", fake.getCalls)
	}
}

func TestRunPut_RefetchFails(t *testing.T) {
	fake := &fakeBackend{
		saveOut: common.SaveOutput{Version: 1},
		getErr:  errors.New("nope"),
	}
	c := newTestEspClient(fake)
	cmd := newCmdWithFlags(t, func(c *cobra.Command) {
		c.Flags().StringP("name", "n", "", "")
		c.Flags().StringP("value", "v", "", "")
		c.Flags().BoolP("secure", "s", false, "")
	}, []string{"--name", "foo", "--value", "v"})

	if err := runPut(cmd, c, testConfig()); err == nil {
		t.Fatal("expected error, got nil")
	}
}
