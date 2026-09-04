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
	// variables holds custom-property declarations in stylesheet source
	// order. Their selectors remain attached so callers can resolve only the
	// scopes that are present instead of leaking an unrelated theme or
	// component's values into every document.
	variables []Declaration
}

// contextualRule is one declaration together with what a document must contain
// for its rule to apply.
type contextualRule struct {
	declaration Declaration
	// requires is the set of classes named anywhere in the selector. A rule
	// whose subject is an element (".book img") requires only the classes.
	requires []string
	// subject is the classes of the selector's last compound — the element
	// the declarations actually land on. Nil when the subject is not a class
	// compound (".book img" lands on the img, which no class identifies).
	subject []string
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
		if strings.HasPrefix(declaration.Property, "--") {
			index.variables = append(index.variables, declaration)
			continue
		}
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
		for _, requirement := range selectorRequirements(declaration.Selector) {
			index.contextual = append(index.contextual, contextualRule{
				declaration: declaration,
				requires:    requirement.requires,
				subject:     requirement.subject,
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

// subjectRequirement is what one selector-list part asks of a document, and
// where its declarations land.
type subjectRequirement struct {
	requires []string
	subject  []string
}

// selectorRequirements reports, for each subject in a selector list, the
// classes a document must contain for that rule to apply, and the classes of
// the compound the declarations land on.
//
// State and structural pseudo-selectors are excluded entirely: ":hover" and
// ":nth-child(2)" describe a moment or a position, and a document containing
// the classes does not mean the rule is in effect.
func selectorRequirements(selector string) []subjectRequirement {
	var out []subjectRequirement
	for _, part := range strings.Split(selector, ",") {
		part = strings.TrimSpace(part)
		if part == "" || strings.ContainsAny(part, ":[") {
			continue
		}
		var requires []string
		for _, match := range classToken.FindAllStringSubmatch(part, -1) {
			if !slices.Contains(requires, match[1]) {
				requires = append(requires, match[1])
			}
		}
		if len(requires) == 0 {
			// An element-only rule ("body { ... }") is reachable by every
			// document, which says nothing about any component.
			continue
		}
		out = append(out, subjectRequirement{
			requires: requires,
			subject:  lastCompoundClasses(part),
		})
	}
	return out
}

// lastCompoundClasses returns the classes of a selector part's final compound
// — the element its declarations land on. Nil when that compound carries no
// class, because then no class identifies the element being styled.
func lastCompoundClasses(part string) []string {
	compounds := strings.FieldsFunc(part, func(r rune) bool {
		return r == ' ' || r == '\t' || r == '>' || r == '+' || r == '~'
	})
	if len(compounds) == 0 {
		return nil
	}
	var out []string
	for _, match := range classToken.FindAllStringSubmatch(compounds[len(compounds)-1], -1) {
		if !slices.Contains(out, match[1]) {
			out = append(out, match[1])
		}
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

// Variables returns the unconditional custom properties available to a
// document without opting into a class-scoped context.
//
// These are what make a product's classes readable as design tokens rather
// than as raw paint. Collect declares "--sm-green: rgb(var(--surface-success,
// 69 159 80))", so a class that fills with var(--sm-green) is not stating a
// colour at all — it is naming the platform's success surface, one level
// removed. Following that chain is how a class resolves back to the token a
// design tool should bind to.
//
// Variables is retained as the context-free API. It deliberately excludes
// declarations scoped to classes; use VariablesFor when resolving a theme or
// component context. Within the unconditional scope, a later declaration wins
// as the cascade does at equal specificity.
func (i *Index) Variables() map[string]string {
	return i.VariablesFor()
}

// VariablesFor returns the custom properties available in a document whose
// ancestry and elements collectively carry classes.
//
// Root-level declarations are always included. A class-scoped declaration is
// included only when all classes named by one selector-list branch are
// present. Scoped values layer over inherited root values even when the root
// declaration appears later in the file; declarations within each layer keep
// source-order precedence. Selectors involving state, attributes, IDs, or
// conditional at-rules are excluded because a class list cannot prove those
// contexts are active.
//
// The class set is intentionally flat, matching Reachable: this API answers
// whether a variable can be available in the supplied document context, not
// which exact element inherits it.
func (i *Index) VariablesFor(classes ...string) map[string]string {
	present := make(map[string]struct{}, len(classes))
	for _, class := range classes {
		present[class] = struct{}{}
	}

	out := map[string]string{}
	for _, declaration := range i.variables {
		if declaration.AtRule == "" && hasGlobalVariableSelector(declaration.Selector) {
			out[declaration.Property] = declaration.Value
		}
	}
	for _, declaration := range i.variables {
		if declaration.AtRule == "" && hasMatchingVariableSelector(declaration.Selector, present) {
			out[declaration.Property] = declaration.Value
		}
	}
	return out
}

func hasGlobalVariableSelector(selector string) bool {
	for part := range strings.SplitSeq(selector, ",") {
		switch strings.TrimSpace(part) {
		case ":root", "html", "html:root", "body", "*":
			return true
		}
	}
	return false
}

func hasMatchingVariableSelector(selector string, present map[string]struct{}) bool {
	for part := range strings.SplitSeq(selector, ",") {
		requirements := selectorRequirements(strings.TrimSpace(part))
		for _, requirement := range requirements {
			if len(requirement.requires) == 0 {
				continue
			}
			matched := true
			for _, required := range requirement.requires {
				if _, ok := present[required]; !ok {
					matched = false
					break
				}
			}
			if matched {
				return true
			}
		}
	}
	return false
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

// SubjectStyle is what a document's stylesheet declares for the elements
// carrying one compound of classes.
type SubjectStyle struct {
	// Requires is the set of classes an element must carry for these
	// declarations to land on it — the selector's subject compound.
	// ".pack-face.pack-front" requires both; a node carrying only pack-face
	// is a different face and must not receive the front's styling.
	Requires []string
	// Declarations is what lands there, keyed by longhand property, later
	// rules overriding earlier ones as the cascade does at equal footing.
	// Specificity is not modelled; the same stated approximation as
	// Declarations.
	//
	// Longhand rather than as-written, and merged here rather than by the
	// caller, because this is where rule order still exists. "padding" and a
	// later "padding-left" both speak about the left padding; a caller
	// handed a property-keyed map has lost which came second, and any merge
	// it does is a coin toss presenting itself as a cascade.
	Declarations map[string]string
}

// SubjectStyles resolves everything a document carrying these classes
// declares, attributed to the elements it lands on.
//
// This is the difference between knowing what a document paints and knowing
// what each of its elements paints. Reachable answers the first; a design
// contract needs the second, because its nodes are elements. A stylesheet
// that writes ".pack-3d .pack-front { width: 256px }" is styling the node
// carrying pack-front — but only in a tree that also carries pack-3d, which
// is why the whole document's classes are the parameter and the subject
// compound is the result.
//
// Rules whose subject no class identifies (".book img") are excluded: their
// declarations land on an element the contract has no way to name.
func (i *Index) SubjectStyles(classes ...string) []SubjectStyle {
	present := make(map[string]struct{}, len(classes))
	for _, class := range classes {
		present[class] = struct{}{}
	}

	merged := map[string]map[string]string{}
	subjects := map[string][]string{}
	var order []string
	for _, rule := range i.contextual {
		if len(rule.subject) == 0 {
			continue
		}
		satisfied := true
		for _, required := range rule.requires {
			if _, ok := present[required]; !ok {
				satisfied = false
				break
			}
		}
		if !satisfied {
			continue
		}
		key := strings.Join(sortedCopy(rule.subject), " ")
		if merged[key] == nil {
			merged[key] = map[string]string{}
			subjects[key] = rule.subject
			order = append(order, key)
		}
		for longhand, value := range ExpandDeclaration(rule.declaration.Property, rule.declaration.Value) {
			merged[key][longhand] = value
		}
	}

	out := make([]SubjectStyle, 0, len(order))
	for _, key := range order {
		out = append(out, SubjectStyle{
			Requires:     subjects[key],
			Declarations: merged[key],
		})
	}
	return out
}

func sortedCopy(values []string) []string {
	out := slices.Clone(values)
	sort.Strings(out)
	return out
}

// MissingContext reports the classes a document would have to gain for rules
// it nearly matches to take effect, counted by how many rules each would
// enable.
//
// It answers a question that looks like a bug in the stylesheet and is usually
// a bug in the caller: markup carrying .trace-slider renders unstyled because
// every rule for it is written ".collect-trace .trace-slider", and the page
// that supplies .collect-trace is not around. A preview or a test harness that
// renders a component outside the context its styles assume shows something
// the product never displays.
//
// Only rules whose subject the document already has are considered, so this
// reports context the markup is missing rather than every rule in the
// stylesheet it fails to match.
func (i *Index) MissingContext(classes ...string) map[string]int {
	present := make(map[string]struct{}, len(classes))
	for _, class := range classes {
		present[class] = struct{}{}
	}

	missing := map[string]int{}
	for _, rule := range i.contextual {
		var absent []string
		var anyPresent bool
		for _, required := range rule.requires {
			if _, ok := present[required]; ok {
				anyPresent = true
				continue
			}
			absent = append(absent, required)
		}
		if !anyPresent || len(absent) == 0 {
			continue
		}
		for _, class := range absent {
			missing[class]++
		}
	}
	return missing
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
