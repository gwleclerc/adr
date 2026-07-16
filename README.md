# ADR Golang

A Golang Architectural decision records (Adrs) CLI

A simple command line written in Golang to manage [Architecture Decision Records](https://adr.github.io/) (ADRs).

# Getting Started

## Installation

### Install script (recommended)

Download the latest prebuilt binary for your platform (Linux/macOS, amd64/arm64/…):

```bash
curl -fsSL https://raw.githubusercontent.com/gwleclerc/adr/main/install.sh | sh
```

You can pin a version or change the install directory:

```bash
ADR_VERSION=v1.0.0 ADR_INSTALL_DIR="$HOME/.local/bin" \
  sh -c "$(curl -fsSL https://raw.githubusercontent.com/gwleclerc/adr/main/install.sh)"
```

Prebuilt archives for every platform (including Windows `.zip`) are attached to each
[GitHub release](https://github.com/gwleclerc/adr/releases).

### From source

```bash
go install github.com/gwleclerc/adr@latest
```

Check the installed version (and build metadata) with:

```bash
adr --version
```

## Initializing

Before creating a new record, you must initialize the configuration with a folder that will contain your ADRs
with the following command:

```bash
adr init docs/adrs
```

It will create a `.adrrc.yml` configuration file with the directory path inside.

## Creating a new record

You can create a new record with the following command:

```bash
adr new decisive decision of architecture
```

You can also add flags to set record's metadata:

```
NAME:
   adr new - Create a new ADR

USAGE:
   adr new [options] <record title...>

OPTIONS:
   --author string, -a string     author of the record
   --status string, -s string     status of the record, allowed: "unknown", "proposed", "accepted", "deprecated", "superseded" or "observed" (default: "accepted")
   --tags string, -t string        tags of the record
   --supersedes string, -r string  record ids superseded by this one
   --template string                body template name (see `adr template list`) (default: "bare")
   --body-file string               read the record body from a file (or - for stdin) instead of the template
   --help, -h                      show help
```

By default `new` scaffolds a minimal (Nygard-style) body. Pass `--template madr` for a
richer [MADR](https://github.com/adr/madr)-lite layout with *Context and Problem
Statement*, *Considered Options*, *Decision Outcome* and *Consequences* sections:

```bash
adr new use urfave/cli over cobra --template madr
```

It will create a new numbered ADR in your ADR folder `001_decisive_decision_of_architecture.md`
with placeholder prose for each section, ready to edit in your preferred editor.

## Templates

Templates define the body structure of a record. Inspect them with:

```bash
adr template list          # available templates (bare, madr, plus your own)
adr template show madr     # print a template's sections and guidance
```

Instead of editing the scaffolded file, you can supply a ready-made body — the CLI wraps
it with the metadata and **validates** that it matches the template's sections (missing
or empty section → error):

```bash
adr new "use urfave/cli over cobra" --template madr --body-file draft.md
cat draft.md | adr new "use urfave/cli over cobra" --template madr --body-file -
```

### Custom templates

Declare a templates directory in `.adrrc.yml`; every `*.tpl` file there becomes a template
named after the file (a custom name equal to a built-in overrides it):

```yaml
directory: docs/adrs
templates_dir: .adr/templates   # relative to the .adrrc.yml (or absolute)
```

A template file is just the body skeleton — the section headings (and optional `>`
guidance). For example `.adr/templates/lightweight.tpl`:

```markdown
## Context
> Why is this decision needed?

## Decision
> What did we decide, and why?
```

Then: `adr new "my decision" --template lightweight`.

## Record statuses

| Status | Meaning |
|---|---|
| `unknown` | status is not determined |
| `proposed` | proposed but not accepted yet by stakeholders |
| `accepted` | accepted by stakeholders |
| `deprecated` | no longer applies |
| `superseded` | replaced by a newer record (set automatically via `new -r`) |
| `observed` | documents a **pre-existing** decision reconstructed after the fact — e.g. while making sense of legacy code you did not write |

`observed` is handy for retrospective ADRs: when re-appropriating an inherited codebase,
record *how things already are and why* (`adr new "..." -s observed`) rather than
pretending the decision is being taken now.

## Updating a record

You can change the metadata of an existing record with `update`. Flags that are not
provided are left untouched; passing an empty value (e.g. `--tags=`) clears the field.

```bash
adr update <record ID> -s deprecated -t design,api
```

## Adding metadata to a record

`add` appends tags or superseders to a record without touching its other metadata:

```bash
adr add <record ID> -t security -r <other record ID>
```

## Listing records

You can list all records using the following command:

```bash
adr list
```

By using flags, you can filter records based on their metadata:

```
NAME:
   adr list - List ADR files

USAGE:
   adr list [options]

OPTIONS:
   --authors string, -a string  filter records by authors
   --status string, -s string   filter records by status
   --tags string, -t string     filter records by tags
   --help, -h                   show help
```

This will display the records as a table.

```bash
+-----------+-----------------------------------+----------+-----------------+-------------+--------------+
|    ID     |              TITLE                |  STATUS  |     AUTHOR      |    DATE     |     TAGS     |
+-----------+-----------------------------------+----------+-----------------+-------------+--------------+
| zl3cUj97R | decisive_decision_of_architecture | accepted | Gwendal Leclerc | 2 hours ago | architecture |
+-----------+-----------------------------------+----------+-----------------+-------------+--------------+
```

# Development

Common tasks are wrapped in the `Makefile`:

```bash
make build          # build the binary into ./build
make test           # unit tests with the race detector
make integration    # end-to-end tests (installs and runs venom)
make lint           # golangci-lint
make release VERSION=v1.2.3 RELEASE=1   # cross-compile archives into ./dist
make install-claude # symlink the Claude Code skill + /adr command into ~/.claude
```

Releases are produced automatically by GitHub Actions: pushing a `v*` tag builds
binaries for Linux, macOS and Windows across amd64/arm64/386/arm and publishes them
to a GitHub release (see `.github/workflows/release.yml`).

## Claude Code integration

This repo ships a Claude Code skill (`.claude/skills/adr`) and an `/adr` command
(`.claude/commands/adr.md`). Run `make install-claude` to symlink them into `~/.claude`
so they're available in every repo (the source stays versioned here).

- The **skill** teaches Claude to author structured ADRs with this CLI, and triggers
  proactively when a decision is made in a session.
- The **`/adr` command** records decisions on demand and routes by intent: capture the
  decision(s) behind the changes just made (from the diff, as `accepted`), document a
  specific decision, or sweep a codebase you're taking over and retro-document it as
  several condensed `observed` ADRs. It never invents rationale — it asks you about the
  *why* whenever it's ambiguous.
