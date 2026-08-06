# Agent orientation

Accessible, server-friendly UI foundations for PlatformKit and other Go/web applications.

## Where things live

- top-level packages per the README

## Working rules

- `go build ./... && go test ./...` (and `make check`/`make ci` where
  defined) must stay green before opening a PR.
- Public contracts live in `pkg/` exported APIs; internal helpers stay
  unexported. Docs and the README are the advertised surface — keep them
  truthful when behavior changes.
