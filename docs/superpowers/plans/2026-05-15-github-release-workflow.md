# GitHub Release Workflow Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add GitHub Actions CI (PR/main gate) and tag-triggered release (linux/amd64, linux/arm64, darwin/arm64) for the `esp` Go CLI using GoReleaser, with build-time version injection.

**Architecture:** Two workflows + one GoReleaser config + a small `cmd/version.go` refactor. `cmd/version.go` exposes overridable `version`/`commit`/`date` package vars; GoReleaser injects them via `-ldflags -X`. `ci.yml` runs vet/test/build/`goreleaser check` on PRs and `main`. `release.yml` fires on `v*` tags and runs GoReleaser to build, archive, checksum, and publish to a GitHub Release.

**Tech Stack:** Go 1.26 (from `go.mod`), GitHub Actions (`actions/checkout@v4`, `actions/setup-go@v5`, `goreleaser/goreleaser-action@v6`), GoReleaser v2 config syntax.

**Spec:** `docs/superpowers/specs/2026-05-15-github-release-workflow-design.md`

---

## File Structure

| File | Action | Responsibility |
|---|---|---|
| `cmd/version.go` | Modify | Hold injectable `version`/`commit`/`date` vars; expose `versionString()` helper for testability; `runVersion()` prints it. |
| `cmd/version_test.go` | Modify | Add a test for `versionString()` default output. |
| `.goreleaser.yaml` | Create | GoReleaser config: builds, archives, checksums, changelog, release. |
| `.github/workflows/ci.yml` | Create | PR + `main` gate. |
| `.github/workflows/release.yml` | Create | Tag-triggered (`v*`) release. |

---

## Task 1: Refactor `cmd/version.go` to use injectable vars

Replace the hardcoded version string with package-level vars that GoReleaser can override via ldflags. Extract a `versionString()` helper so the format is unit-testable. The existing `TestRunVersion` continues to pass; a new test pins the default output.

**Files:**
- Modify: `cmd/version.go`
- Modify: `cmd/version_test.go`

- [ ] **Step 1: Write the failing test**

Add this test to `cmd/version_test.go` (keep the existing `TestRunVersion` intact):

```go
func TestVersionString_Defaults(t *testing.T) {
	got := versionString()
	want := "esp dev (commit none, built unknown)"
	if got != want {
		t.Fatalf("versionString() = %q, want %q", got, want)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./cmd -run TestVersionString_Defaults -v`
Expected: FAIL with `undefined: versionString` (compile error).

- [ ] **Step 3: Implement `versionString()` + new vars**

Replace the contents of `cmd/version.go` with:

```go
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
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

func versionString() string {
	return fmt.Sprintf("esp %s (commit %s, built %s)", version, commit, date)
}

func runVersion() error {
	fmt.Println(versionString())
	return nil
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./cmd -v`
Expected: both `TestRunVersion` and `TestVersionString_Defaults` PASS.

- [ ] **Step 5: Build and run the binary to sanity-check output**

Run: `go build -o esp . && ./esp version`
Expected output: `esp dev (commit none, built unknown)`
Cleanup: `rm esp`

- [ ] **Step 6: Run the full test suite to confirm nothing else broke**

Run: `go test ./...`
Expected: all tests pass.

- [ ] **Step 7: Commit**

```bash
git add cmd/version.go cmd/version_test.go
git commit -m "refactor(version): inject version/commit/date via ldflags"
```

---

## Task 2: Add `.goreleaser.yaml`

Create the GoReleaser config that drives both `ci.yml` (via `goreleaser check`) and `release.yml` (via `goreleaser release`).

**Files:**
- Create: `.goreleaser.yaml`

- [ ] **Step 1: Create the config file**

Write `.goreleaser.yaml` with this exact content:

```yaml
version: 2

before:
  hooks:
    - go mod tidy
    - go test ./...

builds:
  - id: esp
    binary: esp
    main: .
    env:
      - CGO_ENABLED=0
    goos:
      - linux
      - darwin
    goarch:
      - amd64
      - arm64
    ignore:
      - goos: darwin
        goarch: amd64
    ldflags:
      - -s -w
      - -X github.com/AbsolutOD/esp/cmd.version={{.Tag}}
      - -X github.com/AbsolutOD/esp/cmd.commit={{.ShortCommit}}
      - -X github.com/AbsolutOD/esp/cmd.date={{.Date}}

archives:
  - id: esp
    name_template: "esp_{{ .Version }}_{{ .Os }}_{{ .Arch }}"
    formats: [tar.gz]
    files:
      - LICENSE
      - README.md

checksum:
  name_template: "checksums.txt"
  algorithm: sha256

changelog:
  sort: asc
  use: git
  groups:
    - title: Features
      regexp: '^.*?feat(\(.+\))??!?:.+$'
      order: 0
    - title: Fixes
      regexp: '^.*?fix(\(.+\))??!?:.+$'
      order: 1
    - title: Others
      order: 999
  filters:
    exclude:
      - '^docs:'
      - '^test:'
      - '^chore:'
      - merge conflict
      - Merge pull request
      - Merge remote-tracking branch
      - Merge branch

release:
  draft: false
```

- [ ] **Step 2: Validate the config (if `goreleaser` is installed locally)**

Run: `command -v goreleaser >/dev/null && goreleaser check || echo "goreleaser not installed locally — CI will validate"`
Expected: either `[✔] checking '.goreleaser.yaml' ... valid`, or the fallback message. If goreleaser IS installed and reports errors, fix them before moving on. If it's not installed, that's fine — `ci.yml` (Task 3) will run `goreleaser check` in CI.

- [ ] **Step 3: (Optional, only if goreleaser is installed) Dry-run a snapshot release**

Run: `command -v goreleaser >/dev/null && goreleaser release --snapshot --clean --skip=publish || echo "skipped"`
Expected: artifacts produced under `dist/` for the three target platforms, no errors.
Cleanup: `rm -rf dist/` afterwards. If skipped, that's fine.

- [ ] **Step 4: Commit**

```bash
git add .goreleaser.yaml
git commit -m "build: add GoReleaser config for cross-platform release"
```

---

## Task 3: Add `.github/workflows/ci.yml`

Create the CI workflow that runs on PRs and pushes to `main`. Validates code (vet, test, build) and the GoReleaser config itself.

**Files:**
- Create: `.github/workflows/ci.yml`

- [ ] **Step 1: Verify the `.github/workflows` directory does not yet exist**

Run: `ls -la .github 2>/dev/null || echo "no .github dir — will create"`
Expected: directory does not exist (output ends with "will create"). If it does already exist, that's fine — just create the workflows subdir under it.

- [ ] **Step 2: Create the CI workflow file**

Write `.github/workflows/ci.yml` with this exact content:

```yaml
name: CI

on:
  pull_request:
  push:
    branches: [main]

permissions:
  contents: read

jobs:
  test:
    name: Build, vet, test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod
          cache: true

      - name: go vet
        run: go vet ./...

      - name: go test
        run: go test ./...

      - name: go build
        run: go build -o esp .

      - name: Validate GoReleaser config
        uses: goreleaser/goreleaser-action@v6
        with:
          version: "~> v2"
          args: check
```

- [ ] **Step 3: Validate the YAML parses**

Run: `python3 -c 'import yaml; yaml.safe_load(open(".github/workflows/ci.yml"))' && echo OK`
Expected: `OK` printed with no errors.

- [ ] **Step 4: Commit**

```bash
git add .github/workflows/ci.yml
git commit -m "ci: add PR + main gate (vet, test, build, goreleaser check)"
```

---

## Task 4: Add `.github/workflows/release.yml`

Create the release workflow that fires on `v*` tag pushes and runs GoReleaser end-to-end.

**Files:**
- Create: `.github/workflows/release.yml`

- [ ] **Step 1: Create the release workflow file**

Write `.github/workflows/release.yml` with this exact content:

```yaml
name: Release

on:
  push:
    tags:
      - 'v*'

permissions:
  contents: write

jobs:
  release:
    name: GoReleaser
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod
          cache: true

      - name: Run GoReleaser
        uses: goreleaser/goreleaser-action@v6
        with:
          version: "~> v2"
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

- [ ] **Step 2: Validate the YAML parses**

Run: `python3 -c 'import yaml; yaml.safe_load(open(".github/workflows/release.yml"))' && echo OK`
Expected: `OK` printed with no errors.

- [ ] **Step 3: Run the full Go test suite one final time as a sanity check**

Run: `go test ./...`
Expected: all tests pass.

- [ ] **Step 4: Commit**

```bash
git add .github/workflows/release.yml
git commit -m "ci: add tag-triggered release workflow via GoReleaser"
```

---

## Post-Implementation Verification

These are not commit-producing steps — they're how to confirm the workflows actually work once merged.

1. **Open a PR** with the four commits. `CI` workflow should run and pass (vet, test, build, GoReleaser check).
2. **Merge to `main`.** Same `CI` workflow should run on the merge commit and pass.
3. **Cut a test tag** to validate the release path before the real `v0.3.0`:
   ```sh
   git tag v0.0.0-test.1 && git push origin v0.0.0-test.1
   ```
   The `Release` workflow should fire, GoReleaser should produce 3 archives + `checksums.txt`, and a GitHub Pre-release should appear (GoReleaser auto-marks `-test.N`-style tags as prereleases). Download one archive, extract, run `./esp version` — should print `esp v0.0.0-test.1 (commit <sha>, built <date>)`.
4. **If the test release looks right**, delete the test tag and release from GitHub, then push `v0.3.0` (or whatever real first version you want).

---

## Notes for the implementer

- **Commit message style:** This repo's recent commits use conventional prefixes (`feat:`, `fix:`, `docs:`, `ci:`, `build:`, `refactor:`). Match that style.
- **No Co-Authored-By trailer:** The user's auto-memory specifies omitting Claude trailers from commits. The commit messages above already exclude it; keep it that way.
- **Don't push tags yourself.** Tagging is the user's call — they may want to coordinate with a release-notes review.
- **Don't push the branch yourself.** Open a PR is a user action.
