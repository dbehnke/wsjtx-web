# 10 — Waterfall reliability

Status: `todo`  
Priority: P2  
Depends on: `02-optional-audio.md`

## Goal

Make the audio → FFT → canvas path correct and stable (queue consumption, sample rate, scroll direction).

## Why

`WaterfallDisplay.vue` has incomplete / conflicting approaches (watch on queue vs RAF loop; comments about race conditions). Backend pushes raw PCM binary frames.

## Steps

1. Document contract: Int16 mono PCM @ 12000 Hz binary WS messages.
2. Single consumer loop for `audioQueue` (no deep watch races).
3. Stable canvas scroll (newest line position matches WSJT-X mental model — pick one and document).
4. Pause when `-audio` off or tab hidden.
5. Manual check with known tone or radio audio.

## Acceptance

- [ ] No unbounded queue growth
- [ ] Visible scrolling spectrum when audio enabled
- [ ] Clean idle CPU when audio disabled

## Key files

- `wsjtx-web-ui/src/components/WaterfallDisplay.vue`
- `wsjtx-web-ui/src/stores/wsjtx.ts`
- audio capture path in Go
