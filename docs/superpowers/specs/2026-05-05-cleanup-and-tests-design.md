# Cleanup and Tests Design

**Date:** 2026-05-05
**Status:** Approved
**Scope:** Tidy the codebase and add pure-logic tests, in two sequential PRs, before any new feature work begins.

## Goals

- Remove dead code and obvious bugs that survive in the post-refresh tree.
- Replace ad-hoc `os.Exit` and `panic` with idiomatic cobra `RunE` / `PersistentPreRunE` error flow.
- Fix bugs where `fmt.Errorf` results were discarded, suppressing user-visible error messages.
- Switch `GetParametersByPath` to the v2-idiomatic paginator.
- Add table-driven, stdlib-`testing`-based unit tests for every pure function in the codebase.
- Keep external CLI behavior unchanged on the happy path. On the error path, exit code collapses to `1` (was: `1`, `2`, or `3` in places); stderr messages are preserved or improved.

## Non-Goals

- No dependency injection seam, no interface mocks for the SSM API. Deferred to a follow-up brainstorm ("option C").
- No tests against LocalStack or any AWS-emulating infrastructure.
- No removal of the `cmd/` package-level globals (`esp`, `c`).
- No tightening of the `client.Client` interface beyond adding `error` returns.
- No CI workflow changes.
- No new CLI features.
- No module path change or README rewrite — those live in the rebrand spec (`2026-05-05-rebrand-design.md`).

## PR Plan

Two PRs, branched off `main`, sharing one worktree at `.worktrees/cleanup-and-tests`:

1. **PR 1 — `feature/cleanup-error-flow`** — dead-code removal, error-flow refactor, bug fixes, paginator.
2. **PR 2 — `feature/tests-pure-logic`** — pure-logic tests on top of the cleaner shape from PR 1.

PR 2 is started after PR 1 merges and `main` updates. The worktree is reused: switch the branch to `feature/tests-pure-logic` rebased onto fresh `main`.

## PR 1 — Cleanup

### Dead code removal

- `internal/ssm/utils.go`: delete the `AwsParam` struct and its `isValid()` method. Unreferenced.
- `internal/app/init.go`: delete the empty `setPrefixes` method.

### Error flow: `panic` → returned error

- `internal/client/client.go`: `New(c *app.Config)` becomes `New(c *app.Config) (*EspClient, error)`. Returns `(nil, fmt.Errorf("unsupported backend %q", c.Backend))` for unknown backends instead of panicking.

### Error flow: SSM client returns errors

- `internal/ssm/ssm.go`:
  - `Init()` returns `error` instead of `os.Exit(1)` on `LoadDefaultConfig` failure.
  - All five client methods (`Save`, `GetOne`, `GetMany`, `Copy`, `Delete`) gain an `error` return, replacing the in-method `handleAwsErr` that called `os.Exit`.
  - `client.Client` interface (`internal/client/client.go`) is updated to match the new return signatures.
  - Internal helper `getNextParams` is replaced by `ssm.NewGetParametersByPathPaginator` (see "Paginator" below).
- `internal/ssm/utils.go`: delete `handleAwsErr` (no longer needed; subcommands now receive errors directly).
- `internal/ssm/errors.go`: keep the per-action error mapping functions, but they become pure functions returning `(error, bool)` or just `error` — input goes in, mapped error comes out, no side effects. The "fall through to `smithy.APIError`" pattern is preserved.

### Error flow: cmd subcommands use `RunE`

- All `Run: func(...)` blocks in `cmd/copy.go`, `cmd/delete.go`, `cmd/get.go`, `cmd/list.go`, `cmd/move.go`, `cmd/put.go`, `cmd/init.go`, `cmd/version.go` flip to `RunE: func(cmd *cobra.Command, args []string) error`.
- First line of every `RunE` body: `cmd.SilenceUsage = true`. This is the cobra-idiomatic way to get usage on misuse (argument/flag errors, which fire before `RunE`) but suppress usage on runtime errors.
- Each `RunE` returns `nil` on success and the underlying error otherwise. No `os.Exit` calls inside subcommands.

### Error flow: env-var checks move to `PersistentPreRunE`

- `cmd/root.go`: env-var checks (`AWS_DEFAULT_REGION`, `AWS_PROFILE`), `client.New(esp)` construction, and `MarkFlagRequired` calls move from `cobra.OnInitialize(initConfig)` into `rootCmd.PersistentPreRunE`. All return `error` instead of `os.Exit`.
- `slog` configuration stays in `cobra.OnInitialize` — no I/O, no failure mode, fine where it is.
- `--help` continues to work without AWS env vars: cobra short-circuits before `PersistentPreRunE` for help, just as it does for `OnInitialize` today.

### Error flow: top-level `Execute()`

- The single remaining `os.Exit(1)` in `cmd/root.go::Execute()` stays. It is the top-level error handler for cobra, the idiomatic place for the one process-exiting call.

### Bug audit

- `cmd/copy.go` previously had a `fmt.Errorf(...) + os.Exit(1)` pattern that discarded the error message. That was already fixed during the dependency refresh (it is now `fmt.Fprintln(os.Stderr, ...) + os.Exit(1)`). After the `RunE` flip, both args-empty checks become `return errors.New(...)`.
- Audit `cmd/put.go`, `cmd/move.go`, `cmd/delete.go`, and other subcommands for any remaining `fmt.Errorf`-then-discard pattern, or any other in-line stderr-print + `os.Exit` patterns. Convert them to `return error` as part of the `RunE` migration.

### Paginator

- Replace the manual recursive `getNextParams` in `internal/ssm/ssm.go` with `ssm.NewGetParametersByPathPaginator`. Iteration in `GetMany` becomes the standard `for paginator.HasMorePages() { ... }` loop. Behavior unchanged externally.

### Stub test files

- `cmd/get_test.go`, `cmd/put_test.go`, `internal/app/config_test.go`, `internal/client/client_test.go` are currently empty stubs.
- Decision per file: if PR 2 will fill them with real tests, leave them. If a file ends up unused after PR 2 (i.e. its package has no testable pure logic), delete it then. This decision is finalized in PR 2 — PR 1 does not touch them.

### CLAUDE.md update

- The "AWS_DEFAULT_REGION and AWS_PROFILE must be set..." paragraph is updated to reflect the new flow: env-var checks now live in `PersistentPreRunE`, errors now flow through cobra rather than `os.Exit`, and the exit code is `1` on any failure.

### PR 1 commit cadence

Each commit ends with `go build ./... && go vet ./... && go test ./...` clean.

1. `chore: remove dead code (AwsParam, setPrefixes)`
2. `refactor: client.New returns error instead of panic`
3. `refactor: ssm methods return errors; remove handleAwsErr`
4. `refactor: cmd subcommands use RunE with SilenceUsage`
5. `refactor: env-var validation moves to PersistentPreRunE`
6. `refactor: use ssm paginator for GetParametersByPath`
7. `chore: update CLAUDE.md error-flow notes`

## PR 2 — Pure-Logic Tests

### Style

- Stdlib `testing`, no testify or other test libraries.
- Table-driven where the function has multiple distinct cases.
- One `_test.go` file per package, placed alongside the file under test.
- Test names follow Go convention: `TestFuncName_Case` for table entries, `TestFuncName` for single-case tests.

### Coverage stance

No numeric coverage target. The bar is: every pure function reachable from a subcommand has a happy-path test plus its meaningful edge cases. Functions that exist only for side effects against AWS are out of scope (option-C territory).

### Tests by package

**`internal/app/config_test.go`** (currently empty stub — fill it):
- `TestGetAppPath` — table covering org/env/app combinations and edge cases (empty fields, leading slashes already present in inputs).
- `TestGetAppParamPath` — table covering the path-builder with various param-name shapes.

**`internal/app/init_test.go`** (currently has `TestWriteConfig` — extend):
- `TestUpdateWithInput` — verify the comma-separated env list (`"dev,test,prod"`) splits correctly; whitespace handling.
- `TestCreateEspFile` — `Config` → `espFile` field-by-field equality, table-driven.

**`cmd/put_test.go`** (currently empty stub — fill it):
- `TestEspName` — table covering: already-prefixed inputs (unchanged), lowercase inputs (uppercased and prefixed), hyphenated inputs (hyphens to underscores), mixed-case inputs (the actual rule pinned by the test).

**`cmd/list_test.go`** (currently has `TestGetPathWithFullPath` and `TestGetPathEnvVarName` — extend):
- `TestGetPathRelative` — short name input → resolved through `GetAppParamPath`.

**`cmd/get_test.go`** — only fill if `cmd/get.go` contains pure logic distinct from `cmd/list.go::getPath`. If after audit there is no testable pure logic in `get.go`, delete the empty stub file.

**`internal/ssm/utils_test.go`** (new file):
- `TestSelectType` — `true` → `ssmtypes.ParameterTypeSecureString`; `false` → `ssmtypes.ParameterTypeString`.
- `TestConvertToEspParam` — table: input `ssmtypes.Parameter` (constructed in test) → expected `common.EspParam`. Specifically asserts `Secure: true` is set when input type is `SecureString`.

**`internal/ssm/errors_test.go`** (new file):
- Table-driven against each `checkSSM*Error` function. Inputs are constructed typed errors (`&ssmtypes.ParameterNotFound{}`, `&ssmtypes.InternalServerError{}`, etc.) plus a generic `smithy.APIError` fallback case. Asserts that matched types pass through (function returns the error) and unmatched non-API errors return `nil`.

**`internal/client/client_test.go`** (currently empty stub — fill it after PR 1's panic→error change):
- `TestNewUnsupportedBackend` — `&app.Config{Backend: "vault"}` → returns a non-nil error.
- The `"ssm"` happy path is *not* tested here: it requires `LoadDefaultConfig` and AWS credentials. A comment in the test file notes why.

### PR 2 commit cadence

1. `test: cover internal/app path builders and init`
2. `test: cover cmd path/name helpers`
3. `test: cover internal/ssm utils and errors`
4. `test: cover internal/client.New unsupported-backend path`
5. `chore: delete unused stub test files` (only if the audit during PR 2 finds any)

## Verification (both PRs)

Per commit and at end of each PR:

- `go build ./...` — exit 0, no output.
- `go vet ./...` — exit 0, no output.
- `go test ./...` — all packages report `ok` or `[no test files]`. New tests pass.
- `esp --help` — prints usage without any AWS env vars set.
- After PR 1: `esp <subcommand>` with missing `AWS_DEFAULT_REGION` exits 1 with a clear stderr message (was: exit 1 from `os.Exit(1)` direct print). Same for missing `AWS_PROFILE` (was: exit 2; now: exit 1).

## Risks

- **Cobra error-printing format change.** With `RunE`, cobra's default behavior prints `Error: <msg>` to stderr followed by the command usage. We mitigate by setting `cmd.SilenceUsage = true` as the first line of every `RunE` — usage prints only on misuse (which fires before `RunE`), not on runtime errors. The `Error:` prefix on the message is a small UX change from the current bare-`fmt.Println` style; acceptable.
- **Public-API change to `client.New`'s signature.** Anyone importing `github.com/pinpt/esp/internal/client` as a library would break. The `internal/` directory restricts that to the module itself, so the practical blast radius is zero. Flagged for completeness.
- **Exit-code collapse.** Codes `2` (missing `AWS_PROFILE`) and `3` (`MarkFlagRequired` failure) collapse to `1`. These were undocumented and look incidental rather than designed. If anything depends on `$?`, it would already be brittle.
