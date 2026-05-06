# Rebrand Design

**Date:** 2026-05-05
**Status:** Approved
**Scope:** Remove Pinpoint Software branding from this fork, establish ownership under Matt O'Donnell, refresh the README to reflect what `esp` actually is today.

## Background

`esp` was originally created by Matt O'Donnell while employed at Pinpoint Software, Inc., which legally held the copyright through its operating period (2019–2020). Pinpoint has since gone out of business. This fork (`github.com/AbsolutOD/esp`) is the active continuation of the project, maintained by the original author.

The post-Pinpoint copyright period for the fork is **2021–2026**.

## Goals

- Change the Go module path to match the actual repo location.
- Rewrite the README to reflect the current state of the tool.
- Add explicit copyright attribution for the fork.
- Remove the legacy Pinpoint logo.

## Non-Goals

- No code changes. This PR touches only the module path, import statements, README, LICENSE, and the logo file.
- No license change. The project remains Apache License 2.0.
- No edit to the historical `docs/superpowers/specs/2026-05-03-dependency-refresh-design.md` or its plan — those reference the old module path because they were written under it; preserved as-is.
- No CI workflow changes.

## PR Plan

One PR off `main`: branch `feature/rebrand`, worktree `.worktrees/rebrand`. To be done in a separate session from PRs 1 and 2 of the cleanup-and-tests spec.

## Module Path Change

- `go.mod` line 1: `module github.com/pinpt/esp` → `module github.com/AbsolutOD/esp`.

## Import Statement Updates

All Go files importing `github.com/pinpt/esp/...` are updated to `github.com/AbsolutOD/esp/...`. Affected files (17 import statements across 12 files):

- `main.go` (1)
- `cmd/root.go` (2)
- `cmd/copy.go` (1)
- `cmd/delete.go` (1)
- `cmd/get.go` (1)
- `cmd/list.go` (2)
- `cmd/move.go` (1)
- `cmd/put.go` (1)
- `internal/app/config.go` (1)
- `internal/client/client.go` (3)
- `internal/ssm/ssm.go` (2)
- `internal/ssm/utils.go` (1)

Capitalization is preserved: `AbsolutOD` (matches the GitHub organization's actual casing).

## LICENSE Update

The current `LICENSE` is unmodified Apache 2.0 boilerplate with no actual copyright holder filled in.

Change: prepend a single line at the very top of the file, before the existing Apache header:

```
Copyright 2021-2026 Matt O'Donnell
```

A blank line separates the new copyright line from the unchanged Apache License text. The license terms themselves are not altered.

## README Rewrite

The current `README.md` is ~20 lines, references a Pinpoint logo, links to a sample config that no longer exists, and includes a stale Pinpoint copyright line. Full rewrite covering:

- One-paragraph project description: what `esp` is, what problem it solves.
- Install / build instructions.
- Required environment variables (`AWS_DEFAULT_REGION`, `AWS_PROFILE`) and their behavior.
- The `.espFile` (per-project config) format with a small example.
- Command summary: every subcommand with one line of description.
- A short "common workflows" section showing typical usage.
- Footer line: `Released under the Apache License 2.0 — see LICENSE.`

No logo `<img>` tag at the top. The README should look like a fresh project's README, not a refresh of an old one.

## Logo Deletion

- Delete `.github/logo.svg`. No replacement.

## Verification

- `go build ./...` — exit 0.
- `go vet ./...` — exit 0.
- `go test ./...` — all packages pass; the import-path change does not affect test logic.
- `grep -rn 'pinpt\|Pinpoint\|PINPT' --include='*.go' --include='*.md' --include='go.mod' .` — returns zero matches outside `docs/superpowers/specs/2026-05-03-*` and `docs/superpowers/plans/2026-05-03-*` (preserved historical artifacts).
- README renders correctly on GitHub with no broken image link.
- `LICENSE` first line reads `Copyright 2021-2026 Matt O'Donnell`.

## Risks

- **Module-path change is a breaking change for any external importer.** The repo's public surface is CLI-only and the codebase uses the `internal/` directory throughout (which prevents external imports of internal packages by Go's rules). Practical blast radius for external consumers is zero.
- **Original copyright provenance.** The 2019–2020 period belonged to Pinpoint. This rewrite explicitly does not claim those years. Apache 2.0 does not require preserving prior copyright lines in the LICENSE file when the original was unfilled boilerplate.
