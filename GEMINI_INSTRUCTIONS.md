# Project Instructions

> **Canonical agent entrypoint:** [`AGENTS.md`](./AGENTS.md)  
> **Work queue:** [`todo/`](./todo/) (start at [`todo/00-overview.md`](./todo/00-overview.md))

This file remains for tools that look for `GEMINI_INSTRUCTIONS.md`. Prefer `AGENTS.md` when both are available.

## Overview

WSJT-X Web Interface is a remote control and monitoring interface for WSJT-X. Go backend bridges WSJT-X UDP to WebSockets; Vue UI consumes that stream.

## Architecture

- **Backend**: Go — UDP, WebSocket hub, embeds `dist/`, optional audio (`malgo`).
- **Frontend**: Vue 3 + TypeScript + Vite in `wsjtx-web-ui/`.
- **Protocol**: `pkg/wsjtx`.

## Development

- Backend: `go run main.go` (use `-udp :2238` if GridTracker holds `:2237`)
- Frontend: `cd wsjtx-web-ui && npm install && npm run dev`
- Details and pitfalls: `AGENTS.md` + README (GridTracker section)

## Testing

- `go test ./...`
- `cd wsjtx-web-ui && npm run test:unit -- --run`

## Deployment

- `cd wsjtx-web-ui && npm run build` then `go build` (embeds `dist/`)

## Git Workflow

1. Branch per task: `type/description` (`feat/…`, `fix/…`, `refactor/…`, `docs/…`)
2. Commit with clear messages; prefer one `todo/*.md` plan per PR
3. Push to `origin` and open a PR — do not push features directly to `main`
4. Merge after review
5. Mark the matching todo **done** and check it off in `todo/00-overview.md`
