# platformkit-ui

Accessible, server-friendly UI foundation packages for PlatformKit and other Go/web apps.

This repository is intentionally small. It starts with the public foundation extracted from `platformkit-frontend-kit`:

- `accessibility/`: ARIA helpers for server-rendered HTML.
- `contracts/`: renderer-neutral component prop contracts.

Runtime renderers, product components, private registry wiring, Storybook infrastructure, and PlatformKit app composition remain in `platformkit-frontend-kit`.

## Verify

```bash
make verify   # go test + go vet + pinned staticcheck + race
```
