// Implements: REQ-011.
// Per: ADR-0031.
// Discipline: C-14.

package surface

// Provider exposes one renderer-neutral surface contribution. Implementations
// return a valid, stable contribution and must not expose mutable internal
// state to callers.
type Provider interface {
	SurfaceContribution() Contribution
}
