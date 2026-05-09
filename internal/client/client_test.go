package client

import (
	"errors"
	"strings"
	"testing"

	"github.com/AbsolutOD/esp/internal/app"
	"github.com/AbsolutOD/esp/internal/common"
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

// fakeBackend implements Client by recording calls and returning canned responses.
type fakeBackend struct {
	saveIn    common.EspParamInput
	saveOut   common.SaveOutput
	saveErr   error
	saveCalls int

	getIn    common.GetOneInput
	getOut   common.EspParam
	getErr   error
	getCalls int

	manyIn  common.ListParamInput
	manyOut []common.EspParam
	manyErr error

	copyIn    common.CopyCommand
	copyOut   common.SaveOutput
	copyErr   error
	copyCalls int

	delIn    common.DeleteInput
	delOut   string
	delErr   error
	delCalls int

	// scripted: when set, replaces single fields above. Indexed by call.
	getOuts []common.EspParam
	getErrs []error
	getIdx  int
}

func (f *fakeBackend) Save(p common.EspParamInput) (common.SaveOutput, error) {
	f.saveIn = p
	f.saveCalls++
	return f.saveOut, f.saveErr
}
func (f *fakeBackend) GetOne(p common.GetOneInput) (common.EspParam, error) {
	f.getIn = p
	f.getCalls++
	if f.getIdx < len(f.getOuts) {
		out := f.getOuts[f.getIdx]
		var err error
		if f.getIdx < len(f.getErrs) {
			err = f.getErrs[f.getIdx]
		}
		f.getIdx++
		return out, err
	}
	return f.getOut, f.getErr
}
func (f *fakeBackend) GetMany(p common.ListParamInput) ([]common.EspParam, error) {
	f.manyIn = p
	return f.manyOut, f.manyErr
}
func (f *fakeBackend) Copy(cc common.CopyCommand) (common.SaveOutput, error) {
	f.copyIn = cc
	f.copyCalls++
	return f.copyOut, f.copyErr
}
func (f *fakeBackend) Delete(p common.DeleteInput) (string, error) {
	f.delIn = p
	f.delCalls++
	return f.delOut, f.delErr
}

func TestEspClient_GetParam(t *testing.T) {
	fake := &fakeBackend{getOut: common.EspParam{Name: "/n", Value: "v"}}
	c := &EspClient{Backend: "ssm", Client: fake}

	got, err := c.GetParam(common.GetOneInput{Name: "/n", Decrypt: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Value != "v" {
		t.Errorf("Value = %q, want v", got.Value)
	}
	if fake.getIn.Name != "/n" || !fake.getIn.Decrypt {
		t.Errorf("GetOne input = %+v, want Name=/n Decrypt=true", fake.getIn)
	}
}

func TestEspClient_ListParams(t *testing.T) {
	want := []common.EspParam{{Name: "/a/1"}, {Name: "/a/2"}}
	fake := &fakeBackend{manyOut: want}
	c := &EspClient{Backend: "ssm", Client: fake}

	got, err := c.ListParams(common.ListParamInput{Path: "/a/", Recursive: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 || got[0].Name != "/a/1" {
		t.Errorf("got %+v, want %+v", got, want)
	}
	if fake.manyIn.Path != "/a/" || !fake.manyIn.Recursive {
		t.Errorf("ListParams input = %+v", fake.manyIn)
	}
}

func TestEspClient_Save(t *testing.T) {
	fake := &fakeBackend{saveOut: common.SaveOutput{Version: 3}}
	c := &EspClient{Backend: "ssm", Client: fake}

	out, err := c.Save(common.EspParamInput{Name: "/n", Value: "v", Secure: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Version != 3 {
		t.Errorf("Version = %d, want 3", out.Version)
	}
	if !fake.saveIn.Secure {
		t.Error("Save input.Secure = false, want true")
	}
}

func TestEspClient_Delete(t *testing.T) {
	fake := &fakeBackend{delOut: "/n"}
	c := &EspClient{Backend: "ssm", Client: fake}

	got, err := c.Delete(common.DeleteInput{Name: "/n"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "/n" {
		t.Errorf("got %q, want /n", got)
	}
}

func TestEspClient_Copy_Success(t *testing.T) {
	fake := &fakeBackend{
		copyOut: common.SaveOutput{Version: 1},
		getOut:  common.EspParam{Name: "/dest", Value: "v"},
	}
	c := &EspClient{Backend: "ssm", Client: fake}

	got, err := c.Copy(common.CopyCommand{Source: "/src", Destination: "/dest"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Name != "/dest" {
		t.Errorf("got Name=%q, want /dest", got.Name)
	}
	if fake.copyCalls != 1 {
		t.Errorf("Copy called %d times, want 1", fake.copyCalls)
	}
	if fake.getCalls != 1 {
		t.Errorf("GetOne called %d times, want 1 (re-fetch)", fake.getCalls)
	}
	if !fake.getIn.Decrypt {
		t.Error("re-fetch did not request decrypt; want Decrypt=true")
	}
}

func TestEspClient_Copy_CopyFailsNoRefetch(t *testing.T) {
	fake := &fakeBackend{copyErr: errors.New("nope")}
	c := &EspClient{Backend: "ssm", Client: fake}

	_, err := c.Copy(common.CopyCommand{Source: "/src", Destination: "/dest"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if fake.getCalls != 0 {
		t.Errorf("GetOne was called %d times after Copy failure; want 0", fake.getCalls)
	}
}

func TestEspClient_Copy_RefetchFails(t *testing.T) {
	fake := &fakeBackend{
		copyOut: common.SaveOutput{Version: 1},
		getErr:  errors.New("not found"),
	}
	c := &EspClient{Backend: "ssm", Client: fake}

	_, err := c.Copy(common.CopyCommand{Source: "/src", Destination: "/dest"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestEspClient_Move_Success(t *testing.T) {
	fake := &fakeBackend{
		getOut:  common.EspParam{Name: "/src", Value: "v", Secure: true},
		saveOut: common.SaveOutput{Version: 1},
		delOut:  "/src",
	}
	c := &EspClient{Backend: "ssm", Client: fake}

	mc, err := c.Move(common.MoveCommand{Source: "/src", Destination: "/dest"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mc.Source != "/src" || mc.Destination != "/dest" {
		t.Errorf("got %+v, want Source=/src Destination=/dest", mc)
	}
	if !fake.saveIn.Secure {
		t.Error("source's Secure flag did not carry to Save input")
	}
	if fake.saveIn.Name != "/dest" {
		t.Errorf("Save name = %q, want /dest", fake.saveIn.Name)
	}
	if fake.delIn.Name != "/src" {
		t.Errorf("Delete name = %q, want /src", fake.delIn.Name)
	}
}

func TestEspClient_Move_GetFailsNoSaveNoDelete(t *testing.T) {
	fake := &fakeBackend{getErr: errors.New("not found")}
	c := &EspClient{Backend: "ssm", Client: fake}

	_, err := c.Move(common.MoveCommand{Source: "/src", Destination: "/dest"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if fake.saveCalls != 0 || fake.delCalls != 0 {
		t.Errorf("Save called %d, Delete called %d after Get failure; want both 0", fake.saveCalls, fake.delCalls)
	}
}

func TestEspClient_Move_SaveFailsNoDelete(t *testing.T) {
	fake := &fakeBackend{
		getOut:  common.EspParam{Name: "/src"},
		saveErr: errors.New("limit"),
	}
	c := &EspClient{Backend: "ssm", Client: fake}

	_, err := c.Move(common.MoveCommand{Source: "/src", Destination: "/dest"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if fake.delCalls != 0 {
		t.Errorf("Delete called %d times after Save failure; want 0", fake.delCalls)
	}
}

func TestEspClient_Move_DeleteFailsSurfacesError(t *testing.T) {
	fake := &fakeBackend{
		getOut:  common.EspParam{Name: "/src"},
		saveOut: common.SaveOutput{Version: 1},
		delErr:  errors.New("perms"),
	}
	c := &EspClient{Backend: "ssm", Client: fake}

	_, err := c.Move(common.MoveCommand{Source: "/src", Destination: "/dest"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
