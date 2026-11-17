# 07 — Tests and CI

Status: `todo`  
Priority: P1

## Goal

Catch protocol and build regressions before merge.

## Why

No `.github/workflows` today. Protocol code is the highest-value test surface. Embed/`dist` and npm drift already bit local setup.

## Steps

1. Table-driven tests for `pkg/wsjtx` encode/decode (heartbeat, status, decode, reply, halt).
2. Optional hub unit test with `httptest` + gorilla websocket.
3. GitHub Actions on PR:
   - `go test ./...`
   - `go build` (with placeholder or built dist)
   - `cd wsjtx-web-ui && npm ci && npm run build` (and unit tests)
4. Document in `AGENTS.md`.

## Acceptance

- [ ] `go test ./...` covers core message types
- [ ] CI green on PR to `main`
- [ ] Scaffold HelloWorld test gone or replaced (see `03`)

## Key files

- `pkg/wsjtx/*_test.go`
- `.github/workflows/ci.yml`

## Risks

- CGO/audio in CI — build/test with audio disabled (`02`) so CI stays simple.
