package client

import (
	"strings"
	"testing"

	"github.com/AbsolutOD/esp/internal/app"
)

// TestNewUnsupportedBackend covers the only branch of client.New
// reachable without AWS credentials: an unsupported Backend value
// must return a non-nil error and a nil EspClient.
//
// The "ssm" happy path is not tested here. It calls
// config.LoadDefaultConfig, which would either reach out to the
// environment for AWS credentials or fail with MissingRegion — both
// of which make the test environment-dependent. Exercising it
// requires either a dependency-injection seam or an integration
// harness, both deferred to future work (see the cleanup-and-tests
// design doc, "Non-Goals").
func TestNewUnsupportedBackend(t *testing.T) {
	c, err := New(&app.Config{Backend: "vault"})
	if err == nil {
		t.Fatalf("New with unsupported backend returned nil error; want error, got client = %#v", c)
	}
	if c != nil {
		t.Errorf("New with unsupported backend returned non-nil client = %#v; want nil", c)
	}
	if !strings.Contains(err.Error(), "vault") {
		t.Errorf("error message %q does not mention the offending backend %q", err.Error(), "vault")
	}
}
