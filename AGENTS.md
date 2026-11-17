# Agent Guide — wsjtx-web

This repository is set up for AI agents. Read this file first, then work from `todo/`.

## What this project is

Remote monitor/control UI for [WSJT-X](https://wsjt.sourceforge.io/). A Go process receives WSJT-X UDP traffic, bridges it over WebSockets, and serves (or proxies to) a web UI.

| Layer | Location | Role |
|-------|----------|------|
| Backend | `main.go`, `pkg/wsjtx/` | UDP ↔ WebSocket bridge, optional audio capture, embeds `dist/` |
| Protocol | `pkg/wsjtx/protocol.go` | WSJT-X binary UDP message encode/decode |
| Frontend | `wsjtx-web-ui/` | Vue 3 + TypeScript + Vite + Pinia + Tailwind |
| Built UI | `dist/` | Produced by Vite (`outDir: ../dist`); placeholder `dist/index.html` kept for `go:embed` |

Fork of `dbehnke/wsjtx-web`. Push to `origin` (`justsomed4d/wsjtx-web`). Upstream is `dbehnke/wsjtx-web`.

## How to run (agent checklist)

```bash
# Backend (needs dist/ for embed — placeholder already in repo)
go run main.go

# If UDP :2237 is taken (often GridTracker):
go run main.go -udp :2238 -http :8080

# Optional: enable malgo capture for waterfall (default is off)
go run main.go -audio

# Frontend dev server
cd wsjtx-web-ui && npm install && npm run dev
```

- Backend HTTP: `http://localhost:8080`
- Vite UI: `http://localhost:5173`
- WSJT-X must send UDP to whatever `-udp` address is in use; enable **Accept UDP requests** for reply/halt.
- Waterfall needs `-audio`; without it, `getAudioDevices` / `setAudioDevice` return an empty list / no-op.

See README **Using with GridTracker** for coexistence.

## Source of truth for work

1. **`todo/README.md`** — how the queue works  
2. **`todo/00-overview.md`** — ordered roadmap  
3. Next unchecked item in that overview — open its plan file and execute  

When finishing a todo: mark it done in the plan file and in `00-overview.md`, note follow-ups.

## Hard constraints

- Do **not** push directly to `main` for feature work — use `type/description` branches and PRs (`feat/…`, `fix/…`, `refactor/…`, `docs/…`).
- Do **not** commit secrets, `.env`, or full `dist/assets` build output (only tracked placeholder `dist/index.html`).
- Do **not** `npm audit fix --force` — it may downgrade Vite and reintroduce real advisories. Remaining picomatch highs may be unpatched upstream.
- Prefer small PRs that match one `todo/*.md` plan.
- Keep changes scoped; no drive-by refactors outside the active todo unless required to land it.

## Preferred architecture direction

Documented in `todo/` (do not invent a conflicting long-term plan):

1. Split Go monolith into packages; make audio optional.  
2. Purge Vue scaffold leftovers; fix local-dev proxy/embed story.  
3. Shared/typed WS protocol; harden WebSocket hub.  
4. Tests + CI.  
5. Longer term: consider Go-served UI (HTMX/vanilla) to reduce Node — only via its todo, not ad hoc.

## Commands

| Task | Command |
|------|---------|
| Go test | `go test ./...` |
| UI unit tests | `cd wsjtx-web-ui && npm run test:unit -- --run` |
| UI build | `cd wsjtx-web-ui && npm run build` |
| UI lint | `cd wsjtx-web-ui && npm run lint` |
| Prod binary | `cd wsjtx-web-ui && npm run build && cd .. && go build` |

## Key pitfalls agents have already hit

- `//go:embed dist` fails if `dist/` is missing — keep placeholder `dist/index.html`.  
- GridTracker binds `:2237` — use `-udp` and document WSJT-X port changes.  
- Audio is opt-in (`-audio`); without it `audio.Manager` is nil — guard before use.  
- Message types on the wire are stringified ints (`"1"` status, `"2"` decode) plus named commands (`reply`, `halt`, `getAudioDevices`, …).

## Cursor rules

Project rules live in `.cursor/rules/`. They reinforce this guide; if they conflict with `AGENTS.md` or an active todo, follow `AGENTS.md` + the todo.
