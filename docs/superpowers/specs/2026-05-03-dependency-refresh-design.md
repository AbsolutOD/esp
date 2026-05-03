# Dependency Refresh Design

**Date:** 2026-05-03
**Status:** Approved
**Scope:** Full refresh of `esp` Go toolchain and all direct dependencies, including replacement of deprecated libraries.

## Goals

- Bring `esp` onto a current Go toolchain.
- Update every direct dependency to its latest stable release in the same major, except where the library is deprecated.
- Replace deprecated direct dependencies with current alternatives.
- Keep the user-facing CLI behavior identical: same subcommands, same flags, same espFile format, same SSM path conventions.

## Non-Goals

- No new tests. Existing real tests (`cmd/list_test.go`, `internal/app/init_test.go`) must continue to pass; empty stub test files are left alone.
- No refactoring beyond what the upgrades force.
- No CI workflow changes (`.github/workflows/`).
- No module path change. `github.com/pinpt/esp` stays.
- No README updates.
- No new backends; the SSM backend is the only one touched.

## Go Toolchain

- `go.mod`: `go 1.13` → `go 1.26`
- `go.mod`: add `toolchain go1.26.2` directive (matches the developer's local toolchain).

## Direct Dependency Changes

| Package | Current | Action | Target |
|---|---|---|---|
| `github.com/aws/aws-sdk-go` | v1.44.70 | **remove** | replaced by aws-sdk-go-v2 modules |
| `github.com/aws/aws-sdk-go-v2/config` | — | **add** | latest stable |
| `github.com/aws/aws-sdk-go-v2/service/ssm` | — | **add** | latest stable |
| `github.com/aws/aws-sdk-go-v2/aws` | — | **add** (transitive of config) | latest stable |
| `github.com/aws/smithy-go` | — | **add** (for typed error checks) | latest stable |
| `github.com/spf13/cobra` | v1.5.0 | upgrade | latest stable (~v1.10) |
| `github.com/spf13/viper` | v1.12.0 | upgrade | latest stable (~v1.20) |
| `github.com/spf13/jwalterweatherman` | v1.1.0 | **remove** | replaced by stdlib `log/slog` |
| `github.com/AlecAivazis/survey/v2` | v2.3.5 | upgrade | v2.3.7 |
| `github.com/logrusorgru/aurora` | v2.0.3+incompatible | **replace** | `github.com/logrusorgru/aurora/v4` latest |
| `github.com/olekukonko/tablewriter` | v0.0.5 | upgrade | latest stable (v1.x) — **API rewrite required** |
| `gopkg.in/yaml.v2` | v2.4.0 | **replace** | `gopkg.in/yaml.v3` latest |

After all changes, `go mod tidy` is run to prune indirect deps and rewrite `go.sum`.

## AWS SDK v1 → v2 Migration

All changes confined to `internal/ssm/`.

### `ssm.go`

- Replace imports:
  - `github.com/aws/aws-sdk-go/aws` → `github.com/aws/aws-sdk-go-v2/aws`
  - `github.com/aws/aws-sdk-go/aws/session` → removed
  - `github.com/aws/aws-sdk-go/service/ssm` → `github.com/aws/aws-sdk-go-v2/service/ssm` (kept aliased as `awsssm`)
  - Add `github.com/aws/aws-sdk-go-v2/config` for `config.LoadDefaultConfig`
  - Add `github.com/aws/aws-sdk-go-v2/service/ssm/types` (aliased as `ssmtypes`) for input enums and typed errors
- `Service` struct: drop `*session.Session` field; `Svc` becomes `*awsssm.Client`; `Cfg` becomes `aws.Config` (v2 type).
- `New()`: keep the region resolution from `AWS_REGION` env var (default `us-east-1`).
- `Init()`: replace `session.Must(session.NewSession(...))` and `awsssm.New(...)` with:
  ```go
  cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(s.Region))
  // exit on err — same fail-fast behavior as today
  s.Cfg = cfg
  s.Svc = awsssm.NewFromConfig(cfg)
  ```
- Every API call gains `context.Background()` as its first arg. Full context plumbing from cmd-layer is **out of scope**; we use `context.Background()` at the call site so the change is mechanical.
- `Save`/`Delete`/`GetOne`/`GetMany`/`getNextParams` signatures and behavior unchanged at the boundary; only the SDK call shape changes inside.
- Pagination: keep the existing `NextToken` recursion in `GetMany`/`getNextParams`. The v2 paginator (`ssm.NewGetParametersByPathPaginator`) is a nicer pattern but switching is out of scope.

### `utils.go`

- `selectType(t bool) *string` → `selectType(t bool) ssmtypes.ParameterType`. Return `ssmtypes.ParameterTypeSecureString` or `ssmtypes.ParameterTypeString` directly (v2 uses typed enum, not `*string`).
- Update `PutParameterInput.Type` accordingly in `ssm.go`.
- `convertToEspParam(ap *ssmtypes.Parameter)`:
  - `*ap.ARN` → `aws.ToString(ap.ARN)` (or check nil if simpler — but `aws.ToString` is the v2 idiom)
  - `*ap.Name` → `aws.ToString(ap.Name)`
  - `*ap.Type` → `string(ap.Type)` (it's an enum value, not a pointer)
  - `*ap.Value` → `aws.ToString(ap.Value)`
  - `*ap.Version` → `ap.Version` (v2 makes it `int64`, not `*int64`)
  - `*ap.LastModifiedDate` → unchanged (v2's `ssmtypes.Parameter.LastModifiedDate` is still `*time.Time`); assigns into `common.EspParam.LastModifiedDate` (which is already `time.Time` — verified)
- The unused `AwsParam` struct in `utils.go` (declared but never referenced) has a stale `LastModifiedDate float32` field. Leave alone — out of scope. Removing it is a separate cleanup.

### `errors.go`

- Replace `awserr.Error` checks with `errors.As(err, &apiErr)` against `smithy.APIError`, plus typed checks against v2 error types from `ssmtypes` where the existing code uses string-constant error codes:
  - `awsssm.ErrCodeParameterNotFound` → `*ssmtypes.ParameterNotFound`
  - `awsssm.ErrCodeParameterAlreadyExists` → `*ssmtypes.ParameterAlreadyExists`
  - `awsssm.ErrCodeParameterLimitExceeded` → `*ssmtypes.ParameterLimitExceeded`
  - `awsssm.ErrCodeTooManyUpdates` → `*ssmtypes.TooManyUpdates`
  - `awsssm.ErrCodeHierarchyTypeMismatchException` → `*ssmtypes.HierarchyTypeMismatchException`
  - `awsssm.ErrCodeInvalidAllowedPatternException` → `*ssmtypes.InvalidAllowedPatternException`
  - `awsssm.ErrCodeParameterMaxVersionLimitExceeded` → `*ssmtypes.ParameterMaxVersionLimitExceeded`
  - `awsssm.ErrCodeUnsupportedParameterType` → `*ssmtypes.UnsupportedParameterType`
  - `awsssm.ErrCodePoliciesLimitExceededException` → `*ssmtypes.PoliciesLimitExceededException`
  - `awsssm.ErrCodeInvalidPolicyTypeException` → `*ssmtypes.InvalidPolicyTypeException`
  - `awsssm.ErrCodeInvalidPolicyAttributeException` → `*ssmtypes.InvalidPolicyAttributeException`
  - `awsssm.ErrCodeIncompatiblePolicyException` → `*ssmtypes.IncompatiblePolicyException`
  - `awsssm.ErrCodeInternalServerError` → `*ssmtypes.InternalServerError`
  - `awsssm.ErrCodeInvalidKeyId` → `*ssmtypes.InvalidKeyId`
  - `awsssm.ErrCodeInvalidFilterKey` → `*ssmtypes.InvalidFilterKey`
  - `awsssm.ErrCodeInvalidFilterOption` → `*ssmtypes.InvalidFilterOption`
  - `awsssm.ErrCodeInvalidFilterValue` → `*ssmtypes.InvalidFilterValue`
  - `awsssm.ErrCodeInvalidNextToken` → `*ssmtypes.InvalidNextToken`
- `MissingRegion` is no longer a typed error in v2; `config.LoadDefaultConfig` returns a generic error if region resolution fails. `checkRegion` becomes a string-match helper or is folded into the generic fallback. Acceptable since the pre-existing behavior is just "print and exit."
- Preserve the existing per-action switch in `checkSSMError(a action, err error)`.
- Behavior at the boundary (`handleAwsErr` prints `"SSM Error: ..."` and exits 1) is unchanged.

## Logging Swap (jwalterweatherman → log/slog)

### `cmd/root.go`

- Remove the `jww` import.
- Add `log/slog` and `os`.
- Replace the verbose-flag block:
  ```go
  jww.SetStdoutThreshold(jww.LevelInfo)
  ```
  with:
  ```go
  slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo})))
  ```
- Default (non-verbose) is the stdlib default handler, which logs at `LevelInfo` to stderr. To suppress info output by default, set `slog.LevelWarn` as the default:
  ```go
  slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelWarn})))
  ```
  installed before the verbose check. The verbose branch then re-installs at `LevelInfo`.

### `internal/app/config.go`

- Remove the `jww` import; add `log/slog`.
- Three call sites:
  - `jww.INFO.Printf("rendered path: %s", path)` → `slog.Info("rendered path", "path", path)`
  - `jww.INFO.Printf("rendered param path: %s", path)` → `slog.Info("rendered param path", "path", path)`
  - `jww.INFO.Printf("found ENV: %s", e)` → `slog.Info("found ENV", "env", e)`

## tablewriter v1 API Rewrite

Only `cmd/get.go` uses `tablewriter`. The v1.x API replaces `SetHeader`/`Append`/`Render` with a config-builder + `Header`/`Append`/`Render` shape. Concrete rewrite is left to the implementation step but constrained:

- Output must remain a two-column table (`Field`, `Value` columns) printed to `os.Stdout`.
- Column header colors and the bright-yellow row labels (`ID`, `Last_Modified`, `Name`, `Type`, `Value`, `Version`) are preserved.
- The exact visual output may shift slightly (border style, padding) due to v1.x defaults — acceptable.

## yaml.v2 → yaml.v3

- `internal/app/init.go`: `import "gopkg.in/yaml.v2"` → `import "gopkg.in/yaml.v3"`. The `yaml.Marshal` call signature is unchanged.
- `internal/app/init_test.go`: same import swap. The `yaml.Unmarshal` call signature is unchanged.
- No struct-tag changes expected; the existing `EspConfig` uses default tag inference (lowercased field names) which works identically in v3.

## aurora v2 → v4

- All five `cmd/*.go` files importing `github.com/logrusorgru/aurora` switch to `github.com/logrusorgru/aurora/v4`.
- Function calls (`aurora.BrightYellow(...)`, `.String()`) are API-compatible across versions for these specific uses; verify during implementation.

## Verification

The implementation is complete only when all of these pass:

1. `go build ./...` succeeds with no errors or warnings.
2. `go vet ./...` reports no issues.
3. `go test ./...` passes. Specifically:
   - `cmd/list_test.go`: real tests, must pass.
   - `internal/app/init_test.go`: real tests, must pass.
   - Empty stub test files (`cmd/get_test.go`, `cmd/put_test.go`, `internal/app/config_test.go`, `internal/client/client_test.go`) compile (they're package-declaration-only).
4. `./esp --help` produces output (sanity check that cobra is wired up).
5. `./esp version` produces output.

**Manual smoke testing against a real AWS account is the user's responsibility** — the agent cannot perform this.

## Risks and Mitigations

| Risk | Mitigation |
|---|---|
| AWS SDK v2 typed errors are a different shape; some error paths might not be exercised by unit tests | Document which error types map to which v1 codes (table above); user does manual smoke against real AWS for confidence |
| `tablewriter` v1.x output drift | Visual output considered acceptable to shift slightly; functional output (row content) preserved |
| Viper v1.20 behavior changes around config search | Direct usage in this codebase is minimal (`viper.SetConfigName`, `viper.AddConfigPath`, `viper.ReadInConfig`, `viper.Unmarshal`); these have been stable across the v1.x line |
| Cobra v1.10 deprecations | None expected for the simple `cobra.Command` definitions used here |
| Behavioral drift in error messages between v1 `awserr.Code()` strings and v2 typed-error `.Error()` output | Acceptable — current code prints `"SSM Error: %s"` which is not parsed downstream |

## Order of Implementation (preview, full plan in writing-plans phase)

1. Bump go directive + toolchain.
2. Replace jwalterweatherman with slog (smallest change, easy to verify in isolation).
3. Swap yaml.v2 → yaml.v3 (trivial).
4. Bump aurora v2 → v4 (trivial import + go.mod).
5. Bump survey, cobra, viper to latest (trivial).
6. Rewrite `cmd/get.go` for tablewriter v1 API.
7. Migrate `internal/ssm/` from aws-sdk-go v1 → v2 (largest change, last).
8. `go mod tidy`.
9. Run full verification suite.
