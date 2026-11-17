# 03 — Purge Vue scaffold leftovers

Status: `todo`  
Priority: P1

## Goal

Remove create-vue boilerplate so the UI tree only contains product code.

## Why

Scaffold files confuse agents and humans (`HelloWorld`, welcome icons, counter store, About route).

## Scope

**In** — delete or stop routing:

- `src/components/HelloWorld.vue`
- `src/components/TheWelcome.vue`
- `src/components/WelcomeItem.vue`
- `src/components/icons/*` (if unused)
- `src/stores/counter.ts`
- `src/views/AboutView.vue`
- `src/components/__tests__/HelloWorld.spec.ts`
- Router entries / imports that reference the above

**Out**

- Redesigning HomeView layout
- Tailwind restyle

## Steps

1. Grep for imports of scaffold symbols.
2. Delete unused files; fix router (`src/router/index.ts`).
3. Replace HelloWorld test with a minimal smoke test (e.g. store helper) or drop until `07`.
4. `npm run build` + `npm run test:unit -- --run`.

## Acceptance

- [ ] No scaffold components remain referenced
- [ ] Build and unit tests pass
- [ ] App still loads HomeView with BandActivity / ControlPanel / WaterfallDisplay

## Key files

- `wsjtx-web-ui/src/**`
- `wsjtx-web-ui/src/router/index.ts`
