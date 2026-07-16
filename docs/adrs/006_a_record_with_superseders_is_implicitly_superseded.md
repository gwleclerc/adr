---
id: wcEuxULvg
title: A record with superseders is implicitly superseded
author: Gwendal Leclerc
status: accepted
creation_date: 2026-07-16T18:01:52.029721173+02:00
last_update_date: 2026-07-16T18:01:52.029721173+02:00
tags:
  - refactoring
---

# A record with superseders is implicitly superseded

Date: Thu, 16 Jul 2026 18:01:52 CEST

## Context and Problem Statement

`new -r` marked the target record `superseded` and back-linked it, but `add`/`update` only
appended superseder IDs without updating status, and unknown IDs were silently ignored —
leaving inconsistent records (superseders present, yet status not `superseded`).

## Considered Options

- Leave the asymmetry — simplest, but keeps producing inconsistent records.
- Treat superseders as authoritative: any record that gains superseders becomes `superseded`, and unknown IDs are reported.

## Decision Outcome

A record that gains superseders is set to `superseded` in `add`/`update` (unless an explicit
`-s` is given in the same call), matching what `new -r` already did to its targets. Unknown
referenced IDs now emit a warning instead of being silently ignored.

### Consequences

- Consistent supersede semantics across `new`, `add` and `update`; no silent no-ops.
- Adding a superseder now also changes status, which callers must be aware of.
