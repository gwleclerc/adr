---
id: zl3cUj97R
title: change_add_command_to_new
author: Gwendal Leclerc
status: accepted
creation_date: 2022-06-06T00:35:32.033552+02:00
last_update_date: 2022-06-07T00:45:04.70907+02:00
tags:
  - refactor
---

# Change Add Command To New

Date: Mon, 06 Jun 2022 00:35:32 CEST

## Context

The add command is ambiguous. It can represent the addition of a new record or the addition of information about an existing record.

## Decision

The "add" command will be replaced by a "new" command.

## Concequences

A new command "add" can be added in the future.
