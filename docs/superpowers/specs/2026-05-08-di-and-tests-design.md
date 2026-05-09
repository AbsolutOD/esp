# Dependency Injection and Tests Design

**Date:** 2026-05-08
**Status:** Approved
**Scope:** Introduce dependency injection at the AWS SDK boundary (`internal/ssm`) and the cobra subcommand boundary (`cmd/`), then add unit tests covering everything those seams unlock. Two sequential PRs.
**Predecessor:** [`2026-05-05-cleanup-and-tests-design.md`](2026-05-05-cleanup-and-tests-design.md). The prior spec listed DI as a deferred non-goal ("option C territory"). This design picks it up.

## Goals

- Introduce DI at two layers so behaviorally meaningful code becomes unit-testable: the AWS SDK boundary inside `internal/ssm`, and the cobra subcommand boundary inside `cmd/`.
- Add unit tests covering every reachable function gated by those two seams. Combined with the pure-logic tests already merged (PRs #26–#28), this should cover essentially all production behavior except `config.LoadDefaultConfig` and `survey.Ask`.
- Restore `go test ./...` to green. Five test files still reference the pre-rebrand `github.com/pinpt/esp/...` import path; they need to flip to `github.com/AbsolutOD/esp/...` before any other work.

## Non-Goals

- No abstraction of `survey.Ask` or `os.WriteFile` in `internal/app/init.go`. That was Approach C in the brainstorm; deferred.
- No new backends, no CLI feature work.
- No CI workflow changes.
- No integration tests against LocalStack or real AWS.
- No coverage-percentage target. Bar is: every function reachable from a subcommand has a happy-path test plus its meaningful error paths.
- No removal of pure helpers in `cmd/` (`formatParamName`, `getFullPath`, `getParamPath`, `getPath`) — they keep their existing tests, updated for the new explicit `*app.Config` parameter.

## PR Plan

Two PRs, branched off `main`, sharing one worktree at `.worktrees/di-and-tests`:

1. **PR 1 — `feature/di-refactor`** — restore broken tests, extract DI seams, refactor `cmd/` to constructor pattern. No new tests beyond restoring the existing suite. No behavior change.
2. **PR 2 — `feature/di-tests`** — new unit tests against the seams from PR 1.

PR 2 is started after PR 1 merges and `main` updates. The worktree is reused: switch the branch to `feature/di-tests` rebased onto fresh `main`.

## PR 1 — DI Refactor

### Commit 1 — Restore green test suite

`chore: fix stale pinpt/esp imports in tests`

The rebrand (PR #29) flipped the module path to `github.com/AbsolutOD/esp` but left five test files referencing `github.com/pinpt/esp/...`:

- `cmd/get_test.go`
- `cmd/put_test.go`
- `cmd/list_test.go`
- `internal/client/client_test.go`
- `internal/ssm/utils_test.go`

Find-replace: `github.com/pinpt/esp/` → `github.com/AbsolutOD/esp/`.

Verification: `go test ./...` is green at the end of this commit.

### Commit 2 — `internal/ssm` DI seam

`refactor(ssm): extract ssmAPI interface, collapse New+Init`

Add an unexported interface in `internal/ssm/ssm.go`:

```go
type ssmAPI interface {
    PutParameter(ctx context.Context, in *awsssm.PutParameterInput, optFns ...func(*awsssm.Options)) (*awsssm.PutParameterOutput, error)
    GetParameter(ctx context.Context, in *awsssm.GetParameterInput, optFns ...func(*awsssm.Options)) (*awsssm.GetParameterOutput, error)
    DeleteParameter(ctx context.Context, in *awsssm.DeleteParameterInput, optFns ...func(*awsssm.Options)) (*awsssm.DeleteParameterOutput, error)
    GetParametersByPath(ctx context.Context, in *awsssm.GetParametersByPathInput, optFns ...func(*awsssm.Options)) (*awsssm.GetParametersByPathOutput, error)
}
```

`*awsssm.Client` satisfies this naturally; `awsssm.NewGetParametersByPathPaginator`'s argument type (`GetParametersByPathAPIClient`) is satisfied by the fourth method, so the same `ssmAPI` value is passed to the paginator.

`Service` becomes:

```go
type Service struct {
    api    ssmAPI
    Region string
}
```

The `Svc *awsssm.Client` and `Cfg aws.Config` fields are removed (unused outside `Init`).

`New()` collapses `New` + `Init` into a single constructor:

```go
func New() (*Service, error) {
    region := utils.GetEnv("AWS_REGION", "us-east-1")
    cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
    if err != nil {
        return nil, err
    }
    return &Service{api: awsssm.NewFromConfig(cfg), Region: region}, nil
}
```

`Service.Init()` is deleted.

The five `Service` methods (`Save`, `GetOne`, `GetMany`, `Copy`, `Delete`) replace `s.Svc.X(...)` with `s.api.X(...)`. The `GetMany` paginator becomes `awsssm.NewGetParametersByPathPaginator(s.api, si)`.

Test seam: an unexported helper `newWithAPI(api ssmAPI) *Service` in a new `internal/ssm/ssm_test_helpers_test.go` (or inline in `ssm_test.go`) lets PR 2 construct a `Service` around a fake without touching the AWS SDK.

### Commit 3 — `internal/client` adapts to new ssm.New signature

`refactor(client): adapt to ssm.New error return`

```go
func New(c *app.Config) (*EspClient, error) {
    if c.Backend != "ssm" {
        return nil, fmt.Errorf("unsupported backend %q", c.Backend)
    }
    svc, err := ssm.New()
    if err != nil {
        return nil, err
    }
    return &EspClient{Backend: c.Backend, Client: svc}, nil
}
```

The prior two-step `svc := ssm.New(); svc.Init()` collapses.

### Commit 4 — `cmd/` constructor pattern + App holder

`refactor(cmd): introduce App holder; convert subcommands to constructor pattern`

New `App` struct in `cmd/root.go`:

```go
type App struct {
    Config  *app.Config
    Client  *client.EspClient
    Verbose bool
}
```

`Config` is populated at construction. `Client` is `nil` until `PersistentPreRunE` populates it.

The package-level `esp *app.Config`, `c *client.EspClient`, and `verbose bool` globals in `cmd/root.go` are deleted. `App.Verbose` replaces the `verbose` global.

Each subcommand file replaces its `var fooCmd = &cobra.Command{...}` + `init()` pair with an exported constructor:

```go
// cmd/get.go
func newGetCmd(a *App) *cobra.Command {
    cmd := &cobra.Command{
        Use:  "get [path]",
        Args: cobra.MinimumNArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            cmd.SilenceUsage = true
            return runGet(cmd, args, a.Client, a.Config)
        },
    }
    cmd.Flags().BoolP("details", "t", false, "...")
    cmd.Flags().BoolP("decrypt", "d", false, "...")
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

Pattern applies to: `get`, `put`, `list`, `copy`, `move`, `delete`, `init`, `version`. (`list` and `init` are already constructors today; their RunE bodies still need the `runX` extraction.)

Pure helpers that read the package-level config gain an explicit `*app.Config` parameter:

- `getParamPath(p string)` → `getParamPath(cfg *app.Config, p string)` (in `cmd/get.go`)
- `formatParamName(n string)` → `formatParamName(cfg *app.Config, n string)` (in `cmd/put.go`)
- `getFullPath(n string)` → `getFullPath(cfg *app.Config, n string)` (in `cmd/put.go`)
- `getPath(a []string)` → `getPath(cfg *app.Config, a []string)` (in `cmd/list.go`)
- `buildEspParamInputFromCmd(cmd)` → `buildEspParamInputFromCmd(cfg *app.Config, cmd)` (in `cmd/put.go`)

`cmd/root.go::Execute()` becomes:

```go
func Execute() {
    a := &App{Config: app.New(false)}
    a.Config.Backend = "ssm"
    root := newRootCmd(a)
    root.AddCommand(
        newGetCmd(a), newPutCmd(a), newListCmd(a), newCopyCmd(a),
        newMoveCmd(a), newDeleteCmd(a), newInitCmd(a), newVersionCmd(a),
    )
    if err := root.Execute(); err != nil {
        os.Exit(1)
    }
}
```

`newRootCmd(a *App) *cobra.Command` builds the cobra root, registers persistent flags against `&a.Config.Env`, `&a.Config.Backend`, `&a.Verbose`, and sets `PersistentPreRunE` to a closure capturing `a` that performs:

1. `cmd.SilenceUsage = true`
2. Env-var checks (`AWS_DEFAULT_REGION`, `AWS_PROFILE`)
3. `client.New(a.Config)` — assigns to `a.Client`
4. Viper config name/path setup
5. `viper.ReadInConfig` — sets `a.Config.IsEspProject = true` if found
6. If `IsEspProject`: `viper.Unmarshal(a.Config)` and `cmd.Root().MarkFlagRequired("env")`

`configureLogging` (the `cobra.OnInitialize` callback) keeps its current behavior; it now reads from `a.Verbose` via a captured closure.

`main.go` is unchanged — still calls `cmd.Execute()`.

### Commit 5 — Update existing tests to App holder pattern

`refactor(cmd): update existing tests for App holder pattern`

The `withOrgPrefix` (`cmd/put_test.go`) and `withAppConfig` (`cmd/get_test.go`) helpers are deleted. Tests construct a per-test `*app.Config` directly and pass it to the helpers under test:

```go
func TestEspName(t *testing.T) {
    cfg := &app.Config{OrgPrefix: "ACME"}
    tests := []struct{...}{...}
    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            got := formatParamName(cfg, tc.in)
            // ...
        })
    }
}
```

Existing test cases stay the same. Only the call signatures and helper plumbing change.

### Commit cadence

Each commit ends with `go build ./... && go vet ./... && go test ./...` clean.

## PR 2 — Tests

One commit per package boundary. All tests use stdlib `testing`, no testify. Table-driven where the function has multiple distinct cases. Fakes are defined locally per test file.

### Commit 1 — `test(ssm): cover Service methods via fake ssmAPI`

`internal/ssm/ssm_test.go` (new file). A local `fakeSSMAPI` implements the `ssmAPI` interface, recording inputs and returning canned outputs/errors.

- **`TestService_Save`** — table:
  - success returns `Version` from the SDK output;
  - `Secure: true` → `PutParameterInput.Type == ParameterTypeSecureString`;
  - AWS returns `*ssmtypes.ParameterAlreadyExists` → mapped error returned (verifies `mapErr(Save, …)` integration).
- **`TestService_GetOne`** — table:
  - success returns the `EspParam` produced by `convertToEspParam`;
  - `Decrypt: true` → `GetParameterInput.WithDecryption == true`;
  - AWS returns `*ssmtypes.ParameterNotFound` → mapped error.
- **`TestService_GetMany`** — table:
  - paginator drives multiple pages, all params concatenated;
  - error on second page returns mapped error mid-iteration;
  - empty pages return empty slice, no error;
  - `Recursive: true` and `Decrypt: true` flags reach `GetParametersByPathInput`.
- **`TestService_Delete`** — table:
  - success returns the input name;
  - AWS returns `*ssmtypes.ParameterNotFound` → mapped error.
- **`TestService_Copy`** — table:
  - success returns `SaveOutput`;
  - `GetOne` fails → `Save` is never invoked;
  - `Save` fails after successful `GetOne` → error surfaces.

### Commit 2 — `test(client): cover EspClient wrapper methods`

Extends `internal/client/client_test.go`. A local `fakeBackend` implements `client.Client`.

- **`TestEspClient_GetParam` / `_ListParams` / `_Save` / `_Delete`** — pass-through assertions: input goes in, recorded call on the fake matches, output flows back unchanged.
- **`TestEspClient_Copy`** — composition:
  - success → `Client.Copy` invoked, then re-fetches destination with `Decrypt: true`;
  - `Client.Copy` fails → no re-fetch, error returned;
  - re-fetch fails → error returned.
- **`TestEspClient_Move`** — three-step composition (`GetParam` → `Save` → `Delete`):
  - all succeed → returns `MoveCommand`;
  - each step's failure returns early, subsequent steps not invoked;
  - source's `Secure` flag is carried into the `Save` input.

`TestNewUnsupportedBackend` already exists from PR #27.

### Commit 3 — `test(cmd): cover runX functions via fake EspClient`

A new `cmd/testutil_test.go` (test-only file, build-tagged via `_test.go` suffix only — no separate build tag) provides:

- `fakeBackend` — implements `client.Client`. Records calls, returns canned responses keyed by method.
- `newTestApp(t *testing.T, fake client.Client, cfg app.Config) *App` — builds an `*App` with the fake client wrapped in a real `*client.EspClient{Backend: "ssm", Client: fake}`.

Tests per subcommand:

- **`TestRunGet`** — table:
  - literal path (leading `/`) passes through unchanged;
  - relative path resolved through `cfg.GetAppParamPath`;
  - `--decrypt` and `--details` flags reach the underlying call;
  - `c.GetParam` error surfaces.
- **`TestRunPut`** — table:
  - `formatParamName` wiring (`foo` → `<PREFIX>_FOO`, `foo-bar` → `<PREFIX>_FOO_BAR`);
  - `Save` failure returns early (re-fetch never invoked);
  - re-fetch failure surfaces.
- **`TestRunList`** — table:
  - no args → uses `cfg.GetAppPath()`;
  - with arg → uses arg verbatim;
  - `--decrypt` reaches `ListParamInput`;
  - `ListParams` error surfaces.
- **`TestRunCopy`** — table:
  - happy path (both args set);
  - empty source → `errors.New("source can not be empty")`;
  - empty destination → `errors.New("destination can not be empty")`;
  - `c.Copy` failure surfaces.
- **`TestRunMove`** — happy path; `c.Move` failure surfaces.
- **`TestRunDelete`** — happy path; `c.Delete` failure surfaces.
- **`TestRunVersion`** — happy path (prints version constant).

### Commit 4 — `test(cmd): cover persistentPreRunE env-var validation`

New `cmd/root_test.go`. Uses `t.Setenv` for isolation.

- Missing `AWS_DEFAULT_REGION` → returns error mentioning `AWS_DEFAULT_REGION`.
- Missing `AWS_PROFILE` (with region set) → returns error mentioning `AWS_PROFILE`.
- Both set + `Backend: "vault"` (unsupported) → returns the unsupported-backend error from `client.New`.

The "both env vars set + backend ssm" path is **not** tested here — `client.New("ssm")` invokes `config.LoadDefaultConfig` which is environment-coupled. That gap is documented (see "Coverage gaps" below).

### Commit cadence

Each commit ends with `go build ./... && go vet ./... && go test ./...` clean.

## Coverage Gaps (Explicit, Accepted)

- **`ssm.New()`** — calls `config.LoadDefaultConfig`, which loads from environment. Not tested.
- **`app.InitQuestions` and `cmd/init.go::runInit`** — `survey.Ask` is interactive. Approach C territory; deferred.
- **`cmd/get.go::display*` and `cmd/list.go::displayParams`** — print to stdout via `fmt.Printf`. Pure formatting, no logic. Testing requires `io.Writer` injection; out of scope.
- **`cmd/root.go::configureLogging`** — sets the slog default. No I/O, no failure path. Not worth testing.
- **`main.go`** — single line calling `cmd.Execute()`. Not tested.

## Risks

- **`cmd/` global removal** inverts a non-goal of the prior cleanup-and-tests spec ("No removal of the cmd/ package-level globals"). Deliberate trade for testability under the constructor-injection pattern. Internal-only; no public API impact.
- **`ssm.Service` field removal** (`Svc *awsssm.Client`, `Cfg aws.Config`) is a breaking change for any caller reaching into those fields. The `internal/` directory restricts callers to this module; the only current caller (`internal/client/`) invokes `ssm.New()` and the `Client` interface methods only. Practical blast radius is zero.
- **Constructor-pattern churn in `cmd/`.** Every subcommand file is rewritten. Reviewers verify cobra wiring is preserved (flags, aliases, args constraints, examples). Mitigation: PR 1 carries no behavior change — `esp <subcommand>` invocations should produce identical output before and after.
- **`ssmAPI` interface drift.** If the AWS SDK changes a method signature in a future bump, the interface and fake both update. Cost is small (4 methods); SDK signatures are stable.
- **Test-helper duplication.** `fakeBackend` (implements `client.Client`) is defined in both `internal/client/client_test.go` and `cmd/testutil_test.go`. Duplication is acceptable for two callers; promoting to a shared `internal/client/clienttest` package would be over-abstraction.

## Verification (Both PRs, Every Commit)

- `go build ./...` — exit 0, no output.
- `go vet ./...` — exit 0, no output.
- `go test ./...` — all packages report `ok` or `[no test files]`.
- `esp --help` — prints usage without AWS env vars set.
- After PR 1, `esp get /some/path` (with `AWS_DEFAULT_REGION` and `AWS_PROFILE` set) behaves identically to pre-refactor. PR 2 adds tests, no behavior change.
