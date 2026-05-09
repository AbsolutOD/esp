# DI and Tests Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Introduce dependency injection at the AWS SDK and cobra subcommand boundaries, then add unit tests covering everything those seams unlock.

**Architecture:** Two PRs. PR 1 extracts an `ssmAPI` interface inside `internal/ssm`, collapses `Service.New`/`Init` into one constructor, and converts `cmd/*.go` from package-level globals to a constructor pattern threaded through an `App` holder. PR 2 adds unit tests against those seams using local fakes.

**Tech Stack:** Go 1.26, cobra, viper, AWS SDK v2 (`github.com/aws/aws-sdk-go-v2/service/ssm`), stdlib `testing`.

**Spec:** [`docs/superpowers/specs/2026-05-08-di-and-tests-design.md`](../specs/2026-05-08-di-and-tests-design.md)

**Commit-boundary note:** The spec's PR 1 Commits 2 and 3 (`refactor(ssm): collapse New+Init` and `refactor(client): adapt`) are merged into a single commit in this plan (Task 2). They cannot be separate commits while preserving the spec's `go build ./... clean` invariant — deleting `Service.Init()` while `client.New` still calls it would break the build between commits.

---

## File Inventory

**Modified by PR 1:**

| File | Task | Change |
|---|---|---|
| `cmd/get_test.go` | 1, 4 | T1: import path fix. T4: drop `withAppConfig` helper, switch to direct `*app.Config`. |
| `cmd/put_test.go` | 1, 4 | T1: import path fix. T4: drop `withOrgPrefix` helper, switch to direct `*app.Config`. |
| `cmd/list_test.go` | 1, 4 | T1: import path fix. T4: switch to direct `*app.Config`. |
| `internal/client/client_test.go` | 1 | T1: import path fix. |
| `internal/ssm/utils_test.go` | 1 | T1: import path fix. |
| `internal/ssm/ssm.go` | 2 | Add `ssmAPI` interface; `Service` holds `api ssmAPI`; `New() (*Service, error)` collapses `New`+`Init`; methods use `s.api.*`. |
| `internal/client/client.go` | 2 | `client.New` adapts to `(*ssm.Service, error)` return. |
| `cmd/root.go` | 3 | Add `App` struct; replace package globals; `Execute()` builds `App` + root; `newRootCmd(*App)` registers flags + `PersistentPreRunE`. |
| `cmd/get.go` | 3 | Convert to `newGetCmd(*App)`; extract `runGet(cmd, args, *EspClient, *app.Config)`. |
| `cmd/put.go` | 3 | Convert to `newPutCmd(*App)`; extract `runPut`. |
| `cmd/list.go` | 3 | Convert to `newListCmd(*App)`; extract `runList`. |
| `cmd/copy.go` | 3 | Convert to `newCopyCmd(*App)`; extract `runCopy`. |
| `cmd/move.go` | 3 | Convert to `newMoveCmd(*App)`; extract `runMove`. |
| `cmd/delete.go` | 3 | Convert to `newDeleteCmd(*App)`; extract `runDelete`. |
| `cmd/init.go` | 3 | Convert to `newInitCmd(*App)`. |
| `cmd/version.go` | 3 | Convert to `newVersionCmd()`; extract `runVersion`. |

**Created by PR 2:**

| File | Task | Purpose |
|---|---|---|
| `internal/ssm/ssm_test.go` | 7 | `TestService_*` with `fakeSSMAPI`. |
| `cmd/testutil_test.go` | 9 | `fakeBackend` + `newTestApp`. |
| `cmd/copy_test.go` | 9 | `TestRunCopy`. |
| `cmd/move_test.go` | 9 | `TestRunMove`. |
| `cmd/delete_test.go` | 9 | `TestRunDelete`. |
| `cmd/version_test.go` | 9 | `TestRunVersion`. |
| `cmd/root_test.go` | 10 | `TestPersistentPreRunE_*`. |

**Modified by PR 2:**

| File | Task | Change |
|---|---|---|
| `internal/client/client_test.go` | 8 | Add `TestEspClient_*` tests with local `fakeBackend`. |
| `cmd/get_test.go` | 9 | Add `TestRunGet`. |
| `cmd/put_test.go` | 9 | Add `TestRunPut`. |
| `cmd/list_test.go` | 9 | Add `TestRunList`. |

---

## Task 0: Worktree Setup

**Files:** none (creates `.worktrees/di-and-tests/`)

- [ ] **Step 1: Verify clean main**

```bash
git status
git branch --show-current
```

Expected: clean working tree, branch `main`. If not on main, `git checkout main`. If untracked changes, stash or commit before proceeding.

- [ ] **Step 2: Pull latest main**

```bash
git pull --ff-only origin main
```

- [ ] **Step 3: Create worktree on PR 1 branch**

```bash
git worktree add -b feature/di-refactor .worktrees/di-and-tests main
cd .worktrees/di-and-tests
```

Expected: new worktree at `.worktrees/di-and-tests` on branch `feature/di-refactor`.

- [ ] **Step 4: Verify baseline build (will fail)**

```bash
go build ./...
go test ./...
```

Expected: build clean. Tests fail with `no required module provides package github.com/pinpt/esp/...` in `cmd`, `internal/client`, `internal/ssm`. This is the breakage Task 1 fixes.

---

## Task 1: Fix stale pinpt/esp imports

**Commit message:** `chore: fix stale pinpt/esp imports in tests`

**Files:**
- Modify: `cmd/get_test.go`
- Modify: `cmd/put_test.go`
- Modify: `cmd/list_test.go`
- Modify: `internal/client/client_test.go`
- Modify: `internal/ssm/utils_test.go`

- [ ] **Step 1: Replace import paths**

Run from worktree root:

```bash
grep -rl 'github.com/pinpt/esp/' --include='*.go' .
```

Expected output (5 files):
```
cmd/get_test.go
cmd/put_test.go
cmd/list_test.go
internal/client/client_test.go
internal/ssm/utils_test.go
```

Replace in each:

```bash
find . -name '*_test.go' -not -path './.worktrees/*' -exec sed -i.bak 's|github.com/pinpt/esp/|github.com/AbsolutOD/esp/|g' {} \;
find . -name '*.bak' -not -path './.worktrees/*' -delete
```

- [ ] **Step 2: Confirm no remaining occurrences**

```bash
grep -r 'github.com/pinpt/esp/' --include='*.go' . || echo "OK: no matches"
```

Expected: `OK: no matches`.

- [ ] **Step 3: Run tests — expect green**

```bash
go build ./...
go vet ./...
go test ./...
```

Expected: all packages `ok` or `[no test files]`. No `setup failed`.

- [ ] **Step 4: Commit**

```bash
git add cmd/get_test.go cmd/put_test.go cmd/list_test.go internal/client/client_test.go internal/ssm/utils_test.go
git commit -m "chore: fix stale pinpt/esp imports in tests"
```

---

## Task 2: ssmAPI interface + collapse Service.New/Init + adapt client.New

**Commit message:** `refactor(ssm,client): introduce ssmAPI interface and collapse Service constructor`

**Files:**
- Modify: `internal/ssm/ssm.go`
- Modify: `internal/client/client.go`

- [ ] **Step 1: Replace `internal/ssm/ssm.go` body**

Open `internal/ssm/ssm.go` and replace the entire file contents with:

```go
package ssm

import (
	"context"

	"github.com/AbsolutOD/esp/internal/common"
	"github.com/AbsolutOD/esp/internal/utils"
	"github.com/aws/aws-sdk-go-v2/config"
	awsssm "github.com/aws/aws-sdk-go-v2/service/ssm"
)

type action string

const (
	Get     action = "get"
	GetMany action = "getMany"
	Save    action = "save"
	Delete  action = "delete"
)

// ssmAPI is the subset of the AWS SSM client used by Service.
// The concrete *awsssm.Client satisfies it; tests inject a fake.
// The four-method shape is required by NewGetParametersByPathPaginator,
// whose first argument is awsssm.GetParametersByPathAPIClient.
type ssmAPI interface {
	PutParameter(ctx context.Context, in *awsssm.PutParameterInput, optFns ...func(*awsssm.Options)) (*awsssm.PutParameterOutput, error)
	GetParameter(ctx context.Context, in *awsssm.GetParameterInput, optFns ...func(*awsssm.Options)) (*awsssm.GetParameterOutput, error)
	DeleteParameter(ctx context.Context, in *awsssm.DeleteParameterInput, optFns ...func(*awsssm.Options)) (*awsssm.DeleteParameterOutput, error)
	GetParametersByPath(ctx context.Context, in *awsssm.GetParametersByPathInput, optFns ...func(*awsssm.Options)) (*awsssm.GetParametersByPathOutput, error)
}

// Service is the SSM-backed implementation of client.Client.
type Service struct {
	api    ssmAPI
	Region string
}

// mapErr applies the per-action error mapper. If the mapper returns
// nil (the error wasn't a recognized AWS error type), return the raw
// error so the caller still sees the failure.
func mapErr(a action, err error) error {
	if mapped := checkSSMError(a, err); mapped != nil {
		return mapped
	}
	return err
}

// New builds a Service backed by a real AWS SSM client. Returns an
// error if AWS config loading fails.
func New() (*Service, error) {
	region := utils.GetEnv("AWS_REGION", "us-east-1")
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		return nil, err
	}
	return &Service{api: awsssm.NewFromConfig(cfg), Region: region}, nil
}

// Save a single param for a given path
func (s *Service) Save(p common.EspParamInput) (common.SaveOutput, error) {
	pi := &awsssm.PutParameterInput{
		Type:  selectType(p.Secure),
		Name:  &p.Name,
		Value: &p.Value,
	}
	param, err := s.api.PutParameter(context.Background(), pi)
	if err != nil {
		return common.SaveOutput{}, mapErr(Save, err)
	}
	return common.SaveOutput{Version: param.Version}, nil
}

// Delete a single param for a given path
func (s *Service) Delete(p common.DeleteInput) (string, error) {
	dpi := &awsssm.DeleteParameterInput{
		Name: &p.Name,
	}
	_, err := s.api.DeleteParameter(context.Background(), dpi)
	if err != nil {
		return "", mapErr(Delete, err)
	}
	return p.Name, nil
}

// GetOne gets a single param for a given path
func (s *Service) GetOne(p common.GetOneInput) (common.EspParam, error) {
	si := &awsssm.GetParameterInput{
		Name:           &p.Name,
		WithDecryption: &p.Decrypt,
	}
	resp, err := s.api.GetParameter(context.Background(), si)
	if err != nil {
		return common.EspParam{}, mapErr(Get, err)
	}
	return convertToEspParam(*resp.Parameter), nil
}

// GetMany recursively gets parameters from a given path
func (s *Service) GetMany(p common.ListParamInput) ([]common.EspParam, error) {
	si := &awsssm.GetParametersByPathInput{
		Path:           &p.Path,
		WithDecryption: &p.Decrypt,
		Recursive:      &p.Recursive,
	}
	paginator := awsssm.NewGetParametersByPathPaginator(s.api, si)

	var espParams []common.EspParam
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, mapErr(GetMany, err)
		}
		for _, v := range page.Parameters {
			espParams = append(espParams, convertToEspParam(v))
		}
	}
	return espParams, nil
}

// Copy method copies the given parameter to a new location
func (s *Service) Copy(cc common.CopyCommand) (common.SaveOutput, error) {
	sparam, err := s.GetOne(common.GetOneInput{Name: cc.Source, Decrypt: true})
	if err != nil {
		return common.SaveOutput{}, err
	}
	return s.Save(common.EspParamInput{
		Name:   cc.Destination,
		Secure: sparam.Secure,
		Value:  sparam.Value,
	})
}
```

Notes:
- `aws.String()`/`aws.Bool()`/`aws.ToString()` calls in the prior file used `github.com/aws/aws-sdk-go-v2/aws`. The new file uses `&p.Field` directly because the input fields are already `string`/`bool`. Drop the `aws` import.
- `Service.Cfg aws.Config` is removed (was unused outside `Init`).
- `Service.Init()` is gone.

- [ ] **Step 2: Update `internal/client/client.go::New`**

Edit `internal/client/client.go`. Replace the existing `New` function with:

```go
// New creates a new instance of the Client for esp
func New(c *app.Config) (*EspClient, error) {
	if c.Backend != "ssm" {
		return nil, fmt.Errorf("unsupported backend %q", c.Backend)
	}
	svc, err := ssm.New()
	if err != nil {
		return nil, err
	}
	return &EspClient{
		Backend: c.Backend,
		Client:  svc,
	}, nil
}
```

The change is two lines: `svc := ssm.New()` becomes `svc, err := ssm.New()` with an immediate error check, and the `svc.Init()` call is deleted.

- [ ] **Step 3: Build, vet, test**

```bash
go build ./...
go vet ./...
go test ./...
```

Expected: all green. The existing `TestNewUnsupportedBackend` in `internal/client/client_test.go` continues to pass (the unsupported-backend branch returns before reaching `ssm.New`).

- [ ] **Step 4: Commit**

```bash
git add internal/ssm/ssm.go internal/client/client.go
git commit -m "refactor(ssm,client): introduce ssmAPI interface and collapse Service constructor"
```

---

## Task 3: Convert cmd/ to constructor pattern with App holder

**Commit message:** `refactor(cmd): introduce App holder and convert subcommands to constructor pattern`

This is the largest task. All `cmd/*.go` files plus `cmd/root.go` change in lockstep — they cannot be split because the package-level globals being removed are referenced by every subcommand file. Single commit.

**Files (all modified):**
- `cmd/root.go`
- `cmd/get.go`
- `cmd/put.go`
- `cmd/list.go`
- `cmd/copy.go`
- `cmd/move.go`
- `cmd/delete.go`
- `cmd/init.go`
- `cmd/version.go`

The existing tests (`cmd/get_test.go`, `cmd/put_test.go`, `cmd/list_test.go`) **will break compile-wise** during this task — they reference the removed `esp` global and call helpers (`getParamPath`, `formatParamName`, `getPath`) with the old signatures. We fix them in Task 4. Build will be red after Task 3, green again after Task 4. **Therefore Tasks 3 and 4 are committed together** (the Task 4 commit is the verification gate).

Actually no — the spec lists Commits 4 and 5 separately. To honor the green-build invariant per commit, we'll do this differently: in Task 3 we update both the production code and the existing test files in a single commit. The split between "refactor production code" (Commit 4) and "update test signatures" (Commit 5) becomes a logical split inside one commit message. Let me revise — the plan keeps them as one commit:

**Revised commit message for combined Task 3 + Task 4:** `refactor(cmd): introduce App holder, convert subcommands to constructor pattern, update tests`

So Tasks 3 and 4 below are sequential steps but produce a **single commit** at the end of Task 4.

- [ ] **Step 3.1: Replace `cmd/root.go` body**

Open `cmd/root.go` and replace entire contents with:

```go
package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/AbsolutOD/esp/internal/app"
	"github.com/AbsolutOD/esp/internal/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// App holds the runtime dependencies threaded into every subcommand.
// Config is populated at construction; Client is populated by the root
// command's PersistentPreRunE after env-var validation succeeds.
type App struct {
	Config  *app.Config
	Client  *client.EspClient
	Verbose bool
}

// Execute is the entry point invoked by main. It builds the App,
// constructs the cobra tree around it, and runs.
func Execute() {
	a := &App{Config: app.New(false)}
	a.Config.Backend = "ssm"
	root := newRootCmd(a)
	root.AddCommand(
		newGetCmd(a),
		newPutCmd(a),
		newListCmd(a),
		newCopyCmd(a),
		newMoveCmd(a),
		newDeleteCmd(a),
		newInitCmd(a),
		newVersionCmd(),
	)
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

// newRootCmd builds the root cobra command bound to the given App.
func newRootCmd(a *App) *cobra.Command {
	root := &cobra.Command{
		Use:   "esp",
		Short: "A utility to browse and export SSM Parameter values into different formats.",
	}
	cobra.OnInitialize(func() { configureLogging(a.Verbose) })
	root.PersistentPreRunE = persistentPreRunE(a)
	root.PersistentFlags().StringVarP(&a.Config.Env, "env", "e", "", "Declare the env to work on.")
	root.PersistentFlags().StringVarP(&a.Config.Backend, "backend", "b", "ssm", "Set which backend to use.")
	root.PersistentFlags().BoolVar(&a.Verbose, "verbose", false, "Show more output")
	return root
}

// configureLogging sets the slog default handler. No I/O, no failure
// path. Runs via cobra.OnInitialize, which is skipped for --help.
func configureLogging(verbose bool) {
	level := slog.LevelWarn
	if verbose {
		level = slog.LevelInfo
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})))
}

// persistentPreRunE returns the closure that validates AWS env vars,
// constructs the backend client, and reads the .espFile. Cobra
// short-circuits before PersistentPreRunE for --help, so help still
// works without AWS credentials.
func persistentPreRunE(a *App) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		if _, ok := os.LookupEnv("AWS_DEFAULT_REGION"); !ok {
			return fmt.Errorf("AWS_DEFAULT_REGION environment variable is not set")
		}
		if _, ok := os.LookupEnv("AWS_PROFILE"); !ok {
			return fmt.Errorf("AWS_PROFILE environment variable is not set")
		}

		c, err := client.New(a.Config)
		if err != nil {
			return err
		}
		a.Client = c

		viper.SetConfigName(a.Config.Filename)
		viper.AddConfigPath(a.Config.Path)

		if err := viper.ReadInConfig(); err == nil {
			a.Config.IsEspProject = true
		}

		if a.Config.IsEspProject {
			if err := viper.Unmarshal(a.Config); err != nil {
				return fmt.Errorf("parsing %s: %w", a.Config.Filename, err)
			}
			if err := cmd.Root().MarkFlagRequired("env"); err != nil {
				return fmt.Errorf("marking --env required: %w", err)
			}
		}
		return nil
	}
}
```

Removed: package-level `var verbose bool`, `var esp *app.Config`, `var c *client.EspClient`, `init()`, the old `Execute()`, the old `persistentPreRun`.

- [ ] **Step 3.2: Replace `cmd/get.go` body**

```go
package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/AbsolutOD/esp/internal/app"
	"github.com/AbsolutOD/esp/internal/client"
	"github.com/AbsolutOD/esp/internal/common"
	"github.com/logrusorgru/aurora/v4"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func display(p common.EspParam, detail bool) {
	if detail {
		detailDisplay(p)
	} else {
		displayParam(p)
	}
}

func displayParam(p common.EspParam) {
	name := aurora.BrightYellow(p.Name)
	fmt.Printf("%s: %s\n", name, p.Value)
}

func detailDisplay(p common.EspParam) {
	data := [][]string{
		{aurora.BrightYellow("ID").String(), p.Id},
		{aurora.BrightYellow("Last_Modified").String(), p.LastModifiedDate.String()},
		{aurora.BrightYellow("Name").String(), p.Name},
		{aurora.BrightYellow("Type").String(), p.Type},
		{aurora.BrightYellow("Value").String(), p.Value},
		{aurora.BrightYellow("Version").String(), strconv.FormatInt(p.Version, 10)},
	}
	table := tablewriter.NewTable(os.Stdout)
	table.Header("Keys", "Value")
	if err := table.Bulk(data); err != nil {
		fmt.Fprintf(os.Stderr, "table render error: %v\n", err)
		return
	}
	if err := table.Render(); err != nil {
		fmt.Fprintf(os.Stderr, "table render error: %v\n", err)
	}
}

// getParamPath resolves an argument to a full SSM path. Leading "/"
// means the caller passed a literal path; everything else is routed
// through the project-aware GetAppParamPath.
func getParamPath(cfg *app.Config, p string) string {
	if strings.HasPrefix(p, "/") {
		return p
	}
	return cfg.GetAppParamPath(p)
}

func newGetCmd(a *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [path]",
		Short: "Query path for SSM",
		Long:  `Allows you to get a specific ssm parameter with an exact path or recursively get params.`,
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			return runGet(cmd, args, a.Client, a.Config)
		},
	}
	cmd.Flags().BoolP("details", "t", false, "Show all of the attributes of a parameter.")
	cmd.Flags().BoolP("decrypt", "d", false, "Decrypt SSM secure strings.")
	return cmd
}

func runGet(cmd *cobra.Command, args []string, c *client.EspClient, cfg *app.Config) error {
	decrypt, _ := cmd.Flags().GetBool("decrypt")
	details, _ := cmd.Flags().GetBool("details")
	param, err := c.GetParam(common.GetOneInput{
		Name:    getParamPath(cfg, args[0]),
		Decrypt: decrypt,
	})
	if err != nil {
		return err
	}
	display(param, details)
	return nil
}
```

- [ ] **Step 3.3: Replace `cmd/put.go` body**

```go
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/AbsolutOD/esp/internal/app"
	"github.com/AbsolutOD/esp/internal/client"
	"github.com/AbsolutOD/esp/internal/common"
	"github.com/spf13/cobra"
)

// formatParamName ensures the param ends up as a valid env-var
// identifier under the org prefix. Inputs already prefixed pass
// through; everything else is uppercased and hyphens become
// underscores.
func formatParamName(cfg *app.Config, n string) string {
	if strings.HasPrefix(n, cfg.OrgPrefix) {
		return n
	}
	normalized := strings.ReplaceAll(strings.ToUpper(n), "-", "_")
	return cfg.OrgPrefix + "_" + normalized
}

// getFullPath resolves a put-target name. Leading "/" is a literal
// path; otherwise normalize the name and route through GetAppParamPath.
func getFullPath(cfg *app.Config, n string) string {
	if strings.HasPrefix(n, "/") {
		return n
	}
	return cfg.GetAppParamPath(formatParamName(cfg, n))
}

func buildEspParamInputFromCmd(cfg *app.Config, cmd *cobra.Command) common.EspParamInput {
	name, _ := cmd.Flags().GetString("name")
	secure, _ := cmd.Flags().GetBool("secure")
	value, _ := cmd.Flags().GetString("value")
	return common.EspParamInput{
		Name:   getFullPath(cfg, name),
		Secure: secure,
		Value:  value,
	}
}

func newPutCmd(a *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "put",
		Aliases: []string{"add", "create"},
		Short:   "Creates an SSM parameter with the given value.",
		Long:    `Simple command to add values to SSM parameter store.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			return runPut(cmd, a.Client, a.Config)
		},
	}
	cmd.Flags().StringP("name", "n", "", "The name for your parameter.")
	cmd.Flags().StringP("value", "v", "", "The value to be stored in the SSM.")
	cmd.Flags().BoolP("secure", "s", false, "Sets the SSM parameter type to 'SecureString'.")
	if err := cobra.MarkFlagRequired(cmd.Flags(), "name"); err != nil {
		fmt.Fprintln(os.Stderr, "can't set flag --name as required")
	}
	if err := cobra.MarkFlagRequired(cmd.Flags(), "value"); err != nil {
		fmt.Fprintln(os.Stderr, "can't set flag --value as required")
	}
	return cmd
}

func runPut(cmd *cobra.Command, c *client.EspClient, cfg *app.Config) error {
	param := buildEspParamInputFromCmd(cfg, cmd)
	if _, err := c.Save(param); err != nil {
		return err
	}
	saved, err := c.GetParam(common.GetOneInput{Name: param.Name})
	if err != nil {
		return err
	}
	detailDisplay(saved)
	return nil
}
```

The previous `MarkFlagRequired` error handler used `fmt.Print` (stdout, no newline); the rewrite uses `fmt.Fprintln(os.Stderr, ...)` for proper stderr routing.

- [ ] **Step 3.4: Replace `cmd/list.go` body**

```go
package cmd

import (
	"fmt"

	"github.com/AbsolutOD/esp/internal/app"
	"github.com/AbsolutOD/esp/internal/client"
	"github.com/AbsolutOD/esp/internal/common"
	"github.com/logrusorgru/aurora/v4"
	"github.com/spf13/cobra"
)

func displayParams(ps []common.EspParam) {
	for _, p := range ps {
		name := aurora.BrightYellow(p.Name)
		fmt.Printf("%s: %s\n", name, p.Value)
	}
}

// getPath returns the SSM path to list. No args means "the project's
// base path"; one arg means "this exact path" (literal or short).
func getPath(cfg *app.Config, args []string) string {
	if len(args) == 0 {
		return cfg.GetAppPath()
	}
	return args[0]
}

func newListCmd(a *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list [path]",
		Aliases: []string{"ls"},
		Short:   "Recursively list a SSM path if given.",
		Long: `The list command gives you an easy way to recursively get all SSM parameters with a base path.
If you have a .espFile.yaml in the current directory this command will list all params under the project path.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			return runList(cmd, args, a.Client, a.Config)
		},
	}
	cmd.Flags().BoolP("decrypt", "d", false, "Decrypt SSM secure strings.")
	cmd.Flags().BoolP("path", "p", false, "Path to list parameters.")
	return cmd
}

func runList(cmd *cobra.Command, args []string, c *client.EspClient, cfg *app.Config) error {
	decrypt, _ := cmd.Flags().GetBool("decrypt")
	params, err := c.ListParams(common.ListParamInput{
		Path:      getPath(cfg, args),
		Decrypt:   decrypt,
		Recursive: true,
	})
	if err != nil {
		return err
	}
	displayParams(params)
	return nil
}
```

- [ ] **Step 3.5: Replace `cmd/copy.go` body**

```go
package cmd

import (
	"errors"

	"github.com/AbsolutOD/esp/internal/client"
	"github.com/AbsolutOD/esp/internal/common"
	"github.com/spf13/cobra"
)

func newCopyCmd(a *App) *cobra.Command {
	return &cobra.Command{
		Use:     "copy [OPTIONS] SRC_SSM_PATH DEST_SSM_PATH",
		Aliases: []string{"cp"},
		Short:   "Copy a SSM Param from its current path to a new SSM Path",
		Long:    "Copy SSM value from an existing path to a new path.\n",
		Args:    cobra.ExactArgs(2),
		Example: "esp cp /ssm/path/key /ssm/new/path/key",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			return runCopy(args, a.Client)
		},
	}
}

func runCopy(args []string, c *client.EspClient) error {
	if args[0] == "" {
		return errors.New("source can not be empty")
	}
	if args[1] == "" {
		return errors.New("destination can not be empty")
	}
	_, err := c.Copy(common.CopyCommand{
		Source:      args[0],
		Destination: args[1],
	})
	return err
}
```

- [ ] **Step 3.6: Replace `cmd/move.go` body**

```go
package cmd

import (
	"fmt"

	"github.com/AbsolutOD/esp/internal/client"
	"github.com/AbsolutOD/esp/internal/common"
	"github.com/logrusorgru/aurora/v4"
	"github.com/spf13/cobra"
)

func newMoveCmd(a *App) *cobra.Command {
	return &cobra.Command{
		Use:     "move [path]",
		Aliases: []string{"mv"},
		Short:   "move a parameter by path in SSM",
		Long:    `Allows you to move a specific ssm parameter with an exact path.`,
		Args:    cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			return runMove(args, a.Client)
		},
	}
}

func runMove(args []string, c *client.EspClient) error {
	p, err := c.Move(common.MoveCommand{Source: args[0], Destination: args[1]})
	if err != nil {
		return err
	}
	fmt.Printf("%s => %s\n", aurora.BrightYellow(p.Source), aurora.BrightYellow(p.Destination))
	return nil
}
```

- [ ] **Step 3.7: Replace `cmd/delete.go` body**

```go
package cmd

import (
	"fmt"

	"github.com/AbsolutOD/esp/internal/client"
	"github.com/AbsolutOD/esp/internal/common"
	"github.com/logrusorgru/aurora/v4"
	"github.com/spf13/cobra"
)

func newDeleteCmd(a *App) *cobra.Command {
	return &cobra.Command{
		Use:     "delete [path]",
		Aliases: []string{"rm"},
		Short:   "Delete a parameter by path in SSM",
		Long:    `Allows you to delete a specific ssm parameter with an exact path.`,
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			return runDelete(args, a.Client)
		},
	}
}

func runDelete(args []string, c *client.EspClient) error {
	name, err := c.Delete(common.DeleteInput{Name: args[0]})
	if err != nil {
		return err
	}
	fmt.Printf("Deleted: %s\n", aurora.BrightYellow(name))
	return nil
}
```

- [ ] **Step 3.8: Replace `cmd/init.go` body**

```go
package cmd

import "github.com/spf13/cobra"

func newInitCmd(a *App) *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initializes the current directory to be an ESP based application.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			a.Config.InitQuestions()
			return nil
		},
	}
}
```

- [ ] **Step 3.9: Replace `cmd/version.go` body**

```go
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Version of esp",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			return runVersion()
		},
	}
}

func runVersion() error {
	fmt.Println("ESP version 0.2.0")
	return nil
}
```

- [ ] **Step 3.10: Build — expect compile errors in test files only**

```bash
go build ./...
```

Expected: build clean. Tests next.

```bash
go test ./...
```

Expected: failures in `cmd` package — `undefined: esp` and similar in `cmd/get_test.go`, `cmd/put_test.go`, `cmd/list_test.go`. Fix in Task 4. **Do not commit yet.**

---

## Task 4: Update existing cmd tests for new signatures

Same commit as Task 3. After this task, build is fully green, then commit.

**Files:**
- `cmd/get_test.go`
- `cmd/put_test.go`
- `cmd/list_test.go`

- [ ] **Step 4.1: Replace `cmd/get_test.go` body**

```go
package cmd

import (
	"testing"

	"github.com/AbsolutOD/esp/internal/app"
)

// TestGetParamPath pins getParamPath's rule: a leading "/" means the
// caller already wrote a literal SSM path, so it passes through; any
// other shape is routed through cfg.GetAppParamPath.
func TestGetParamPath(t *testing.T) {
	cfg := &app.Config{
		OrgName: "acme",
		Env:     "dev",
		AppName: "billing",
	}

	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "leading slash → literal path returned unchanged",
			in:   "/acme/dev/billing/DB_URL",
			want: "/acme/dev/billing/DB_URL",
		},
		{
			name: "absolute path with different segments → unchanged",
			in:   "/some/other/place",
			want: "/some/other/place",
		},
		{
			name: "short name → resolved through GetAppParamPath",
			in:   "DB_URL",
			want: "/acme/dev/billing/DB_URL",
		},
		{
			name: "lowercase short name passes through verbatim (no upcasing here)",
			in:   "secret",
			want: "/acme/dev/billing/secret",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := getParamPath(cfg, tc.in)
			if got != tc.want {
				t.Errorf("getParamPath(cfg, %q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}
```

- [ ] **Step 4.2: Replace `cmd/put_test.go` body**

```go
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
```

- [ ] **Step 4.3: Replace `cmd/list_test.go` body**

```go
package cmd

import (
	"testing"

	"github.com/AbsolutOD/esp/internal/app"
)

func TestGetPathWithFullPath(t *testing.T) {
	cfg := &app.Config{}
	testPath := "/corpa/dev/foo_app/"
	if got := getPath(cfg, []string{testPath}); got != testPath {
		t.Errorf("want: %s | got %s", testPath, got)
	}
}

func TestGetPathEnvVarName(t *testing.T) {
	cfg := &app.Config{}
	envVar := "TEST_VAR"
	if got := getPath(cfg, []string{envVar}); got != envVar {
		t.Errorf("want: %s | got %s", envVar, got)
	}
}

// TestGetPathRelative pins getPath's empty-args branch.
func TestGetPathRelative(t *testing.T) {
	cfg := &app.Config{
		OrgName: "acme",
		Env:     "dev",
		AppName: "billing",
	}
	got := getPath(cfg, nil)
	want := "/acme/dev/billing/"
	if got != want {
		t.Errorf("getPath(cfg, nil) = %q, want %q", got, want)
	}
}
```

- [ ] **Step 4.4: Verify all green**

```bash
go build ./...
go vet ./...
go test ./...
```

Expected: all packages `ok` or `[no test files]`.

- [ ] **Step 4.5: Smoke-test the binary**

```bash
go build -o /tmp/esp .
AWS_DEFAULT_REGION=us-east-1 AWS_PROFILE=default /tmp/esp --help
```

Expected: standard cobra help output, no AWS errors (cobra short-circuits before `PersistentPreRunE` for help).

```bash
unset_pair() { unset AWS_DEFAULT_REGION; unset AWS_PROFILE; }
unset_pair
/tmp/esp --help
```

Expected: still works (help short-circuits).

```bash
/tmp/esp version 2>&1 || true
```

Expected: `Error: AWS_DEFAULT_REGION environment variable is not set`, exit 1.

```bash
rm /tmp/esp
```

- [ ] **Step 4.6: Commit (combined Task 3 + 4)**

```bash
git add cmd/
git commit -m "refactor(cmd): introduce App holder, convert subcommands to constructor pattern, update tests"
```

---

## Task 5: Open PR 1

- [ ] **Step 1: Push branch**

```bash
git push -u origin feature/di-refactor
```

- [ ] **Step 2: Open PR**

```bash
gh pr create --title "PR 1: DI refactor (ssmAPI interface + cmd App holder)" --body "$(cat <<'EOF'
## Summary

PR 1 of 2 implementing the DI and tests spec at \`docs/superpowers/specs/2026-05-08-di-and-tests-design.md\`.

This PR is purely structural — no behavior changes. It introduces the dependency-injection seams that PR 2 will exploit for unit tests.

- **Commit 1:** restore green test suite by fixing stale \`pinpt/esp\` imports left over from the rebrand (PR #29).
- **Commit 2:** extract unexported \`ssmAPI\` interface inside \`internal/ssm\`; collapse \`Service.New\`/\`Init\` into one constructor; \`client.New\` adapts.
- **Commit 3:** introduce \`cmd.App\` holder; convert all subcommands to \`newXCmd(*App)\` constructor pattern with extracted \`runX\` functions; remove the package-level \`esp\` and \`c\` globals; update existing tests for new pure-helper signatures.

PR 2 will add unit tests against the new seams.

## Test plan

- [ ] \`go build ./...\` clean
- [ ] \`go vet ./...\` clean
- [ ] \`go test ./...\` green (existing tests preserved)
- [ ] \`esp --help\` works without AWS env vars set
- [ ] \`esp version\` with missing \`AWS_DEFAULT_REGION\` returns exit 1 with the env-var error
- [ ] \`esp\` invocations behave identically to pre-refactor on the happy path
EOF
)"
```

- [ ] **Step 3: Note the PR URL** for follow-up.

After merge, return to this plan and continue with Task 6.

---

## Task 6: Rebase worktree onto fresh main for PR 2

- [ ] **Step 1: From the worktree, fetch and check out new branch**

```bash
git fetch origin
git checkout -b feature/di-tests origin/main
```

Expected: branch `feature/di-tests` based on the latest `main` (which now contains the merged PR 1).

- [ ] **Step 2: Verify baseline green**

```bash
go build ./... && go vet ./... && go test ./...
```

Expected: all green.

---

## Task 7: SSM Service tests via fakeSSMAPI

**Commit message:** `test(ssm): cover Service methods via fake ssmAPI`

**Files:**
- Create: `internal/ssm/ssm_test.go`

- [ ] **Step 1: Create the test file**

```go
package ssm

import (
	"context"
	"errors"
	"testing"

	"github.com/AbsolutOD/esp/internal/common"
	awsssm "github.com/aws/aws-sdk-go-v2/service/ssm"
	ssmtypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

// fakeSSMAPI implements ssmAPI by returning canned responses.
// Each method records its last input for assertion.
type fakeSSMAPI struct {
	putIn    *awsssm.PutParameterInput
	putOut   *awsssm.PutParameterOutput
	putErr   error
	getIn    *awsssm.GetParameterInput
	getOut   *awsssm.GetParameterOutput
	getErr   error
	delIn    *awsssm.DeleteParameterInput
	delOut   *awsssm.DeleteParameterOutput
	delErr   error
	pathIn   *awsssm.GetParametersByPathInput
	pathOuts []*awsssm.GetParametersByPathOutput
	pathErrs []error
	pathIdx  int
}

func (f *fakeSSMAPI) PutParameter(_ context.Context, in *awsssm.PutParameterInput, _ ...func(*awsssm.Options)) (*awsssm.PutParameterOutput, error) {
	f.putIn = in
	return f.putOut, f.putErr
}

func (f *fakeSSMAPI) GetParameter(_ context.Context, in *awsssm.GetParameterInput, _ ...func(*awsssm.Options)) (*awsssm.GetParameterOutput, error) {
	f.getIn = in
	return f.getOut, f.getErr
}

func (f *fakeSSMAPI) DeleteParameter(_ context.Context, in *awsssm.DeleteParameterInput, _ ...func(*awsssm.Options)) (*awsssm.DeleteParameterOutput, error) {
	f.delIn = in
	return f.delOut, f.delErr
}

func (f *fakeSSMAPI) GetParametersByPath(_ context.Context, in *awsssm.GetParametersByPathInput, _ ...func(*awsssm.Options)) (*awsssm.GetParametersByPathOutput, error) {
	f.pathIn = in
	i := f.pathIdx
	f.pathIdx++
	if i >= len(f.pathOuts) {
		return &awsssm.GetParametersByPathOutput{}, nil
	}
	var err error
	if i < len(f.pathErrs) {
		err = f.pathErrs[i]
	}
	return f.pathOuts[i], err
}

func newServiceWithAPI(api ssmAPI) *Service {
	return &Service{api: api, Region: "us-east-1"}
}

func strPtr(s string) *string { return &s }

// --- Save ---

func TestService_Save_Success(t *testing.T) {
	fake := &fakeSSMAPI{putOut: &awsssm.PutParameterOutput{Version: 7}}
	s := newServiceWithAPI(fake)

	out, err := s.Save(common.EspParamInput{Name: "/x/y", Value: "v", Secure: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Version != 7 {
		t.Errorf("Version = %d, want 7", out.Version)
	}
	if fake.putIn.Type != ssmtypes.ParameterTypeString {
		t.Errorf("Type = %q, want String", fake.putIn.Type)
	}
	if got := *fake.putIn.Name; got != "/x/y" {
		t.Errorf("Name = %q, want /x/y", got)
	}
	if got := *fake.putIn.Value; got != "v" {
		t.Errorf("Value = %q, want v", got)
	}
}

func TestService_Save_SecureFlagSetsSecureString(t *testing.T) {
	fake := &fakeSSMAPI{putOut: &awsssm.PutParameterOutput{Version: 1}}
	s := newServiceWithAPI(fake)
	if _, err := s.Save(common.EspParamInput{Name: "n", Value: "v", Secure: true}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fake.putIn.Type != ssmtypes.ParameterTypeSecureString {
		t.Errorf("Type = %q, want SecureString", fake.putIn.Type)
	}
}

func TestService_Save_AlreadyExistsErrorPropagates(t *testing.T) {
	awsErr := &ssmtypes.ParameterAlreadyExists{}
	fake := &fakeSSMAPI{putErr: awsErr}
	s := newServiceWithAPI(fake)

	_, err := s.Save(common.EspParamInput{Name: "n", Value: "v"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var typed *ssmtypes.ParameterAlreadyExists
	if !errors.As(err, &typed) {
		t.Errorf("error %v is not *ParameterAlreadyExists", err)
	}
}

// --- GetOne ---

func TestService_GetOne_Success(t *testing.T) {
	fake := &fakeSSMAPI{getOut: &awsssm.GetParameterOutput{
		Parameter: &ssmtypes.Parameter{
			Name:    strPtr("/x/y"),
			Value:   strPtr("v"),
			Type:    ssmtypes.ParameterTypeString,
			Version: 1,
		},
	}}
	s := newServiceWithAPI(fake)

	got, err := s.GetOne(common.GetOneInput{Name: "/x/y", Decrypt: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Name != "/x/y" {
		t.Errorf("Name = %q, want /x/y", got.Name)
	}
	if got.Value != "v" {
		t.Errorf("Value = %q, want v", got.Value)
	}
	if !*fake.getIn.WithDecryption {
		t.Error("WithDecryption = false, want true")
	}
}

func TestService_GetOne_NotFoundErrorPropagates(t *testing.T) {
	fake := &fakeSSMAPI{getErr: &ssmtypes.ParameterNotFound{}}
	s := newServiceWithAPI(fake)

	_, err := s.GetOne(common.GetOneInput{Name: "/x/y"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var typed *ssmtypes.ParameterNotFound
	if !errors.As(err, &typed) {
		t.Errorf("error %v is not *ParameterNotFound", err)
	}
}

// --- GetMany ---

func TestService_GetMany_MultiplePages(t *testing.T) {
	fake := &fakeSSMAPI{
		pathOuts: []*awsssm.GetParametersByPathOutput{
			{
				Parameters: []ssmtypes.Parameter{
					{Name: strPtr("/a/1"), Value: strPtr("v1"), Type: ssmtypes.ParameterTypeString},
				},
				NextToken: strPtr("page2"),
			},
			{
				Parameters: []ssmtypes.Parameter{
					{Name: strPtr("/a/2"), Value: strPtr("v2"), Type: ssmtypes.ParameterTypeString},
				},
			},
		},
	}
	s := newServiceWithAPI(fake)

	params, err := s.GetMany(common.ListParamInput{Path: "/a/", Recursive: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(params) != 2 {
		t.Fatalf("len(params) = %d, want 2", len(params))
	}
	if params[0].Name != "/a/1" || params[1].Name != "/a/2" {
		t.Errorf("got names %q,%q; want /a/1,/a/2", params[0].Name, params[1].Name)
	}
}

func TestService_GetMany_ErrorMidIteration(t *testing.T) {
	fake := &fakeSSMAPI{
		pathOuts: []*awsssm.GetParametersByPathOutput{
			{
				Parameters: []ssmtypes.Parameter{{Name: strPtr("/a/1"), Type: ssmtypes.ParameterTypeString}},
				NextToken:  strPtr("page2"),
			},
			{},
		},
		pathErrs: []error{nil, &ssmtypes.InternalServerError{}},
	}
	s := newServiceWithAPI(fake)

	_, err := s.GetMany(common.ListParamInput{Path: "/a/"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var typed *ssmtypes.InternalServerError
	if !errors.As(err, &typed) {
		t.Errorf("error %v is not *InternalServerError", err)
	}
}

func TestService_GetMany_EmptyPages(t *testing.T) {
	fake := &fakeSSMAPI{
		pathOuts: []*awsssm.GetParametersByPathOutput{{}},
	}
	s := newServiceWithAPI(fake)
	params, err := s.GetMany(common.ListParamInput{Path: "/x/"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(params) != 0 {
		t.Errorf("got %d params, want 0", len(params))
	}
}

// --- Delete ---

func TestService_Delete_Success(t *testing.T) {
	fake := &fakeSSMAPI{delOut: &awsssm.DeleteParameterOutput{}}
	s := newServiceWithAPI(fake)

	name, err := s.Delete(common.DeleteInput{Name: "/x/y"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if name != "/x/y" {
		t.Errorf("name = %q, want /x/y", name)
	}
	if got := *fake.delIn.Name; got != "/x/y" {
		t.Errorf("input name = %q, want /x/y", got)
	}
}

func TestService_Delete_NotFoundErrorPropagates(t *testing.T) {
	fake := &fakeSSMAPI{delErr: &ssmtypes.ParameterNotFound{}}
	s := newServiceWithAPI(fake)

	_, err := s.Delete(common.DeleteInput{Name: "/x/y"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var typed *ssmtypes.ParameterNotFound
	if !errors.As(err, &typed) {
		t.Errorf("error %v is not *ParameterNotFound", err)
	}
}

// --- Copy ---

func TestService_Copy_Success(t *testing.T) {
	// First call: GetOne returns the source param.
	// Second call: PutParameter.
	fake := &fakeSSMAPI{
		getOut: &awsssm.GetParameterOutput{
			Parameter: &ssmtypes.Parameter{
				Name:  strPtr("/src"),
				Value: strPtr("vv"),
				Type:  ssmtypes.ParameterTypeSecureString,
			},
		},
		putOut: &awsssm.PutParameterOutput{Version: 9},
	}
	s := newServiceWithAPI(fake)

	out, err := s.Copy(common.CopyCommand{Source: "/src", Destination: "/dest"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Version != 9 {
		t.Errorf("Version = %d, want 9", out.Version)
	}
	if fake.putIn.Type != ssmtypes.ParameterTypeSecureString {
		t.Errorf("dest Type = %q, want SecureString (carried from source)", fake.putIn.Type)
	}
	if got := *fake.putIn.Name; got != "/dest" {
		t.Errorf("dest Name = %q, want /dest", got)
	}
}

func TestService_Copy_GetOneFailsSaveNotCalled(t *testing.T) {
	fake := &fakeSSMAPI{getErr: &ssmtypes.ParameterNotFound{}}
	s := newServiceWithAPI(fake)

	_, err := s.Copy(common.CopyCommand{Source: "/src", Destination: "/dest"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if fake.putIn != nil {
		t.Error("PutParameter was called despite GetOne failure")
	}
}

func TestService_Copy_SaveFails(t *testing.T) {
	fake := &fakeSSMAPI{
		getOut: &awsssm.GetParameterOutput{
			Parameter: &ssmtypes.Parameter{
				Name: strPtr("/src"), Value: strPtr("v"), Type: ssmtypes.ParameterTypeString,
			},
		},
		putErr: &ssmtypes.ParameterLimitExceeded{},
	}
	s := newServiceWithAPI(fake)

	_, err := s.Copy(common.CopyCommand{Source: "/src", Destination: "/dest"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var typed *ssmtypes.ParameterLimitExceeded
	if !errors.As(err, &typed) {
		t.Errorf("error %v is not *ParameterLimitExceeded", err)
	}
}

```

- [ ] **Step 2: Run tests**

```bash
go test ./internal/ssm/ -v -run TestService_
```

Expected: every `TestService_*` PASS.

- [ ] **Step 3: Run full suite**

```bash
go build ./... && go vet ./... && go test ./...
```

Expected: all green.

- [ ] **Step 4: Commit**

```bash
git add internal/ssm/ssm_test.go
git commit -m "test(ssm): cover Service methods via fake ssmAPI"
```

---

## Task 8: EspClient wrapper tests

**Commit message:** `test(client): cover EspClient wrapper methods`

**Files:**
- Modify: `internal/client/client_test.go`

- [ ] **Step 1: Append to `internal/client/client_test.go`**

After the existing `TestNewUnsupportedBackend` function, append:

```go
// fakeBackend implements Client by recording calls and returning canned responses.
type fakeBackend struct {
	saveIn   common.EspParamInput
	saveOut  common.SaveOutput
	saveErr  error
	saveCalls int

	getIn   common.GetOneInput
	getOut  common.EspParam
	getErr  error
	getCalls int

	manyIn   common.ListParamInput
	manyOut  []common.EspParam
	manyErr  error

	copyIn   common.CopyCommand
	copyOut  common.SaveOutput
	copyErr  error
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
```

Add the import:

```go
import (
	"errors"  // add if not already imported
	"strings"
	"testing"

	"github.com/AbsolutOD/esp/internal/app"
	"github.com/AbsolutOD/esp/internal/common"  // add
)
```

- [ ] **Step 2: Add `TestEspClient_GetParam`**

```go
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
```

- [ ] **Step 3: Add `TestEspClient_ListParams`**

```go
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
```

- [ ] **Step 4: Add `TestEspClient_Save`**

```go
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
```

- [ ] **Step 5: Add `TestEspClient_Delete`**

```go
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
```

- [ ] **Step 6: Add `TestEspClient_Copy`**

```go
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
```

- [ ] **Step 7: Add `TestEspClient_Move`**

```go
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
```

- [ ] **Step 8: Run tests**

```bash
go test ./internal/client/ -v
```

Expected: all `TestEspClient_*` and `TestNewUnsupportedBackend` PASS.

- [ ] **Step 9: Run full suite**

```bash
go build ./... && go vet ./... && go test ./...
```

Expected: all green.

- [ ] **Step 10: Commit**

```bash
git add internal/client/client_test.go
git commit -m "test(client): cover EspClient wrapper methods"
```

---

## Task 9: cmd runX tests

**Commit message:** `test(cmd): cover runX functions via fake EspClient`

**Files:**
- Create: `cmd/testutil_test.go`
- Modify: `cmd/get_test.go`
- Modify: `cmd/put_test.go`
- Modify: `cmd/list_test.go`
- Create: `cmd/copy_test.go`
- Create: `cmd/move_test.go`
- Create: `cmd/delete_test.go`
- Create: `cmd/version_test.go`

- [ ] **Step 1: Create `cmd/testutil_test.go`**

```go
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
```

- [ ] **Step 2: Add `TestRunGet` cases to `cmd/get_test.go`**

Update the import block at the top of `cmd/get_test.go` so it includes all of:

```go
import (
	"errors"
	"testing"

	"github.com/AbsolutOD/esp/internal/app"
	"github.com/AbsolutOD/esp/internal/common"
	"github.com/spf13/cobra"
)
```

Then append the following test functions to the end of the file:

```go
func TestRunGet_LiteralPath(t *testing.T) {
	fake := &fakeBackend{getOut: common.EspParam{Name: "/x/y", Value: "v"}}
	c := newTestEspClient(fake)
	cfg := testConfig()
	cmd := newCmdWithFlags(t, func(c *cobra.Command) {
		c.Flags().BoolP("decrypt", "d", false, "")
		c.Flags().BoolP("details", "t", false, "")
	}, []string{})

	if err := runGet(cmd, []string{"/x/y"}, c, cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fake.getIn.Name != "/x/y" {
		t.Errorf("GetParam called with %q, want /x/y", fake.getIn.Name)
	}
}

func TestRunGet_RelativePathResolved(t *testing.T) {
	fake := &fakeBackend{getOut: common.EspParam{Name: "/acme/dev/billing/DB"}}
	c := newTestEspClient(fake)
	cfg := testConfig()
	cmd := newCmdWithFlags(t, func(c *cobra.Command) {
		c.Flags().BoolP("decrypt", "d", false, "")
		c.Flags().BoolP("details", "t", false, "")
	}, []string{})

	if err := runGet(cmd, []string{"DB"}, c, cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fake.getIn.Name != "/acme/dev/billing/DB" {
		t.Errorf("got %q, want /acme/dev/billing/DB", fake.getIn.Name)
	}
}

func TestRunGet_DecryptFlagPropagates(t *testing.T) {
	fake := &fakeBackend{getOut: common.EspParam{Name: "/x"}}
	c := newTestEspClient(fake)
	cmd := newCmdWithFlags(t, func(c *cobra.Command) {
		c.Flags().BoolP("decrypt", "d", false, "")
		c.Flags().BoolP("details", "t", false, "")
	}, []string{"--decrypt"})

	if err := runGet(cmd, []string{"/x"}, c, testConfig()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !fake.getIn.Decrypt {
		t.Error("Decrypt = false, want true")
	}
}

func TestRunGet_GetParamErrorSurfaces(t *testing.T) {
	fake := &fakeBackend{getErr: errors.New("nope")}
	c := newTestEspClient(fake)
	cmd := newCmdWithFlags(t, func(c *cobra.Command) {
		c.Flags().BoolP("decrypt", "d", false, "")
		c.Flags().BoolP("details", "t", false, "")
	}, []string{})

	err := runGet(cmd, []string{"/x"}, c, testConfig())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
```

- [ ] **Step 3: Add `TestRunPut` cases to `cmd/put_test.go`**

Update the import block to include all of:

```go
import (
	"errors"
	"testing"

	"github.com/AbsolutOD/esp/internal/app"
	"github.com/AbsolutOD/esp/internal/common"
	"github.com/spf13/cobra"
)
```

Then append the following test functions to the end of the file:

```go
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
```

- [ ] **Step 4: Add `TestRunList` cases to `cmd/list_test.go`**

Update the import block to include all of:

```go
import (
	"errors"
	"testing"

	"github.com/AbsolutOD/esp/internal/app"
	"github.com/AbsolutOD/esp/internal/common"
	"github.com/spf13/cobra"
)
```

Then append the following test functions to the end of the file:

```go
func TestRunList_NoArgsUsesAppPath(t *testing.T) {
	fake := &fakeBackend{manyOut: []common.EspParam{{Name: "/acme/dev/billing/X", Value: "v"}}}
	c := newTestEspClient(fake)
	cfg := testConfig()
	cmd := newCmdWithFlags(t, func(c *cobra.Command) {
		c.Flags().BoolP("decrypt", "d", false, "")
		c.Flags().BoolP("path", "p", false, "")
	}, []string{})

	if err := runList(cmd, nil, c, cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fake.manyIn.Path != "/acme/dev/billing/" {
		t.Errorf("Path = %q, want /acme/dev/billing/", fake.manyIn.Path)
	}
	if !fake.manyIn.Recursive {
		t.Error("Recursive = false, want true")
	}
}

func TestRunList_WithArgUsesArg(t *testing.T) {
	fake := &fakeBackend{}
	c := newTestEspClient(fake)
	cmd := newCmdWithFlags(t, func(c *cobra.Command) {
		c.Flags().BoolP("decrypt", "d", false, "")
		c.Flags().BoolP("path", "p", false, "")
	}, []string{})

	if err := runList(cmd, []string{"/elsewhere/"}, c, testConfig()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fake.manyIn.Path != "/elsewhere/" {
		t.Errorf("Path = %q, want /elsewhere/", fake.manyIn.Path)
	}
}

func TestRunList_DecryptFlagPropagates(t *testing.T) {
	fake := &fakeBackend{}
	c := newTestEspClient(fake)
	cmd := newCmdWithFlags(t, func(c *cobra.Command) {
		c.Flags().BoolP("decrypt", "d", false, "")
		c.Flags().BoolP("path", "p", false, "")
	}, []string{"--decrypt"})

	if err := runList(cmd, []string{"/x/"}, c, testConfig()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !fake.manyIn.Decrypt {
		t.Error("Decrypt = false, want true")
	}
}

func TestRunList_ErrorSurfaces(t *testing.T) {
	fake := &fakeBackend{manyErr: errors.New("api err")}
	c := newTestEspClient(fake)
	cmd := newCmdWithFlags(t, func(c *cobra.Command) {
		c.Flags().BoolP("decrypt", "d", false, "")
		c.Flags().BoolP("path", "p", false, "")
	}, []string{})

	if err := runList(cmd, []string{"/x/"}, c, testConfig()); err == nil {
		t.Fatal("expected error, got nil")
	}
}
```

- [ ] **Step 5: Create `cmd/copy_test.go`**

```go
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
```

- [ ] **Step 6: Create `cmd/move_test.go`**

```go
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
```

- [ ] **Step 7: Create `cmd/delete_test.go`**

```go
package cmd

import (
	"errors"
	"testing"
)

func TestRunDelete_HappyPath(t *testing.T) {
	fake := &fakeBackend{delOut: "/x"}
	if err := runDelete([]string{"/x"}, newTestEspClient(fake)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fake.delIn.Name != "/x" {
		t.Errorf("Delete name = %q, want /x", fake.delIn.Name)
	}
}

func TestRunDelete_BackendError(t *testing.T) {
	fake := &fakeBackend{delErr: errors.New("nope")}
	if err := runDelete([]string{"/x"}, newTestEspClient(fake)); err == nil {
		t.Fatal("expected error, got nil")
	}
}
```

- [ ] **Step 8: Create `cmd/version_test.go`**

```go
package cmd

import "testing"

func TestRunVersion(t *testing.T) {
	if err := runVersion(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
```

- [ ] **Step 9: Run cmd tests**

```bash
go test ./cmd/ -v
```

Expected: all `Test*` PASS, including the existing `TestEspName`, `TestGetParamPath`, `TestGetPath*`, plus all new `TestRun*` tests.

- [ ] **Step 10: Run full suite**

```bash
go build ./... && go vet ./... && go test ./...
```

Expected: all green.

- [ ] **Step 11: Commit**

```bash
git add cmd/testutil_test.go cmd/get_test.go cmd/put_test.go cmd/list_test.go cmd/copy_test.go cmd/move_test.go cmd/delete_test.go cmd/version_test.go
git commit -m "test(cmd): cover runX functions via fake EspClient"
```

---

## Task 10: persistentPreRunE env-var tests

**Commit message:** `test(cmd): cover persistentPreRunE env-var validation`

**Files:**
- Create: `cmd/root_test.go`

- [ ] **Step 1: Create `cmd/root_test.go`**

```go
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
```

- [ ] **Step 2: Run tests**

```bash
go test ./cmd/ -v -run TestPersistentPreRunE
```

Expected: all three tests PASS.

- [ ] **Step 3: Run full suite**

```bash
go build ./... && go vet ./... && go test ./...
```

Expected: all green.

- [ ] **Step 4: Commit**

```bash
git add cmd/root_test.go
git commit -m "test(cmd): cover persistentPreRunE env-var validation"
```

---

## Task 11: Open PR 2

- [ ] **Step 1: Push branch**

```bash
git push -u origin feature/di-tests
```

- [ ] **Step 2: Open PR**

```bash
gh pr create --title "PR 2: DI tests (Service, EspClient, runX, persistentPreRunE)" --body "$(cat <<'EOF'
## Summary

PR 2 of 2 implementing the DI and tests spec at \`docs/superpowers/specs/2026-05-08-di-and-tests-design.md\`.

This PR is purely additive — new tests against the seams introduced in PR 1. No production code changes.

- **Commit 1:** \`Service.{Save,GetOne,GetMany,Delete,Copy}\` tests via \`fakeSSMAPI\`, including AWS error mapping.
- **Commit 2:** \`EspClient.{GetParam,ListParams,Save,Delete,Copy,Move}\` tests via \`fakeBackend\`, including composition order in Copy and Move.
- **Commit 3:** \`runGet/Put/List/Copy/Move/Delete/Version\` tests via fake-wrapped \`EspClient\`.
- **Commit 4:** \`persistentPreRunE\` env-var validation tests.

Coverage gaps explicitly accepted: \`ssm.New()\` (LoadDefaultConfig — environment-coupled), \`InitQuestions\` (interactive survey), \`display*\` print helpers, \`configureLogging\`, \`main.go\`. See spec.

## Test plan

- [ ] \`go build ./...\` clean
- [ ] \`go vet ./...\` clean
- [ ] \`go test ./...\` green; new tests pass
EOF
)"
```

- [ ] **Step 3: Note the PR URL** for follow-up.

---

## Self-Review Checklist

After completing all tasks, verify:

- [ ] `go build ./...` — clean.
- [ ] `go vet ./...` — clean.
- [ ] `go test ./...` — every package green.
- [ ] `esp --help` works without AWS env vars.
- [ ] `esp version` with missing `AWS_DEFAULT_REGION` exits 1 with the env-var error.
- [ ] All five spec PR 1 commits represented (rebrand-fix is its own; ssm+client are merged into Task 2's commit; cmd refactor + test signature update are merged into Task 3+4's commit). Verify the commit log:
  ```bash
  git log --oneline main..feature/di-refactor
  ```
  Expected three commits.
- [ ] All four spec PR 2 commits represented:
  ```bash
  git log --oneline main..feature/di-tests
  ```
  Expected four commits.
