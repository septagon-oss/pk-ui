# pk-ui Charter

## Purpose

Canonical OSS UI contracts and generic server-rendered primitives for
PlatformKit and other Go/web applications.

## In Scope

- Accessibility helpers: ARIA attributes and keyboard-navigation contracts
- Renderer-neutral component property and slot contracts
- Canonical component identity, atomic-design tiers, ownership metadata, and
  typed renderer contributions
- Renderer-neutral surface, route, navigation, rich-page, presenter, template,
  preview, and route-ownership contracts
- Generic gomponents renderers for public contracts
- Conformance helpers for downstream providers

## Out of Scope

- Product-specific components, shells, or registries
- PlatformKit application composition and dependency-injection wiring
- Client, tenant, or brand-specific UI
- Framework-specific renderers such as React, Astro, or Svelte
- Design tokens or themes (handled by pk-design)
- Product CSS, templates, or static assets

## Extension Rule

- OSS defines identity, contracts, validation, cloning, generic rendering, and
  conformance.
- Downstream runtimes bind those contracts to concrete components and shells.
- Applications register product and business contributions explicitly.
- Downstream packages must not fork or flatten the canonical component and page
  models into parallel vocabularies.
- This repository must never import a private `github.com/septagon-dev/...`
  package.

## Dependencies

- `maragu.dev/gomponents` — HTML element construction helpers
- `github.com/septagon-oss/styleengine` — typed generic CSS emission
- `github.com/septagon-oss/tw` — typed utility-class vocabulary
