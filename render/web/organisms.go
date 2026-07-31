// Implements: REQ-011.
// Per: ADR-0031.
// Discipline: C-14.

package web

// organisms.go begins the organism layer: self-contained page sections
// composed from atoms and molecules through explicit named slots. An organism
// owns the order and layout of its parts without embedding another
// component's data contract inside its own props.

import (
	"strconv"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"github.com/septagon-oss/pk-ui/contracts/organisms"
)

// DataGridSlots is the typed Go composition surface corresponding exactly to
// the public design-delivery slot contract.
type DataGridSlots struct {
	Search     g.Node
	Filters    []g.Node
	Actions    []g.Node
	Table      g.Node
	Status     []g.Node
	Pagination g.Node
}

// WindowedCollectionSlots carries the already bounded item projection and
// its cursor navigation without putting domain rows into portable props.
type WindowedCollectionSlots struct {
	Items    []g.Node
	Controls g.Node
}

// DataGrid arranges an explicitly composed toolbar, table, status region, and
// pagination. It never reconstructs child components from opaque nested props.
func DataGrid(p organisms.DataGridProps, slots DataGridSlots) g.Node {
	var toolbar []g.Node
	if slots.Search != nil {
		toolbar = append(toolbar, slots.Search)
	}
	toolbar = append(toolbar, slots.Filters...)
	if len(slots.Actions) > 0 {
		toolbar = append(toolbar, h.Div(
			h.Class(clGridActions.Compile()),
			g.Group(slots.Actions),
		))
	}

	section := baseAttrs(p.ComponentProps, htmxAttrs(p.HTMXProps)...)
	section = append(section, classes(clGridSection.Compile(), p.Class))
	if len(toolbar) > 0 {
		bar := []g.Node{h.Class(clGridToolbar.Compile())}
		if slots.Search != nil {
			bar = append(bar, h.Role("search"))
		}
		section = append(section, h.Div(append(bar, toolbar...)...))
	}
	if slots.Table != nil {
		section = append(section, slots.Table)
	}
	section = append(section, slots.Status...)
	if slots.Pagination != nil {
		section = append(section, slots.Pagination)
	}
	return h.Section(section...)
}

// WindowedCollection renders one bounded collection window and its resilient
// states. The caller owns item projection; the shell refuses over-budget
// content before placing it in the DOM.
func WindowedCollection(
	p organisms.WindowedCollectionProps,
	slots WindowedCollectionSlots,
) g.Node {
	if p.MaxItems == 0 {
		p.MaxItems = 100
	}
	state := p.State
	if state == "" {
		state = "ready"
	}
	contractError := p.ContractError ||
		p.MaxItems < 1 ||
		p.MaxItems > 500 ||
		p.ItemCount < 0 ||
		p.ItemCount > p.MaxItems
	if contractError || (state != "ready" && state != "loading" && state != "error") {
		state = "error"
		contractError = true
		slots.Items = nil
	}

	collectionLabel := fallbackText(p.CollectionLabel, "Results")
	section := baseAttrs(p.ComponentProps)
	section = append(
		section,
		classes(clWindow.Compile(), p.Class),
		g.Attr("data-component", "windowed-collection"),
		g.Attr("data-windowed-state", state),
		g.Attr("data-windowed-item-count", strconv.Itoa(p.ItemCount)),
		g.Attr("data-windowed-max-items", strconv.Itoa(p.MaxItems)),
		g.Attr("aria-label", collectionLabel),
		g.Attr("aria-busy", strconv.FormatBool(state == "loading")),
	)
	if len(slots.Items) > 0 {
		section = append(section, h.Div(
			h.Class(clWindowItems.Compile()),
			g.Attr("data-windowed-items", ""),
			g.Group(slots.Items),
		))
	}

	switch {
	case state == "loading":
		section = append(section, windowedStatus(
			"data-windowed-loading",
			"status",
			clWindowLoading.Compile(),
			fallbackText(p.LoadingLabel, "Loading results"),
			"",
		))
	case state == "error":
		section = append(section, windowedError(p, contractError))
	case p.ItemCount == 0:
		section = append(section, windowedStatus(
			"data-windowed-empty",
			"status",
			clWindowEmpty.Compile(),
			fallbackText(p.EmptyTitle, "No results"),
			fallbackText(p.EmptyDescription, "There are no items to show."),
		))
	case p.NavigationUnavailable != "":
		section = append(section, h.Div(
			h.Class(clWindowNavigationError.Compile()),
			g.Attr("data-windowed-navigation-error", ""),
			g.Attr("role", "alert"),
			g.Text(p.NavigationUnavailable),
		))
	case slots.Controls != nil:
		section = append(section, h.Footer(
			h.Class(clWindowFooter.Compile()),
			g.Attr("data-windowed-controls", ""),
			slots.Controls,
		))
	}
	return h.Section(section...)
}

func windowedError(
	p organisms.WindowedCollectionProps,
	contractError bool,
) g.Node {
	nodes := []g.Node{
		h.Class(clWindowError.Compile()),
		g.Attr("data-windowed-error", ""),
		g.Attr("role", "alert"),
	}
	if contractError {
		nodes = append(nodes, g.Attr("data-windowed-error-kind", "contract"))
	}
	nodes = append(
		nodes,
		h.P(
			h.Class(clWindowTitle.Compile()),
			g.Text(fallbackText(p.ErrorTitle, "Unable to load results")),
		),
		h.P(
			h.Class(clWindowDescription.Compile()),
			g.Text(fallbackText(p.ErrorDescription, "Try again in a moment.")),
		),
	)
	if p.RetryURL != "" {
		enhancement := p.HTMXProps
		if hasHTMXEnhancement(enhancement) {
			enhancement.Get = p.RetryURL
		}
		retry := []g.Node{
			h.Class(clWindowRetry.Compile()),
			h.Href(p.RetryURL),
			g.Attr("data-windowed-retry", ""),
		}
		retry = append(retry, htmxAttrs(enhancement)...)
		retry = append(retry, g.Text(fallbackText(p.RetryLabel, "Try again")))
		nodes = append(nodes, h.A(retry...))
	}
	return h.Div(nodes...)
}

func windowedStatus(
	attribute string,
	role string,
	class string,
	title string,
	description string,
) g.Node {
	nodes := []g.Node{
		h.Class(class),
		g.Attr(attribute, ""),
		g.Attr("role", role),
	}
	if attribute == "data-windowed-loading" {
		nodes = append(
			nodes,
			g.Attr("aria-live", "polite"),
			g.Attr("aria-atomic", "true"),
		)
	}
	nodes = append(nodes, h.P(h.Class(clWindowTitle.Compile()), g.Text(title)))
	if description != "" {
		nodes = append(nodes, h.P(
			h.Class(clWindowDescription.Compile()),
			g.Text(description),
		))
	}
	return h.Div(nodes...)
}

func fallbackText(value string, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}
