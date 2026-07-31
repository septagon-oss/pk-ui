// Implements: REQ-011.
// Per: ADR-0031, ADR-0076.
// Discipline: C-14.

package layout

// roundtrip.go carries an arrangement change from the design tool back to the
// product, using the protocol already proven for tokens and component styles:
// a content digest, a change set pinned to the digest it was observed against,
// and an all-or-nothing apply. Three payloads, one set of conflict semantics —
// which is what lets an operator reason about any of them.
//
// What a designer may change is deliberately narrow: the order of siblings,
// a wrapper's classes and attributes, and whether a section appears. What they
// may not change is the set of sections a page can draw from, because that is
// the boundary between arrangement and behaviour. A change set that introduces
// a section name is refused rather than applied, since the Go function behind
// it would not exist.

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"time"
)

// ChangeSetSchema identifies an arrangement change set.
const ChangeSetSchema = "pk.ui.page-layout-changes.v1"

// Provenance records who proposed an arrangement and why. It is excluded from
// the digest so the same arrangement always identifies the same.
type Provenance struct {
	Source      string    `json:"source"`
	GeneratedAt time.Time `json:"generatedAt"`
	Author      string    `json:"author,omitempty"`
	Reason      string    `json:"reason,omitempty"`
}

// Digest is the arrangement's content identity.
func (l Layout) Digest() (string, error) {
	payload := struct {
		SchemaVersion string `json:"schemaVersion"`
		Page          string `json:"page"`
		Root          Node   `json:"root"`
	}{SchemaVersion, l.Page, l.Root}
	raw, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("digest layout: %w", err)
	}
	sum := sha256.Sum256(raw)
	return "sha256:" + hex.EncodeToString(sum[:]), nil
}

// ChangeSet is a proposed replacement arrangement, pinned to the one it was
// derived from.
//
// The payload is the whole arrangement rather than a list of edits. A page has
// one tree and reordering is most of what changes about it; expressing that as
// per-node operations would be a merge algorithm nobody needs, while the whole
// tree is small, reviewable, and unambiguous.
type ChangeSet struct {
	SchemaVersion string     `json:"schemaVersion"`
	Page          string     `json:"page"`
	BaseDigest    string     `json:"baseDigest"`
	Provenance    Provenance `json:"provenance"`
	Proposed      Layout     `json:"proposed"`
}

// Diff produces the change set carrying base to proposed, refusing any change
// that is not an arrangement change.
func Diff(base, proposed Layout, provenance Provenance) (ChangeSet, error) {
	if base.Page != proposed.Page {
		return ChangeSet{}, fmt.Errorf(
			"layout diff: arrangements describe different pages (%q and %q)", base.Page, proposed.Page)
	}
	baseSections := SectionNames(base)
	proposedSections := SectionNames(proposed)

	// A section that appears is behaviour the product does not have; one that
	// disappears is content silently dropped. Removal is legitimate, so it is
	// permitted — but introduction never is.
	for _, name := range proposedSections {
		if !slices.Contains(baseSections, name) {
			return ChangeSet{}, fmt.Errorf(
				"layout diff: %q introduces section %q, which has no renderer behind it", proposed.Page, name)
		}
	}

	baseDigest, err := base.Digest()
	if err != nil {
		return ChangeSet{}, err
	}
	proposedDigest, err := proposed.Digest()
	if err != nil {
		return ChangeSet{}, err
	}
	if baseDigest == proposedDigest {
		return ChangeSet{}, fmt.Errorf("layout diff: the arrangement did not change")
	}
	return ChangeSet{
		SchemaVersion: ChangeSetSchema,
		Page:          base.Page,
		BaseDigest:    baseDigest,
		Provenance:    provenance,
		Proposed:      proposed,
	}, nil
}

// Validate rejects a malformed change set before it reaches an arrangement.
func (c ChangeSet) Validate() error {
	if c.SchemaVersion != ChangeSetSchema {
		return fmt.Errorf("unsupported layout change-set schema %q", c.SchemaVersion)
	}
	if !strings.HasPrefix(c.BaseDigest, "sha256:") || len(c.BaseDigest) != len("sha256:")+64 {
		return fmt.Errorf("baseDigest %q is not a sha256 digest", c.BaseDigest)
	}
	if strings.TrimSpace(c.Provenance.Source) == "" {
		return fmt.Errorf("layout change set must record its source")
	}
	if strings.TrimSpace(c.Page) == "" {
		return fmt.Errorf("layout change set does not name its page")
	}
	if c.Proposed.Page != c.Page {
		return fmt.Errorf("layout change set proposes an arrangement for a different page")
	}
	return nil
}

// Apply replaces an arrangement with a proposed one, refusing on any drift.
//
// The guards compose: the change set must be well formed, must have been
// observed against the arrangement it is being applied to, must not introduce
// a section, and the result must render against the product's real renderers.
// The last check is what makes this safe to automate — an arrangement that
// cannot draw is refused before it can replace one that can.
func Apply(current Layout, changes ChangeSet, sections Sections) (Layout, error) {
	if err := changes.Validate(); err != nil {
		return Layout{}, err
	}
	if current.Page != changes.Page {
		return Layout{}, fmt.Errorf(
			"layout apply: change set is for %q, not %q", changes.Page, current.Page)
	}
	digest, err := current.Digest()
	if err != nil {
		return Layout{}, err
	}
	if digest != changes.BaseDigest {
		return Layout{}, fmt.Errorf(
			"layout %s has changed since the change set was exported (%s, expected %s); re-export and retry",
			current.Page, digest, changes.BaseDigest)
	}
	for _, name := range SectionNames(changes.Proposed) {
		if !slices.Contains(SectionNames(current), name) {
			return Layout{}, fmt.Errorf(
				"layout apply: %q introduces section %q, which has no renderer behind it", changes.Page, name)
		}
	}
	if err := Validate(changes.Proposed, sections); err != nil {
		return Layout{}, err
	}
	next := changes.Proposed
	next.SchemaVersion = SchemaVersion
	return next, nil
}
