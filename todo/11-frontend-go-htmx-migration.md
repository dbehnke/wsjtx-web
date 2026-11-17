# 11 — Frontend migration (Go + HTMX / vanilla)

Status: `todo`  
Priority: P3 (direction, not urgent)  
Depends on: foundation todos `01`–`07` strongly recommended first

## Goal

Evaluate and (if approved) replace the Vue/npm SPA with a Go-served UI: templates + HTMX or light vanilla JS, plus a small static `waterfall.js`.

## Why

Maintainers prefer less Node. This UI is mostly live table + controls + one canvas — a strong fit for server-rendered HTML + WebSocket/HTMX.

## Decision gate (do this before coding)

Write a short ADR in `docs/adr/001-frontend-stack.md`:

- Stay on Vue, or migrate?
- HTMX vs pure WS DOM updates?
- CSS approach without a Node pipeline (checked-in CSS, etc.)

Only proceed to implementation if ADR says migrate.

## Migration sketch (if yes)

1. Add `internal/web` templates for shell, status, band table.
2. WS still pushes decode/status JSON; small JS patches DOM (or HTMX WebSocket extension).
3. Port waterfall to `static/waterfall.js`.
4. Remove `wsjtx-web-ui` when parity reached; simplify Makefile/CI.
5. Keep protocol types from `05`.

## Acceptance

- [ ] ADR merged with explicit choice
- [ ] If migrate: feature parity for monitor + reply + halt; single `go build` artifact; no npm required for prod
- [ ] If stay: ADR closes this todo as **done** with rationale

## Key files

- New `internal/web`, `static/`
- Eventually delete `wsjtx-web-ui/`

## Risks

- Large rewrite — must not block P0/P1 todos
- Waterfall/WebAudio still needs browser JS either way
