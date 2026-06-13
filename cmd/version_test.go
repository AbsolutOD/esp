package cmd

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/AbsolutOD/esp/internal/app"
)

func TestVersionCmd_NoAWSEnvVars(t *testing.T) {
	unsetEnv(t, "AWS_DEFAULT_REGION")
	unsetEnv(t, "AWS_PROFILE")

	a := &App{Config: app.New(false)}
	a.Config.Backend = "ssm"
	root := newRootCmd(a)
	root.AddCommand(newVersionCmd())
	root.SetArgs([]string{"version"})

	// Capture stdout — runVersion uses fmt.Println, which writes to os.Stdout.
	r, w, _ := os.Pipe()
	origStdout := os.Stdout
	os.Stdout = w
	t.Cleanup(func() { os.Stdout = origStdout })

	err := root.Execute()
	_ = w.Close()
	out, _ := io.ReadAll(r)

	if err != nil {
		t.Fatalf("Execute() returned error %v; expected nil (version should bypass AWS checks)", err)
	}
	got := string(bytes.TrimRight(out, "\n"))
	want := "esp dev (commit none, built unknown)"
	if got != want {
		t.Errorf("stdout = %q, want %q", got, want)
	}
}

func TestRunVersion(t *testing.T) {
	if err := runVersion(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestVersionString_Defaults(t *testing.T) {
	got := versionString()
	want := "esp dev (commit none, built unknown)"
	if got != want {
		t.Fatalf("versionString() = %q, want %q", got, want)
	}
}
