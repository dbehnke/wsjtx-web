# 05 — Shared / typed WS protocol

Status: `todo`  
Priority: P1

## Goal

Replace magic string types (`"1"`, `"2"`) and `any` payloads with a documented, typed contract used by Go and the UI.

## Why

Fragile and hostile to agents. Easy to break reply fields (`DeltaFrequency` vs `DeltaFreq`).

## Options (pick one in implementation)

1. **Lightweight:** Go constants + JSON `type` string enums (`status`, `decode`, …) and matching TS union types in `wsjtx-web-ui/src/types/ws.ts`.
2. **Heavier:** JSON Schema or `quicktype` codegen from a single schema file.

Prefer (1) unless codegen is already justified.

## Steps

1. Inventory all WS server→client and client→server messages.
2. Change server wrapper to use stable string names (keep numeric aliases briefly if needed).
3. Add TS types; tighten Pinia store.
4. Document in `AGENTS.md` or `docs/protocol.md`.

## Acceptance

- [ ] No `msg.type === "2"` in UI (use named types)
- [ ] Go and TS agree on field names for `reply` / `halt`
- [ ] Short protocol doc exists

## Key files

- `main.go` / `internal/server`
- `wsjtx-web-ui/src/stores/wsjtx.ts`
- `pkg/wsjtx/protocol.go`

## Risks

- Breaking existing clients — this repo is the only client; still bump carefully.
