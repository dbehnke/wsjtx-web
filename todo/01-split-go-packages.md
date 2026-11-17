# 01 — Split Go packages

Status: `done`  
Priority: P0 (foundation)

## Goal

Break `main.go` into focused packages so UDP, WebSocket hub, HTTP, and audio can evolve independently and be tested.

## Why

`main.go` mixes transport, hub, static file serving, and audio. That blocks optional audio, hub hardening, and unit tests.

## Scope

**In**

- Introduce `internal/hub`, `internal/udp` (or keep UDP loop near main), `internal/audio`, `internal/server` (HTTP + WS routes)
- `main` only parses flags and wires dependencies
- Preserve current behavior and CLI flags

**Out**

- Frontend changes
- Protocol redesign (see later todos)

## Steps

1. Map current types/funcs in `main.go` (Hub, AudioManager, handlers).
2. Move Hub + broadcast types → `internal/hub`.
3. Move HTTP/WS registration → `internal/server` (inject hub, UDP writer, audio).
4. Move malgo code → `internal/audio`.
5. Leave `pkg/wsjtx` as public protocol package.
6. `go test ./...` and manual smoke: UDP decode → WS, reply/halt.

## Acceptance

- [x] `main.go` is thin wiring (< ~80 lines ideal)
- [x] `go build` and `go test ./...` pass
- [x] Existing flags still work
- [x] No behavior change for decode broadcast / reply / halt

Note: `main.go` landed ~108 lines (flags + UDP loop + optional audio wiring); hub/audio/HTTP live under `internal/`.

## Key files

- `main.go`
- `pkg/wsjtx/`
- new `internal/**`

## Risks

- Circular imports — depend inward toward `pkg/wsjtx` and small interfaces for “send UDP” / “list devices”.
