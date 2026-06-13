# GitHub Actions CI + Release Workflow Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add automated build/test on PRs and `main`, and tag-triggered binary releases for `esp` via GoReleaser.

**Architecture:** Two GitHub Actions workflows. `ci.yml` runs `go vet`, `go test`, `go build`, and `goreleaser check` on every PR and push to `main`. `release.yml` triggers on `v*` SemVer tags, invokes GoReleaser with the repo's default `GITHUB_TOKEN`, and produces a GitHub Release containing three cross-compiled binaries (`linux/amd64`, `linux/arm64`, `darwin/arm64`) with version/commit/date injected via ldflags. `cmd/version.go` was already updated to read these from package-level vars.

**Tech Stack:** GitHub Actions, GoReleaser v2, Go 1.26.

**Spec:** [`docs/superpowers/specs/2026-05-15-github-release-workflow-design.md`](../specs/2026-05-15-github-release-workflow-design.md) (commit `ab3554c`).

---

## Pre-flight (one-time, optional but recommended)

Local validation is faster than waiting for CI. Install GoReleaser if you want to validate the config and do snapshot builds locally:

```bash
brew install goreleaser
```

If you skip this, the `goreleaser check` step in `ci.yml` will validate the config on the PR — you just won't see failures until CI runs. Plan steps that say "run locally" will show CI as a fallback.

---

### Task 1: Add `.goreleaser.yaml`

**Files:**
- Create: `.goreleaser.yaml`

- [ ] **Step 1: Create `.goreleaser.yaml`**

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

- [ ] **Step 2: Validate config syntax**

If GoReleaser is installed locally:
```bash
goreleaser check
```
Expected output: `• checking config file • config is valid`

If GoReleaser is not installed: skip this step. The `goreleaser check` step in `ci.yml` (Task 2) will validate it on the PR.

- [ ] **Step 3: (Optional) Snapshot build to prove ldflags injection works**

Only if GoReleaser is installed locally:
```bash
goreleaser release --snapshot --clean
```
Expected: builds complete, `dist/` directory populated. Inspect a binary:
```bash
ls dist/                          # find the linux_amd64 subdirectory
dist/esp_linux_amd64_v1/esp version   # path may differ slightly by GoReleaser version
```
Expected output shape: `esp <version> (commit <sha>, built <timestamp>)` where `<sha>` is a 7-char hex string and `<timestamp>` is an RFC3339 date — confirms ldflags are reaching the binary. (The `<version>` string in snapshot mode is GoReleaser-synthesized, not a real tag — that's fine; we're verifying injection, not the exact value.)

Clean up after:
```bash
rm -rf dist/
```

If GoReleaser is not installed: skip. Real verification happens when the first `v*` tag is pushed.

- [ ] **Step 4: Commit**

```bash
git add .goreleaser.yaml
git commit -m "ci: add GoReleaser config for v* tag releases"
```

---

### Task 2: Add `.github/workflows/ci.yml`

**Files:**
- Create: `.github/workflows/ci.yml`

- [ ] **Step 1: Create the workflows directory**

```bash
mkdir -p .github/workflows
```

- [ ] **Step 2: Create `ci.yml`**

```yaml
name: CI

on:
  pull_request:
  push:
    branches: [main]

jobs:
  test:
    name: vet, test, build
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod
          cache: true

      - name: Vet
        run: go vet ./...

      - name: Test
        run: go test ./...

      - name: Build
        run: go build -o esp .

      - name: GoReleaser config check
        uses: goreleaser/goreleaser-action@v6
        with:
          version: latest
          args: check
```

- [ ] **Step 3: Local sanity check — same commands the workflow runs**

```bash
go vet ./...
go test ./...
go build -o esp .
```
Expected: all three succeed. The `esp` binary is produced in the working directory.

Clean up:
```bash
rm esp
```

- [ ] **Step 4: Validate YAML syntax**

```bash
python3 -c "import yaml; yaml.safe_load(open('.github/workflows/ci.yml'))" && echo "ok"
```
Expected output: `ok`. (If `python3` isn't available, any YAML linter or just visual inspection works — the real validation happens when GitHub parses the file on push.)

- [ ] **Step 5: Commit**

```bash
git add .github/workflows/ci.yml
git commit -m "ci: add CI workflow for PRs and main (vet, test, build, goreleaser check)"
```

---

### Task 3: Add `.github/workflows/release.yml`

**Files:**
- Create: `.github/workflows/release.yml`

- [ ] **Step 1: Create `release.yml`**

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

      - name: Release
        uses: goreleaser/goreleaser-action@v6
        with:
          version: latest
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

- [ ] **Step 2: Validate YAML syntax**

```bash
python3 -c "import yaml; yaml.safe_load(open('.github/workflows/release.yml'))" && echo "ok"
```
Expected output: `ok`.

- [ ] **Step 3: Commit**

```bash
git add .github/workflows/release.yml
git commit -m "ci: add release workflow triggered by v* SemVer tags"
```

---

### Task 4: Open PR and verify CI runs green

**Files:** none modified — verification only.

- [ ] **Step 1: Push the branch and open a PR**

```bash
git push -u origin HEAD
gh pr create --title "ci: add CI + release workflows" --body "$(cat <<'EOF'
## Summary
- Adds `.github/workflows/ci.yml` — runs vet/test/build/goreleaser-check on PRs and pushes to main
- Adds `.github/workflows/release.yml` — fires on `v*` SemVer tags, releases via GoReleaser
- Adds `.goreleaser.yaml` — builds linux/amd64, linux/arm64, darwin/arm64 with version/commit/date injected via ldflags

Spec: `docs/superpowers/specs/2026-05-15-github-release-workflow-design.md`

## Test plan
- [ ] CI workflow runs green on this PR (vet, test, build, goreleaser check)
- [ ] After merge, cut a throwaway tag (e.g. `v0.0.0-test.1`) to exercise the release workflow end-to-end, then delete the test release/tag before the first real release
EOF
)"
```

- [ ] **Step 2: Watch CI**

```bash
gh pr checks --watch
```
Expected: the `CI / vet, test, build` check passes. If `goreleaser check` fails, the `.goreleaser.yaml` has a syntax error — fix it and push.

- [ ] **Step 3: (Post-merge, manual) Smoke test the release workflow**

Out of scope for this plan, but worth doing once: after merging the PR, push a throwaway pre-release tag:

```bash
git checkout main && git pull
git tag v0.0.0-test.1
git push origin v0.0.0-test.1
gh run watch
```
Expected: the `Release` workflow runs and produces a GitHub Release with three `.tar.gz` archives + `checksums.txt`. Inspect a binary:

```bash
gh release download v0.0.0-test.1 -p '*linux_amd64*' -D /tmp/esp-test
tar -xzf /tmp/esp-test/esp_*_linux_amd64.tar.gz -C /tmp/esp-test
/tmp/esp-test/esp version
```
Expected output: `esp v0.0.0-test.1 (commit <sha>, built <date>)`.

Clean up:
```bash
gh release delete v0.0.0-test.1 --yes
git push --delete origin v0.0.0-test.1
git tag -d v0.0.0-test.1
```

---

## Notes for the implementer

- **Do not modify `cmd/version.go`** — it's already in the desired state (package-level `version`/`commit`/`date` vars + `PersistentPreRunE` no-op so `esp version` doesn't require AWS env vars).
- **Do not add tests for `versionString()`** — the existing `TestRunVersion` covers the no-error contract, which is sufficient. The format is verified end-to-end when the release smoke test runs.
- **Action versions** in this plan (`actions/checkout@v4`, `actions/setup-go@v5`, `goreleaser/goreleaser-action@v6`) are current as of writing. If a newer major is out by execution time, prefer the newer one and update the plan accordingly.
- **`{{.Tag}}` vs `{{.Version}}`** in ldflags: `{{.Tag}}` preserves the leading `v` (so `esp version` prints `esp v0.3.0`). `{{.Version}}` strips it. We want `{{.Tag}}` here.
- **`PersistentPreRunE` no-op on `versionCmd`**: the root command's `PersistentPreRunE` enforces `AWS_DEFAULT_REGION`/`AWS_PROFILE`. Overriding it on `versionCmd` (even with a no-op) replaces the inherited check, so `esp version` works without those env vars. This was added by the user when implementing the spec.
