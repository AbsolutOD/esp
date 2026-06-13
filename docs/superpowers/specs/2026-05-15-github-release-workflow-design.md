# GitHub Actions CI + Release Workflow for `esp`

**Date:** 2026-05-15
**Status:** Approved (pending implementation)

## Goal

Automate build, test, and binary publication for the `esp` Go CLI. PRs and pushes to `main` get a fast feedback gate; SemVer tag pushes produce a GitHub Release with prebuilt binaries for the three target platforms.

## Non-goals

- Publishing a Docker image (GHCR). Can be added later as a separate workflow.
- Windows builds. Skipped intentionally — re-add if a user needs it.
- darwin/amd64 builds. Skipped intentionally per current usage; users on Intel Macs can `go install`.
- Publishing to "GitHub Packages" in any other sense (npm/Maven/etc.). Go modules are distributed via VCS tags; `go install github.com/AbsolutOD/esp@<tag>` will work as soon as a tag exists.
- Code signing / notarization. Out of scope for v1.

## Build targets

GoReleaser will produce exactly three binaries per release:

- `linux/amd64`
- `linux/arm64`
- `darwin/arm64`

`CGO_ENABLED=0` for all targets (pure-Go build, statically linked).

## Versioning

- **Tag convention:** SemVer, `vX.Y.Z` (e.g. `v0.3.0`). Pre-release suffixes (`-rc.1`, `-beta.2`) are supported by GoReleaser; tags matching those automatically produce a GitHub pre-release.
- **Version injection:** the binary learns its own version via `-ldflags -X` at build time. Three values are injected:
  - `version` — the SemVer string from the tag (e.g. `v0.3.0`)
  - `commit` — the short Git commit SHA
  - `date` — the build timestamp (RFC3339)
- **Local builds:** running `go build` outside CI leaves the defaults intact, so `esp version` prints `esp dev (commit none, built unknown)`. No tooling required to develop locally.

### `cmd/version.go` modification

Replace the hardcoded version string with overridable package-level variables:

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

func runVersion() error {
	fmt.Printf("esp %s (commit %s, built %s)\n", version, commit, date)
	return nil
}
```

The existing `cmd/version_test.go` only asserts that `runVersion()` returns no error — it does not pin the output string — so no test changes are required.

## File layout

```
.github/
  workflows/
    ci.yml          # PRs + push to main: vet, test, build, validate goreleaser config
    release.yml     # tag push (v*): full GoReleaser release
.goreleaser.yaml    # build/archive/checksum/changelog/release config
cmd/version.go      # modified: ldflags-injectable vars
```

## `.github/workflows/ci.yml`

**Triggers:** `pull_request` (any branch), `push` to `main`.

**Single job** on `ubuntu-latest`:

1. `actions/checkout@v4`
2. `actions/setup-go@v5` with `go-version-file: go.mod` (single source of truth for Go version — currently `go 1.26`)
3. `go vet ./...`
4. `go test ./...`
5. `go build -o esp .` (smoke-build the binary)
6. `goreleaser/goreleaser-action@v6` with `args: check` — validates `.goreleaser.yaml` syntax/options so a broken config is caught on PR, not at tag time.

Standard `GITHUB_TOKEN` is sufficient — no extra secrets.

## `.github/workflows/release.yml`

**Trigger:** `push` on tags matching `v*`.

**Permissions:**
```yaml
permissions:
  contents: write   # required to create a GitHub Release
```

**Single job** on `ubuntu-latest`:

1. `actions/checkout@v4` with `fetch-depth: 0` (full history needed for GoReleaser's changelog generation)
2. `actions/setup-go@v5` with `go-version-file: go.mod`
3. `goreleaser/goreleaser-action@v6` with `args: release --clean`, passing `GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}` in env

## `.goreleaser.yaml`

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
  # GoReleaser auto-marks tags matching common pre-release patterns
  # (e.g. v0.3.0-rc.1) as prereleases.
```

Notes:
- The `before:` hooks run `go mod tidy` (defensive — ensures `go.sum` is in sync) and `go test ./...` (gate: a failing test aborts the release before any artifact is produced).
- `{{.Tag}}` (e.g. `v0.3.0`) is injected as the version string so `esp version` displays the leading `v`. (`{{.Version}}` would strip it.)
- Archive `formats` is plural (modern GoReleaser v2 syntax).

## Release flow (operational)

1. Bump anything that needs bumping, merge to `main`.
2. Tag: `git tag v0.3.0 && git push origin v0.3.0`.
3. `release.yml` fires; GoReleaser runs:
   - `go mod tidy` + `go test ./...`
   - Cross-compiles 3 binaries with injected ldflags
   - Produces `.tar.gz` archives + `checksums.txt`
   - Generates changelog from commits since previous tag
   - Creates GitHub Release with all artifacts attached
4. Users install via:
   - `curl` + `tar` from the Release page, or
   - `go install github.com/AbsolutOD/esp@v0.3.0` (works because the tag exists in VCS)

## Failure modes & handling

| Failure | Outcome |
|---|---|
| Test fails on PR | `ci.yml` red, PR blocked. |
| `goreleaser check` fails on PR | `ci.yml` red. Bad config caught before tag time. |
| Test fails on tag push | GoReleaser `before:` hook aborts; no Release created; workflow shows red. |
| Build fails for one arch on tag push | GoReleaser aborts the whole release atomically; no partial artifacts uploaded. |
| Tag pushed to branch that hasn't merged latest code | Whatever's at the tag is what gets built. Operator concern, not a workflow concern. |

## Verification plan

- Local: `go build && ./esp version` confirms default-string output.
- PR validation: the PR introducing these files exercises `ci.yml` end-to-end (including `goreleaser check`).
- Release validation: cut a `v0.0.0-test.1` tag against a throwaway commit and inspect the resulting GitHub Release. If acceptable, delete the test tag/release and proceed with the real first SemVer release.

## Out of scope but worth noting

- **Docker / GHCR image:** GoReleaser supports `dockers:` and multi-arch manifests cleanly; this can be added in a follow-up without restructuring the workflow.
- **Homebrew tap:** GoReleaser supports `brews:` for auto-publishing a formula. Worth considering if macOS users need an easier install path.
- **SBOM / provenance:** Sigstore + SLSA provenance can be layered in via GoReleaser's `sboms:` and GitHub's attestation actions if supply-chain attestation matters later.
