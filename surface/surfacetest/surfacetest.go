// Package surfacetest provides the driver conformance suite for the
// contribution contract and a reference Static provider usable as a test fake.
package surfacetest

// surfacetest.go — surface provider conformance suite.
//
// Implements: REQ-011.
// Per: ADR-0031.
// Discipline: C-14.

import (
	"encoding/json"
	"testing"

	"github.com/septagon-oss/pk-ui/surface"
)

// ProviderConformance runs the behavioral contract every surface.Provider
// must satisfy: a valid contribution, stable across calls, and isolated
// from caller mutation.
func ProviderConformance(t *testing.T, newProvider func(t *testing.T) surface.Provider) {
	t.Helper()

	t.Run("contribution is valid", func(t *testing.T) {
		p := newProvider(t)
		if err := surface.ValidateContribution(p.SurfaceContribution()); err != nil {
			t.Fatalf("SurfaceContribution is invalid: %v", err)
		}
	})

	t.Run("contribution is stable across calls", func(t *testing.T) {
		p := newProvider(t)
		first := snapshot(t, p.SurfaceContribution())
		second := snapshot(t, p.SurfaceContribution())
		if first != second {
			t.Fatalf("SurfaceContribution changed between calls:\n%s\n---\n%s", first, second)
		}
	})

	t.Run("caller mutation does not leak back", func(t *testing.T) {
		p := newProvider(t)
		before := snapshot(t, p.SurfaceContribution())
		got := p.SurfaceContribution()
		if len(got.Routes) > 0 {
			got.Routes[0].ID = "mutated-by-caller"
		}
		got.Widgets = append(got.Widgets, surface.Widget{ID: "injected", Kind: "stats"})
		after := snapshot(t, p.SurfaceContribution())
		if before != after {
			t.Fatal("mutating a returned contribution changed the provider's state; providers must return isolated copies or treat-as-read-only documented values")
		}
	})
}

// snapshot renders a contribution to canonical JSON for comparison.
func snapshot(t *testing.T, c surface.Contribution) string {
	t.Helper()
	b, err := json.Marshal(c)
	if err != nil {
		t.Fatalf("marshal contribution: %v", err)
	}
	return string(b)
}

// Static is a reference surface.Provider returning a deep copy of a fixed
// contribution on every call.
type Static struct {
	Contribution surface.Contribution
}

var _ surface.Provider = Static{}

// SurfaceContribution implements surface.Provider.
func (s Static) SurfaceContribution() surface.Contribution {
	return clone(s.Contribution)
}

// clone deep-copies a contribution via JSON round-trip; reference-driver
// simplicity is preferred over hand-written copy code here.
func clone(c surface.Contribution) surface.Contribution {
	b, err := json.Marshal(c)
	if err != nil {
		return surface.Contribution{}
	}
	var out surface.Contribution
	if err := json.Unmarshal(b, &out); err != nil {
		return surface.Contribution{}
	}
	return out
}
