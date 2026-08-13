# pk-ui

Accessible, server-friendly UI foundations for PlatformKit and other Go/web
applications.

`pk-ui` is the canonical OSS UI pillar:

- `accessibility/`: ARIA helpers for server-rendered HTML.
- `contracts/`: renderer-neutral component property and slot contracts.
- `component/`: canonical component identity, atomic-design tiers, ownership,
  and typed renderer contributions.
- `surface/`: canonical routes, navigation, rich page contracts, entity-route
  publication, hypermedia protocol names, complete error-document rendering,
  admin-page canonicalization, namespaced composition, section rendering,
  provider contracts, route ownership, and conformance helpers.
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

## Adaptive data semantics

`WindowedCollection` is the canonical bounded list/grid shell. Its state axis
is `ready`, `loading`, `error`, or `offline`; an empty view is derived from
`itemCount`. Offline projections keep already verified items visible, use
`offlineTitle` and `offlineDescription` for localized recovery copy, and do
not expose cursor controls until the collection returns to `ready`.

`DetailList` is the canonical label/value section. `title` and `description`
are visible copy. `semanticRole` is an optional stable, non-localized machine
key (for example `identity`, `preferences`, or `activity`) that adaptive
renderers may use to preserve information architecture without interpreting
translated item labels. It is never projected as an HTML or ARIA role.

## Verify

```bash
make verify   # go test + go vet + pinned staticcheck + race
```
