---
name: adr
description: >-
  Create and manage Architecture Decision Records (ADRs) with the `adr` CLI.
  Use whenever the user wants to record, document, list, update, or supersede
  an architecture/design decision, or mentions "ADR" / "decision record".
---

# Managing ADRs with the `adr` CLI

This project ships `adr`, a CLI to manage [Architecture Decision Records](https://adr.github.io/).
Prefer these commands over hand-writing decision files, so numbering, metadata and
the markdown template stay consistent.

## Prerequisites

- The `adr` binary must be on `PATH` (`go install github.com/gwleclerc/adr@latest`,
  the `install.sh` script, or `make build`).
- A `.adrrc.yml` must exist. If commands fail to find records, initialize once:
  ```bash
  adr init docs/adrs
  ```
  This stores the records directory in `.adrrc.yml` at the repo root.

## Commands

### Create a record
```bash
adr new <title...> [-a author] [-s status] [-t tag1,tag2] [-r <id>...]
```
- `-s/--status` defaults to `accepted`. Allowed: `unknown`, `proposed`, `accepted`,
  `deprecated`, `superseded`, `observed`.
- `-r/--supersedes <id>` marks the referenced record as `superseded` and links back to it.
- The command creates a numbered file (e.g. `001_my_title.md`) and prints the new record ID.

### List / filter records
```bash
adr list [-a authors] [-s status] [-t tags]
```
Filters are comma-separated and combine (author AND status AND tag). Output is a table.

### Update a record's metadata
```bash
adr update <id> [-a author] [-s status] [-t tags] [-r superseders]
```
- Only the flags you pass are changed; the body is preserved.
- Pass an **empty** value to clear a field: `--tags=` removes all tags, `--superseders=` clears them.

### Add tags / superseders (append, non-destructive)
```bash
adr add <id> [-t tags] [-r superseders]
```
Unlike `update`, `add` only appends and never removes existing metadata.

## Guidance

- After `adr new`, open the generated file and fill in the template sections
  (**Context**, **Decision**, **Implications**) — the CLI only scaffolds metadata + headings.
- To supersede a decision, use `adr new "<new title>" -r <old-id>` rather than editing the
  old file by hand, so the back-link and status update happen automatically.
- Record IDs are short opaque strings (printed on creation and shown by `adr list`), not the
  numeric filename prefix. Use the ID for `update`/`add`/`-r`.
- Error/diagnostic output goes to **stderr**; success messages and tables go to **stdout**.
