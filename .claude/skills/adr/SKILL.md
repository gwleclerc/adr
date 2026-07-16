---
name: adr
description: >-
  Capture and retro-document architecture decisions as structured ADRs using the
  `adr` CLI. Trigger proactively when a non-trivial technical decision is made in a
  session (choosing a library, pattern, protocol, data model, trade-off, or reversing
  an earlier decision); when the user asks to record/document a decision or mentions
  "ADR"; and when the user wants to retro-document the decisions embodied in an
  existing/legacy repository (see Mode B, also reachable via the /adr command). The CLI
  stays a light scaffolder; this skill supplies the MADR-lite structure and workflow.
---

# Authoring ADRs with the `adr` CLI

`adr` is a lightweight scaffolder: it handles numbering, metadata and the file, and
nothing else. The *structure and rigor* live here, in the workflow. The goal is
MADR-quality records (context, options, decision, consequences) without full-MADR
ceremony.

## Two modes

- **Mode A — capture a decision** taken *now*: as it happens in the session
  (proactively), on request, or **from the changes just made** (reconstructed from the
  working `git diff` / recent commits). Status `accepted`.
- **Mode B — retro-document an existing repository**: sweep a codebase you're taking
  over and emit several *condensed* `observed` ADRs.

The `/adr` command routes across these based on intent: "record the modifications we
just made" → Mode A from the diff; "document this repo" → Mode B.

## Prerequisites

- `adr` on `PATH` (`go install github.com/gwleclerc/adr@latest`, `install.sh`, or `make build`).
- A `.adrrc.yml`. If records aren't found, initialize once: `adr init docs/adrs`.

## Mode A — capture a decision taken now

Trigger proactively when a meaningful decision is reached mid-session, on request, or
when asked to document the changes just made. Do **not** create an ADR for trivial/local
choices (variable names, one-off refactors).

1. **Get the contract from the tool.** Run `adr template show madr` to see the exact
   sections and their guidance — this is the single source of truth for what a body must
   contain (don't hardcode the structure). Use `madr` unless the user prefers another
   template (`adr template list`).
2. **Draft the body — do not create the file yet.** Gather the reasoning from the session
   context and, when documenting recent work, from the current `git diff` / recent
   commits (reconstruct *what* changed and *why*). Write one markdown body that fills every
   section of the contract; fill what you can and ask the user only for genuine gaps.
   Name the implicit "status quo / do nothing" option under *Considered Options* when it
   applies.
3. **Validate with the user.** Show the draft and confirm it's correct and complete.
4. **Create the record in one shot, directly as `accepted`,** passing the validated body:
   ```bash
   adr new "<title>" --template madr -s accepted --body-file <draft> [-a author] [-t tags] [-r <superseded-id>...]
   ```
   The CLI validates the body against the template and rejects it if a section is missing
   or empty, so there is no scaffold-then-overwrite step. (Omit `--body-file` to instead
   scaffold the empty template for a human to fill.)

## Mode B — retro-document an existing repository (`/adr`)

Goal: capture the significant decisions **already embodied** in a codebase you did not
author, as several **condensed** `observed` ADRs — asking the user *why* when the
rationale is not evident. This is for making sense of and re-appropriating legacy code.

1. **Setup.** Ensure `.adrrc.yml` exists (offer `adr init docs/adrs` if not). Check
   existing ADRs with `adr list` so you don't duplicate.
2. **Survey the repo for decision signals** (read, don't guess). Look at:
   - build/deps manifests → language, framework and library choices;
   - directory layout & module boundaries → architectural style (layers, hexagonal, monorepo…);
   - data layer → database, schema/migrations, storage;
   - interfaces → API style (REST/gRPC/GraphQL), CLI framework;
   - infra/CI → Dockerfile, workflows, deployment target;
   - cross-cutting → auth, logging, error handling, concurrency model;
   - README, docs, notable comments and commit history → any *stated* rationale.
3. **Cluster into a SMALL set of decisions.** Group related choices; aim for roughly the
   5–12 decisions that actually shaped the system. Condensed means *several focused
   records*, not one per file and not one giant document. Skip the trivial.
4. **For each candidate, prepare a condensed draft:** the decision embodied, a
   reconstructed context, the alternatives that plausibly existed, the rationale, and
   the consequences visible in today's code. **Cite the evidence** (file/signal) each
   decision is inferred from.
5. **Handle the "why" honestly (key requirement).** Never invent rationale. For every
   point where the *why* is ambiguous from code/docs/history, mark it as a question.
   **Batch all ambiguous questions and ask the user once**, then fold the answers in.
   In the body, label reasoning as *inferred* vs *confirmed by <source/user>*.
6. **Confirm the shortlist with the user** (which decisions, what granularity) before
   writing files.
7. **Create one condensed `observed` record per confirmed decision,** passing the drafted
   body (get the section contract from `adr template show madr` first):
   ```bash
   adr new "<decision title>" --template madr -s observed --body-file <draft> -t <area> -a "<original team|unknown>"
   ```
   Keep each record terse: a few lines per section, one decision each.

## Status rules

- In Mode A, **never write an ADR as `proposed`.** Validation *is* acceptance: once the
  user confirms the draft, create it as `accepted`. An ADR should already be `accepted`
  by the time it lands in a PR — a reviewed PR counts as acceptance, so avoid the
  `proposed → accepted` git ping-pong after approval.
- In Mode B, records are `observed` (pre-existing decisions reconstructed after the fact).
- Later transitions use `adr update <id> -s <status>`: `deprecated` when it no longer
  applies, `superseded` when replaced (prefer the automatic path below).
- Run `adr new --help` for the meaning of every status.

## Superseding a decision

Create the new record referencing the old ID — this sets the old record to `superseded`
and back-links it automatically. Do not hand-edit the old file's status.
```bash
adr new "<new title>" --template madr -s accepted -r <old-id>
```

## Templates

The body structure is owned by templates, not by this skill — always learn a template's
sections from the tool rather than hardcoding them:

- `adr template list` — available templates (`bare`, `madr`, plus any custom ones).
- `adr template show <name>` — the sections + guidance a body must fill (the contract).

Projects can add their own templates: set `templates_dir: <path>` in `.adrrc.yml`; every
`*.tpl` file there becomes a template named after the file (overriding a built-in of the
same name). When a project ships custom templates, prefer them over `madr`.

## CLI reference

| Command | Purpose |
|---|---|
| `adr new <title> [--template <name>] [--body-file <f\|->] [-s] [-a] [-t] [-r]` | create a record (prints its ID); `--body-file` supplies a validated body, else the template is scaffolded |
| `adr template list` / `adr template show <name>` | discover templates / print a template's section contract |
| `adr list [-a] [-s] [-t] [--json]` | list/filter records (filters are AND across types); **use `--json`** to parse records reliably |
| `adr show <id> [--json]` | print one record (raw file, or metadata as JSON with `--json`) |
| `adr update <id> [-a] [-s] [-t] [-r]` | change metadata; only passed flags change; `--tags=` clears |
| `adr add <id> [-t] [-r]` | append tags/superseders (adding superseders marks the record `superseded`) |

- Statuses: `unknown`, `proposed`, `accepted`, `deprecated`, `superseded`, `observed`
  (`observed` is a CLI extension for retrospective records; not part of MADR/Nygard).
- Record IDs are short opaque strings (shown on creation and by `adr list`), not the
  numeric filename prefix. Use the ID for `update`/`add`/`-r`.
- Diagnostics go to **stderr**; success output and tables go to **stdout**.
