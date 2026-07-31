// Implements: REQ-011.
// Per: ADR-0031, ADR-0076.
// Discipline: C-14.

// Package layout renders a page's arrangement from data.
//
// It exists to close the last gap in the design round trip. Tokens and
// component styles already travel from a design tool back into code, because
// both are *values* — and a value can be edited safely by anyone. A page's
// arrangement is equally safe to edit (reordering sections touches no data),
// but it could not travel, because the arrangement lived in Go source. Applying
// a designer's reorder would have meant rewriting a function body, which is
// structural code generation: the kind of change that fails silently and that
// this system refuses to make.
//
// Making arrangement data removes that obstacle without weakening anything.
// A page becomes an ordered tree of elements and named sections; the sections
// themselves stay ordinary Go functions, and the data they receive stays
// ordinary Go. What a design tool may change is the arrangement — order,
// wrappers, and which sections appear — and nothing else.
//
// The boundary is enforced rather than documented: a section name with no Go
// function behind it is refused at render time, so an arrangement can never
// invent behaviour.
package layout

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

// SchemaVersion identifies the arrangement document.
const SchemaVersion = "pk.ui.page-layout.v1"

// Layout is one page's arrangement.
type Layout struct {
	SchemaVersion string `json:"schemaVersion"`
	// Page names the arrangement ("collect/stickify"), and is the coordinate
	// a change set addresses.
	Page string `json:"page"`
	Root Node   `json:"root"`
}

// Node is one element or section in the arrangement.
//
// A node is either a section reference or an element, never both: mixing them
// would make it ambiguous whether children belong to the element or to the
// section's own output.
type Node struct {
	// Element is the HTML element to emit ("section", "div").
	Element string `json:"element,omitempty"`
	// Classes is the element's class attribute, verbatim.
	Classes string `json:"classes,omitempty"`
	// Attrs are literal attributes, emitted in sorted order so the same
	// arrangement always renders the same bytes.
	Attrs map[string]string `json:"attrs,omitempty"`
	// Section names a registered renderer. When set, the node contributes
	// that renderer's output and may not carry an element or children.
	Section string `json:"section,omitempty"`
	// Children are nested nodes, in order. Order is the thing a designer
	// most often changes, and the thing this package exists to let them
	// change.
	Children []Node `json:"children,omitempty"`
}

// Sections maps a section name to the renderer that produces it. The map is
// the allowlist: an arrangement may only reference what a product has
// implemented.
type Sections map[string]func() g.Node

// Render walks an arrangement, resolving section references against the
// supplied renderers.
//
// Every failure is a refusal to guess. An unknown section is not skipped and
// not rendered empty, because either would let a page silently lose content
// that a reviewer would then have to notice by eye.
func Render(page Layout, sections Sections) (g.Node, error) {
	if page.SchemaVersion != "" && page.SchemaVersion != SchemaVersion {
		return nil, fmt.Errorf("layout: unsupported schema %q", page.SchemaVersion)
	}
	return renderNode(page.Root, sections, page.Page)
}

func renderNode(node Node, sections Sections, page string) (g.Node, error) {
	section := strings.TrimSpace(node.Section)
	element := strings.TrimSpace(node.Element)

	if section != "" {
		if element != "" || len(node.Children) > 0 {
			return nil, fmt.Errorf(
				"layout %s: section %q also declares an element or children; a node is one or the other",
				page, section)
		}
		render, known := sections[section]
		if !known {
			return nil, fmt.Errorf(
				"layout %s: section %q has no renderer; an arrangement cannot introduce one",
				page, section)
		}
		if render == nil {
			return nil, fmt.Errorf("layout %s: section %q is registered as nil", page, section)
		}
		return render(), nil
	}

	if element == "" {
		return nil, fmt.Errorf("layout %s: a node declares neither an element nor a section", page)
	}

	parts := make([]g.Node, 0, len(node.Children)+len(node.Attrs)+1)
	if node.Classes != "" {
		parts = append(parts, h.Class(node.Classes))
	}
	for _, key := range slices.Sorted(maps.Keys(node.Attrs)) {
		parts = append(parts, g.Attr(key, node.Attrs[key]))
	}
	for _, child := range node.Children {
		rendered, err := renderNode(child, sections, page)
		if err != nil {
			return nil, err
		}
		parts = append(parts, rendered)
	}
	return g.El(element, parts...), nil
}

// Validate checks an arrangement without rendering it, so a change set can be
// refused before anything is drawn.
func Validate(page Layout, sections Sections) error {
	if strings.TrimSpace(page.Page) == "" {
		return fmt.Errorf("layout: arrangement does not name its page")
	}
	if page.SchemaVersion != "" && page.SchemaVersion != SchemaVersion {
		return fmt.Errorf("layout %s: unsupported schema %q", page.Page, page.SchemaVersion)
	}
	return validateNode(page.Root, sections, page.Page)
}

func validateNode(node Node, sections Sections, page string) error {
	section := strings.TrimSpace(node.Section)
	element := strings.TrimSpace(node.Element)
	switch {
	case section != "" && (element != "" || len(node.Children) > 0):
		return fmt.Errorf(
			"layout %s: section %q also declares an element or children; a node is one or the other",
			page, section)
	case section != "":
		if _, known := sections[section]; !known {
			return fmt.Errorf(
				"layout %s: section %q has no renderer; an arrangement cannot introduce one",
				page, section)
		}
		return nil
	case element == "":
		return fmt.Errorf("layout %s: a node declares neither an element nor a section", page)
	}
	for _, child := range node.Children {
		if err := validateNode(child, sections, page); err != nil {
			return err
		}
	}
	return nil
}

// SectionNames lists the sections an arrangement references, sorted. A change
// set uses it to prove it neither introduced nor dropped a section.
func SectionNames(page Layout) []string {
	seen := map[string]struct{}{}
	collectSections(page.Root, seen)
	return slices.Sorted(maps.Keys(seen))
}

func collectSections(node Node, into map[string]struct{}) {
	if name := strings.TrimSpace(node.Section); name != "" {
		into[name] = struct{}{}
		return
	}
	for _, child := range node.Children {
		collectSections(child, into)
	}
}
