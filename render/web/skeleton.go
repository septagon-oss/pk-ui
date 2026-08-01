package web

// skeleton.go renders the loading state of the library. A skeleton is the
// loading rendering of a component, not a separate component: primitives hold
// the geometry of pending content, twins mirror Table and Card exactly (same
// class lists), and DeferredSlot is the HTMX seam that swaps the finished
// fragment in.
//
// Accessibility contract: each placeholder is aria-hidden (it carries no
// information), while the enclosing DeferredSlot announces the pending state
// once via aria-busy. Motion is AnimatePulse — a decorative opacity fade whose
// pk-pulse keyframes tw/emission emits in Base().

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"github.com/septagon-oss/pk-ui/contracts/atoms"
	"github.com/septagon-oss/pk-ui/contracts/molecules"
	"github.com/septagon-oss/tw"
)

// Skeleton renders atoms.SkeletonProps. Shapes: "block" (full-width rectangle,
// the default), "text" (one or more prose lines, last line short), "circle"
// (avatar-sized disc). Unknown shape or size strings fall back to the
// defaults, matching the contracts' "data schema, not behavior" stance.
func Skeleton(p atoms.SkeletonProps) g.Node {
	size := p.Size
	if size == "" {
		size = "md"
	}
	switch p.Shape {
	case "circle":
		cl := clSkeleton.Merge(variantOr(clSkeletonCircleSize, size, "md"))
		return h.Div(append(baseAttrs(p.ComponentProps),
			classes(cl.Compile(), p.Class), g.Attr("aria-hidden", "true"))...)
	case "text":
		lines := max(p.Lines, 1)
		if lines == 1 {
			cl := skeletonLine(size, false)
			return h.Div(append(baseAttrs(p.ComponentProps),
				classes(cl.Compile(), p.Class), g.Attr("aria-hidden", "true"))...)
		}
		nodes := append(baseAttrs(p.ComponentProps),
			classes(clSkeletonText.Compile(), p.Class), g.Attr("aria-hidden", "true"))
		for i := 0; i < lines; i++ {
			nodes = append(nodes, h.Div(h.Class(skeletonLine(size, i == lines-1).Compile())))
		}
		return h.Div(nodes...)
	default:
		cl := clSkeleton.Merge(variantOr(clSkeletonBlockSize, size, "md"))
		return h.Div(append(baseAttrs(p.ComponentProps),
			classes(cl.Compile(), p.Class), g.Attr("aria-hidden", "true"))...)
	}
}

// skeletonLine is one pulsing prose line; the last line of a multi-line text
// placeholder stops short so the paragraph reads as prose, not a slab.
func skeletonLine(size string, last bool) tw.ClassList {
	width := clSkeletonLine
	if last {
		width = clSkeletonLineLast
	}
	return clSkeleton.Merge(width).Merge(variantOr(clSkeletonLineSize, size, "md"))
}

// DeferredSlot renders atoms.DeferredSlotProps around a placeholder. HTMX
// fetches Get as soon as the element lands (hx-trigger="load") and replaces
// the whole slot with the server-rendered fragment (hx-swap="outerHTML"), so
// the fragment route returns the finished component and no client code runs.
// Both defaults yield to explicit Trigger/Swap values for other moments
// (e.g. "revealed" for below-the-fold slots).
func DeferredSlot(p atoms.DeferredSlotProps, placeholder ...g.Node) g.Node {
	if p.Trigger == "" {
		p.Trigger = "load"
	}
	if p.Swap == "" {
		p.Swap = "outerHTML"
	}
	nodes := baseAttrs(p.ComponentProps, htmxAttrs(p.HTMXProps)...)
	if p.Class != "" {
		nodes = append(nodes, h.Class(p.Class))
	}
	nodes = append(nodes, g.Attr("aria-busy", "true"))
	nodes = append(nodes, placeholder...)
	return h.Div(nodes...)
}

// TableSkeleton renders molecules.TableSkeletonProps: the loading rendering of
// Table, built from Table's own class lists so the swap-in causes no layout
// shift. Defaults: 4 columns, 3 rows.
func TableSkeleton(p molecules.TableSkeletonProps) g.Node {
	cols := p.Columns
	if cols < 1 {
		cols = 4
	}
	rows := p.Rows
	if rows < 1 {
		rows = 3
	}
	tdClass := clTableTd
	if p.Compact {
		tdClass = clTableTdC
	}
	cellLine := skeletonLine("sm", false)

	headCells := make([]g.Node, 0, cols)
	for i := 0; i < cols; i++ {
		headCells = append(headCells, h.Th(
			h.Class(clTableTh.Compile()), g.Attr("scope", "col"),
			h.Div(h.Class(cellLine.Compile())),
		))
	}

	bodyRows := make([]g.Node, 0, rows)
	for i := 0; i < rows; i++ {
		cells := []g.Node{h.Class(clTableRow.Compile())}
		for j := 0; j < cols; j++ {
			cells = append(cells, h.Td(h.Class(tdClass.Compile()),
				h.Div(h.Class(cellLine.Compile()))))
		}
		bodyRows = append(bodyRows, h.Tr(cells...))
	}

	nodes := baseAttrs(p.ComponentProps)
	nodes = append(nodes, classes(clTableWrap.Compile(), p.Class),
		g.Attr("aria-hidden", "true"),
		h.Table(
			h.Class(clTable.Compile()),
			h.THead(h.Class(clTableHead.Compile()), h.Tr(headCells...)),
			h.TBody(bodyRows...),
		),
	)
	return h.Div(nodes...)
}

// CardSkeleton renders molecules.CardSkeletonProps: the loading rendering of
// Card — a short title line and a text-shaped body inside the real card
// frame. Default body: 3 lines.
func CardSkeleton(p molecules.CardSkeletonProps) g.Node {
	lines := p.Lines
	if lines < 1 {
		lines = 3
	}
	nodes := baseAttrs(p.ComponentProps)
	nodes = append(nodes, classes(clCard.Compile(), p.Class),
		g.Attr("aria-hidden", "true"),
		h.Div(h.Class(skeletonLine("lg", true).Compile())),
		Skeleton(atoms.SkeletonProps{Shape: "text", Lines: lines}),
	)
	return h.Div(nodes...)
}
