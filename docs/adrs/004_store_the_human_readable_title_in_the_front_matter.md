---
id: ZcEubULDg
title: Store the human-readable title in the front matter
author: Gwendal Leclerc
status: accepted
creation_date: 2026-07-16T18:01:52.008176813+02:00
last_update_date: 2026-07-16T18:01:52.008176813+02:00
tags:
  - refactoring
---

# Store the human-readable title in the front matter

Date: Thu, 16 Jul 2026 18:01:52 CEST

## Context and Problem Statement

`new` stored the filename slug in the `title:` metadata, so `adr list` showed
`modernize_the_cli_and_build_tooling` instead of a readable title, and the `#` heading was
force-title-cased, mangling acronyms (e.g. "Cli").

## Considered Options

- Keep the slug as the title — simple, but unreadable listings and mangled acronyms.
- Store the verbatim human title in the metadata and heading, keep the slug only for the filename.

## Decision Outcome

Store the human title verbatim in `title:` and the `#` heading; derive the filename slug
separately. Titles are used as typed, so casing and acronyms are preserved.

### Consequences

- `adr list` and `adr show` are readable, and acronyms survive.
- Records created before this change keep their slug-style stored title (not migrated).
