package cmd

import (
	"testing"

	"github.com/AbsolutOD/esp/internal/app"
	"github.com/AbsolutOD/esp/internal/client"
	"github.com/AbsolutOD/esp/internal/common"
	"github.com/spf13/cobra"
)

// fakeBackend implements client.Client by recording calls and
// returning canned responses. A copy of internal/client's fakeBackend
// — duplication is acceptable for two callers (see spec risks).
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

// newTestEspClient wraps a fakeBackend in a real *EspClient so runX
// functions (which take *EspClient) accept it.
func newTestEspClient(fake *fakeBackend) *client.EspClient {
	return &client.EspClient{Backend: "ssm", Client: fake}
}

// newCmdWithFlags builds a bare cobra.Command with the given flags
// pre-defined and parsed. Tests call runX directly on it.
func newCmdWithFlags(t *testing.T, flagSetup func(*cobra.Command), argv []string) *cobra.Command {
	t.Helper()
	cmd := &cobra.Command{Use: "test"}
	flagSetup(cmd)
	if err := cmd.ParseFlags(argv); err != nil {
		t.Fatalf("ParseFlags(%v): %v", argv, err)
	}
	return cmd
}

// testConfig returns a Config with sensible defaults for cmd tests.
func testConfig() *app.Config {
	return &app.Config{
		OrgName:   "acme",
		OrgPrefix: "ACME",
		AppName:   "billing",
		Env:       "dev",
		Filename:  ".espFile",
	}
}
