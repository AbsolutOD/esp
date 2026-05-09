package cmd

import (
	"errors"
	"testing"

	"github.com/AbsolutOD/esp/internal/common"
)

func TestRunCopy_HappyPath(t *testing.T) {
	fake := &fakeBackend{
		copyOut: common.SaveOutput{Version: 1},
		getOut:  common.EspParam{Name: "/dest"},
	}
	c := newTestEspClient(fake)

	if err := runCopy([]string{"/src", "/dest"}, c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fake.copyIn.Source != "/src" || fake.copyIn.Destination != "/dest" {
		t.Errorf("got %+v", fake.copyIn)
	}
}

func TestRunCopy_EmptySource(t *testing.T) {
	fake := &fakeBackend{}
	err := runCopy([]string{"", "/dest"}, newTestEspClient(fake))
	if err == nil || err.Error() != "source can not be empty" {
		t.Errorf("got %v, want \"source can not be empty\"", err)
	}
	if fake.copyCalls != 0 {
		t.Errorf("Copy called %d times, want 0", fake.copyCalls)
	}
}

func TestRunCopy_EmptyDestination(t *testing.T) {
	fake := &fakeBackend{}
	err := runCopy([]string{"/src", ""}, newTestEspClient(fake))
	if err == nil || err.Error() != "destination can not be empty" {
		t.Errorf("got %v, want \"destination can not be empty\"", err)
	}
	if fake.copyCalls != 0 {
		t.Errorf("Copy called %d times, want 0", fake.copyCalls)
	}
}

func TestRunCopy_BackendError(t *testing.T) {
	fake := &fakeBackend{copyErr: errors.New("perm")}
	if err := runCopy([]string{"/src", "/dest"}, newTestEspClient(fake)); err == nil {
		t.Fatal("expected error, got nil")
	}
}
