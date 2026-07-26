// Implements: REQ-011.
// Per: ADR-0031.
// Discipline: C-14.

// Package component defines renderer-neutral component identity, atomic-design
// tiers, and contribution metadata. Rendering runtimes bind their own concrete
// implementation type to Contribution without leaking that implementation
// into the OSS contract.
package component

import (
	"fmt"
	"regexp"
	"strings"
)

// ID is the stable wire identity used to select a component renderer.
type ID string

// Tier identifies a component's atomic-design composition level.
type Tier string

// Atomic-design tiers.
const (
	TierAtom     Tier = "atom"
	TierMolecule Tier = "molecule"
	TierOrganism Tier = "organism"
	TierTemplate Tier = "template"
	TierPage     Tier = "page"
)

var canonicalIDPattern = regexp.MustCompile(`^[A-Z][A-Za-z0-9]*(?:[./:][A-Za-z][A-Za-z0-9]*)*$`)

// Descriptor is the renderer-neutral identity carried by every component
// contribution. Owner names the package or module responsible for the
// contract; Version is the contract version, not an implementation build.
type Descriptor struct {
	ID      ID     `json:"id"`
	Tier    Tier   `json:"tier"`
	Owner   string `json:"owner"`
	Version string `json:"version,omitempty"`
}

// Validate rejects ambiguous or incomplete component identities.
func (d Descriptor) Validate() error {
	id := strings.TrimSpace(string(d.ID))
	if id == "" {
		return fmt.Errorf("component id is required")
	}
	if id != string(d.ID) || !canonicalIDPattern.MatchString(id) {
		return fmt.Errorf("component id %q is not canonical", d.ID)
	}
	if !d.Tier.Valid() {
		return fmt.Errorf("component %q has invalid atomic tier %q", d.ID, d.Tier)
	}
	if strings.TrimSpace(d.Owner) == "" {
		return fmt.Errorf("component %q owner is required", d.ID)
	}
	if strings.TrimSpace(d.Version) != d.Version {
		return fmt.Errorf("component %q version contains surrounding whitespace", d.ID)
	}
	return nil
}

// Valid reports whether tier is a supported atomic-design level.
func (t Tier) Valid() bool {
	switch t {
	case TierAtom, TierMolecule, TierOrganism, TierTemplate, TierPage:
		return true
	default:
		return false
	}
}

// IsCanonicalID reports whether value is a valid component wire identity.
// Namespaced extension IDs are supported without adding them to a central
// vocabulary, for example "Acme/InvoiceTimeline".
func IsCanonicalID(value string) bool {
	return value == strings.TrimSpace(value) && canonicalIDPattern.MatchString(value)
}

// Contribution binds an OSS component descriptor to a renderer-specific
// implementation. T belongs to the consuming runtime: gomponents, React,
// UIKit, A2UI, or another target can use the same contract independently.
type Contribution[T any] struct {
	Descriptor Descriptor
	Value      T
}

// NewContribution validates descriptor before binding it to value.
func NewContribution[T any](descriptor Descriptor, value T) (Contribution[T], error) {
	if err := descriptor.Validate(); err != nil {
		return Contribution[T]{}, err
	}
	return Contribution[T]{Descriptor: descriptor, Value: value}, nil
}

// MustContribution returns a validated contribution or panics during static
// application composition.
func MustContribution[T any](descriptor Descriptor, value T) Contribution[T] {
	contribution, err := NewContribution(descriptor, value)
	if err != nil {
		panic(err)
	}
	return contribution
}
