// Implements: REQ-011.
// Per: ADR-0031, ADR-0076.
// Discipline: C-14.

// Package componentgen turns an authored component document into the pk-ui
// artifacts that today are written and maintained by hand: the typed Props
// contract, the tw class lists, the web renderer, and the delivery-catalog
// entry.
//
// It exists because those artifacts are already four views of one truth with
// nothing enforcing their agreement — a duplication that has measurably
// drifted (see the Badge deviations in badge.go). Generating them from one
// source removes the drift by construction, and is the precondition for
// designing a component in Figma and having it land as verified code.
//
// Two documents, deliberately separate:
//
//   - Source is renderer-neutral: identity, props, roles, examples. It is the
//     future pk.design.component-source.v1 and belongs to pk-design, because
//     an iOS or Android renderer would consume exactly this and nothing more.
//   - WebStyle is the web projection: which tw utilities each part carries and
//     how the renderer assembles them. It is the future pk.ui.web-style.v1.
//
// Keeping them apart is what stops web utility classes from leaking into a
// contract that other renderers must implement.
package componentgen

// Source is the renderer-neutral component document.
type Source struct {
	// ID is the PascalCase component type ("Badge"). It determines the Go
	// identifiers, the wire identity, and the stable delivery ID.
	ID string
	// Tier is atom, molecule, organism, template, or page. Figma has no
	// equivalent, so it is authored rather than derived.
	Tier string
	// Presentational marks a component whose renderer is fully generated.
	// Behavioral components (sorting tables, paginators) generate everything
	// except the renderer, which is hand-written against the generated class
	// lists.
	Presentational bool
	// Description is the one-line contract description.
	Description string
	// LibraryPage and LibrarySection place the component in the design tool
	// ("02 Atoms", "Status").
	LibraryPage    string
	LibrarySection string
	// Props are declared in authored order; the generator emits them in that
	// order and lets the delivery layer sort its own view.
	Props []Prop
	// Examples are the reviewed fixtures. Exactly one must be canonical.
	Examples []Example
}

// Prop is one component property.
type Prop struct {
	// Name is the JSON/wire name ("variant"); GoName is the struct field
	// ("Variant"). They are related by convention but not derivable in both
	// directions, so both are authored.
	Name   string
	GoName string
	// Type is "string" or "bool". The delivery layer derives its own PropType
	// from the generated struct by reflection, so this only shapes the Go
	// field.
	Type string
	// Role is the pk-design PropRole: content, variant, tone, size, state,
	// modifier, or "" for a plain enum. This is the irreducibly human
	// decision — it decides which enums become native design-tool variant
	// axes — so it is authored, never inferred.
	Role string
	// Required drops omitempty from the JSON tag, which is what makes the
	// delivery layer treat the prop as required.
	Required bool
	// Description is the delivery-contract description.
	Description string
	// Enum lists the legal values for an enum prop; Default names the
	// fallback the renderer applies when the value is unset or unknown.
	Enum    []string
	Default string
}

// Example is one reviewed fixture.
type Example struct {
	Name      string
	Canonical bool
	Props     map[string]any
}

// WebStyle is the web projection of a component.
type WebStyle struct {
	// ID matches Source.ID.
	ID string
	// Element is the HTML element the root renders as.
	Element string
	// Base carries the utilities every instance has.
	//
	// Invariant: no utility here may declare a CSS property that any variant
	// also declares. Two single-class utilities have equal specificity, so a
	// collision is resolved by alphabetical sheet order — silently. The
	// generator hoists colliding properties out of Base into every variant
	// rather than emitting the collision.
	Base []Utility
	// VariantProp names the prop that selects a variant list, and Variants
	// maps each of that prop's values to its utilities.
	VariantProp string
	Variants    map[string][]Utility
	// Parts are named sub-elements (a badge's status dot) with their own
	// utilities.
	Parts map[string][]Utility
	// Body is the renderer's child sequence, in order.
	Body []Step
}

// Utility is one typed tw builder call: Method is the ClassList method name,
// Arg is the Go constant name passed to it ("Display", "DisplayInlineFlex").
// Zero-argument utilities leave Arg empty.
//
// Storing the constant name rather than the compiled class is deliberate: the
// generated code must read as hand-written typed Go, and the *Raw escape
// hatches are unresolvable by design.
type Utility struct {
	Method string
	Arg    string
}

// StepKind enumerates the renderer's child sequence.
type StepKind string

const (
	// StepBaseAttrs emits the shared id/hidden/disabled/attrs prelude.
	StepBaseAttrs StepKind = "base-attrs"
	// StepClasses emits the merged base+variant class attribute, appending
	// the caller's own Class.
	StepClasses StepKind = "classes"
	// StepPartIfBool emits a named part element when a bool prop is set.
	StepPartIfBool StepKind = "part-if-bool"
	// StepIconIfSet emits the icon placeholder when a string prop is set.
	StepIconIfSet StepKind = "icon-if-set"
	// StepText emits a string prop as the element's text.
	StepText StepKind = "text"
)

// Step is one entry in the renderer's child sequence.
type Step struct {
	Kind StepKind
	// Prop is the Go field this step reads.
	Prop string
	// Part names the entry in WebStyle.Parts, for StepPartIfBool.
	Part string
	// Element is the part's HTML element, for StepPartIfBool.
	Element string
	// Fallback is the sample text the design tool shows for a text node when
	// no content is bound, for StepText.
	Fallback string
	// Attrs are literal attributes on the part — an aria-hidden on a
	// decorative dot, for instance.
	Attrs [][2]string
}
