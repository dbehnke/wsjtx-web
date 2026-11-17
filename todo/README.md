# Todo queue

Plans for features and refactors. Agents and humans work these in order unless a plan says otherwise.

## How to use

1. Open `00-overview.md` — pick the first item with status **todo**.
2. Open that plan file — follow **Goal**, **Steps**, **Acceptance**.
3. Implement on a feature branch; open a PR.
4. When merged (or landed), set status to **done**, check the box in `00-overview.md`, and list any new follow-ups.

## Status values

| Status | Meaning |
|--------|---------|
| `todo` | Not started |
| `in_progress` | Actively being worked |
| `blocked` | Waiting on decision/dependency — say why |
| `done` | Complete; keep the file for history |

## Plan file template

New plans should include:

- Title, status, priority
- Goal / why
- Scope (in / out)
- Steps
- Acceptance criteria
- Key files
- Risks / notes

## Naming

`NN-short-slug.md` — numbers define default order (`00` = roadmap index).
