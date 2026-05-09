package cmd

import (
	"errors"
	"testing"

	"github.com/AbsolutOD/esp/internal/common"
)

func TestRunMove_HappyPath(t *testing.T) {
	fake := &fakeBackend{
		getOut:  common.EspParam{Name: "/src", Value: "v"},
		saveOut: common.SaveOutput{Version: 1},
		delOut:  "/src",
	}
	c := newTestEspClient(fake)

	if err := runMove([]string{"/src", "/dest"}, c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunMove_BackendError(t *testing.T) {
	fake := &fakeBackend{getErr: errors.New("not found")}
	if err := runMove([]string{"/src", "/dest"}, newTestEspClient(fake)); err == nil {
		t.Fatal("expected error, got nil")
	}
}
