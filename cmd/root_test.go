package cmd

import (
	"os"
	"strings"
	"testing"

	"github.com/AbsolutOD/esp/internal/app"
	"github.com/spf13/cobra"
)

// rootForPreRunTest constructs a minimal root cobra.Command bound to
// the App, so the persistentPreRunE closure can call cmd.Root() if
// it reaches the IsEspProject branch.
func rootForPreRunTest(a *App) *cobra.Command {
	root := &cobra.Command{Use: "esp"}
	root.PersistentFlags().StringVarP(&a.Config.Env, "env", "e", "", "")
	return root
}

// unsetEnv unsets the variable for this test and registers a cleanup
// that restores the prior value. t.Setenv only sets; "must be unset"
// cases need this.
func unsetEnv(t *testing.T, key string) {
	t.Helper()
	prev, had := os.LookupEnv(key)
	if err := os.Unsetenv(key); err != nil {
		t.Fatalf("os.Unsetenv(%q): %v", key, err)
	}
	t.Cleanup(func() {
		if had {
			_ = os.Setenv(key, prev)
		} else {
			_ = os.Unsetenv(key)
		}
	})
}

func TestPersistentPreRunE_MissingRegion(t *testing.T) {
	unsetEnv(t, "AWS_DEFAULT_REGION")
	t.Setenv("AWS_PROFILE", "default")

	a := &App{Config: app.New(false)}
	a.Config.Backend = "ssm"
	pre := persistentPreRunE(a)

	err := pre(rootForPreRunTest(a), nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "AWS_DEFAULT_REGION") {
		t.Errorf("error %q does not mention AWS_DEFAULT_REGION", err.Error())
	}
}

func TestPersistentPreRunE_MissingProfile(t *testing.T) {
	t.Setenv("AWS_DEFAULT_REGION", "us-east-1")
	unsetEnv(t, "AWS_PROFILE")

	a := &App{Config: app.New(false)}
	a.Config.Backend = "ssm"
	pre := persistentPreRunE(a)

	err := pre(rootForPreRunTest(a), nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "AWS_PROFILE") {
		t.Errorf("error %q does not mention AWS_PROFILE", err.Error())
	}
}

func TestPersistentPreRunE_UnsupportedBackend(t *testing.T) {
	t.Setenv("AWS_DEFAULT_REGION", "us-east-1")
	t.Setenv("AWS_PROFILE", "default")

	a := &App{Config: app.New(false)}
	a.Config.Backend = "vault"
	pre := persistentPreRunE(a)

	err := pre(rootForPreRunTest(a), nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "vault") {
		t.Errorf("error %q does not mention the backend name", err.Error())
	}
}
