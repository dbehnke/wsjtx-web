# 06 — WebSocket hub hardening

Status: `todo`  
Priority: P1  
Depends on: `01-split-go-packages.md` (ideal)

## Goal

Make the hub safe under slow clients, reconnect storms, and basic abuse.

## Why

Current hub can block on `WriteMessage`; UI reconnects every 1s forever; `CheckOrigin` allows all; dead sleep-loop goroutine exists in `main`.

## Steps

1. Delete the useless sleep goroutine in `main`.
2. Per-client write timeout; drop client on failure (already partially there).
3. Non-blocking broadcast or per-client buffer with drop-oldest.
4. Server ping/pong; client backoff reconnect (cap + jitter).
5. Configurable allowed origins (flag/env); default localhost in dev docs.

## Acceptance

- [ ] No empty background sleep loop
- [ ] Slow client cannot stall all broadcasts
- [ ] Reconnect uses backoff
- [ ] Origin policy is configurable

## Key files

- Hub package / `main.go`
- `wsjtx-web-ui/src/stores/wsjtx.ts`
