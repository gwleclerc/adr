---
id: S9MFFQYvR
title: Add repo-maintenance and lifecycle commands
author: Gwendal Leclerc
status: accepted
creation_date: 2026-07-17T01:09:44.203573363+02:00
last_update_date: 2026-07-17T01:09:44.203573363+02:00
tags:
  - cli
  - tooling
---

# Add repo-maintenance and lifecycle commands

Date: Fri, 17 Jul 2026 01:09:44 CEST

## Context and Problem Statement

The CLI could create and query records but offered nothing to *maintain* a growing set:
no reader-facing index, no way to catch inconsistencies (dangling superseder references,
duplicate numbers, invalid statuses), and routine status transitions required verbose
`update` / `add -r` invocations.

## Considered Options

- Leave maintenance to external scripts and manual review — no new CLI surface, but repo
  health goes unchecked and there is no index for readers.
- Add focused, deterministic maintenance commands to the CLI — a little more surface, but
  structural operations belong in the tool (the skill stays for authoring judgement).

## Decision Outcome

Add `adr toc` (a markdown index), `adr lint` (consistency checks that exit non-zero on
issues, for CI), and the lifecycle shortcuts `adr deprecate <id>` and
`adr supersede <old> <new>`. Index rendering and linting are pure functions (unit-tested),
and linting reuses the already-indexed records. These are deterministic, structural
operations, so they live in the CLI rather than the skill.

### Consequences

- Readers get an index; CI can gate on `adr lint`; common transitions are a single command.
- More CLI surface to maintain, but each command is thin and covered by unit + venom tests.
