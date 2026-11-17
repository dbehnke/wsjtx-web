# Roadmap overview

Work top to bottom. Skip only if blocked and note why.

## Refactors (foundation)

- [x] `01-split-go-packages.md` — Split `main.go` into packages
- [x] `02-optional-audio.md` — Gate audio behind a flag; default off
- [ ] `03-purge-vue-scaffold.md` — Remove Vue CLI/scaffold leftovers
- [ ] `04-dev-proxy-and-embed.md` — Clean local-dev vs production embed story
- [ ] `05-shared-protocol-types.md` — Typed WS message contract
- [ ] `06-websocket-hardening.md` — Hub reliability and safety
- [ ] `07-tests-and-ci.md` — Protocol tests + GitHub Actions

## Features

- [ ] `08-decode-filters-and-ux.md` — Band activity filters / usability
- [ ] `09-status-and-controls-parity.md` — Richer status + Enable Tx / config commands
- [ ] `10-waterfall-reliability.md` — Fix streaming FFT/waterfall path

## Ops

- [ ] `12-containerize.md` — Dockerfile / Compose multi-stage image (UI build + Go embed)

## Longer-term direction

- [ ] `11-frontend-go-htmx-migration.md` — Optional move off Node SPA toward Go-served UI

## Recently landed (context, not queue)

- Configurable `-udp` / `-http`; GridTracker docs in README
- `dist/index.html` placeholder for `go:embed`
- Audio nil-guards; Vite high-severity dep upgrades (picomatch still unpatched upstream)
