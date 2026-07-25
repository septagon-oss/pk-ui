// Implements: REQ-011.
// Per: ADR-0031.
// Discipline: C-14.

// Package web renders pk-ui component contracts to HTML with gomponents,
// styled entirely through tw class lists against the PlatformKit design
// system's role variables.
//
// Every renderer follows the same contract:
//
//   - Props in, gomponents Node out — no hidden state, no template files.
//   - Styling comes only from tw ClassLists declared in classlists.go, so an
//     application derives its stylesheet with tw/emission.For(web.ClassLists()...)
//     instead of scanning source. If a renderer used a class its list does not
//     declare, TestRenderedClassesAreDeclared fails.
//   - Accessibility is by construction: labels are associated, icons are
//     hidden from assistive tech, and states carry their ARIA attributes.
//
// The set implemented here is the working subset the PlatformKit admin and
// module pages compose. Contracts without a renderer yet are listed in
// Unimplemented, and the completeness test keeps that list honest — removing
// an entry without adding the renderer fails.
package web

import (
	"sort"
	"strconv"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"github.com/septagon-oss/pk-ui/contracts"
)

// baseAttrs renders the shared ComponentProps onto any element.
func baseAttrs(p contracts.ComponentProps, extra ...g.Node) []g.Node {
	var nodes []g.Node
	if p.ID != "" {
		nodes = append(nodes, h.ID(p.ID))
	}
	if p.Hidden {
		nodes = append(nodes, g.Attr("hidden"))
	}
	if p.Disabled {
		nodes = append(nodes, h.Disabled())
	}
	// Deterministic attribute order for stable golden tests.
	keys := make([]string, 0, len(p.Attrs))
	for k := range p.Attrs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		nodes = append(nodes, g.Attr(k, p.Attrs[k]))
	}
	return append(nodes, extra...)
}

// attrPairs renders a ComponentProps.Attrs map in deterministic order. Form
// controls apply it directly to the control element (not the field wrapper),
// so contract attributes like minlength, autocomplete, or data-* land where
// browsers and scripts read them.
func attrPairs(attrs map[string]string) []g.Node {
	keys := make([]string, 0, len(attrs))
	for k := range attrs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	nodes := make([]g.Node, 0, len(keys))
	for _, k := range keys {
		nodes = append(nodes, g.Attr(k, attrs[k]))
	}
	return nodes
}

// htmxAttrs renders HTMXProps as hx-* attributes; zero values emit nothing,
// so plain server-rendered pages carry no HTMX residue.
func htmxAttrs(p contracts.HTMXProps) []g.Node {
	var nodes []g.Node
	add := func(name, v string) {
		if v != "" {
			nodes = append(nodes, g.Attr(name, v))
		}
	}
	add("hx-get", p.Get)
	add("hx-post", p.Post)
	add("hx-put", p.Put)
	add("hx-patch", p.Patch)
	add("hx-delete", p.Delete)
	add("hx-target", p.Target)
	add("hx-swap", p.Swap)
	add("hx-trigger", p.Trigger)
	add("hx-confirm", p.Confirm)
	add("hx-ext", p.Ext)
	return nodes
}

// classes joins a declared base ClassList's compiled form with the caller's
// ComponentProps.Class escape hatch.
func classes(compiled, extra string) g.Node {
	if extra != "" {
		compiled += " " + extra
	}
	return h.Class(compiled)
}

// icon renders a decorative icon placeholder span. pk-ui ships no icon font;
// the name lands in a data attribute for the application's icon system, and
// the node is hidden from assistive technology because adjacent text carries
// the meaning.
func icon(name string) g.Node {
	if name == "" {
		return nil
	}
	return h.Span(
		h.Class(clIcon.Compile()),
		g.Attr("data-pk-icon", name),
		g.Attr("aria-hidden", "true"),
	)
}

func itoa(n int) string { return strconv.Itoa(n) }
