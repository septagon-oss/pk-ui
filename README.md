# pk-ui

Accessible, server-friendly UI foundations for PlatformKit and other Go/web
applications.

`pk-ui` is the canonical OSS UI pillar:

- `accessibility/`: ARIA helpers for server-rendered HTML.
- `contracts/`: renderer-neutral component property and slot contracts.
- `component/`: canonical component identity, atomic-design tiers, ownership,
  and typed renderer contributions.
- `surface/`: canonical routes, navigation, rich page contracts, route
  ownership, section rendering, preview providers, and conformance helpers.
- `render/web/`: generic gomponents renderers for the public contracts.

PlatformKit extends these foundations downstream with its concrete component
catalog, shells, Storybook runtime, and application composition. The OSS
pillar never imports private PlatformKit packages.

## Dependency direction

```text
pk-ui contracts + component + surface + generic web renderers
                              ↓
             platformkit-frontend-kit runtime
                              ↓
                platformkit-apps composition
                              ↑
          business-module surface contributions
```

`surface.PageContract` is the single page-composition model. Downstream shells
bind it directly instead of translating it into a private parallel vocabulary.
Product and client UI stays downstream.

## Verify

```bash
make verify   # go test + go vet + pinned staticcheck + race
```
