---
description: Record architecture decisions with the adr CLI — capture the decision(s) behind recent changes, document a specific decision, or retro-document the whole repo
---

Record architecture decision(s) using the **`adr` skill** (load it for the methodology,
CLI usage and status rules). First pick the mode that fits my arguments and the session
context — if it is genuinely unclear, ask me before proceeding:

- **Capture recent changes (default when there is recent work to document).** If I refer
  to "the changes/modifs just made", or there is a relevant working diff or recent
  commits from this session: inspect `git diff` / `git log` and the session context, and
  draft an ADR for the decision(s) those changes embody. Status `accepted` (these are
  decisions we are making now). This is Mode A of the skill.
- **A specific decision.** If my arguments describe a decision directly, draft that one
  (Mode A), status `accepted`.
- **Retro-document the repository.** If I ask to document/sweep the existing codebase
  (e.g. a repo I'm taking over), run Mode B: survey the code and emit several **condensed**
  `observed` ADRs for the decisions it already embodies.

In every mode:
- **Draft first**, don't create files until I've validated the shortlist.
- **Never invent the rationale.** Where the *why* is ambiguous, mark it and ask me the
  ambiguous questions **batched together**; label reasoning as *inferred* vs *confirmed*.
- Create records with `adr new "<title>" --template madr -s <status>` — `accepted` for
  decisions taken now, `observed` for pre-existing ones reconstructed after the fact —
  then fill the body with the validated content.

$ARGUMENTS
