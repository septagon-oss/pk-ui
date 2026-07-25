// Implements: REQ-011.
// Per: ADR-0031.
// Discipline: C-14.

package web

// atoms.go renders the atom contracts. Each renderer takes its Props struct
// and returns a gomponents Node; unknown variant or size strings fall back to
// the documented defaults rather than failing, matching the contracts'
// "data schema, not behavior" stance.

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"github.com/septagon-oss/pk-ui/contracts/atoms"
	"github.com/septagon-oss/tw"
)

func variantOr(m map[string]tw.ClassList, key, fallback string) tw.ClassList {
	if cl, ok := m[key]; ok {
		return cl
	}
	return m[fallback]
}

// Button renders atoms.ButtonProps as a <button>.
func Button(p atoms.ButtonProps) g.Node {
	cl := clButtonBase.
		Merge(variantOr(clButtonVariant, p.Variant, "primary")).
		Merge(variantOr(clButtonSize, p.Size, "medium"))
	if p.FullWidth {
		cl = cl.Merge(clButtonFull)
	}
	typ := p.Type
	if typ == "" {
		typ = "button"
	}
	var children []g.Node
	children = append(children, baseAttrs(p.ComponentProps, htmxAttrs(p.HTMXProps)...)...)
	children = append(children, classes(cl.Compile(), p.Class), h.Type(typ))
	if p.Loading {
		children = append(children,
			g.Attr("aria-busy", "true"),
			Spinner(atoms.SpinnerProps{Size: "small", Label: ""}),
		)
	} else if p.Icon != "" && p.IconPosition != "right" {
		children = append(children, icon(p.Icon))
	}
	children = append(children, g.Text(p.Text))
	if !p.Loading && p.Icon != "" && p.IconPosition == "right" {
		children = append(children, icon(p.Icon))
	}
	return h.Button(children...)
}

// Badge renders atoms.BadgeProps.
func Badge(p atoms.BadgeProps) g.Node {
	cl := clBadgeBase.Merge(variantOr(clBadgeVariant, p.Variant, "default"))
	var children []g.Node
	children = append(children, baseAttrs(p.ComponentProps)...)
	children = append(children, classes(cl.Compile(), p.Class))
	if p.Dot {
		children = append(children, h.Span(h.Class(clBadgeDot.Compile()), g.Attr("aria-hidden", "true")))
	}
	if p.Icon != "" {
		children = append(children, icon(p.Icon))
	}
	children = append(children, g.Text(p.Text))
	return h.Span(children...)
}

// Alert renders atoms.AlertProps with role="alert" for danger/warning and
// role="status" otherwise, so severity maps to interruption behavior.
func Alert(p atoms.AlertProps) g.Node {
	cl := clAlertBase.Merge(variantOr(clAlertVariant, p.Variant, "info"))
	role := "status"
	if p.Variant == "error" || p.Variant == "danger" || p.Variant == "warning" {
		role = "alert"
	}
	body := []g.Node{h.Class(clAlertBody.Compile())}
	if p.Title != "" {
		body = append(body, h.P(h.Class(clAlertTitle.Compile()), g.Text(p.Title)))
	}
	body = append(body, h.P(h.Class(clAlertMessage.Compile()), g.Text(p.Message)))

	var children []g.Node
	children = append(children, baseAttrs(p.ComponentProps)...)
	children = append(children, classes(cl.Compile(), p.Class), h.Role(role))
	if p.Icon != "" {
		children = append(children, icon(p.Icon))
	}
	children = append(children, h.Div(body...))
	return h.Div(children...)
}

// Input renders atoms.InputProps as a labelled form field. When Error is set
// the input carries aria-invalid and is described by the error element.
func Input(p atoms.InputProps) g.Node {
	id := p.ID
	if id == "" && p.Name != "" {
		id = "pk-input-" + p.Name
	}
	typ := p.Type
	if typ == "" {
		typ = "text"
	}
	cl := clInput.Merge(clInputNormal)
	if p.Error != "" {
		cl = clInput.Merge(clInputError)
	}

	input := []g.Node{
		classes(cl.Compile(), p.Class),
		h.ID(id), h.Name(p.Name), h.Type(typ),
	}
	input = append(input, attrPairs(p.Attrs)...)
	input = append(input, htmxAttrs(p.HTMXProps)...)
	if p.Value != "" {
		input = append(input, h.Value(p.Value))
	}
	if p.Placeholder != "" {
		input = append(input, h.Placeholder(p.Placeholder))
	}
	if p.Required {
		input = append(input, h.Required())
	}
	if p.ReadOnly {
		input = append(input, h.ReadOnly())
	}
	if p.AutoFocus {
		input = append(input, h.AutoFocus())
	}
	if p.Disabled {
		input = append(input, h.Disabled())
	}
	if p.Min != "" {
		input = append(input, h.Min(p.Min))
	}
	if p.Max != "" {
		input = append(input, h.Max(p.Max))
	}
	if p.Pattern != "" {
		input = append(input, h.Pattern(p.Pattern))
	}
	if p.Error != "" {
		input = append(input,
			g.Attr("aria-invalid", "true"),
			g.Attr("aria-describedby", id+"-error"),
		)
	} else if p.HelpText != "" {
		input = append(input, g.Attr("aria-describedby", id+"-help"))
	}

	field := []g.Node{h.Class(clFieldWrap.Compile())}
	if p.Label != "" {
		field = append(field, Label(atoms.LabelProps{Text: p.Label, For: id, Required: p.Required}))
	}
	field = append(field, h.Input(input...))
	if p.Error != "" {
		field = append(field, h.P(h.ID(id+"-error"), h.Class(clFieldErr.Compile()), g.Text(p.Error)))
	} else if p.HelpText != "" {
		field = append(field, h.P(h.ID(id+"-help"), h.Class(clHelp.Compile()), g.Text(p.HelpText)))
	}
	return h.Div(field...)
}

// Select renders atoms.SelectProps as a labelled native <select>, styled as
// an input-family control. A Placeholder renders as an empty leading option
// so an untouched control has no accidental value; when the field is
// Required the placeholder is how "nothing chosen yet" stays expressible.
func Select(p atoms.SelectProps) g.Node {
	id := p.ID
	if id == "" && p.Name != "" {
		id = "pk-select-" + p.Name
	}
	cl := clInput.Merge(clInputNormal)
	if p.Error != "" {
		cl = clInput.Merge(clInputError)
	}

	var options []g.Node
	if p.Placeholder != "" || !p.Required {
		label := p.Placeholder
		if label == "" {
			label = "Choose…"
		}
		options = append(options, h.Option(h.Value(""), g.Text(label)))
	}
	for _, o := range p.Options {
		opt := []g.Node{h.Value(o.Value), g.Text(o.Label)}
		if o.Value == p.Value && p.Value != "" {
			opt = append(opt, h.Selected())
		}
		if o.Disabled {
			opt = append(opt, h.Disabled())
		}
		options = append(options, h.Option(opt...))
	}

	sel := []g.Node{
		classes(cl.Compile(), p.Class),
		h.ID(id), h.Name(p.Name),
	}
	sel = append(sel, attrPairs(p.Attrs)...)
	sel = append(sel, htmxAttrs(p.HTMXProps)...)
	if p.Required {
		sel = append(sel, h.Required())
	}
	if p.Disabled {
		sel = append(sel, h.Disabled())
	}
	if p.Error != "" {
		sel = append(sel, g.Attr("aria-invalid", "true"), g.Attr("aria-describedby", id+"-error"))
	} else if p.HelpText != "" {
		sel = append(sel, g.Attr("aria-describedby", id+"-help"))
	}
	sel = append(sel, options...)

	field := []g.Node{h.Class(clFieldWrap.Compile())}
	if p.Label != "" {
		field = append(field, Label(atoms.LabelProps{Text: p.Label, For: id, Required: p.Required}))
	}
	field = append(field, h.Select(sel...))
	if p.Error != "" {
		field = append(field, h.P(h.ID(id+"-error"), h.Class(clFieldErr.Compile()), g.Text(p.Error)))
	} else if p.HelpText != "" {
		field = append(field, h.P(h.ID(id+"-help"), h.Class(clHelp.Compile()), g.Text(p.HelpText)))
	}
	return h.Div(field...)
}

// Textarea renders atoms.TextareaProps as a labelled multi-line field.
func Textarea(p atoms.TextareaProps) g.Node {
	id := p.ID
	if id == "" && p.Name != "" {
		id = "pk-textarea-" + p.Name
	}
	cl := clInput.Merge(clInputNormal)
	if p.ErrorMessage != "" {
		cl = clInput.Merge(clInputError)
	}
	rows := p.Rows
	if rows == 0 {
		rows = 4
	}
	area := []g.Node{
		classes(cl.Compile(), p.Class),
		h.ID(id), h.Name(p.Name), h.Rows(itoa(rows)),
	}
	area = append(area, attrPairs(p.Attrs)...)
	area = append(area, htmxAttrs(p.HTMXProps)...)
	if p.Placeholder != "" {
		area = append(area, h.Placeholder(p.Placeholder))
	}
	if p.Required {
		area = append(area, h.Required())
	}
	if p.ReadOnly {
		area = append(area, h.ReadOnly())
	}
	if p.MaxLength > 0 {
		area = append(area, h.MaxLength(itoa(p.MaxLength)))
	}
	if p.ErrorMessage != "" {
		area = append(area, g.Attr("aria-invalid", "true"), g.Attr("aria-describedby", id+"-error"))
	}
	area = append(area, g.Text(p.Value))

	field := []g.Node{h.Class(clFieldWrap.Compile())}
	if p.Label != "" {
		field = append(field, Label(atoms.LabelProps{Text: p.Label, For: id, Required: p.Required}))
	}
	field = append(field, h.Textarea(area...))
	if p.ErrorMessage != "" {
		field = append(field, h.P(h.ID(id+"-error"), h.Class(clFieldErr.Compile()), g.Text(p.ErrorMessage)))
	} else if p.HelperText != "" {
		field = append(field, h.P(h.Class(clHelp.Compile()), g.Text(p.HelperText)))
	}
	return h.Div(field...)
}

// Checkbox renders atoms.CheckboxProps as a labelled checkbox row.
func Checkbox(p atoms.CheckboxProps) g.Node {
	id := p.ID
	if id == "" && p.Name != "" {
		id = "pk-checkbox-" + p.Name
	}
	box := []g.Node{
		h.Class(clCheckbox.Compile()), h.ID(id), h.Name(p.Name), h.Type("checkbox"),
	}
	box = append(box, attrPairs(p.Attrs)...)
	if p.Value != "" {
		box = append(box, h.Value(p.Value))
	}
	if p.Checked {
		box = append(box, h.Checked())
	}
	if p.Required {
		box = append(box, h.Required())
	}
	if p.Disabled {
		box = append(box, h.Disabled())
	}
	if p.HelpText != "" {
		box = append(box, g.Attr("aria-describedby", id+"-help"))
	}
	row := []g.Node{
		h.Class(clCheckRow.Compile()),
		h.Input(box...),
		h.Label(h.For(id), h.Class(clLabel.Compile()), g.Text(p.Label)),
	}
	if p.HelpText == "" {
		return h.Div(row...)
	}
	return h.Div(
		h.Class(clFieldWrap.Compile()),
		h.Div(row...),
		h.P(h.ID(id+"-help"), h.Class(clHelp.Compile()), g.Text(p.HelpText)),
	)
}

// Label renders atoms.LabelProps; required fields carry a visible marker the
// screen reader skips (the input's required attribute carries the semantics).
func Label(p atoms.LabelProps) g.Node {
	children := []g.Node{h.For(p.For), classes(clLabel.Compile(), p.Class), g.Text(p.Text)}
	if p.Required {
		children = append(children, h.Span(
			h.Class(clRequired.Compile()), g.Attr("aria-hidden", "true"), g.Text(" *"),
		))
	}
	return h.Label(children...)
}

// Text renders atoms.TextProps as a paragraph.
func Text(p atoms.TextProps) g.Node {
	cl := tw.New()
	if sz, ok := clTextSize[p.Size]; ok {
		cl = cl.FontSize(sz)
	}
	if w, ok := clTextWeight[p.Weight]; ok {
		cl = cl.FontWeight(w)
	}
	if c, ok := clTextColor[p.Color]; ok {
		cl = cl.TextColor(c)
	}
	if p.Truncate {
		cl = cl.Merge(clTruncate)
	}
	var children []g.Node
	children = append(children, baseAttrs(p.ComponentProps)...)
	children = append(children, classes(cl.Compile(), p.Class), g.Text(p.Content))
	return h.P(children...)
}

// Heading renders atoms.HeadingProps at the given level (clamped 1..6) in the
// design system's display face.
func Heading(p atoms.HeadingProps) g.Node {
	level := p.Level
	if level < 1 || level > 6 {
		level = 2
	}
	cl := clHeadingBase.Merge(clHeadingLevel[level])
	if p.Truncate {
		cl = cl.Merge(clTruncate)
	}
	var children []g.Node
	children = append(children, baseAttrs(p.ComponentProps)...)
	if p.Anchor != "" {
		children = append(children, h.ID(p.Anchor))
	}
	children = append(children, classes(cl.Compile(), p.Class), g.Text(p.Text))
	switch level {
	case 1:
		return h.H1(children...)
	case 3:
		return h.H3(children...)
	case 4:
		return h.H4(children...)
	case 5:
		return h.H5(children...)
	case 6:
		return h.H6(children...)
	default:
		return h.H2(children...)
	}
}

// Divider renders atoms.DividerProps as an <hr> (horizontal) or a separator
// span (vertical).
func Divider(p atoms.DividerProps) g.Node {
	if p.Orientation == "vertical" {
		return h.Span(
			classes(clDividerV.Compile(), p.Class),
			h.Role("separator"), g.Attr("aria-orientation", "vertical"),
		)
	}
	return h.Hr(classes(clDividerH.Compile(), p.Class))
}

// Spinner renders atoms.SpinnerProps; the label is announced, the rotation is
// decoration.
func Spinner(p atoms.SpinnerProps) g.Node {
	cl := clSpinner.Merge(variantOr(clSpinnerSize, p.Size, "medium"))
	label := p.Label
	if label == "" {
		label = "Loading"
	}
	return h.Span(
		classes(cl.Compile(), p.Class),
		h.Role("status"), g.Attr("aria-label", label),
	)
}

// EmptyState renders atoms.EmptyStateProps.
func EmptyState(p atoms.EmptyStateProps) g.Node {
	cl := clEmpty.Merge(clEmptyPad)
	if p.Compact {
		cl = clEmpty.Merge(clEmptyCompact)
	}
	if p.Bordered {
		cl = cl.Merge(clEmptyBordered)
	}
	var children []g.Node
	children = append(children, baseAttrs(p.ComponentProps)...)
	children = append(children, classes(cl.Compile(), p.Class))
	if p.Icon != "" {
		children = append(children, icon(p.Icon))
	}
	children = append(children, h.P(h.Class(clEmptyTitle.Compile()), g.Text(p.Title)))
	if p.Description != "" {
		children = append(children, h.P(h.Class(clEmptyDesc.Compile()), g.Text(p.Description)))
	}
	for _, a := range p.Actions {
		variant := a.Variant
		if variant == "" {
			variant = "primary"
		}
		cl := clButtonBase.
			Merge(variantOr(clButtonVariant, variant, "primary")).
			Merge(variantOr(clButtonSize, "medium", "medium"))
		action := []g.Node{h.Class(cl.Compile()), h.Href(a.Href)}
		if a.Icon != "" {
			action = append(action, icon(a.Icon))
		}
		action = append(action, g.Text(a.Label))
		children = append(children, h.A(action...))
	}
	return h.Div(children...)
}

// Kbd renders atoms.KbdProps as a sequence of <kbd> elements.
func Kbd(p atoms.KbdProps) g.Node {
	var children []g.Node
	children = append(children, baseAttrs(p.ComponentProps)...)
	children = append(children, h.Class(clCheckRow.Compile()))
	for _, k := range p.Keys {
		children = append(children, h.Kbd(classes(clKbd.Compile(), p.Class), g.Text(k)))
	}
	return h.Span(children...)
}

// Link renders atoms.LinkProps; external links open safely.
func Link(p atoms.LinkProps) g.Node {
	var children []g.Node
	children = append(children, baseAttrs(p.ComponentProps, htmxAttrs(p.HTMXProps)...)...)
	children = append(children, classes(clLink.Compile(), p.Class), h.Href(p.Href))
	target, rel := p.Target, p.Rel
	if p.External {
		if target == "" {
			target = "_blank"
		}
		if rel == "" {
			rel = "noopener noreferrer"
		}
	}
	if target != "" {
		children = append(children, h.Target(target))
	}
	if rel != "" {
		children = append(children, h.Rel(rel))
	}
	children = append(children, g.Text(p.Text))
	if p.External {
		children = append(children, h.Span(g.Attr("aria-hidden", "true"), g.Text(" ↗")))
	}
	return h.A(children...)
}

// Tag renders atoms.TagProps; removable tags post to OnRemoveURL.
func Tag(p atoms.TagProps) g.Node {
	cl := clTagBase.Merge(clTagIdle)
	if p.Selected {
		cl = clTagBase.Merge(clTagSelected)
	}
	var children []g.Node
	children = append(children, baseAttrs(p.ComponentProps)...)
	children = append(children, classes(cl.Compile(), p.Class))
	if p.Icon != "" {
		children = append(children, icon(p.Icon))
	}
	children = append(children, g.Text(p.Text))
	if p.Removable && p.OnRemoveURL != "" {
		children = append(children, h.Button(
			h.Type("button"), h.Class(clLink.Compile()),
			g.Attr("hx-delete", p.OnRemoveURL),
			g.Attr("aria-label", "Remove "+p.Text),
			g.Text("×"),
		))
	}
	return h.Span(children...)
}
