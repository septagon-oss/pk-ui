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
	"strings"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"github.com/septagon-oss/pk-ui/contracts/atoms"
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

// DashboardWidgetSlots carries rich widget content without serializing
// arbitrary child markup into the portable property contract.
type DashboardWidgetSlots struct {
	Content []g.Node
	Footer  []g.Node
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
	section = append(section,
		classes(clGridSection.Compile(), p.Class),
		g.Attr("data-component", "data-grid"),
	)
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

// DashboardWidget renders the canonical dashboard metric/content card.
func DashboardWidget(p organisms.DashboardWidgetProps, slots DashboardWidgetSlots) g.Node {
	span := strings.TrimSpace(p.Span)
	if _, ok := clDashboardWidgetSpan[span]; !ok {
		span = "1"
	}
	trend := strings.ToLower(strings.TrimSpace(p.Trend))
	if trend != "up" && trend != "down" && trend != "flat" {
		trend = "flat"
	}

	htmx := p.HTMXProps
	if htmx.Get == "" {
		htmx.Get = p.RefreshURL
	}
	if htmx.Get != "" {
		if htmx.Swap == "" {
			htmx.Swap = "innerHTML"
		}
		if htmx.Trigger == "" {
			var triggers []string
			if p.RefreshSeconds > 0 && p.RefreshSeconds <= 86400 {
				triggers = append(triggers, "every "+strconv.Itoa(p.RefreshSeconds)+"s")
			} else {
				triggers = append(triggers, "load")
			}
			if event := strings.TrimSpace(p.RefreshOn); event != "" {
				triggers = append(triggers, event+" from:body")
			}
			htmx.Trigger = strings.Join(triggers, ", ")
		}
	}

	root := baseAttrs(p.ComponentProps, htmxAttrs(htmx)...)
	root = append(root,
		classes(clDashboardWidget.Merge(clDashboardWidgetSpan[span]).Compile(), p.Class),
		g.Attr("data-component", "dashboard-widget"),
		g.Attr("data-widget-type", fallbackText(p.Type, "stat")),
	)

	var title g.Node = g.Text(p.Title)
	if p.DetailURL != "" {
		title = h.A(h.Href(p.DetailURL), g.Text(p.Title))
	}
	header := []g.Node{h.Div(
		h.Class(clDashboardWidgetCopy.Compile()),
		h.P(h.Class(clDashboardWidgetTitle.Compile()), title),
		g.If(p.Subtitle != "", h.P(
			h.Class(clDashboardWidgetSubtitle.Compile()),
			g.Text(p.Subtitle),
		)),
	)}
	if p.Icon != "" {
		header = append([]g.Node{h.Div(
			h.Class(clDashboardWidgetIcon.Compile()),
			Icon(atoms.IconProps{Name: p.Icon, Size: "md", Tone: "brand"}),
		)}, header...)
	}
	root = append(root, h.Header(
		h.Class(clDashboardWidgetHeader.Compile()),
		g.Group(header),
	))

	if len(slots.Content) > 0 {
		root = append(root, h.Div(h.Class(clDashboardWidgetBody.Compile()), g.Group(slots.Content)))
	} else {
		body := []g.Node{h.P(
			h.Class(clDashboardWidgetValue.Compile()),
			g.Attr("data-dashboard-widget-value", ""),
			g.Text(p.Value),
		)}
		if p.Change != "" {
			arrow := map[string]string{"up": "↑ ", "down": "↓ ", "flat": "→ "}[trend]
			body = append(body, h.Div(
				h.Class(clDashboardWidgetTrend.Compile()),
				g.Attr("data-dashboard-widget-trend", trend),
				h.Span(
					h.Class(clDashboardWidgetChange.Merge(clDashboardWidgetTrendTone[trend]).Compile()),
					g.Text(arrow+p.Change),
				),
				g.If(p.PreviousValue != "", h.Span(
					h.Class(clDashboardWidgetPrevious.Compile()),
					g.Text("from "+p.PreviousValue),
				)),
			))
		}
		root = append(root, h.Div(h.Class(clDashboardWidgetBody.Compile()), g.Group(body)))
	}
	if len(slots.Footer) > 0 {
		root = append(root, h.Footer(
			h.Class(clDashboardWidgetFooter.Compile()),
			g.Group(slots.Footer),
		))
	}
	return h.Article(root...)
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
	layout := p.Layout
	if layout != "grid" {
		layout = "list"
	}
	contractError := p.ContractError ||
		p.MaxItems < 1 ||
		p.MaxItems > 500 ||
		p.ItemCount < 0 ||
		p.ItemCount > p.MaxItems
	if contractError ||
		(state != "ready" && state != "loading" && state != "error" && state != "offline") {
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
		g.Attr("data-windowed-layout", layout),
		g.Attr("data-windowed-state", state),
		g.Attr("data-windowed-item-count", strconv.Itoa(p.ItemCount)),
		g.Attr("data-windowed-max-items", strconv.Itoa(p.MaxItems)),
		g.Attr("aria-label", collectionLabel),
		g.Attr("aria-busy", strconv.FormatBool(state == "loading")),
	)
	if p.Title != "" || p.Description != "" {
		section = append(section, h.Header(
			h.Class(clWindowHeader.Compile()),
			h.Span(
				h.Class(clWindowAccent.Compile()),
				g.Attr("aria-hidden", "true"),
			),
			g.If(p.Title != "", h.H2(
				h.Class(clWindowHeading.Compile()),
				g.Text(p.Title),
			)),
			g.If(p.Description != "", h.P(
				h.Class(clWindowIntro.Compile()),
				g.Text(p.Description),
			)),
		))
	}
	if len(slots.Items) > 0 {
		itemsClass := clWindowItems
		if layout == "grid" {
			itemsClass = clWindowItemsGrid
		}
		section = append(section, h.Div(
			h.Class(itemsClass.Compile()),
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
	case state == "offline":
		section = append(section, windowedRecoverableStatus(
			p,
			"offline",
			clWindowOffline.Compile(),
			fallbackText(p.OfflineTitle, "You are offline"),
			fallbackText(p.OfflineDescription, "Reconnect to load this window, then try again."),
			false,
		))
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
	return windowedRecoverableStatus(
		p,
		"error",
		clWindowError.Compile(),
		fallbackText(p.ErrorTitle, "Unable to load results"),
		fallbackText(p.ErrorDescription, "Try again in a moment."),
		contractError,
	)
}

func windowedRecoverableStatus(
	p organisms.WindowedCollectionProps,
	kind string,
	class string,
	title string,
	description string,
	contractError bool,
) g.Node {
	nodes := []g.Node{
		h.Class(class),
		g.Attr("data-windowed-"+kind, ""),
		g.Attr("role", "alert"),
	}
	if contractError {
		nodes = append(nodes, g.Attr("data-windowed-error-kind", "contract"))
	}
	nodes = append(
		nodes,
		h.P(
			h.Class(clWindowTitle.Compile()),
			g.Text(title),
		),
		h.P(
			h.Class(clWindowDescription.Compile()),
			g.Text(description),
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
