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

`AWS_DEFAULT_REGION` and `AWS_PROFILE` must be set before invoking any subcommand — `cmd/root.go`'s `persistentPreRun` (wired as `rootCmd.PersistentPreRunE`) returns an error if either is missing, which cobra surfaces as `Error: <message>` on stderr; `Execute()` then exits 1. (Previously, missing env vars exited 1 or 2 directly via `os.Exit`; the codes now collapse to a uniform 1.) They are NOT required for `esp --help` (cobra short-circuits before `PersistentPreRunE`) or for `go test ./...`.

## Architecture

The code is organized in three layers:

1. **`cmd/`** — Cobra subcommands (`get`, `list`/`ls`, `put`/`add`/`create`, `delete`/`rm`, `copy`/`cp`, `move`/`mv`, `init`, `version`). Each file registers itself onto `rootCmd` via `init()`. The package-level globals `esp *app.Config` and `c *client.EspClient` (in `cmd/root.go`) are constructed once in `root.go`'s `init()` and used by every subcommand.

2. **`internal/client/`** — A thin facade over a backend `Client` interface (`Save`, `GetOne`, `GetMany`, `Copy`, `Delete`); each method returns `(value, error)`. `client.New` switches on `Config.Backend`; today the only valid value is `"ssm"`, and any other value returns `fmt.Errorf("unsupported backend %q", c.Backend)`. Add new backends here.

3. **`internal/ssm/`** — The AWS SSM implementation of the `Client` interface, plus parameter-type conversion (`utils.go`) and AWS-error mapping (`errors.go`).

`internal/common/` defines the wire types (`EspParam`, `EspParamInput`, `GetOneInput`, `ListParamInput`, `CopyCommand`, `MoveCommand`, etc.) shared between `cmd`, `client`, and `ssm` to avoid import cycles.

`internal/app/` owns the per-project `Config` and the `init` flow:
- `Config` is loaded by Viper from a `.espFile` (YAML) in the working directory. If found, `IsEspProject` is set true and `--env` becomes required.
- `Config.GetAppPath()` and `GetAppParamPath(name)` produce SSM paths using the convention **`/<OrgName>/<Env>/<AppName>/<param>`**. Subcommands call these to build paths from short names.
- `cmd/put.go` additionally normalizes parameter names: anything not already starting with `OrgPrefix` is rewritten to `<OrgPrefix>_<UPPER_NAME>`. `cmd/get.go` and `cmd/list.go` treat any argument starting with `/` as a literal SSM path and otherwise pass it through `GetAppParamPath`.
- `init.go::InitQuestions` uses `AlecAivazis/survey` to prompt for backend/org/prefix/app/envs and writes the `.espFile` via `WriteConfig`.

## Adding a subcommand

Create a new file under `cmd/`, define a `*cobra.Command`, register it in an `init()` with `rootCmd.AddCommand(...)`, and use the package-level `c` (`*client.EspClient`) and `esp` (`*app.Config`) for backend operations and path resolution. Mirror the existing pattern of building a `common.*Input` struct and passing it to the client method.

## Adding a backend

Implement the `client.Client` interface (`internal/client/client.go`) in a new package under `internal/`, then extend the switch in `client.New` to construct it when `Config.Backend` matches.

## Testing notes

Tests use stdlib `testing` only (no testify). Table-driven where the function has multiple distinct cases. One `_test.go` file per package, alongside the file under test. The bar is: every pure function reachable from a subcommand has a happy-path test plus its meaningful edge cases.

Functions that exist only for side effects against AWS — `Service.Save` / `GetOne` / `GetMany` / `Copy` / `Delete`, the `EspClient` wrappers in `internal/client/`, and `Service.Init` — are intentionally not unit-tested. Exercising them requires either real AWS credentials or a dependency-injection seam that doesn't exist yet. The `client.New` test (`internal/client/client_test.go`) covers only the unsupported-backend path for that reason.

The `cmd` package shares two package-level globals (`esp *app.Config` and `c *client.EspClient`) across subcommands. Tests that depend on them use the local `withOrgPrefix` and `withAppConfig` helpers in `cmd/put_test.go` and `cmd/get_test.go` to swap values in and restore them via `t.Cleanup`.
