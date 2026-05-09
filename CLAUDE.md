# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

`esp` is a Go CLI utility for managing AWS SSM Parameter Store entries as environment variables. It is built on Cobra/Viper and has a pluggable backend layer (only the `ssm` backend exists today).

## Common commands

```sh
# Build the binary
go build -o esp .

# Run all tests
go test ./...

# Run a single package's tests
go test ./internal/app

# Run a single test
go test ./internal/app -run TestWriteConfig

# Run with verbose logging (sets slog default level to INFO)
./esp --verbose <subcommand>
```

`AWS_DEFAULT_REGION` and `AWS_PROFILE` must be set before invoking any subcommand — `cmd/root.go`'s `persistentPreRunE` returns an error if either is missing, which cobra surfaces as `Error: <message>` on stderr; `Execute()` then exits 1. (Previously, missing env vars exited 1 or 2 directly via `os.Exit`; the codes now collapse to a uniform 1.) They are NOT required for `esp --help` (cobra short-circuits before `PersistentPreRunE`) or for `go test ./...`.

## Architecture

The code is organized in three layers:

1. **`cmd/`** — Cobra subcommands (`get`, `list`/`ls`, `put`/`add`/`create`, `delete`/`rm`, `copy`/`cp`, `move`/`mv`, `init`, `version`). Each file exports a `newXCmd(*App)` constructor (or `newVersionCmd()` for version, which doesn't depend on App). `cmd.Execute()` builds a single `*App` holder containing `Config *app.Config` and `Client *client.EspClient`, threads it through each constructor, and adds the resulting commands to the root. `Config` is populated at construction time; `Client` is `nil` until the root's `PersistentPreRunE` validates AWS env vars and calls `client.New(a.Config)` to populate it. Each subcommand's `RunE` is a thin closure that calls a free function `runX(cmd, args, c, cfg)` — those `runX` functions are the unit-testable seam.

2. **`internal/client/`** — A thin facade over a backend `Client` interface (`Save`, `GetOne`, `GetMany`, `Copy`, `Delete`); each method returns `(value, error)`. `client.New` switches on `Config.Backend`; today the only valid value is `"ssm"`, and any other value returns `fmt.Errorf("unsupported backend %q", c.Backend)`. Add new backends here.

3. **`internal/ssm/`** — The AWS SSM implementation of the `Client` interface, plus parameter-type conversion (`utils.go`) and AWS-error mapping (`errors.go`). `Service` holds an unexported `ssmAPI` interface (covering the four AWS SDK methods used: `PutParameter`, `GetParameter`, `DeleteParameter`, `GetParametersByPath`); the concrete `*awsssm.Client` satisfies it in production, and tests inject a fake. `ssm.New() (*Service, error)` calls `config.LoadDefaultConfig` and `awsssm.NewFromConfig` itself — the prior two-step `New`+`Init` is gone.

`internal/common/` defines the wire types (`EspParam`, `EspParamInput`, `GetOneInput`, `ListParamInput`, `CopyCommand`, `MoveCommand`, etc.) shared between `cmd`, `client`, and `ssm` to avoid import cycles.

`internal/app/` owns the per-project `Config` and the `init` flow:
- `Config` is loaded by Viper from a `.espFile` (YAML) in the working directory. If found, `IsEspProject` is set true and `--env` becomes required.
- `Config.GetAppPath()` and `GetAppParamPath(name)` produce SSM paths using the convention **`/<OrgName>/<Env>/<AppName>/<param>`**. Subcommands call these to build paths from short names.
- `cmd/put.go` additionally normalizes parameter names: anything not already starting with `OrgPrefix` is rewritten to `<OrgPrefix>_<UPPER_NAME>`. `cmd/get.go` and `cmd/list.go` treat any argument starting with `/` as a literal SSM path and otherwise pass it through `GetAppParamPath`.
- `init.go::InitQuestions` uses `AlecAivazis/survey` to prompt for backend/org/prefix/app/envs and writes the `.espFile` via `WriteConfig`.

## Adding a subcommand

Create a new file under `cmd/` and export a constructor `newFooCmd(a *App) *cobra.Command` returning the configured cobra command. Inside `RunE`, delegate to a free function `runFoo(cmd, args, a.Client, a.Config)` (or whatever subset of dependencies the command needs) — keep the closure tiny so the body is testable in isolation. Wire the constructor into `cmd/root.go::Execute()`'s `root.AddCommand(...)` call. Mirror the existing pattern of building a `common.*Input` struct and passing it to the client method.

## Adding a backend

Implement the `client.Client` interface (`internal/client/client.go`) in a new package under `internal/`, then extend the switch in `client.New` to construct it when `Config.Backend` matches.

## Testing notes

Tests use stdlib `testing` only (no testify). Table-driven where the function has multiple distinct cases. One `_test.go` file per package, alongside the file under test. The bar is: every pure function reachable from a subcommand has a happy-path test plus its meaningful edge cases.

Two DI seams enable unit-testing the AWS-coupled code without real credentials:
- **`internal/ssm/`** — the unexported `ssmAPI` interface on `Service.api`. Tests construct a `Service` with a fake `ssmAPI` implementation and assert behavior of `Save`/`GetOne`/`GetMany`/`Delete`/`Copy` including AWS error-mapping.
- **`internal/client/`** — the existing `Client` interface on `EspClient.Client`. Tests construct an `EspClient{Backend: "ssm", Client: fake}` and assert wrapper composition (notably `Move` and `Copy`).

In `cmd/`, tests construct a per-test `*app.Config` and a fake-backed `*client.EspClient`, then invoke the relevant `runX` function directly. There are no package-level globals in `cmd/` — all state lives on the `*App` holder threaded through subcommand constructors.

Coverage gaps that remain (deliberately): `ssm.New()` (calls `config.LoadDefaultConfig`, environment-coupled), `app.InitQuestions` and `cmd/init.go` (interactive `survey.Ask`), and the `cmd/get.go::display*` print helpers (formatting only, no logic).
