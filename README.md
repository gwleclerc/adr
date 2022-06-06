# ADR Golang

A Golang Architectural decision records (Adrs) CLI

A simple command line written in Golang to manage [Architecture Decision Records](https://adr.github.io/) (ADRs).

# Getting Started

## Installation

You need to install it manually with the following command:

```bash
go install github.com/marouni/adr@main
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
adr new decisive decision of articheture
```

You can also add flags to set record's metadata:

```bash
Create a new architecture decision record.
It will be created in the directory defined in the nearest .adrrc.yml configuration file.

Usage:
  adr new [flags] <record title...>

Flags:
  -a, --author string   author of the record
  -h, --help            help for new
  -s, --status status   status of the record, allowed: "unknown", "proposed", "accepted", "deprecated" or "superseded" (default accepted)
  -t, --tags strings    tags of the record
```

It will create a new numbered ADR in your ADR folder `001_decisive_decision_of_articheture.md`.

Then you will be have to open the file in your preferred editor and starting editing the ADR.

The template contains placeholders to indicate the purpose of each section.

## Listing records

You can list all records using the following command:

```bash
adr list
```

By using flags, you can filter records based on their metadata:

```bash
List ADR files present in directory stored in .adrrc.yml configuration file.

Usage:
  adr list [flags]

Flags:
  -a, --authors strings   filter records by authors
  -h, --help              help for list
  -s, --status strings    filter records by status
  -t, --tags strings      filter records by tags
```

This will display the records as a table.

```bash
+-----------+----------------------------------+----------+-----------------+-------------+--------------+
|    ID     |              TITLE               |  STATUS  |     AUTHOR      |    DATE     |     TAGS     |
+-----------+----------------------------------+----------+-----------------+-------------+--------------+
| zl3cUj97R | decisive_decision_of_articheture | accepted | Gwendal Leclerc | 2 hours ago | architecture |
+-----------+----------------------------------+----------+-----------------+-------------+--------------+
```
