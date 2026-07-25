// Implements: REQ-011.
// Per: ADR-0031.
// Discipline: C-14.

package web

// molecules.go renders the molecule contracts the admin and module pages
// compose: tabular data, cards, and the navigation set (breadcrumb,
// pagination, tabs). Interactive open/close behavior (drawers, modals,
// dropdown menus) is deliberately not in this set yet: those need a
// progressive-enhancement story decided once, not five ad-hoc scripts.

import (
	"fmt"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"github.com/septagon-oss/pk-ui/contracts/atoms"
	"github.com/septagon-oss/pk-ui/contracts/molecules"
)

// Table renders molecules.TableProps. Cell values render via fmt.Sprint;
// rows are keyed by column order. An empty Rows slice renders EmptyText.
func Table(p molecules.TableProps) g.Node {
	sortable := func(c molecules.TableColumn) bool { return p.Sortable && c.Sortable }

	head := []g.Node{h.Class(clTableHead.Compile())}
	var headCells []g.Node
	if p.Selectable {
		headCells = append(headCells, h.Th(
			h.Class(clTableTh.Compile()), g.Attr("scope", "col"),
			h.Input(h.Class(clCheckbox.Compile()), h.Type("checkbox"),
				g.Attr("data-pk-select", "all"), g.Attr("aria-label", "Select all rows")),
		))
	}
	for _, c := range p.Columns {
		if sortable(c) {
			// A real button inside the th: keyboard operable, and the page
			// script owns cycling aria-sort none → ascending → descending.
			headCells = append(headCells, h.Th(
				h.Class(clTableThSort.Compile()), g.Attr("scope", "col"),
				g.Attr("aria-sort", "none"),
				h.Button(h.Class(clTableSortBtn.Compile()), h.Type("button"),
					g.Attr("data-pk-sort", c.Key),
					g.Text(c.Label),
					h.Span(g.Attr("aria-hidden", "true"), g.Attr("data-pk-sort-icon", ""), g.Text("↕")),
				),
			))
			continue
		}
		headCells = append(headCells, h.Th(
			h.Class(clTableTh.Compile()), g.Attr("scope", "col"), g.Text(c.Label),
		))
	}
	head = append(head, h.Tr(headCells...))

	tdClass := clTableTd
	if p.Compact {
		tdClass = clTableTdC
	}
	var bodyRows []g.Node
	for i, r := range p.Rows {
		rowClass := clTableRow
		if p.Striped && i%2 == 1 {
			rowClass = clTableRow.Merge(clTableRowAlt)
		}
		cells := []g.Node{h.Class(rowClass.Compile())}
		if r.ID != "" {
			cells = append(cells, g.Attr("data-pk-row", r.ID))
		}
		if p.Selectable {
			cells = append(cells, h.Td(h.Class(tdClass.Compile()),
				h.Input(h.Class(clCheckbox.Compile()), h.Type("checkbox"),
					g.Attr("data-pk-select", r.ID), g.Attr("aria-label", "Select row")),
			))
		}
		for _, c := range p.Columns {
			v := ""
			if raw, ok := r.Cells[c.Key]; ok && raw != nil {
				v = fmt.Sprint(raw)
			}
			cell := tdClass
			if c.Primary {
				cell = cell.Merge(clTableTdStrong)
			}
			td := []g.Node{h.Class(cell.Compile())}
			switch c.Align {
			case "center":
				td = append(td, g.Attr("style", "text-align:center"))
			case "right":
				td = append(td, g.Attr("style", "text-align:right"))
			}
			bodyCell := h.Td(append(td, g.Text(v))...)
			cells = append(cells, bodyCell)
		}
		bodyRows = append(bodyRows, h.Tr(cells...))
	}
	if len(bodyRows) == 0 {
		empty := p.EmptyText
		if empty == "" {
			empty = "Nothing to show yet."
		}
		span := len(p.Columns)
		if p.Selectable {
			span++
		}
		bodyRows = append(bodyRows, h.Tr(h.Td(
			h.Class(clTableTd.Compile()), h.ColSpan(itoa(span)),
			h.Class(clHelp.Compile()), g.Text(empty),
		)))
	}

	wrap := baseAttrs(p.ComponentProps, htmxAttrs(p.HTMXProps)...)
	wrap = append(wrap, classes(clTableWrap.Compile(), p.Class),
		h.Table(
			h.Class(clTable.Compile()),
			h.THead(head...),
			h.TBody(bodyRows...),
		),
	)
	return h.Div(wrap...)
}

// Card renders molecules.CardProps; a clickable card with an Href becomes a
// block anchor so the whole surface is one focusable target.
func Card(p molecules.CardProps, children ...g.Node) g.Node {
	cl := clCard
	if p.Clickable {
		cl = cl.Merge(clCardClickable)
	}
	nodes := baseAttrs(p.ComponentProps, htmxAttrs(p.HTMXProps)...)
	nodes = append(nodes, classes(cl.Compile(), p.Class))
	if p.Title != "" {
		nodes = append(nodes, h.P(h.Class(clCardTitle.Compile()), g.Text(p.Title)))
	}
	if p.Description != "" {
		nodes = append(nodes, h.P(h.Class(clCardDesc.Compile()), g.Text(p.Description)))
	}
	nodes = append(nodes, children...)
	if p.Clickable && p.Href != "" {
		return h.A(append(nodes, h.Href(p.Href))...)
	}
	return h.Div(nodes...)
}

// Breadcrumb renders molecules.BreadcrumbProps as an aria-labelled trail; the
// current page is text, not a link, and carries aria-current.
func Breadcrumb(p molecules.BreadcrumbProps) g.Node {
	sep := p.Separator
	if sep == "" {
		sep = "/"
	}
	items := []g.Node{h.Class(clBreadcrumb.Compile())}
	for i, it := range p.Items {
		if i > 0 {
			items = append(items, h.Li(
				h.Class(clBreadcrumbSep.Compile()), g.Attr("aria-hidden", "true"), g.Text(sep),
			))
		}
		current := it.Current || (it.Href == "" && i == len(p.Items)-1)
		if current {
			items = append(items, h.Li(
				h.Class(clBreadcrumbCur.Compile()), g.Attr("aria-current", "page"), g.Text(it.Label),
			))
			continue
		}
		items = append(items, h.Li(Link(atoms.LinkProps{Text: it.Label, Href: it.Href})))
	}
	nav := baseAttrs(p.ComponentProps)
	nav = append(nav, g.Attr("aria-label", "Breadcrumb"), h.Ol(items...))
	return h.Nav(nav...)
}

// Pagination renders molecules.PaginationProps as previous/next plus a
// sibling window around the current page. Page links append ?page=N to
// BaseURL; when HTMX props are set they ride along on every link.
func Pagination(p molecules.PaginationProps) g.Node {
	if p.TotalPages == 0 {
		return paginationCursor(p)
	}
	if p.TotalPages <= 1 {
		return g.Text("")
	}
	siblings := p.Siblings
	if siblings <= 0 {
		siblings = 1
	}
	pageHref := func(n int) string {
		return p.BaseURL + "?page=" + itoa(n)
	}
	pageLink := func(n int, label string, current bool, ariaLabel string) g.Node {
		cl := clPageBtn.Merge(clPageIdle)
		if current {
			cl = clPageBtn.Merge(clPageCur)
		}
		nodes := []g.Node{h.Class(cl.Compile()), h.Href(pageHref(n))}
		nodes = append(nodes, htmxAttrs(p.HTMXProps)...)
		if current {
			nodes = append(nodes, g.Attr("aria-current", "page"))
		}
		if ariaLabel != "" {
			nodes = append(nodes, g.Attr("aria-label", ariaLabel))
		}
		return h.A(append(nodes, g.Text(label))...)
	}

	items := []g.Node{h.Class(clPagination.Compile())}
	if p.CurrentPage > 1 {
		items = append(items, pageLink(p.CurrentPage-1, "‹", false, "Previous page"))
	}
	lo, hi := p.CurrentPage-siblings, p.CurrentPage+siblings
	if lo < 1 {
		lo = 1
	}
	if hi > p.TotalPages {
		hi = p.TotalPages
	}
	if lo > 1 {
		items = append(items, pageLink(1, "1", p.CurrentPage == 1, ""))
		if lo > 2 {
			items = append(items, h.Span(h.Class(clBreadcrumbSep.Compile()), g.Text("…")))
		}
	}
	for n := lo; n <= hi; n++ {
		items = append(items, pageLink(n, itoa(n), n == p.CurrentPage, ""))
	}
	if hi < p.TotalPages {
		if hi < p.TotalPages-1 {
			items = append(items, h.Span(h.Class(clBreadcrumbSep.Compile()), g.Text("…")))
		}
		items = append(items, pageLink(p.TotalPages, itoa(p.TotalPages), false, ""))
	}
	if p.CurrentPage < p.TotalPages {
		items = append(items, pageLink(p.CurrentPage+1, "›", false, "Next page"))
	}

	nav := baseAttrs(p.ComponentProps)
	nav = append(nav, g.Attr("aria-label", "Pagination"))
	nav = append(nav, items...)
	return h.Nav(nav...)
}

// paginationCursor renders offset pagination for the common API shape that
// reports no total count: previous/next controls around a page label. The
// elements carry data-pk-pagination hooks — prev, next, label — because this
// mode is a progressive-enhancement shell: the application's script owns
// fetching and enables, disables, and relabels as pages load. Server-side,
// previous starts disabled on the first page and next starts enabled.
func paginationCursor(p molecules.PaginationProps) g.Node {
	page := p.CurrentPage
	if page < 1 {
		page = 1
	}
	prev := []g.Node{
		h.Class(clPageBtn.Merge(clPageIdle).Compile()), h.Type("button"),
		g.Attr("data-pk-pagination", "prev"), g.Attr("aria-label", "Previous page"),
		g.Text("← Previous"),
	}
	if page == 1 {
		prev = append(prev, h.Disabled())
	}
	nav := baseAttrs(p.ComponentProps)
	nav = append(nav, g.Attr("aria-label", "Pagination"), h.Class(clPagination.Compile()),
		h.Button(prev...),
		h.Span(h.Class(clPageLabel.Compile()), g.Attr("data-pk-pagination", "label"),
			g.Attr("aria-live", "polite"), g.Text("Page "+itoa(page))),
		h.Button(h.Class(clPageBtn.Merge(clPageIdle).Compile()), h.Type("button"),
			g.Attr("data-pk-pagination", "next"), g.Attr("aria-label", "Next page"),
			g.Text("Next →")),
	)
	return h.Nav(nav...)
}

// SearchBar renders molecules.SearchBarProps as a labelled search field. The
// wrapper reads as the control; the input inside is borderless. Behavior —
// instant search, clearing — belongs to the application's script or HTMX
// attributes; the markup carries the affordances.
func SearchBar(p molecules.SearchBarProps) g.Node {
	label := p.Label
	if label == "" {
		label = "Search"
	}
	input := []g.Node{
		h.Class(clSearchInput.Compile()), h.Type("search"),
		g.Attr("aria-label", label), h.AutoComplete("off"),
	}
	if p.ID != "" {
		input = append(input, h.ID(p.ID))
	}
	input = append(input, attrPairs(p.Attrs)...)
	input = append(input, htmxAttrs(p.HTMXProps)...)
	if p.Placeholder != "" {
		input = append(input, h.Placeholder(p.Placeholder))
	}
	if p.SearchURL != "" {
		input = append(input, g.Attr("data-pk-search-url", p.SearchURL))
	}
	return h.Label(
		classes(clSearchWrap.Compile(), p.Class),
		h.Span(g.Attr("aria-hidden", "true"), g.Text("⌕")),
		h.Input(input...),
	)
}

// Tabs renders molecules.TabsProps as link-style tabs. Items with a URL
// navigate (or hx-get when the app enhances); static Content is not rendered
// here — the page owns the panel.
func Tabs(p molecules.TabsProps) g.Node {
	items := []g.Node{h.Class(clTabList.Compile()), h.Role("tablist")}
	for _, it := range p.Items {
		active := it.Key == p.ActiveTab
		cl := clTab.Merge(clTabIdle)
		if active {
			cl = clTab.Merge(clTabActive)
		}
		node := []g.Node{h.Class(cl.Compile()), h.Role("tab")}
		if active {
			node = append(node, g.Attr("aria-selected", "true"))
		} else {
			node = append(node, g.Attr("aria-selected", "false"))
		}
		if it.Icon != "" {
			node = append(node, icon(it.Icon))
		}
		node = append(node, g.Text(it.Label))
		if it.URL != "" {
			items = append(items, h.A(append(node, h.Href(it.URL))...))
			continue
		}
		items = append(items, h.Button(append(node, h.Type("button"))...))
	}
	nav := baseAttrs(p.ComponentProps)
	nav = append(nav, items...)
	return h.Nav(nav...)
}
