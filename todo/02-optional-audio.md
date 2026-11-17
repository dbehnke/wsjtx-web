# 02 — Optional audio

Status: `done`  
Priority: P0  
Depends on: `01-split-go-packages.md` (ideal), can be partial before

## Goal

Audio capture (`malgo` / CGO) is opt-in via `-audio` (default **off**). Core bridge runs without audio devices or noisy PulseAudio logs.

## Why

Audio complicates startup, adds CGO surface, and is only needed for waterfall. Most users want decode table + controls.

## Scope

**In**

- Flag `-audio` (bool, default false)
- Skip `initAudioManager` / `Start` when false
- WS commands `getAudioDevices` / `setAudioDevice` return a clear error or no-op when disabled
- README note

**Out**

- Rewriting waterfall DSP
- Removing malgo from go.mod entirely (optional follow-up)

## Steps

1. Add flag; wire through server/audio package.
2. Ensure nil/disabled paths never panic.
3. Confirm `go run main.go` with GridTracker-free port starts cleanly without audio init spam.
4. Document waterfall needing `-audio`.

## Acceptance

- [x] Default run does not init malgo
- [x] `-audio` restores previous capture behavior
- [x] No nil pointer paths on device commands when disabled

## Key files

- `main.go` / `internal/audio` / `internal/server`
- `README.md`
- `wsjtx-web-ui/src/stores/wsjtx.ts` (tolerate missing devices)

## Risks

- Frontend may assume devices always arrive — handle empty list gracefully.
