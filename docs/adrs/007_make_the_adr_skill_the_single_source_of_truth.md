---
id: TNaqIwLDg
title: Make the ADR skill the single source of truth
author: Gwendal Leclerc
status: accepted
creation_date: 2026-07-17T00:10:27.220008595+02:00
last_update_date: 2026-07-17T00:10:27.220008595+02:00
tags:
  - claude
  - process
---

# Make the ADR skill the single source of truth

Date: Fri, 17 Jul 2026 00:10:27 CEST

## Context and Problem Statement

The project ships both a Claude Code skill and an `/adr` command, both named `adr` and
both describing the ADR workflow. The duplicated instructions were ambiguous — it was
unclear which was authoritative — and the two copies risked drifting out of sync.

## Considered Options

- Keep both with the workflow duplicated in each — status quo; ambiguous and drift-prone.
- Drop the command, skill only — one entity, but relies on skills being slash-invokable
  and loses an explicit, discoverable command.
- Skill as the single source of truth + command as a thin delegate (superpowers-style) —
  one "brain", one explicit trigger, no duplicated methodology.

## Decision Outcome

Adopt the skill-centric model (inspired by the superpowers skills): the `adr` skill
holds all methodology and composes with design skills (it references `brainstorming` /
`writing-plans` by name to graft in as the downstream "record the decision" step), while
the `/adr` command becomes a thin entry point that only delegates to it. The dependency
direction is one-way: command → skill → other skills, never the reverse.

### Consequences

- One authoritative workflow; no duplicated instructions to keep in sync.
- The skill grafts cleanly onto design/planning flows as their downstream capture step.
- `/adr` still works as an explicit trigger. A same-name command and skill still coexist
  in listings, but their roles are now disjoint (trigger vs. methodology).
