# Dependency Refresh Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Bring `esp` onto Go 1.26 and refresh every direct dependency, replacing deprecated libs (aws-sdk-go v1 → v2, jwalterweatherman → log/slog, yaml.v2 → yaml.v3, aurora v2 → v4) without changing user-facing CLI behavior.

**Architecture:** Mechanical upgrades grouped by dependency, smallest-impact first, biggest (AWS SDK) last. Each task ends with `go build ./...`, `go vet ./...`, `go test ./...`, and a commit so the tree is bisectable.

**Tech Stack:** Go 1.26, AWS SDK for Go v2, Cobra, Viper, AlecAivazis/survey, logrusorgru/aurora, olekukonko/tablewriter, gopkg.in/yaml.v3, log/slog.

**Spec:** `docs/superpowers/specs/2026-05-03-dependency-refresh-design.md`

---

## Conventions used in every task

- All commands run from repo root (`/Users/modonnell/workspaces/instride/esp`).
- Verification commands at the end of every task:
  - `go build ./...` → expect: no output, exit 0
  - `go vet ./...` → expect: no output, exit 0
  - `go test ./...` → expect: PASS for `cmd` and `internal/app` packages; other packages "ok" or "no test files"
- Commits use `git add <file paths>` (never `-A`/`.`) to avoid staging unintended files.
- `AWS_DEFAULT_REGION` and `AWS_PROFILE` must be set in the environment when invoking the binary (not for `go build`/`go test`). The agent does NOT need real AWS credentials for this plan.

---

## Task 1: Bump Go directive and toolchain

**Files:**
- Modify: `go.mod` (lines 1–3)

- [ ] **Step 1: Edit `go.mod`**

Replace:
```
module github.com/pinpt/esp

go 1.13
```

With:
```
module github.com/pinpt/esp

go 1.26

toolchain go1.26.2
```

- [ ] **Step 2: Verify build still works (no dep changes yet)**

```
go build ./...
go vet ./...
go test ./...
```

Expected: all pass. The Go directive bump alone should not break anything; the existing deps still satisfy their go.mod requirements.

- [ ] **Step 3: Commit**

```
git add go.mod
git commit -m "chore: bump go directive to 1.26 with toolchain go1.26.2"
```

---

## Task 2: Replace `jwalterweatherman` with `log/slog`

**Files:**
- Modify: `cmd/root.go`
- Modify: `internal/app/config.go`
- Modify: `go.mod` (jwalterweatherman line removed by `go mod tidy` at end)

- [ ] **Step 1: Edit `cmd/root.go` imports**

Replace lines 3–12 (the import block):

```go
import (
	"fmt"
	"github.com/pinpt/esp/internal/app"
	"github.com/pinpt/esp/internal/client"
	jww "github.com/spf13/jwalterweatherman"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)
```

With:

```go
import (
	"fmt"
	"log/slog"
	"os"

	"github.com/pinpt/esp/internal/app"
	"github.com/pinpt/esp/internal/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)
```

- [ ] **Step 2: Replace verbose-flag block in `cmd/root.go` `initConfig`**

Locate (around line 60–62):

```go
if verbose {
    jww.SetStdoutThreshold(jww.LevelInfo)
}
```

Replace with:

```go
level := slog.LevelWarn
if verbose {
    level = slog.LevelInfo
}
slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})))
```

- [ ] **Step 3: Edit `internal/app/config.go` imports**

Replace lines 3–7:

```go
import (
	"fmt"
	"github.com/pinpt/esp/internal/utils"
	jww "github.com/spf13/jwalterweatherman"
)
```

With:

```go
import (
	"fmt"
	"log/slog"

	"github.com/pinpt/esp/internal/utils"
)
```

- [ ] **Step 4: Replace the three `jww.INFO.Printf` call sites in `internal/app/config.go`**

Line 34 — replace:
```go
jww.INFO.Printf("rendered path: %s", path)
```
With:
```go
slog.Info("rendered path", "path", path)
```

Line 41 — replace:
```go
jww.INFO.Printf("rendered param path: %s", path)
```
With:
```go
slog.Info("rendered param path", "path", path)
```

Line 48 — replace:
```go
jww.INFO.Printf("found ENV: %s", e)
```
With:
```go
slog.Info("found ENV", "env", e)
```

- [ ] **Step 5: Verify**

```
go build ./...
go vet ./...
go test ./...
```

Expected: all pass. (`jwalterweatherman` is still in `go.mod`/`go.sum` but no longer imported; that's fine — `go mod tidy` in Task 8 will remove it.)

- [ ] **Step 6: Commit**

```
git add cmd/root.go internal/app/config.go
git commit -m "refactor: replace jwalterweatherman with log/slog"
```

---

## Task 3: Swap `gopkg.in/yaml.v2` → `gopkg.in/yaml.v3`

**Files:**
- Modify: `internal/app/init.go` (line 6)
- Modify: `internal/app/init_test.go` (line 4)
- Modify: `go.mod`

- [ ] **Step 1: Add yaml.v3 and remove yaml.v2 from go.mod**

```
go get gopkg.in/yaml.v3@latest
```

(This adds yaml.v3. The yaml.v2 line stays until tidy, but that's fine.)

- [ ] **Step 2: Update `internal/app/init.go` import**

Replace line 6:
```go
	"gopkg.in/yaml.v2"
```

With:
```go
	"gopkg.in/yaml.v3"
```

- [ ] **Step 3: Update `internal/app/init_test.go` import**

Replace line 4:
```go
	"gopkg.in/yaml.v2"
```

With:
```go
	"gopkg.in/yaml.v3"
```

- [ ] **Step 4: Verify**

```
go build ./...
go vet ./...
go test ./...
```

Expected: all pass. `TestWriteConfig` in `internal/app/init_test.go` exercises `yaml.Marshal` and `yaml.Unmarshal` — both are API-compatible between v2 and v3 for this struct.

- [ ] **Step 5: Commit**

```
git add internal/app/init.go internal/app/init_test.go go.mod go.sum
git commit -m "deps: upgrade yaml.v2 to yaml.v3"
```

---

## Task 4: Swap `logrusorgru/aurora` v2 → v4

**Files:** (4 files, all in `cmd/`)
- Modify: `cmd/get.go` (import line)
- Modify: `cmd/list.go` (import line)
- Modify: `cmd/delete.go` (import line)
- Modify: `cmd/move.go` (import line)
- Modify: `go.mod`

- [ ] **Step 1: Add aurora v4 to go.mod**

```
go get github.com/logrusorgru/aurora/v4@latest
```

- [ ] **Step 2: Update `cmd/get.go` import**

Replace:
```go
	"github.com/logrusorgru/aurora"
```

With:
```go
	"github.com/logrusorgru/aurora/v4"
```

- [ ] **Step 3: Update `cmd/list.go` import**

Same replacement as Step 2.

- [ ] **Step 4: Update `cmd/delete.go` import**

Same replacement as Step 2.

- [ ] **Step 5: Update `cmd/move.go` import**

Same replacement as Step 2.

- [ ] **Step 6: Verify**

```
go build ./...
go vet ./...
go test ./...
```

Expected: all pass. The function calls used (`aurora.BrightYellow(x)`, `.String()`) are unchanged across versions. The package name is still `aurora` after `/v4` since Go's "v2+" rule uses the directory but the package name stays.

If the build fails because the package name is actually `aurora` from a `v4` subpackage (it should be, but verify), check the package's repo for the exact name.

- [ ] **Step 7: Commit**

```
git add cmd/get.go cmd/list.go cmd/delete.go cmd/move.go go.mod go.sum
git commit -m "deps: upgrade logrusorgru/aurora v2 to v4"
```

---

## Task 5: Rewrite `cmd/get.go` for tablewriter v1 API

**Files:**
- Modify: `cmd/get.go` (the `detailDisplay` function, lines 28–41)
- Modify: `go.mod`

**Note:** The tablewriter v1.x API replaced `SetHeader`/`AppendBulk`/`Render` with a different shape (e.g. `Header`/`Bulk`/`Render`). Before editing, verify the exact method names against the package's current docs.

- [ ] **Step 1: Upgrade tablewriter**

```
go get github.com/olekukonko/tablewriter@latest
```

- [ ] **Step 2: Check the v1.x API**

Run:
```
go doc github.com/olekukonko/tablewriter
```

Look for: a constructor (likely `NewTable` or `NewWriter`), how to set headers, how to append rows in bulk, how to render.

If the package doc isn't clear, check `pkg.go.dev/github.com/olekukonko/tablewriter` or the source in `~/go/pkg/mod/github.com/olekukonko/tablewriter@<version>/`.

- [ ] **Step 3: Rewrite `detailDisplay` in `cmd/get.go`**

Replace the function body (lines 28–41) so it produces a 2-column table with header `["Keys", "Value"]` and one row per data entry, written to `os.Stdout`. Preserve the `aurora.BrightYellow(...)` styling on the row labels.

The shape will be approximately:

```go
func detailDisplay(p common.EspParam) {
	data := [][]string{
		{aurora.BrightYellow("ID").String(), p.Id},
		{aurora.BrightYellow("Last_Modified").String(), p.LastModifiedDate.String()},
		{aurora.BrightYellow("Name").String(), p.Name},
		{aurora.BrightYellow("Type").String(), p.Type},
		{aurora.BrightYellow("Value").String(), p.Value},
		{aurora.BrightYellow("Version").String(), strconv.FormatInt(p.Version, 10)},
	}
	table := tablewriter.NewTable(os.Stdout)  // or NewWriter — verify in Step 2
	table.Header("Keys", "Value")             // or table.Header([]string{"Keys", "Value"}) — verify
	if err := table.Bulk(data); err != nil {  // method name may be Bulk / AppendBulk — verify
		fmt.Fprintf(os.Stderr, "table render error: %v\n", err)
		return
	}
	if err := table.Render(); err != nil {
		fmt.Fprintf(os.Stderr, "table render error: %v\n", err)
	}
}
```

Adjust method names to match what `go doc` showed in Step 2. If `Render()` returns no error in v1, drop the error handling for it.

- [ ] **Step 4: Verify**

```
go build ./...
go vet ./...
go test ./...
```

Expected: all pass.

- [ ] **Step 5: Smoke test the table output (optional but recommended)**

Run:
```
AWS_DEFAULT_REGION=us-east-1 AWS_PROFILE=dummy ./esp --help
```

Expected: `--help` output for `esp` (this verifies cobra still wires up; the real table render needs AWS, which is out of scope here).

- [ ] **Step 6: Commit**

```
git add cmd/get.go go.mod go.sum
git commit -m "deps: upgrade tablewriter to v1.x and rewrite detailDisplay"
```

---

## Task 6: Migrate AWS SDK v1 → v2 (`internal/ssm/`)

**Files:**
- Modify: `internal/ssm/ssm.go` (full rewrite of imports + body)
- Modify: `internal/ssm/utils.go` (full rewrite)
- Modify: `internal/ssm/errors.go` (full rewrite)
- Modify: `go.mod`

This is the largest change. Done as one task because the three files are interlocked — partial edits won't compile.

- [ ] **Step 1: Add AWS SDK v2 modules**

```
go get github.com/aws/aws-sdk-go-v2/aws@latest
go get github.com/aws/aws-sdk-go-v2/config@latest
go get github.com/aws/aws-sdk-go-v2/service/ssm@latest
go get github.com/aws/smithy-go@latest
```

- [ ] **Step 2: Replace `internal/ssm/utils.go` entirely**

Write the file as:

```go
package ssm

import (
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	ssmtypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/pinpt/esp/internal/common"
)

// ParamType sets the base type for SSM parameter types
type ParamType string

// Defines the SSM types
const (
	String       ParamType = "string"
	SecureString ParamType = "SecureString"
	StringList   ParamType = "Stringlist"
)

// AwsParam represents an individual SSM parameter
type AwsParam struct {
	Arn              string
	Name             string
	Type             ParamType
	Value            string
	Version          int
	LastModifiedDate float32
}

func (p *AwsParam) isValid() error {
	switch p.Type {
	case String, SecureString, StringList:
		return nil
	}
	return errors.New("invalid SSM Parameter Type")
}

func selectType(t bool) ssmtypes.ParameterType {
	if t {
		return ssmtypes.ParameterTypeSecureString
	}
	return ssmtypes.ParameterTypeString
}

func convertToEspParam(ap ssmtypes.Parameter) common.EspParam {
	param := common.EspParam{
		Id:               aws.ToString(ap.ARN),
		Name:             aws.ToString(ap.Name),
		Type:             string(ap.Type),
		Value:            aws.ToString(ap.Value),
		Version:          ap.Version,
		LastModifiedDate: aws.ToTime(ap.LastModifiedDate),
	}

	if param.Type == string(ssmtypes.ParameterTypeSecureString) {
		param.Secure = true
	}
	return param
}

// handleAwsErr it will for all of the AWS API errors and exit if exists
func handleAwsErr(a action, err error) {
	awsErr := checkSSMError(a, err)
	if awsErr != nil {
		fmt.Printf("SSM Error: %s\n", awsErr.Error())
		os.Exit(1)
	}
}
```

Notes:
- `convertToEspParam` now takes `ssmtypes.Parameter` by value (v2 returns slices of values, not pointers). The single-call `GetOne` will need to dereference its `*ssmtypes.Parameter` result; we'll handle that in Step 3.
- The `AwsParam` struct is preserved as-is even though it's unused dead code (out of scope to remove).

- [ ] **Step 3: Replace `internal/ssm/ssm.go` entirely**

Write the file as:

```go
package ssm

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	awsssm "github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/pinpt/esp/internal/common"
	"github.com/pinpt/esp/internal/utils"
)

type action string

// Constants to represent actions to take against SSM
const (
	Get     action = "get"
	GetMany action = "getMany"
	//Put     action = "put"
	Save    action = "save"
	Delete  action = "delete"
)

// Service the struct representing AWS service/Session
type Service struct {
	Svc    *awsssm.Client
	Region string
	Cfg    aws.Config
}

// Save a single param for a given path
func (s *Service) Save(p common.EspParamInput) common.SaveOutput {
	pi := &awsssm.PutParameterInput{
		Type:  selectType(p.Secure),
		Name:  aws.String(p.Name),
		Value: aws.String(p.Value),
	}
	param, err := s.Svc.PutParameter(context.Background(), pi)
	if err != nil {
		handleAwsErr(Save, err)
	}
	return common.SaveOutput{Version: param.Version}
}

// Delete a single param for a given path
func (s *Service) Delete(p common.DeleteInput) string {
	dpi := &awsssm.DeleteParameterInput{
		Name: aws.String(p.Name),
	}
	_, err := s.Svc.DeleteParameter(context.Background(), dpi)
	if err != nil {
		handleAwsErr(Delete, err)
	}
	return p.Name
}

// GetOne gets a single param for a given path
func (s *Service) GetOne(p common.GetOneInput) common.EspParam {
	si := &awsssm.GetParameterInput{
		Name:           aws.String(p.Name),
		WithDecryption: aws.Bool(p.Decrypt),
	}
	resp, err := s.Svc.GetParameter(context.Background(), si)
	if err != nil {
		handleAwsErr(Get, err)
	}
	return convertToEspParam(*resp.Parameter)
}

// GetMany recursively gets parameters from a given path
func (s *Service) GetMany(p common.ListParamInput) []common.EspParam {
	si := &awsssm.GetParametersByPathInput{
		Path:           aws.String(p.Path),
		WithDecryption: aws.Bool(p.Decrypt),
		Recursive:      aws.Bool(p.Recursive),
	}
	params, err := s.Svc.GetParametersByPath(context.Background(), si)
	if err != nil {
		handleAwsErr(GetMany, err)
	}

	var espParams []common.EspParam
	for _, v := range params.Parameters {
		espParams = append(espParams, convertToEspParam(v))
	}

	if params.NextToken != nil {
		si.NextToken = params.NextToken
		moreParams := s.getNextParams(si)
		espParams = append(espParams, moreParams...)
	}
	return espParams
}

// getNextParams uses the NextToken to get more params
func (s *Service) getNextParams(pi *awsssm.GetParametersByPathInput) []common.EspParam {
	params, err := s.Svc.GetParametersByPath(context.Background(), pi)
	if err != nil {
		handleAwsErr(GetMany, err)
	}

	var espParams []common.EspParam
	for _, v := range params.Parameters {
		espParams = append(espParams, convertToEspParam(v))
	}

	if params.NextToken != nil {
		pi.NextToken = params.NextToken
		moreParams := s.getNextParams(pi)
		espParams = append(espParams, moreParams...)
	}
	return espParams
}

// Copy method copies the given parameter to a new location
func (s *Service) Copy(cc common.CopyCommand) common.SaveOutput {
	input := common.GetOneInput{
		Name:    cc.Source,
		Decrypt: true,
	}
	sparam := s.GetOne(input)

	dparam := common.EspParamInput{
		Name:   cc.Destination,
		Secure: sparam.Secure,
		Value:  sparam.Value,
	}
	return s.Save(dparam)
}

// New Create a new SSM service
func New() *Service {
	svc := new(Service)
	svc.Region = utils.GetEnv("AWS_REGION", "us-east-1")
	return svc
}

// Init create the actual session to talk to the AWS API
func (s *Service) Init() {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(s.Region))
	if err != nil {
		fmt.Printf("AWS config load error: %s\n", err.Error())
		os.Exit(1)
	}
	s.Cfg = cfg
	s.Svc = awsssm.NewFromConfig(cfg)
}
```

Key changes from v1:
- Dropped `*session.Session` field; v2 uses `aws.Config` only.
- `New()` no longer pre-populates `aws.Config` (it has to wait for `LoadDefaultConfig`).
- `Init()` calls `config.LoadDefaultConfig` and exits 1 on failure (matches existing fail-fast behavior).
- All API calls take `context.Background()` as first arg.
- `param.Version` (was `*int64`) → `param.Version` (now `int64`).
- `*resp.Parameter` dereferenced because `convertToEspParam` now takes value, not pointer.

- [ ] **Step 4: Replace `internal/ssm/errors.go` entirely**

Write the file as:

```go
package ssm

import (
	"errors"

	ssmtypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/aws/smithy-go"
)

// checkSSMError is the entry point for check all of the based and call specific errors
func checkSSMError(a action, err error) error {
	awsErr := checkBaseSSMErrors(err)
	switch a {
	case Get:
		return checkSSMGetParameterError(err)
	case Save:
		return checkSSMPutParameterError(err)
	case Delete:
		return checkDeleteParameterError(err)
	}
	return awsErr
}

// checkBaseSSMErrors checks for the common errors all SSM API calls might return
func checkBaseSSMErrors(err error) error {
	var invalidKey *ssmtypes.InvalidKeyId
	if errors.As(err, &invalidKey) {
		return err
	}
	var internalErr *ssmtypes.InternalServerError
	if errors.As(err, &internalErr) {
		return err
	}
	// Generic API error fallback (covers MissingRegion-style errors and any unmapped types).
	var apiErr smithy.APIError
	if errors.As(err, &apiErr) {
		return err
	}
	return nil
}

func checkDeleteParameterError(err error) error {
	var notFound *ssmtypes.ParameterNotFound
	if errors.As(err, &notFound) {
		return err
	}
	var apiErr smithy.APIError
	if errors.As(err, &apiErr) {
		return err
	}
	return nil
}

// checkSSMGetParameterError checks for errors the GetParameter API call might return
func checkSSMGetParameterError(err error) error {
	var notFound *ssmtypes.ParameterNotFound
	if errors.As(err, &notFound) {
		return err
	}
	var apiErr smithy.APIError
	if errors.As(err, &apiErr) {
		return err
	}
	return nil
}

// checkSSMPutParameterError checks for errors the PutParameter API call might return
func checkSSMPutParameterError(err error) error {
	var limitExceeded *ssmtypes.ParameterLimitExceeded
	if errors.As(err, &limitExceeded) {
		return err
	}
	var tooMany *ssmtypes.TooManyUpdates
	if errors.As(err, &tooMany) {
		return err
	}
	var hierarchyMismatch *ssmtypes.HierarchyTypeMismatchException
	if errors.As(err, &hierarchyMismatch) {
		return err
	}
	var invalidPattern *ssmtypes.InvalidAllowedPatternException
	if errors.As(err, &invalidPattern) {
		return err
	}
	var maxVersion *ssmtypes.ParameterMaxVersionLimitExceeded
	if errors.As(err, &maxVersion) {
		return err
	}
	var unsupportedType *ssmtypes.UnsupportedParameterType
	if errors.As(err, &unsupportedType) {
		return err
	}
	var policyLimit *ssmtypes.PoliciesLimitExceededException
	if errors.As(err, &policyLimit) {
		return err
	}
	var invalidPolicyType *ssmtypes.InvalidPolicyTypeException
	if errors.As(err, &invalidPolicyType) {
		return err
	}
	var invalidPolicyAttr *ssmtypes.InvalidPolicyAttributeException
	if errors.As(err, &invalidPolicyAttr) {
		return err
	}
	var incompatible *ssmtypes.IncompatiblePolicyException
	if errors.As(err, &incompatible) {
		return err
	}
	var alreadyExists *ssmtypes.ParameterAlreadyExists
	if errors.As(err, &alreadyExists) {
		return err
	}
	var apiErr smithy.APIError
	if errors.As(err, &apiErr) {
		return err
	}
	return nil
}

// checkSSMByPathError checks for errors the GetParameterByPath API call might return
func checkSSMByPathError(err error) error {
	var internalErr *ssmtypes.InternalServerError
	if errors.As(err, &internalErr) {
		return err
	}
	var invalidFilterKey *ssmtypes.InvalidFilterKey
	if errors.As(err, &invalidFilterKey) {
		return err
	}
	var invalidFilterOption *ssmtypes.InvalidFilterOption
	if errors.As(err, &invalidFilterOption) {
		return err
	}
	var invalidFilterValue *ssmtypes.InvalidFilterValue
	if errors.As(err, &invalidFilterValue) {
		return err
	}
	var invalidKey *ssmtypes.InvalidKeyId
	if errors.As(err, &invalidKey) {
		return err
	}
	var invalidNextToken *ssmtypes.InvalidNextToken
	if errors.As(err, &invalidNextToken) {
		return err
	}
	return nil
}
```

Key changes from v1:
- `awserr.Error` interface gone; replaced with `errors.As` against typed errors and `smithy.APIError` as the generic fallback.
- The `MissingRegion` v1 awserr code is no longer surfaced as a typed error in v2 — `LoadDefaultConfig` returns it directly, so `Init()` (in `ssm.go`) handles it.
- `checkRegion` helper deleted (no longer needed).
- All checks now return the original `err` (rather than the awserr wrapper) so `handleAwsErr` prints the v2 error message verbatim.

- [ ] **Step 5: Verify**

```
go build ./...
go vet ./...
go test ./...
```

Expected: all pass. The existing tests (`cmd/list_test.go`, `internal/app/init_test.go`) don't touch the SSM package, so they'll still pass without an AWS connection.

If a build error mentions a typed error like `*ssmtypes.SomeError` not existing, check the v2 SSM types package — names sometimes differ slightly from v1 codes (e.g. v2 may use `ParameterAlreadyExists` vs v1's `ErrCodeParameterAlreadyExists`). Use `go doc github.com/aws/aws-sdk-go-v2/service/ssm/types | grep Error` to find the correct names and adjust.

- [ ] **Step 6: Commit**

```
git add internal/ssm/ssm.go internal/ssm/utils.go internal/ssm/errors.go go.mod go.sum
git commit -m "deps: migrate aws-sdk-go v1 to v2"
```

---

## Task 7: Bump cobra, viper, survey to latest

**Files:**
- Modify: `go.mod`, `go.sum`

- [ ] **Step 1: Upgrade**

```
go get github.com/spf13/cobra@latest
go get github.com/spf13/viper@latest
go get github.com/AlecAivazis/survey/v2@latest
```

- [ ] **Step 2: Verify**

```
go build ./...
go vet ./...
go test ./...
```

Expected: all pass. The codebase uses only the most basic API surface of all three (`cobra.Command`, `cobra.MinimumNArgs`, `viper.SetConfigName`, `viper.AddConfigPath`, `viper.ReadInConfig`, `viper.Unmarshal`, `survey.Ask`) — all stable across recent releases.

If the build fails, read the error and fix the breaking change inline. Report any non-trivial breakage to the user before continuing.

- [ ] **Step 3: Commit**

```
git add go.mod go.sum
git commit -m "deps: upgrade cobra, viper, and survey/v2 to latest"
```

---

## Task 8: Final tidy and verification

**Files:**
- Modify: `go.mod`, `go.sum`

- [ ] **Step 1: Tidy**

```
go mod tidy
```

This removes `jwalterweatherman`, `aws-sdk-go` v1, `yaml.v2`, and `aurora` v2 from `go.mod`/`go.sum` along with any orphaned indirects.

- [ ] **Step 2: Confirm the deprecated direct deps are gone**

```
grep -E "jwalterweatherman|aws-sdk-go [^v]|aws-sdk-go$|yaml.v2|aurora v2" go.mod
```

Expected: no output (none found).

- [ ] **Step 3: Final verification**

```
go build ./...
go vet ./...
go test ./...
```

Expected: all pass.

- [ ] **Step 4: Smoke test the binary**

```
go build -o esp .
AWS_DEFAULT_REGION=us-east-1 AWS_PROFILE=dummy ./esp --help
AWS_DEFAULT_REGION=us-east-1 AWS_PROFILE=dummy ./esp version
```

Expected: `--help` shows the usage; `version` prints the version string. Neither requires real AWS credentials.

If either fails with a runtime error other than the AWS env-var checks, debug before committing.

- [ ] **Step 5: Commit**

```
git add go.mod go.sum
git commit -m "deps: tidy go.mod after dependency refresh"
```

- [ ] **Step 6: Show the final dep state to the user**

Run:
```
go list -m -mod=mod all | grep -E "^(github.com/aws/aws-sdk-go|github.com/spf13|github.com/AlecAivazis|github.com/logrusorgru|github.com/olekukonko|gopkg.in/yaml)" | sort
```

Then summarize: Go version, AWS SDK v2, all direct deps current, deprecated libs removed.

---

## Out of scope (do NOT do these)

- Adding new tests
- Refactoring `cmd/` to use a shared `context.Context`
- Switching `GetParametersByPath` to `ssm.NewGetParametersByPathPaginator`
- Removing the unused `AwsParam` struct in `internal/ssm/utils.go`
- Removing the empty `setPrefixes` method in `internal/app/init.go`
- Touching `README.md`, `.github/workflows/`, or the module path
