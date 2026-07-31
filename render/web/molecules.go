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
	"net/url"
	"strings"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"github.com/septagon-oss/pk-ui/contracts"
	"github.com/septagon-oss/pk-ui/contracts/atoms"
	"github.com/septagon-oss/pk-ui/contracts/molecules"
)

// TableSlots is the trusted Go composition seam for rich web cells and
// server-driven sorting. Portable table data remains in TableProps; callers
// only opt into these callbacks when a cell needs real markup.
type TableSlots struct {
	Cell           func(molecules.TableRow, molecules.TableColumn) g.Node
	SortURL        func(molecules.TableColumn) string
	SortState      func(molecules.TableColumn) string
	SelectAllLabel string
	SelectRowLabel func(molecules.TableRow) string
}

// Table renders molecules.TableProps. Cell values render via fmt.Sprint;
// rows are keyed by column order. An empty Rows slice renders EmptyText.
func Table(p molecules.TableProps) g.Node {
	return TableWithSlots(p, TableSlots{})
}

// TableWithSlots renders the canonical table while allowing trusted Go
// composition to project rich cell nodes without creating a second table
// renderer.
func TableWithSlots(p molecules.TableProps, slots TableSlots) g.Node {
	sortable := func(c molecules.TableColumn) bool { return p.Sortable && c.Sortable }

	head := []g.Node{h.Class(clTableHead.Compile())}
	var headCells []g.Node
	if p.Selectable {
		label := fallbackText(slots.SelectAllLabel, "Select all rows")
		headCells = append(headCells, h.Th(
			h.Class(clTableTh.Compile()), g.Attr("scope", "col"),
			h.Input(h.Class(clCheckbox.Compile()), h.Type("checkbox"),
				g.Attr("data-pk-select", "all"), g.Attr("aria-label", label)),
		))
	}
	for _, c := range p.Columns {
		if sortable(c) {
			sortState := "none"
			if slots.SortState != nil {
				switch state := slots.SortState(c); state {
				case "ascending", "descending":
					sortState = state
				}
			}
			glyph := "↕"
			switch sortState {
			case "ascending":
				glyph = "↑"
			case "descending":
				glyph = "↓"
			}
			button := []g.Node{
				h.Class(clTableSortBtn.Compile()), h.Type("button"),
				g.Attr("data-pk-sort", c.Key),
				g.Text(c.Label),
				h.Span(g.Attr("aria-hidden", "true"), g.Attr("data-pk-sort-icon", ""), g.Text(glyph)),
			}
			if slots.SortURL != nil {
				if sortURL := slots.SortURL(c); sortURL != "" {
					enhancement := p.HTMXProps
					enhancement.Get = sortURL
					enhancement.Trigger = ""
					button = append(button, htmxAttrs(enhancement)...)
				}
			}
			// A real button inside the th: keyboard operable, and the page
			// script owns cycling aria-sort none → ascending → descending.
			headCells = append(headCells, h.Th(
				h.Class(clTableThSort.Compile()), g.Attr("scope", "col"),
				g.Attr("aria-sort", sortState),
				g.If(c.Width != "", g.Attr("style", "width:"+c.Width)),
				h.Button(button...),
			))
			continue
		}
		cell := []g.Node{h.Class(clTableTh.Compile()), g.Attr("scope", "col")}
		if c.Width != "" {
			cell = append(cell, g.Attr("style", "width:"+c.Width))
		}
		headCells = append(headCells, h.Th(append(cell, g.Text(c.Label))...))
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
			label := "Select row"
			if slots.SelectRowLabel != nil {
				label = fallbackText(slots.SelectRowLabel(r), label)
			}
			cells = append(cells, h.Td(h.Class(tdClass.Compile()),
				h.Input(h.Class(clCheckbox.Compile()), h.Type("checkbox"),
					g.Attr("data-pk-select", r.ID), g.Attr("aria-label", label)),
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
			content := g.Node(g.Text(v))
			if slots.Cell != nil {
				if rich := slots.Cell(r, c); rich != nil {
					content = rich
				}
			}
			bodyCell := h.Td(append(td, content)...)
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
	visible := breadcrumbItems(p.Items, p.MaxItems)
	items := []g.Node{h.Class(clBreadcrumb.Compile())}
	for i, it := range visible {
		if i > 0 {
			items = append(items, h.Li(
				h.Class(clBreadcrumbSep.Compile()), g.Attr("aria-hidden", "true"), g.Text(sep),
			))
		}
		if it.Label == "…" && it.Href == "" && !it.Current {
			items = append(items, h.Li(
				h.Class(clBreadcrumbSep.Compile()),
				g.Attr("aria-hidden", "true"),
				g.Attr("data-breadcrumb-ellipsis", ""),
				g.Text("…"),
			))
			continue
		}
		current := it.Current || (it.Href == "" && i == len(visible)-1)
		if current {
			content := []g.Node{g.Text(it.Label)}
			if it.Icon != "" {
				content = append([]g.Node{icon(it.Icon)}, content...)
			}
			items = append(items, h.Li(append(
				[]g.Node{h.Class(clBreadcrumbCur.Compile()), g.Attr("aria-current", "page")},
				content...,
			)...))
			continue
		}
		var adornment []g.Node
		if it.Icon != "" {
			adornment = []g.Node{icon(it.Icon)}
		}
		items = append(items, h.Li(linkWithSlots(
			atoms.LinkProps{Label: it.Label, Href: it.Href},
			adornment,
		)))
	}
	nav := baseAttrs(p.ComponentProps)
	nav = append(nav, g.Attr("aria-label", "Breadcrumb"), h.Ol(items...))
	return h.Nav(nav...)
}

func breadcrumbItems(items []molecules.BreadcrumbItem, maxItems int) []molecules.BreadcrumbItem {
	if maxItems <= 0 || len(items) <= maxItems || maxItems < 2 {
		return items
	}

	tailCount := max(maxItems-1, 1)
	startOfTail := len(items) - tailCount
	visible := make([]molecules.BreadcrumbItem, 0, maxItems+1)
	visible = append(visible, items[0], molecules.BreadcrumbItem{Label: "…"})
	return append(visible, items[startOfTail:]...)
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
		return paginationPageURL(p.BaseURL, n)
	}
	pageLink := func(n int, label string, current bool, ariaLabel, marker string) g.Node {
		cl := clPageBtn.Merge(clPageIdle)
		if current {
			cl = clPageBtn.Merge(clPageCur)
		}
		nodes := []g.Node{
			h.Class(cl.Compile()),
			h.Href(pageHref(n)),
			g.Attr("data-page", itoa(n)),
		}
		if marker != "" {
			nodes = append(nodes, g.Attr(marker, ""))
		}
		if current {
			nodes = append(nodes, g.Attr("aria-current", "page"))
		}
		if ariaLabel != "" {
			nodes = append(nodes, g.Attr("aria-label", ariaLabel))
		}
		if !current {
			enhancement := p.HTMXProps
			if hasHTMXEnhancement(enhancement) {
				enhancement.Get = pageHref(n)
			}
			nodes = append(nodes, htmxAttrs(enhancement)...)
		}
		return h.A(append(nodes, g.Text(label))...)
	}
	disabledBoundary := func(label, marker, glyph string) g.Node {
		return h.Button(
			h.Class(clPageBtn.Merge(clPageIdle).Compile()),
			h.Type("button"),
			g.Attr(marker, ""),
			g.Attr("aria-label", label),
			h.Disabled(),
			g.Text(glyph),
		)
	}

	items := []g.Node{h.Class(clPagination.Compile())}
	if p.CurrentPage > 1 {
		items = append(items, pageLink(p.CurrentPage-1, "‹", false, "Previous page", "data-pagination-prev"))
	} else {
		items = append(items, disabledBoundary("Previous page", "data-pagination-prev", "‹"))
	}
	lo, hi := p.CurrentPage-siblings, p.CurrentPage+siblings
	if lo < 1 {
		lo = 1
	}
	if hi > p.TotalPages {
		hi = p.TotalPages
	}
	if lo > 1 {
		items = append(items, pageLink(1, "1", p.CurrentPage == 1, "Go to page 1", ""))
		if lo > 2 {
			items = append(items, h.Span(h.Class(clBreadcrumbSep.Compile()), g.Text("…")))
		}
	}
	for n := lo; n <= hi; n++ {
		ariaLabel := "Go to page " + itoa(n)
		if n == p.CurrentPage {
			ariaLabel = "Page " + itoa(n) + ", current page"
		}
		items = append(items, pageLink(n, itoa(n), n == p.CurrentPage, ariaLabel, ""))
	}
	if hi < p.TotalPages {
		if hi < p.TotalPages-1 {
			items = append(items, h.Span(h.Class(clBreadcrumbSep.Compile()), g.Text("…")))
		}
		items = append(items, pageLink(p.TotalPages, itoa(p.TotalPages), false, "Go to page "+itoa(p.TotalPages), ""))
	}
	if p.CurrentPage < p.TotalPages {
		items = append(items, pageLink(p.CurrentPage+1, "›", false, "Next page", "data-pagination-next"))
	} else {
		items = append(items, disabledBoundary("Next page", "data-pagination-next", "›"))
	}

	nav := baseAttrs(p.ComponentProps)
	nav = append(nav, g.Attr("data-component", "pagination"), g.Attr("aria-label", "Pagination"))
	nav = append(nav, items...)
	return h.Nav(nav...)
}

func paginationPageURL(baseURL string, page int) string {
	parsed, err := url.Parse(baseURL)
	if err == nil {
		query := parsed.Query()
		query.Set("page", itoa(page))
		parsed.RawQuery = query.Encode()
		return parsed.String()
	}

	separator := "?"
	if strings.Contains(baseURL, "?") {
		separator = "&"
	}
	return baseURL + separator + "page=" + itoa(page)
}

// paginationCursor renders bounded cursor links when continuation URLs are
// supplied. The older button shell remains available when no URLs are present.
func paginationCursor(p molecules.PaginationProps) g.Node {
	page := p.CurrentPage
	if page < 1 {
		page = 1
	}
	navigationLabel := p.NavigationLabel
	if navigationLabel == "" {
		navigationLabel = "Pagination"
	}
	previousLabel := p.PreviousLabel
	if previousLabel == "" {
		previousLabel = "Previous"
	}
	nextLabel := p.NextLabel
	if nextLabel == "" {
		nextLabel = "Next"
	}
	loadMoreLabel := p.LoadMoreLabel
	if loadMoreLabel == "" {
		loadMoreLabel = "Load more"
	}
	mode := p.CursorMode
	if mode == "" {
		mode = "previous-next"
	}

	if p.PreviousURL != "" || p.NextURL != "" {
		items := baseAttrs(p.ComponentProps)
		items = append(
			items,
			g.Attr("aria-label", navigationLabel),
			g.Attr("data-component", "cursor-pagination"),
			g.Attr("data-cursor-mode", mode),
			classes(clPagination.Compile(), p.Class),
		)
		if mode == "previous-next" && p.PreviousURL != "" {
			items = append(items, paginationCursorLink(
				p,
				p.PreviousURL,
				previousLabel,
				"backward",
				"prev",
			))
		}
		if p.NextURL != "" {
			label := nextLabel
			if mode == "load-more" {
				label = loadMoreLabel
			}
			items = append(items, paginationCursorLink(
				p,
				p.NextURL,
				label,
				"forward",
				"next",
			))
		}
		return h.Nav(items...)
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
	nav = append(nav, g.Attr("aria-label", navigationLabel), h.Class(clPagination.Compile()),
		h.Button(prev...),
		h.Span(h.Class(clPageLabel.Compile()), g.Attr("data-pk-pagination", "label"),
			g.Attr("aria-live", "polite"), g.Text("Page "+itoa(page))),
		h.Button(h.Class(clPageBtn.Merge(clPageIdle).Compile()), h.Type("button"),
			g.Attr("data-pk-pagination", "next"), g.Attr("aria-label", "Next page"),
			g.Text("Next →")),
	)
	return h.Nav(nav...)
}

func paginationCursorLink(
	p molecules.PaginationProps,
	url string,
	label string,
	direction string,
	rel string,
) g.Node {
	enhancement := p.HTMXProps
	if hasHTMXEnhancement(enhancement) {
		enhancement.Get = url
	}
	nodes := []g.Node{
		h.Href(url),
		h.Class(clPageBtn.Merge(clPageIdle).Compile() + " break-words"),
		g.Attr("data-cursor-direction", direction),
		g.Attr("rel", rel),
		g.Attr("aria-label", label),
	}
	if controlledID, ok := simpleIDTarget(p.Target); ok {
		nodes = append(nodes, g.Attr("aria-controls", controlledID))
	}
	nodes = append(nodes, htmxAttrs(enhancement)...)
	return h.A(append(nodes, g.Text(label))...)
}

func hasHTMXEnhancement(p contracts.HTMXProps) bool {
	return p.Get != "" || p.Post != "" || p.Put != "" || p.Patch != "" ||
		p.Delete != "" || p.Target != "" || p.Swap != "" || p.Trigger != "" ||
		p.Confirm != "" || p.Ext != "" || p.Indicator != "" ||
		p.DisabledElt != "" || p.Vals != "" || p.PushURL != "" ||
		p.Select != "" || p.Boost || p.Disable
}

func simpleIDTarget(target string) (string, bool) {
	if !strings.HasPrefix(target, "#") || len(target) < 2 {
		return "", false
	}
	id := strings.TrimPrefix(target, "#")
	for _, char := range id {
		if (char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '_' || char == '-' || char == ':' || char == '.' {
			continue
		}
		return "", false
	}
	return id, true
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
	name := p.Name
	if name == "" {
		name = "search"
	}
	debounce := p.DebounceMS
	if debounce <= 0 {
		debounce = 300
	}
	clearLabel := p.ClearLabel
	if clearLabel == "" {
		clearLabel = "Clear search"
	}
	shortcutKey := p.ShortcutKey
	if shortcutKey == "" {
		shortcutKey = "/"
	}

	enhancement := p.HTMXProps
	if p.Instant && p.SearchURL != "" {
		if enhancement.Get == "" {
			enhancement.Get = p.SearchURL
		}
		if enhancement.Trigger == "" {
			enhancement.Trigger = fmt.Sprintf("keyup changed delay:%dms, search", debounce)
		}
		if enhancement.Indicator == "" {
			enhancement.Indicator = "#" + name + "-indicator"
		}
	}

	input := []g.Node{
		h.Class(clSearchInput.Compile()), h.Type("search"),
		h.Name(name), h.ID(name),
		g.Attr("aria-label", label), h.AutoComplete("off"),
		g.Attr("data-searchbar-input", "true"),
		g.Attr("data-action", "input->searchbar#syncInput"),
	}
	if p.Value != "" {
		input = append(input, h.Value(p.Value))
	}
	if p.Disabled {
		input = append(input, h.Disabled())
	}
	if p.MinChars > 0 {
		input = append(input, g.Attr("data-search-min-chars", itoa(p.MinChars)))
	}
	input = append(input, htmxAttrs(enhancement)...)
	if p.Placeholder != "" {
		input = append(input, h.Placeholder(p.Placeholder))
	}
	if p.SearchURL != "" {
		input = append(input, g.Attr("data-pk-search-url", p.SearchURL))
	}

	children := []g.Node{
		Icon(atoms.IconProps{Name: "search", Size: "md", Tone: "neutral"}),
		h.Input(input...),
	}
	if p.Instant && p.SearchURL != "" {
		children = append(children, h.Span(
			h.ID(name+"-indicator"),
			h.Class("htmx-indicator"),
			g.Attr("data-searchbar-indicator", "true"),
			Spinner(atoms.SpinnerProps{Size: "sm", Label: "Searching"}),
		))
	}
	if p.ShowClear {
		children = append(children, buttonWithSlots(
			atoms.ButtonProps{
				Label: clearLabel, Variant: "ghost", Tone: "neutral", Size: "xs",
				IconOnly: true, AriaLabel: clearLabel,
				ComponentProps: contracts.ComponentProps{
					Disabled: p.Disabled,
					Hidden:   p.Value == "",
					Attrs: map[string]string{
						"data-action":                 "click->searchbar#clear",
						"data-searchbar-clear-button": "true",
					},
				},
			},
			[]g.Node{Icon(atoms.IconProps{Name: "x-mark", Size: "sm", Tone: "neutral"})},
			nil,
		))
	}
	if p.ShowShortcut {
		shortcut := []g.Node{
			h.Class(clKbd.Compile()),
			g.Attr("data-searchbar-shortcut", "true"),
			g.Text(shortcutKey),
		}
		if p.Disabled || p.Value != "" {
			shortcut = append(shortcut, g.Attr("hidden"))
		}
		children = append(children, h.Kbd(shortcut...))
	}

	root := baseAttrs(p.ComponentProps)
	root = append(root,
		classes(clSearchWrap.Compile(), p.Class),
		h.Role("search"),
		g.Attr("aria-label", label),
		g.Attr("data-component", "searchbar"),
		g.Attr("data-controller", "searchbar"),
		g.Attr("data-searchbar-show-clear-value", fmt.Sprintf("%t", p.ShowClear)),
		g.Attr("data-searchbar-show-shortcut-value", fmt.Sprintf("%t", p.ShowShortcut)),
		g.Attr("data-searchbar-trigger-search-on-clear-value", fmt.Sprintf("%t", p.Instant && p.SearchURL != "")),
	)
	return h.Label(append(root, children...)...)
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
