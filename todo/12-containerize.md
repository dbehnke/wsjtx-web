# 12 — Containerize

Status: `todo`  
Priority: P2  
Depends on: `02-optional-audio.md` (strongly recommended), `04-dev-proxy-and-embed.md` (nice)

## Goal

Ship a reproducible container image and a simple run path so agents and operators can start wsjtx-web without installing Go/Node locally.

## Why

Local setup depends on Go, Node, `dist` embed, and UDP port conflicts (GridTracker). A multi-stage image makes demos and deploys predictable.

## Scope

**In**

- Multi-stage `Dockerfile`: Node build of `wsjtx-web-ui` → `dist/`, then Go build with embed
- `.dockerignore` (exclude `node_modules`, `.git`, binaries)
- Optional `compose.yaml` for documented ports/env
- README + `AGENTS.md` run snippets
- Default image runs **without** host audio (`-audio` off / omit malgo init) unless clearly documented otherwise

**Out**

- Kubernetes / Helm
- Publishing to a registry (optional follow-up once image builds)
- Forcing audio devices into the container (Pulse/ALSA device mapping is host-specific)

## Recommended design

```text
stage 1 (node):  npm ci && npm run build  → /out/dist
stage 2 (golang): copy source + dist, CGO_ENABLED=0 go build -o /wsjtx-web
final:           scratch or distroless/static, copy binary
```

- Expose `8080/tcp` (HTTP/WS) and document UDP (e.g. `2237/udp` or `2238/udp`).
- Prefer **host network** or explicit UDP publish; bridge NAT often breaks WSJT-X discovery/reply sourcing — document the chosen mode.
- Entrypoint flags via `CMD`/`command`: e.g. `-http :8080 -udp :2237`.
- Build with audio **disabled by default** in-container so `CGO_ENABLED=0` static binary stays easy. If audio-in-container is required later, use a glibc/Debian runtime stage and document device mounts (separate follow-up).

## Steps

1. Add `.dockerignore`.
2. Add multi-stage `Dockerfile` that always produces a real UI `dist` (not only the placeholder).
3. Verify: `docker build -t wsjtx-web .` and run with published ports; hit `/` and `/ws`.
4. Document GridTracker conflict: e.g. `-udp :2238` + WSJT-X pointing at the container’s UDP port / host IP.
5. Optional: `compose.yaml` with `ports` / `network_mode: host` profiles (`default` vs `host`).
6. Optional: CI job to `docker build` (can extend `07-tests-and-ci.md`).

## Acceptance

- [ ] `docker build` succeeds on a clean machine
- [ ] Container serves embedded UI on `8080`
- [ ] UDP listen address configurable at run time
- [ ] README documents run, ports, and GridTracker coexistence
- [ ] No requirement for Node/Go on the host to *run* the image

## Key files

- `Dockerfile`
- `.dockerignore`
- `compose.yaml` (optional)
- `README.md`, `AGENTS.md`

## Risks

- UDP + Docker networking: reply packets must reach WSJT-X; test on Linux with host network first.
- `go:embed` needs built `dist` in the builder context — don’t copy only the git placeholder into the final build stage.
- Enabling `malgo` in-container pulls CGO and audio libs — keep that out of the default image.
