---
id: hcPuxUYvg
title: Machine-readable output and a show command
author: Gwendal Leclerc
status: accepted
creation_date: 2026-07-16T18:01:52.01873682+02:00
last_update_date: 2026-07-16T18:01:52.01873682+02:00
tags:
  - feature
---

# Machine-readable output and a show command

Date: Thu, 16 Jul 2026 18:01:52 CEST

## Context and Problem Statement

Now that ADRs are often created and read by agents, scraping the ASCII table was fragile,
and there was no way to print or open a single record by ID.

## Considered Options

- Keep table-only output — human-friendly but hard to parse reliably.
- Add JSON output and a `show`/`edit` command — machine-readable and closes the CRUD gap.

## Decision Outcome

Add `adr list --json` and `adr show <id> [--json]` for machine-readable access, plus
`adr edit <id>` and `adr new --edit` to open records in `$EDITOR`. The `Set` type marshals
to a sorted JSON array so tags/superseders are clean lists.

### Consequences

- Agents and scripts can consume records as JSON instead of parsing a table.
- Slightly more CLI surface to maintain.
