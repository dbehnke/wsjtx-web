# 09 — Status and controls parity

Status: `todo`  
Priority: P2  
Depends on: `05-shared-protocol-types.md` (ideal)

## Goal

Expose more WSJT-X status fields and wire real controls beyond Halt (e.g. Enable Tx / free text where protocol allows).

## Why

`ControlPanel` has placeholders (`Enable Tx` logs only). Operators expect parity with common remote panels.

## Steps

1. Audit `StatusMessage` fields vs UI.
2. Map protocol messages already in `pkg/wsjtx` (`HaltTx`, `FreeText`, `Replay`, `Configure`, …).
3. Implement safe controls with confirmations where destructive.
4. Disable controls until a WSJT-X peer address is known (backend already tracks `lastAddr`).

## Acceptance

- [ ] Enable Tx (or documented substitute) works end-to-end against WSJT-X
- [ ] Status panel shows mode, TX state, DX/DE, dial freq accurately
- [ ] No panic when no peer yet — user-visible error

## Key files

- `pkg/wsjtx/protocol.go`
- server WS command switch
- `ControlPanel.vue`
