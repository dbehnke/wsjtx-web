# 04 — Dev proxy and embed story

Status: `todo`  
Priority: P1

## Goal

One clear local-dev path and one clear production path — no conflicting README advice.

## Why

`go:embed dist` fights “run Vite separately”. Agents stumble on missing `dist` and dual origins for `/ws`.

## Proposed model

| Mode | How |
|------|-----|
| **Dev** | Vite on `:5173` proxies `/ws` (and optional `/api`) to Go on `:8080`. Go may skip serving SPA or serve placeholder only. |
| **Prod** | `npm run build` → `../dist` → `go build` embeds `dist`. |

## Steps

1. Add Vite `server.proxy` for `/ws` → `http://127.0.0.1:8080`.
2. Decide: Go always embeds (placeholder OK in dev) vs build tag `embedui` — prefer keep embed + placeholder for simplicity.
3. Update README + `AGENTS.md` to a single recommended `make dev` or two-terminal recipe.
4. Optional: Makefile `dev` / `run` targets.

## Acceptance

- [ ] From Vite origin, WebSocket connects without hard-coding host hacks
- [ ] Fresh clone: `go run` works with tracked placeholder
- [ ] Prod build still one binary with real UI assets

## Key files

- `wsjtx-web-ui/vite.config.ts`
- `wsjtx-web-ui/src/stores/wsjtx.ts` (WS URL uses `window.location.host` — good with proxy)
- `README.md`, `Makefile`, `AGENTS.md`
