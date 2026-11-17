# 08 — Decode filters and band-activity UX

Status: `todo`  
Priority: P2  
Depends on: `03-purge-vue-scaffold.md` (nice), `05-shared-protocol-types.md` (nice)

## Goal

Make the decode list usable on busy bands: filter, search, clear, and clearer CQ / DX highlighting.

## Why

Unfiltered unshift of last 100 decodes is noisy for real operation.

## Scope

**In**

- Filter by call / grid / CQ / SNR threshold (client-side first)
- Clear list control
- Visual cues for CQ and own call (from status DE call)

**Out**

- Full logging / ADIF browser (future)

## Steps

1. Extend Pinia store with filter state + computed visible messages.
2. UI controls above `BandActivity` table.
3. Keep double-click reply behavior.

## Acceptance

- [ ] Filters change visible rows without dropping WS backlog incorrectly (document retention policy)
- [ ] Clear works
- [ ] Reply still works on visible rows

## Key files

- `wsjtx-web-ui/src/components/BandActivity.vue`
- `wsjtx-web-ui/src/stores/wsjtx.ts`
