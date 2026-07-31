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
	ossicon "github.com/septagon-oss/pk-ui/icon"
	"github.com/septagon-oss/tw"
)

func variantOr(m map[string]tw.ClassList, key, fallback string) tw.ClassList {
	if cl, ok := m[key]; ok {
		return cl
	}
	return m[fallback]
}

// Icon renders the OSS provider's vector directly into the document. Product
// and client providers extend the glyph vocabulary behind icon.Resolve while
// this atom retains sizing, semantic tone, and accessibility ownership.
func Icon(p atoms.IconProps) g.Node {
	size := p.Size
	if size == "" {
		size = "md"
	}
	tone := p.Tone
	if tone == "" {
		tone = "neutral"
	}
	cl := clIcon.
		Merge(variantOr(clIconSize, size, "md")).
		Merge(variantOr(clIconTone, tone, "neutral"))
	glyph, known := ossicon.Resolve(p.Name, p.Weight)
	nodes := baseAttrs(p.ComponentProps)
	nodes = append(
		nodes,
		classes(cl.Compile(), p.Class),
		g.Attr("xmlns", "http://www.w3.org/2000/svg"),
		g.Attr("viewBox", glyph.ViewBox),
		g.Attr("fill", "currentColor"),
		g.Attr("focusable", "false"),
		g.Attr("data-pk-icon", p.Name),
	)
	if p.Weight != "" {
		nodes = append(nodes, g.Attr("data-pk-icon-weight", p.Weight))
	}
	if !known {
		nodes = append(nodes, g.Attr("data-pk-icon-fallback", "true"))
	}
	if p.AriaLabel == "" {
		nodes = append(nodes, g.Attr("aria-hidden", "true"))
	} else {
		nodes = append(nodes, g.Attr("role", "img"), g.Attr("aria-label", p.AriaLabel))
	}
	nodes = append(nodes, g.Raw(glyph.Body))
	return g.El("svg", nodes...)
}

// Button renders atoms.ButtonProps as a <button>.
func Button(p atoms.ButtonProps) g.Node {
	return buttonWithSlots(p, nil, nil)
}

func buttonWithSlots(
	p atoms.ButtonProps,
	iconStart []g.Node,
	iconEnd []g.Node,
) g.Node {
	variant := p.Variant
	if variant == "" {
		variant = "primary"
	}
	appearance := variantOr(clButtonVariant, variant, "primary")
	if p.Tone != "" && p.Tone != "neutral" {
		appearance = variantOr(clButtonTone, p.Tone, "brand")
	}
	cl := clButtonBase.
		Merge(appearance).
		Merge(variantOr(clButtonSize, p.Size, "md"))
	if p.FullWidth {
		cl = cl.Merge(clButtonFull)
	}
	if p.IconOnly {
		cl = cl.Merge(clButtonIconOnly)
	}
	typ := p.Type
	if typ == "" {
		typ = "button"
	}
	var children []g.Node
	children = append(children, baseAttrs(p.ComponentProps, htmxAttrs(p.HTMXProps)...)...)
	children = append(children, classes(cl.Compile(), p.Class), h.Type(typ))
	if label := p.AriaLabel; label != "" {
		children = append(children, g.Attr("aria-label", label))
	} else if p.IconOnly && p.Label != "" {
		children = append(children, g.Attr("aria-label", p.Label))
	}
	if p.Loading {
		children = append(children,
			g.Attr("aria-busy", "true"),
			Spinner(atoms.SpinnerProps{Size: "sm", Label: ""}),
		)
	} else if len(iconStart) > 0 {
		children = append(children, iconStart...)
	}
	if !p.IconOnly {
		children = append(children, g.Text(p.Label))
	}
	if !p.Loading {
		if len(iconEnd) > 0 {
			children = append(children, iconEnd...)
		}
	}
	return h.Button(children...)
}

// Badge renders atoms.BadgeProps.
func Badge(p atoms.BadgeProps) g.Node {
	return badgeWithSlots(p, nil, nil)
}

func badgeWithSlots(
	p atoms.BadgeProps,
	iconStart []g.Node,
	iconEnd []g.Node,
) g.Node {
	variant := p.Variant
	if variant == "" {
		variant = "primary"
	}
	appearance := variantOr(clBadgeVariant, variant, "primary")
	if p.Tone != "" && p.Tone != "neutral" {
		appearance = variantOr(clBadgeTone, p.Tone, "neutral")
	}
	cl := clBadgeBase.
		Merge(appearance).
		Merge(variantOr(clBadgeSize, p.Size, "md"))
	var children []g.Node
	children = append(children, baseAttrs(p.ComponentProps)...)
	children = append(children, classes(cl.Compile(), p.Class))
	if p.Dot {
		children = append(children, h.Span(h.Class(clBadgeDot.Compile()), g.Attr("aria-hidden", "true")))
	}
	children = append(children, iconStart...)
	children = append(children, g.Text(p.Label))
	children = append(children, iconEnd...)
	return h.Span(children...)
}

// Alert renders atoms.AlertProps with role="alert" for danger/warning and
// role="status" otherwise, so severity maps to interruption behavior.
func Alert(p atoms.AlertProps) g.Node {
	return alertWithSlots(p, nil, nil)
}

func alertWithSlots(
	p atoms.AlertProps,
	iconStart []g.Node,
	actions []g.Node,
) g.Node {
	tone := p.Tone
	if tone == "" {
		tone = "info"
	}
	cl := clAlertBase.Merge(variantOr(clAlertVariant, tone, "info"))
	role := "status"
	live := "polite"
	if tone == "danger" || tone == "warning" {
		role = "alert"
		live = "assertive"
	}
	body := []g.Node{h.Class(clAlertBody.Compile())}
	if p.Title != "" {
		body = append(body, h.P(h.Class(clAlertTitle.Compile()), g.Text(p.Title)))
	}
	body = append(body, h.P(h.Class(clAlertMessage.Compile()), g.Text(p.Message)))

	var children []g.Node
	children = append(children, baseAttrs(p.ComponentProps)...)
	children = append(
		children,
		classes(cl.Compile(), p.Class),
		h.Role(role),
		g.Attr("aria-live", live),
		g.Attr("aria-atomic", "true"),
		g.Attr("data-component", "alert"),
		g.Attr("data-alert-tone", tone),
	)
	if p.Dismissible {
		children = append(
			children,
			g.Attr("data-controller", "alert"),
			g.Attr("data-alert-dismissible-value", "true"),
		)
	}
	if len(iconStart) > 0 {
		children = append(children, iconStart...)
	}
	children = append(children, h.Div(body...))
	if len(actions) > 0 {
		children = append(children, h.Div(
			h.Class("mt-3 flex flex-wrap items-center gap-3 text-sm"),
			g.Group(actions),
		))
	}
	if p.Dismissible {
		children = append(children, h.Button(
			h.Type("button"),
			g.Attr("data-action", "click->alert#dismiss"),
			g.Attr("data-alert-close", ""),
			g.Attr("aria-label", "Dismiss notification"),
			icon("x-mark"),
		))
	}
	return h.Div(children...)
}

// Input renders atoms.InputProps as a labelled form field. When Error is set
// the input carries aria-invalid and is described by the error element.
func Input(p atoms.InputProps) g.Node {
	return inputWithSlots(p, nil, nil)
}

func inputWithSlots(
	p atoms.InputProps,
	iconStart []g.Node,
	iconEnd []g.Node,
) g.Node {
	id := p.ID
	if id == "" && p.Name != "" {
		id = "pk-input-" + p.Name
	}
	typ := p.Type
	if typ == "" {
		typ = "text"
	}
	cl := clInput.
		Merge(clInputNormal).
		Merge(variantOr(clInputSize, p.Size, "md"))
	if p.Error != "" {
		cl = clInput.
			Merge(clInputError).
			Merge(variantOr(clInputSize, p.Size, "md"))
	}
	if len(iconStart) > 0 {
		cl = cl.Merge(clInputPadStart)
	}
	if len(iconEnd) > 0 {
		cl = cl.Merge(clInputPadEnd)
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
	if p.Step != "" {
		input = append(input, h.Step(p.Step))
	}
	if p.MinLength > 0 {
		input = append(input, g.Attr("minlength", itoa(p.MinLength)))
	}
	if p.MaxLength > 0 {
		input = append(input, h.MaxLength(itoa(p.MaxLength)))
	}
	if p.Pattern != "" {
		input = append(input, h.Pattern(p.Pattern))
	}
	if p.Autocomplete != "" {
		input = append(input, h.AutoComplete(p.Autocomplete))
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
	control := h.Input(input...)
	if len(iconStart) > 0 || len(iconEnd) > 0 {
		controlChildren := []g.Node{h.Class(clInputIconWrap.Compile())}
		if len(iconStart) > 0 {
			controlChildren = append(controlChildren, h.Span(
				h.Class(clInputIconStart.Compile()),
				g.Group(iconStart),
			))
		}
		controlChildren = append(controlChildren, control)
		if len(iconEnd) > 0 {
			controlChildren = append(controlChildren, h.Span(
				h.Class(clInputIconEnd.Compile()),
				g.Group(iconEnd),
			))
		}
		field = append(field, h.Div(controlChildren...))
	} else {
		field = append(field, control)
	}
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
	cl := clInput.Merge(clInputNormal).Merge(variantOr(clInputSize, "md", "md"))
	if p.Error != "" {
		cl = clInput.Merge(clInputError).Merge(variantOr(clInputSize, "md", "md"))
	}

	selectedValues := make(map[string]struct{}, len(p.Values)+1)
	for _, value := range p.Values {
		selectedValues[value] = struct{}{}
	}
	if p.Value != "" {
		selectedValues[p.Value] = struct{}{}
	}

	var options []g.Node
	if !p.Multiple && (p.Placeholder != "" || !p.Required) {
		label := p.Placeholder
		if label == "" {
			label = "Choose…"
		}
		options = append(options, h.Option(h.Value(""), g.Text(label)))
	}
	for _, o := range p.Options {
		opt := []g.Node{h.Value(o.Value), g.Text(o.Label)}
		if _, selected := selectedValues[o.Value]; selected {
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
	if p.Multiple {
		sel = append(sel, g.Attr("multiple", ""))
	}
	if p.VisibleRows > 0 {
		sel = append(sel, g.Attr("size", itoa(p.VisibleRows)))
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
	cl := clInput.Merge(clInputNormal).Merge(variantOr(clInputSize, "md", "md"))
	if p.ErrorMessage != "" {
		cl = clInput.Merge(clInputError).Merge(variantOr(clInputSize, "md", "md"))
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
	if p.MinLength > 0 {
		area = append(area, g.Attr("minlength", itoa(p.MinLength)))
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

// Divider renders atoms.DividerProps as an <hr>, a labelled horizontal
// separator, or a vertical separator.
func Divider(p atoms.DividerProps) g.Node {
	if p.Orientation == "vertical" {
		return h.Span(baseAttrs(
			p.ComponentProps,
			classes(clDividerV.Compile(), p.Class),
			h.Role("separator"), g.Attr("aria-orientation", "vertical"),
		)...)
	}
	if p.Text != "" {
		children := baseAttrs(
			p.ComponentProps,
			classes(clDividerText.Compile(), p.Class),
			h.Role("presentation"),
		)
		children = append(
			children,
			h.Span(
				classes(clDividerTextLine.Compile(), ""),
				g.Attr("aria-hidden", "true"),
			),
			h.Span(classes(clDividerTextLabel.Compile(), ""), g.Text(p.Text)),
			h.Span(
				classes(clDividerTextLine.Compile(), ""),
				g.Attr("aria-hidden", "true"),
			),
		)
		return h.Div(children...)
	}
	return h.Hr(baseAttrs(
		p.ComponentProps,
		classes(clDividerH.Compile(), p.Class),
		h.Role("separator"),
		g.Attr("aria-orientation", "horizontal"),
	)...)
}

// Spinner renders atoms.SpinnerProps; the label is announced, the rotation is
// decoration.
func Spinner(p atoms.SpinnerProps) g.Node {
	cl := clSpinner.
		Merge(variantOr(clSpinnerSize, p.Size, "md")).
		Merge(variantOr(clSpinnerTone, p.Tone, "brand"))
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
	return emptyStateWithSlots(p, nil, nil)
}

func emptyStateWithSlots(
	p atoms.EmptyStateProps,
	iconStart []g.Node,
	actions []g.Node,
) g.Node {
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
	children = append(children, iconStart...)
	children = append(children, h.P(h.Class(clEmptyTitle.Compile()), g.Text(p.Title)))
	if p.Description != "" {
		children = append(children, h.P(h.Class(clEmptyDesc.Compile()), g.Text(p.Description)))
	}
	children = append(children, actions...)
	return h.Div(children...)
}

// Kbd renders atoms.KbdProps as a sequence of <kbd> elements.
func Kbd(p atoms.KbdProps) g.Node {
	var children []g.Node
	children = append(children, baseAttrs(p.ComponentProps)...)
	children = append(children, h.Class(clCheckRow.Compile()))
	size := variantOr(clKbdSize, p.Size, "md")
	for _, k := range p.Keys {
		children = append(children, h.Kbd(classes(clKbd.Merge(size).Compile(), p.Class), g.Text(k)))
	}
	return h.Span(children...)
}

// Link renders atoms.LinkProps; external links open safely.
func Link(p atoms.LinkProps) g.Node {
	return linkWithSlots(p, nil)
}

func linkWithSlots(
	p atoms.LinkProps,
	trailingAdornment []g.Node,
) g.Node {
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
	children = append(children, g.Text(p.Label))
	children = append(children, trailingAdornment...)
	if p.External {
		children = append(children, h.Span(g.Attr("aria-hidden", "true"), g.Text(" ↗")))
	}
	return h.A(children...)
}

// Tag renders atoms.TagProps; removable tags post to OnRemoveURL.
func Tag(p atoms.TagProps) g.Node {
	return tagWithSlots(p, nil)
}

func tagWithSlots(
	p atoms.TagProps,
	iconStart []g.Node,
) g.Node {
	cl := clTagBase.Merge(clTagIdle)
	if p.Tone != "" && p.Tone != "neutral" {
		cl = clTagBase.Merge(variantOr(clTagTone, p.Tone, "neutral"))
	} else if p.Selected {
		cl = clTagBase.Merge(clTagSelected)
	}
	var children []g.Node
	children = append(children, baseAttrs(p.ComponentProps)...)
	children = append(children, classes(cl.Compile(), p.Class))
	children = append(children, iconStart...)
	children = append(children, g.Text(p.Label))
	if p.Removable && p.OnRemoveURL != "" {
		children = append(children, h.Button(
			h.Type("button"), h.Class(clLink.Compile()),
			g.Attr("hx-delete", p.OnRemoveURL),
			g.Attr("aria-label", "Remove "+p.Label),
			g.Text("×"),
		))
	}
	return h.Span(children...)
}
