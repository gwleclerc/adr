---
id: TTUEkAjng
title: create_records_service
author: Gwendal Leclerc
status: accepted
creation_date: 2022-06-13T03:04:30.1997921+02:00
last_update_date: 2022-06-13T03:04:30.1997921+02:00
---

# Create Records Service

Date: Mon, 13 Jun 2022 03:04:30 CEST

## Context

In order to add the 'supersedes' option when creating a new record, 
the `new` command must be able to modify one or more records, which currently leads to a lot of code duplication.

## Decision

All the code concerning the records will be extracted into a dedicated service.

## Implications

We will be able to create atomic functions to `create` and `modify` records easily.
