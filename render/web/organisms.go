// Implements: REQ-011.
// Per: ADR-0031.
// Discipline: C-14.

package web

// organisms.go begins the organism layer: self-contained page sections
// composed from the molecule and atom renderers in this package. An organism
// owns the ORDER of its parts; page-specific regions (live status lines,
// empty-state panels) slot in through the variadic children, which render
// between the table and the pagination — the one point where a data page
// interleaves its own state.

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"github.com/septagon-oss/pk-ui/contracts/atoms"
	"github.com/septagon-oss/pk-ui/contracts/organisms"
)

// DataGrid renders organisms.DataGridProps: a toolbar (search, filters,
// actions, and an optional create link), the table, any page-owned children,
// then pagination. The toolbar carries role="search" when a search control
// is present. CreateURL renders as a primary Button-classed anchor pushed to
// the toolbar's end.
func DataGrid(p organisms.DataGridProps, children ...g.Node) g.Node {
	var toolbar []g.Node
	search := p.Search
	if search.SearchURL == "" && p.SearchURL != "" {
		search.SearchURL = p.SearchURL
	}
	hasSearch := search.ID != "" || search.Label != "" || search.Placeholder != "" || search.SearchURL != ""
	if hasSearch {
		toolbar = append(toolbar, SearchBar(search))
	}
	for _, f := range p.Filters {
		switch f.Type {
		case "select":
			options := make([]atoms.SelectOption, 0, len(f.Options))
			for _, o := range f.Options {
				options = append(options, atoms.SelectOption{Label: o.Label, Value: o.Value})
			}
			toolbar = append(toolbar, Select(atoms.SelectProps{
				Name: f.Key, Placeholder: f.Label, Options: options,
			}))
		default: // text, date
			typ := f.Type
			if typ != "date" {
				typ = "text"
			}
			toolbar = append(toolbar, Input(atoms.InputProps{
				Name: f.Key, Type: typ, Placeholder: f.Label,
			}))
		}
	}
	for _, a := range p.Actions {
		toolbar = append(toolbar, Button(a))
	}
	if p.CreateURL != "" {
		text := p.CreateText
		if text == "" {
			text = "New"
		}
		toolbar = append(toolbar, h.A(
			h.Class(ButtonClasses("primary", "medium").Merge(clGridCreate).Compile()),
			h.Href(p.CreateURL), g.Text(text),
		))
	}

	section := baseAttrs(p.ComponentProps, htmxAttrs(p.HTMXProps)...)
	section = append(section, classes(clGridSection.Compile(), p.Class))
	if len(toolbar) > 0 {
		bar := []g.Node{h.Class(clGridToolbar.Compile())}
		if hasSearch {
			bar = append(bar, h.Role("search"))
		}
		section = append(section, h.Div(append(bar, toolbar...)...))
	}
	section = append(section, Table(p.Table))
	section = append(section, children...)
	section = append(section, Pagination(p.Pagination))
	return h.Section(section...)
}
