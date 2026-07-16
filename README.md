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
   --help, -h                      show help
```

It will create a new numbered ADR in your ADR folder `001_decisive_decision_of_architecture.md`.

Then you will have to open the file in your preferred editor and start editing the ADR.

The template contains placeholders to indicate the purpose of each section.

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
make build        # build the binary into ./build
make test         # unit tests with the race detector
make integration  # end-to-end tests (installs and runs venom)
make lint         # golangci-lint
make release VERSION=v1.2.3 RELEASE=1   # cross-compile archives into ./dist
```

Releases are produced automatically by GitHub Actions: pushing a `v*` tag builds
binaries for Linux, macOS and Windows across amd64/arm64/386/arm and publishes them
to a GitHub release (see `.github/workflows/release.yml`).
