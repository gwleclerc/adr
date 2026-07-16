---
id: 8YLkGyLDR
title: modernize_the_cli_and_build_tooling
author: Gwendal Leclerc
status: accepted
creation_date: 2026-07-16T04:23:44.829012176+02:00
last_update_date: 2026-07-16T04:23:44.829012176+02:00
tags:
  - refactoring
  - tooling
---

# Modernize The Cli And Build Tooling

Date: Thu, 16 Jul 2026 04:23:44 CEST

## Context and Problem Statement

The CLI was built on `spf13/cobra` + `pflag` with dependencies frozen in 2022
(`io/ioutil`, `golang.org/x/exp`), and carried a few latent bugs (a semaphore leak that
could deadlock indexing, a broken format string, diagnostics printed to stdout). It had
no way to ship binaries for anything but the build host. We wanted a lighter, modern and
easier-to-maintain foundation, plus a real cross-platform release path — while keeping
the tool deliberately small.

## Considered Options

- **Only fix the bugs, keep the stack as-is** — minimal churn; but stays on a heavy,
  dated command framework and offers no release story.
- **Refresh the foundation** — migrate the command layer, adopt the standard library,
  and add a release pipeline; more up-front work, modern and lighter afterwards.

## Decision Outcome

Refresh the project:

- Replace cobra/pflag with **urfave/cli v3** — same multi-command UX with far less
  ceremony for a small tool.
- Adopt the standard library: `io/ioutil` → `os`, `golang.org/x/exp` → `slices`/`cmp`;
  bump to Go 1.23; drop `cobra`, `pflag` and `x/exp`.
- Fix the latent bugs and route diagnostics to **stderr** (success output stays on stdout).
- Wire build metadata into `adr --version`.
- Cross-compile release archives from the existing **Makefile** (`make release`) driven by
  GitHub Actions on `v*` tags, plus an `install.sh` — no extra release tooling, to stay light.

Driver: keep the tool light and easy to use and maintain, on a current dependency set.

### Consequences

- Lighter, maintained dependency tree; simpler per-command code; correct stderr/stdout
  split; real binaries for Linux/macOS/Windows across amd64/arm64/386/arm; self-install.
- One-time rewrite of every command and its tests; slightly different `--help` output;
  contributors now need Go 1.23+.
