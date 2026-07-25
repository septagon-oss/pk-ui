// Implements: REQ-011.
// Per: ADR-0031.
// Discipline: C-14.

package web

// layouts.go renders the layout contracts: structural containers whose whole
// job is arranging children. They accept children as trailing gomponents
// nodes because layout without content is meaningless.

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"github.com/septagon-oss/pk-ui/contracts/layouts"
	"github.com/septagon-oss/tw"
)

func gapOr(gap string, fallback tw.Spacing) tw.Spacing {
	if s, ok := clGapScale[gap]; ok {
		return s
	}
	return fallback
}

var alignItems = map[string]tw.Items{
	"start": tw.ItemsStart, "center": tw.ItemsCenter,
	"end": tw.ItemsEnd, "stretch": tw.ItemsStretch,
}

var justifyContent = map[string]tw.Justify{
	"start": tw.JustifyStart, "center": tw.JustifyCenter,
	"end": tw.JustifyEnd, "between": tw.JustifyBetween,
}

// Stack renders layouts.StackProps: a vertical flex column.
func Stack(p layouts.StackProps, children ...g.Node) g.Node {
	cl := clStack.Gap(gapOr(p.Gap, tw.S4))
	if a, ok := alignItems[p.Align]; ok {
		cl = cl.Items(a)
	}
	nodes := baseAttrs(p.ComponentProps)
	nodes = append(nodes, classes(cl.Compile(), p.Class))
	nodes = append(nodes, children...)
	return h.Div(nodes...)
}

// Flex renders layouts.FlexProps.
func Flex(p layouts.FlexProps, children ...g.Node) g.Node {
	cl := clFlex.Gap(gapOr(p.Gap, tw.S4))
	if p.Direction == "col" || p.Direction == "column" {
		cl = cl.FlexDir(tw.FlexCol)
	} else {
		cl = cl.FlexDir(tw.FlexRow)
	}
	if p.Wrap {
		cl = cl.FlexWrap()
	}
	if a, ok := alignItems[p.Align]; ok {
		cl = cl.Items(a)
	}
	if j, ok := justifyContent[p.Justify]; ok {
		cl = cl.Justify(j)
	}
	nodes := baseAttrs(p.ComponentProps)
	nodes = append(nodes, classes(cl.Compile(), p.Class))
	nodes = append(nodes, children...)
	return h.Div(nodes...)
}

// Grid renders layouts.GridProps; Columns outside 1..12 fall back to 1.
func Grid(p layouts.GridProps, children ...g.Node) g.Node {
	cols := 1
	switch p.Columns {
	case "2":
		cols = 2
	case "3":
		cols = 3
	case "4":
		cols = 4
	case "6":
		cols = 6
	case "12":
		cols = 12
	}
	cl := clGrid.GridCols(cols).Gap(gapOr(p.Gap, tw.S4))
	nodes := baseAttrs(p.ComponentProps)
	nodes = append(nodes, classes(cl.Compile(), p.Class))
	nodes = append(nodes, children...)
	return h.Div(nodes...)
}

var containerWidths = map[string]tw.MaxWidth{
	"sm": tw.MaxWSM, "md": tw.MaxWMD, "lg": tw.MaxWLG, "xl": tw.MaxWXL,
	"2xl": tw.MaxW2XL, "4xl": tw.MaxW4XL, "7xl": tw.MaxW7XL, "full": tw.MaxWFull,
}

// Container renders layouts.ContainerProps: a centered max-width column.
func Container(p layouts.ContainerProps, children ...g.Node) g.Node {
	w := tw.MaxW7XL
	if mw, ok := containerWidths[p.MaxWidth]; ok {
		w = mw
	}
	cl := clContainer.MaxWScaled(w)
	nodes := baseAttrs(p.ComponentProps)
	nodes = append(nodes, classes(cl.Compile(), p.Class))
	nodes = append(nodes, children...)
	return h.Div(nodes...)
}
