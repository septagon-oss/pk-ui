// Implements: REQ-011.
// Per: ADR-0031, ADR-0076.
// Discipline: C-14.

package stylesheet

// resolve.go answers "what does this class actually declare" — the question
// that lets two descriptions of a component be compared without rendering
// either one.
//
// A design contract states a component's styling one way and the shipped
// component states it another: the contract says "w-7 h-7 bg-surface-success"
// where the component says "collect-circle-badge" and the stylesheet gives
// that class the same width, height and surface. Comparing the class names
// calls those different. Comparing what they resolve to calls them the same,
// and calls a genuine difference a difference.
//
// What is deliberately not modelled: the cascade in its full generality.
// Only rules that style a class on its own are indexed — no descendant
// selectors, no pseudo-classes, no media queries. Those describe a class in
// some context or some state, and "what does this class look like" has no
// single answer under them. Including them would mean inventing a context to
// resolve against, and an invented context produces confident wrong answers.
// Resolution stays with what a class declares unconditionally.

import (
	"regexp"
	"slices"
	"sort"
	"strings"
)

// Index is what a stylesheet declares, indexed two ways: what each class
// declares on its own, and what a document containing a given set of classes
// could have declared on it.
type Index struct {
	// declarations maps a class to its properties in stylesheet source order.
	declarations map[string][]Declaration
	// order records where each class was first defined, so merging a class
	// list follows source order the way the cascade does.
	order map[string]int
	// contextual holds every unconditional rule with the classes its selector
	// requires, for reachability.
	contextual []contextualRule
}

// contextualRule is one declaration together with what a document must contain
// for its rule to apply.
type contextualRule struct {
	declaration Declaration
	// requires is the set of classes named anywhere in the selector. A rule
	// whose subject is an element (".book img") requires only the classes.
	requires []string
}

// NewIndex parses a stylesheet and indexes what it declares.
func NewIndex(css []byte) (*Index, error) {
	all, err := IndexDeclarations(css)
	if err != nil {
		return nil, err
	}
	index := &Index{
		declarations: map[string][]Declaration{},
		order:        map[string]int{},
	}
	for position, declaration := range all {
		if declaration.AtRule != "" {
			// A media query states what the class looks like at some width.
			// The unconditional declarations are what both targets can be
			// held to; responsive behaviour needs a viewport to compare at.
			continue
		}
		for _, class := range soleClasses(declaration.Selector) {
			if _, seen := index.order[class]; !seen {
				index.order[class] = position
			}
			index.declarations[class] = append(index.declarations[class], declaration)
		}
		for _, requires := range selectorRequirements(declaration.Selector) {
			index.contextual = append(index.contextual, contextualRule{
				declaration: declaration,
				requires:    requires,
			})
		}
	}
	return index, nil
}

// soleClasses reports the classes a selector styles on their own.
//
// A selector list styles each of its subjects independently, so ".card, .tile"
// contributes to both. Within one subject the class may repeat (".pill.pill"
// is a specificity device, not a second subject) but may not name another
// element, combine with a different class, or attach a pseudo-selector — each
// of those describes the class in a context rather than by itself.
func soleClasses(selector string) []string {
	var out []string
	for _, subject := range strings.Split(selector, ",") {
		if class, ok := soleClass(strings.TrimSpace(subject)); ok {
			out = append(out, class)
		}
	}
	return out
}

func soleClass(subject string) (string, bool) {
	if strings.ContainsAny(subject, " >+~:[*") {
		return "", false
	}
	if !strings.HasPrefix(subject, ".") {
		return "", false
	}
	var class string
	for _, part := range strings.Split(strings.TrimPrefix(subject, "."), ".") {
		if part == "" {
			return "", false
		}
		if class == "" {
			class = part
			continue
		}
		if part != class {
			return "", false
		}
	}
	return class, class != ""
}

// selectorRequirements reports, for each subject in a selector list, the
// classes a document must contain for that rule to apply.
//
// State and structural pseudo-selectors are excluded entirely: ":hover" and
// ":nth-child(2)" describe a moment or a position, and a document containing
// the classes does not mean the rule is in effect.
func selectorRequirements(selector string) [][]string {
	var out [][]string
	for _, subject := range strings.Split(selector, ",") {
		subject = strings.TrimSpace(subject)
		if subject == "" || strings.ContainsAny(subject, ":[") {
			continue
		}
		var requires []string
		for _, match := range classToken.FindAllStringSubmatch(subject, -1) {
			if !slices.Contains(requires, match[1]) {
				requires = append(requires, match[1])
			}
		}
		if len(requires) == 0 {
			// An element-only rule ("body { ... }") is reachable by every
			// document, which says nothing about any component.
			continue
		}
		out = append(out, requires)
	}
	return out
}

// classToken matches one class name in a selector.
var classToken = regexp.MustCompile(`\.(-?[_a-zA-Z]+[_a-zA-Z0-9-]*)`)

// Reachable returns every declaration a document containing exactly these
// classes could have applied to it.
//
// This is reachability, not the cascade: it answers "does this stylesheet
// declare this anywhere that could take effect here", which is what comparing
// two descriptions of a component needs. It deliberately does not say which
// element the declaration lands on or which of two conflicting rules wins —
// those need a document tree, and a component's markup is not one.
//
// It exists because a hand-authored stylesheet styles most classes in context.
// Collect declares ".album-book .cover { ... }", never ".cover" alone, so
// asking what "cover" declares on its own returns nothing while the shipped
// component is plainly styled. Judging by that answer would report a
// difference that is not there.
func (i *Index) Reachable(classes ...string) []Declaration {
	present := make(map[string]struct{}, len(classes))
	for _, class := range classes {
		present[class] = struct{}{}
	}

	var out []Declaration
	for _, rule := range i.contextual {
		satisfied := true
		for _, required := range rule.requires {
			if _, ok := present[required]; !ok {
				satisfied = false
				break
			}
		}
		if satisfied {
			out = append(out, rule.declaration)
		}
	}
	return out
}

// Mentions reports whether any rule names this class at all, in any context.
//
// It is the weaker sibling of Defines, and answers a different question: not
// "what does this class do on its own" but "does the stylesheet know about it".
// A class styled only as ".album-book .cover" is mentioned but not defined, and
// telling the two apart is how a caller distinguishes styling it could not
// attribute from a class nothing styles.
func (i *Index) Mentions(class string) bool {
	for _, rule := range i.contextual {
		if slices.Contains(rule.requires, class) {
			return true
		}
	}
	return false
}

// Defines reports whether the stylesheet styles this class unconditionally.
func (i *Index) Defines(class string) bool {
	_, ok := i.declarations[class]
	return ok
}

// Declarations returns what one class declares, property to value. A property
// declared twice for the same class resolves to the last, as the cascade does.
func (i *Index) Declarations(class string) map[string]string {
	out := map[string]string{}
	for _, declaration := range i.declarations[class] {
		out[declaration.Property] = declaration.Value
	}
	return out
}

// Compute resolves a class list to the declarations it produces together.
//
// Classes are applied in stylesheet source order rather than the order they
// appear in the attribute, because that is the order the cascade uses when
// specificity ties. Classes the stylesheet does not define contribute nothing;
// callers that care whether a class was understood should ask Defines.
func (i *Index) Compute(classes ...string) map[string]string {
	known := make([]string, 0, len(classes))
	for _, class := range classes {
		if i.Defines(class) {
			known = append(known, class)
		}
	}
	sort.SliceStable(known, func(a, b int) bool {
		return i.order[known[a]] < i.order[known[b]]
	})

	out := map[string]string{}
	for _, class := range known {
		for property, value := range i.Declarations(class) {
			out[property] = value
		}
	}
	return out
}

// Classes lists every class the index understands, sorted.
func (i *Index) Classes() []string {
	out := make([]string, 0, len(i.declarations))
	for class := range i.declarations {
		out = append(out, class)
	}
	sort.Strings(out)
	return out
}
