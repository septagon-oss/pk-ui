// Implements: REQ-011.
// Per: ADR-0031, ADR-0076.
// Discipline: C-14.

package componentgen

// style_roundtrip.go closes the return path: a change made in the design tool
// becomes a reviewed edit to the product's own stylesheet.
//
// It carries style *values*, never structure. That is the decision the whole
// design rests on: styling changes constantly and cannot break behaviour or
// accessibility, while structure changes rarely and is full of judgment. A
// payload of values therefore works for every component regardless of how its
// renderer branches or loops — which is why this covers all 45 of Collect's
// components where a structural model reaches 13.
//
// The protocol is deliberately the same shape as the token round trip that
// already ships: a normalized snapshot with a content digest, a change set
// pinned to the digest it was observed against, and an all-or-nothing apply
// that refuses on any drift. Reusing that shape is not laziness — it is what
// makes the two loops behave identically under conflict, which is the only way
// an operator can trust either.

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"time"
)

const (
	// StyleSnapshotSchema identifies a snapshot of a stylesheet's declared
	// values.
	StyleSnapshotSchema = "pk.design.component-styles.v1"
	// StyleChangeSetSchema identifies a proposed set of value edits.
	StyleChangeSetSchema = "pk.design.component-style-changes.v1"
)

// StyleValue is one declared value, addressed the way the editor addresses it.
type StyleValue struct {
	Selector string `json:"selector"`
	Property string `json:"property"`
	Value    string `json:"value"`
	// AtRule scopes the declaration ("@media (min-width: 40rem)"), empty at
	// top level. It is part of the coordinate: a responsive override and its
	// base are different declarations.
	AtRule string `json:"atRule,omitempty"`
}

// Coordinate is the stable identity of one declaration.
func (v StyleValue) Coordinate() string {
	if v.AtRule == "" {
		return v.Selector + "|" + v.Property
	}
	return v.AtRule + "|" + v.Selector + "|" + v.Property
}

// StyleProvenance records who proposed a change and why. Evidence, not
// identity: it is excluded from the digest so a rebuild is reproducible.
type StyleProvenance struct {
	Source      string    `json:"source"`
	GeneratedAt time.Time `json:"generatedAt"`
	Author      string    `json:"author,omitempty"`
	Reason      string    `json:"reason,omitempty"`
}

// StyleSnapshot is the declared state of one stylesheet.
type StyleSnapshot struct {
	SchemaVersion string          `json:"schemaVersion"`
	Stylesheet    string          `json:"stylesheet"`
	Provenance    StyleProvenance `json:"provenance"`
	Styles        []StyleValue    `json:"styles"`
}

// SnapshotStylesheet reads the declared state of a stylesheet.
//
// Declarations that appear more than once at the same coordinate are omitted
// rather than recorded, because they cannot be edited unambiguously — the
// editor refuses them, so offering them in a snapshot would invite a change
// that can only fail.
func SnapshotStylesheet(css []byte, stylesheet string, provenance StyleProvenance) (StyleSnapshot, error) {
	declarations, err := IndexDeclarations(css)
	if err != nil {
		return StyleSnapshot{}, err
	}
	counts := map[string]int{}
	for _, declaration := range declarations {
		counts[styleValueOf(declaration).Coordinate()]++
	}
	styles := make([]StyleValue, 0, len(declarations))
	for _, declaration := range declarations {
		value := styleValueOf(declaration)
		if counts[value.Coordinate()] > 1 {
			continue
		}
		styles = append(styles, value)
	}
	slices.SortFunc(styles, func(a, b StyleValue) int {
		return strings.Compare(a.Coordinate(), b.Coordinate())
	})
	return StyleSnapshot{
		SchemaVersion: StyleSnapshotSchema,
		Stylesheet:    stylesheet,
		Provenance:    provenance,
		Styles:        styles,
	}, nil
}

func styleValueOf(declaration Declaration) StyleValue {
	return StyleValue{
		Selector: declaration.Selector,
		Property: declaration.Property,
		Value:    declaration.Value,
		AtRule:   declaration.AtRule,
	}
}

// Digest is the snapshot's content identity. Provenance is excluded so the
// same stylesheet always digests the same, whoever exported it and whenever.
func (s StyleSnapshot) Digest() (string, error) {
	payload := struct {
		SchemaVersion string       `json:"schemaVersion"`
		Stylesheet    string       `json:"stylesheet"`
		Styles        []StyleValue `json:"styles"`
	}{s.SchemaVersion, s.Stylesheet, s.Styles}
	raw, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("digest style snapshot: %w", err)
	}
	sum := sha256.Sum256(raw)
	return "sha256:" + hex.EncodeToString(sum[:]), nil
}

// StyleChange is one proposed value edit.
type StyleChange struct {
	Before StyleValue `json:"before"`
	After  StyleValue `json:"after"`
}

// StyleChangeSet is a set of edits pinned to the state they were observed
// against.
type StyleChangeSet struct {
	SchemaVersion string          `json:"schemaVersion"`
	Stylesheet    string          `json:"stylesheet"`
	BaseDigest    string          `json:"baseDigest"`
	Provenance    StyleProvenance `json:"provenance"`
	Changes       []StyleChange   `json:"changes"`
}

// DiffStyles produces the smallest change set carrying edited to base.
//
// Only value edits are represented. A coordinate that appears or disappears is
// a structural change to the stylesheet — a new rule or a deleted one — and
// this protocol deliberately refuses to make those, so they are reported
// rather than silently applied.
func DiffStyles(base, edited StyleSnapshot, provenance StyleProvenance) (StyleChangeSet, error) {
	if base.Stylesheet != edited.Stylesheet {
		return StyleChangeSet{}, fmt.Errorf(
			"style diff: snapshots describe different stylesheets (%q and %q)", base.Stylesheet, edited.Stylesheet)
	}
	baseIndex := make(map[string]StyleValue, len(base.Styles))
	for _, value := range base.Styles {
		baseIndex[value.Coordinate()] = value
	}
	editedIndex := make(map[string]StyleValue, len(edited.Styles))
	for _, value := range edited.Styles {
		editedIndex[value.Coordinate()] = value
	}

	var added []string
	for coordinate := range editedIndex {
		if _, known := baseIndex[coordinate]; !known {
			added = append(added, coordinate)
		}
	}
	if len(added) > 0 {
		slices.Sort(added)
		return StyleChangeSet{}, fmt.Errorf(
			"style diff: %d declaration(s) exist only in the edited snapshot (first: %s); "+
				"this protocol edits values and cannot create rules", len(added), added[0])
	}

	var changes []StyleChange
	for _, before := range base.Styles {
		after, present := editedIndex[before.Coordinate()]
		if !present {
			return StyleChangeSet{}, fmt.Errorf(
				"style diff: %s is missing from the edited snapshot; this protocol cannot delete rules",
				before.Coordinate())
		}
		if after.Value != before.Value {
			changes = append(changes, StyleChange{Before: before, After: after})
		}
	}
	if len(changes) == 0 {
		return StyleChangeSet{}, fmt.Errorf("style diff: no value changed")
	}
	digest, err := base.Digest()
	if err != nil {
		return StyleChangeSet{}, err
	}
	return StyleChangeSet{
		SchemaVersion: StyleChangeSetSchema,
		Stylesheet:    base.Stylesheet,
		BaseDigest:    digest,
		Provenance:    provenance,
		Changes:       changes,
	}, nil
}

// Validate rejects a malformed change set before it is allowed near a file.
func (c StyleChangeSet) Validate() error {
	if c.SchemaVersion != StyleChangeSetSchema {
		return fmt.Errorf("unsupported style change-set schema %q", c.SchemaVersion)
	}
	if !strings.HasPrefix(c.BaseDigest, "sha256:") || len(c.BaseDigest) != len("sha256:")+64 {
		return fmt.Errorf("baseDigest %q is not a sha256 digest", c.BaseDigest)
	}
	if strings.TrimSpace(c.Provenance.Source) == "" {
		return fmt.Errorf("style change set must record its source")
	}
	if len(c.Changes) == 0 {
		return fmt.Errorf("style change set carries no changes")
	}
	seen := make(map[string]struct{}, len(c.Changes))
	for index, change := range c.Changes {
		if change.Before.Coordinate() != change.After.Coordinate() {
			return fmt.Errorf("change %d moves a declaration between coordinates", index)
		}
		if strings.TrimSpace(change.After.Value) == "" {
			return fmt.Errorf("change %d would write an empty value", index)
		}
		if _, duplicate := seen[change.Before.Coordinate()]; duplicate {
			return fmt.Errorf("change %d repeats coordinate %s", index, change.Before.Coordinate())
		}
		seen[change.Before.Coordinate()] = struct{}{}
	}
	return nil
}

// ApplyStyleChanges applies a change set to a stylesheet, all or nothing.
//
// The base digest is checked first: a change set observed against a stylesheet
// that has since moved on is refused outright rather than applied piecemeal.
// Each edit then re-verifies the value it expected, so the guarantee survives
// even if the digest were somehow stale.
func ApplyStyleChanges(css []byte, changes StyleChangeSet) ([]byte, error) {
	if err := changes.Validate(); err != nil {
		return nil, err
	}
	current, err := SnapshotStylesheet(css, changes.Stylesheet, changes.Provenance)
	if err != nil {
		return nil, err
	}
	digest, err := current.Digest()
	if err != nil {
		return nil, err
	}
	if digest != changes.BaseDigest {
		return nil, fmt.Errorf(
			"stylesheet has changed since the change set was exported (%s, expected %s); re-export and retry",
			digest, changes.BaseDigest)
	}

	// Media-scoped edits are refused rather than mis-applied: the editor
	// addresses top-level declarations, so accepting one here would silently
	// write to the base rule instead of the override.
	for _, change := range changes.Changes {
		if change.Before.AtRule != "" {
			return nil, fmt.Errorf(
				"change to %s is scoped to %s; scoped declarations are not yet editable",
				change.Before.Coordinate(), change.Before.AtRule)
		}
	}

	out := css
	for _, change := range changes.Changes {
		out, err = SetValue(out, change.Before.Selector, change.Before.Property, change.Before.Value, change.After.Value)
		if err != nil {
			return nil, fmt.Errorf("apply %s: %w", change.Before.Coordinate(), err)
		}
	}
	return out, nil
}
