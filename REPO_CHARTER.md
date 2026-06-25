# platformkit-ui Charter

## Purpose

UI component contracts and accessibility helpers for PlatformKit. Renderer-neutral primitives that describe component roles, slot contracts, and ARIA runtime behaviour.

## In Scope

- Accessibility helpers (`pkg/a11y`): ARIA attributes, keyboard navigation contracts
- Component contracts (`pkg/contracts`): renderer-agnostic component prop and slot definitions
- UI primitives shared across rendering targets

## Out of Scope

- Specific renderers (Astro, React, Svelte — live in frontend-kit or downstream adapters)
- Design tokens or themes (handled by pk-design)
- CSS, HTML templates, or static assets
- Marketing or brand-specific UI

## Dependencies

- `maragu.dev/gomponents` — HTML element construction helpers
