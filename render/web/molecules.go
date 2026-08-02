// Implements: REQ-011.
// Per: ADR-0031.
// Discipline: C-14.

package web

// molecules.go renders the molecule contracts the admin and module pages
// compose: tabular data, cards, navigation, and progressively enhanced
// overlays. Interaction contracts stay in the shared runtime controllers so
// downstream products do not need private scripts or duplicate markup.

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"github.com/septagon-oss/pk-ui/contracts"
	"github.com/septagon-oss/pk-ui/contracts/atoms"
	"github.com/septagon-oss/pk-ui/contracts/molecules"
	"github.com/septagon-oss/tw"
)

// CheckboxGroup renders a native multiple-choice fieldset from canonical
// Checkbox atoms. Required is exposed as a group semantic for assistive
// technology and server validation; applying HTML required to every checkbox
// would incorrectly require every option to be selected.
func CheckboxGroup(p molecules.CheckboxGroupProps) g.Node {
	selected := make(map[string]struct{}, len(p.Selected))
	for _, value := range p.Selected {
		selected[value] = struct{}{}
	}
	return choiceGroup(
		"checkbox-group",
		p.ComponentProps,
		p.HTMXProps,
		p.Name,
		p.Label,
		p.Description,
		p.AriaLabel,
		p.Required,
		p.Orientation,
		p.Error,
		len(p.Options),
		func(groupID string, index int) g.Node {
			option := p.Options[index]
			_, checked := selected[option.Value]
			return Checkbox(atoms.CheckboxProps{
				ComponentProps: contracts.ComponentProps{
					ID:       fmt.Sprintf("%s-option-%d", groupID, index),
					Disabled: p.Disabled || option.Disabled,
				},
				Name:     p.Name,
				Label:    choiceOptionLabel(option),
				Value:    option.Value,
				Checked:  checked,
				HelpText: option.Description,
			})
		},
	)
}

// RadioGroup renders a native single-choice fieldset from canonical Radio
// atoms. Browser radio behavior remains authoritative for keyboard navigation,
// selection, required validation, and submitted values.
func RadioGroup(p molecules.RadioGroupProps) g.Node {
	return choiceGroup(
		"radio-group",
		p.ComponentProps,
		p.HTMXProps,
		p.Name,
		p.Label,
		p.Description,
		p.AriaLabel,
		p.Required,
		p.Orientation,
		p.Error,
		len(p.Options),
		func(groupID string, index int) g.Node {
			option := p.Options[index]
			return Radio(atoms.RadioProps{
				ComponentProps: contracts.ComponentProps{
					ID:       fmt.Sprintf("%s-option-%d", groupID, index),
					Disabled: p.Disabled || option.Disabled,
				},
				Name:     p.Name,
				Label:    choiceOptionLabel(option),
				HelpText: option.Description,
				Value:    option.Value,
				Checked:  option.Value == p.Value,
				Required: p.Required,
			})
		},
	)
}

func choiceGroup(
	componentName string,
	componentProps contracts.ComponentProps,
	enhancement contracts.HTMXProps,
	name string,
	label string,
	description string,
	ariaLabel string,
	required bool,
	orientation string,
	errorMessage string,
	optionCount int,
	optionAt func(groupID string, index int) g.Node,
) g.Node {
	groupID := strings.TrimSpace(componentProps.ID)
	if groupID == "" {
		groupID = "pk-" + componentName
		if strings.TrimSpace(name) != "" {
			groupID += "-" + strings.TrimSpace(name)
		}
	}
	orientation = strings.ToLower(strings.TrimSpace(orientation))
	if orientation != "horizontal" {
		orientation = "vertical"
	}
	describedBy := make([]string, 0, 2)
	if strings.TrimSpace(description) != "" {
		describedBy = append(describedBy, groupID+"-description")
	}
	if strings.TrimSpace(errorMessage) != "" {
		describedBy = append(describedBy, groupID+"-error")
	}

	rootProps := componentProps
	rootProps.ID = ""
	rootProps.Class = ""
	rootProps.Disabled = false
	root := baseAttrs(rootProps,
		h.ID(groupID),
		classes(clChoiceGroupRoot.Compile(), componentProps.Class),
		g.Attr("data-component", componentName),
		g.Attr("data-orientation", orientation),
	)
	if componentName == "radio-group" {
		root = append(root, h.Role("radiogroup"))
	}
	if componentProps.Disabled {
		root = append(root, h.Disabled())
	}
	if required {
		root = append(root, g.Attr("aria-required", "true"), g.Attr("data-required", "true"))
	}
	if errorMessage != "" {
		root = append(root, g.Attr("aria-invalid", "true"))
	}
	if len(describedBy) > 0 {
		root = append(root, g.Attr("aria-describedby", strings.Join(describedBy, " ")))
	}
	if strings.TrimSpace(label) == "" {
		if strings.TrimSpace(ariaLabel) == "" {
			ariaLabel = name
		}
		if strings.TrimSpace(ariaLabel) != "" {
			root = append(root, g.Attr("aria-label", ariaLabel))
		}
	}
	if enhancement.Trigger == "" && hasHTMXRequest(enhancement) {
		enhancement.Trigger = "change"
	}
	root = append(root, htmxAttrs(enhancement)...)

	if label != "" {
		legend := []g.Node{h.Class(clChoiceGroupLegend.Compile()), g.Text(label)}
		if required {
			legend = append(legend, h.Span(h.Class(clRequired.Compile()), g.Attr("aria-hidden", "true"), g.Text(" *")))
		}
		root = append(root, h.Legend(legend...))
	}
	if description != "" {
		root = append(root, h.P(
			h.ID(groupID+"-description"),
			h.Class(clChoiceGroupDescription.Compile()),
			g.Text(description),
		))
	}
	optionsClass := clChoiceGroupOptions.Merge(clChoiceGroupVertical)
	if orientation == "horizontal" {
		optionsClass = clChoiceGroupOptions.Merge(clChoiceGroupHorizontal)
	}
	options := make([]g.Node, 0, optionCount+1)
	options = append(options, h.Class(optionsClass.Compile()))
	for index := 0; index < optionCount; index++ {
		options = append(options, optionAt(groupID, index))
	}
	root = append(root, h.Div(options...))
	if errorMessage != "" {
		root = append(root, h.P(
			h.ID(groupID+"-error"),
			h.Class(clChoiceGroupError.Compile()),
			h.Role("alert"),
			g.Text(errorMessage),
		))
	}
	return h.FieldSet(root...)
}

func hasHTMXRequest(p contracts.HTMXProps) bool {
	return p.Get != "" || p.Post != "" || p.Put != "" || p.Patch != "" || p.Delete != ""
}

func choiceOptionLabel(option molecules.Option) string {
	if label := strings.TrimSpace(option.Label); label != "" {
		return label
	}
	return option.Value
}

// ActionMenu renders the canonical progressively enhanced menu-button pattern.
// Navigation, custom action events, and HTMX mutations all use the same item
// contract while the shared action-menu controller owns keyboard and focus
// behavior.
func ActionMenu(p molecules.ActionMenuProps) g.Node {
	align := strings.ToLower(strings.TrimSpace(p.Align))
	if align != "start" {
		align = "end"
	}
	width := strings.ToLower(strings.TrimSpace(p.Width))
	if _, ok := clActionMenuWidth[width]; !ok {
		width = "md"
	}
	triggerLabel := strings.TrimSpace(p.TriggerLabel)
	triggerAriaLabel := triggerLabel
	if triggerAriaLabel == "" {
		triggerAriaLabel = "Actions"
	}
	triggerIcon := strings.TrimSpace(p.TriggerIcon)
	if triggerIcon == "" {
		triggerIcon = "ellipsis-vertical"
	}

	panelID := ""
	if id := strings.TrimSpace(p.ID); id != "" {
		panelID = id + "-panel"
	}
	trigger := []g.Node{
		h.Type("button"), h.Class(clActionMenuTrigger.Compile()),
		g.Attr("data-action-menu-trigger", ""),
		g.Attr("data-action", "click->action-menu#toggle keydown.escape->action-menu#close"),
		g.Attr("aria-haspopup", "menu"), g.Attr("aria-expanded", "false"),
		g.Attr("aria-label", triggerAriaLabel),
	}
	if panelID != "" {
		trigger = append(trigger, g.Attr("aria-controls", panelID))
	}
	if p.Disabled {
		trigger = append(trigger, h.Disabled())
	}
	trigger = append(trigger,
		Icon(atoms.IconProps{Name: triggerIcon, Size: "md", Tone: "neutral"}),
	)
	if triggerLabel != "" {
		trigger = append(trigger, g.Text(triggerLabel))
	}

	sections := make([]molecules.ActionMenuSection, 0, len(p.Sections)+1)
	if len(p.Items) > 0 {
		sections = append(sections, molecules.ActionMenuSection{Items: p.Items})
	}
	sections = append(sections, p.Sections...)
	traceAttrs := make(map[string]string, 3)
	for _, name := range []string{"data-request-id", "data-trace-id", "data-span-id"} {
		if value := p.Attrs[name]; value != "" {
			traceAttrs[name] = value
		}
	}
	sectionNodes := make([]g.Node, 0, len(sections)*2)
	for index, section := range sections {
		if index > 0 {
			sectionNodes = append(sectionNodes, h.Div(
				h.Class(clActionMenuSeparator.Compile()), h.Role("separator"),
				g.Attr("aria-hidden", "true"),
			))
		}
		group := []g.Node{h.Role("group")}
		if label := strings.TrimSpace(section.Label); label != "" {
			group = append(group, g.Attr("aria-label", label), h.Div(
				h.Class(clActionMenuSectionLabel.Compile()), h.Role("presentation"), g.Text(label),
			))
		}
		for _, item := range section.Items {
			group = append(group, actionMenuItem(p.Disabled, item, traceAttrs))
		}
		sectionNodes = append(sectionNodes, h.Div(group...))
	}

	panel := []g.Node{
		h.Class(clActionMenuPanel.
			Merge(clActionMenuAlign[align]).
			Merge(clActionMenuWidth[width]).
			Compile()),
		g.Attr("data-action-menu-panel", ""), h.Role("menu"),
		g.Attr("aria-orientation", "vertical"), g.Attr("tabindex", "-1"),
		g.Attr("hidden"),
	}
	if panelID != "" {
		panel = append(panel, h.ID(panelID))
	}
	panel = append(panel, h.Div(h.Class(clActionMenuPanelInner.Compile()), g.Group(sectionNodes)))

	rootProps := p.ComponentProps
	rootProps.Disabled = false
	rootProps.Class = ""
	root := baseAttrs(rootProps)
	root = append(root,
		classes(clActionMenuRoot.Compile(), p.Class),
		g.Attr("data-component", "action-menu"),
		g.Attr("data-controller", "action-menu"),
		g.Attr("data-action-menu-align-value", align),
		g.Attr("data-state", "closed"),
	)
	if p.Disabled {
		root = append(root, g.Attr("aria-disabled", "true"))
	}
	root = append(root, h.Button(trigger...), h.Div(panel...))
	return h.Div(root...)
}

func actionMenuItem(menuDisabled bool, item molecules.ActionMenuItem, traceAttrs map[string]string) g.Node {
	disabled := menuDisabled || item.Disabled
	tone := strings.ToLower(strings.TrimSpace(item.Tone))
	if item.Danger {
		tone = "danger"
	}
	itemClass := clActionMenuItem.Merge(clActionMenuItemTone["neutral"])
	if disabled {
		itemClass = clActionMenuItem.Merge(clActionMenuItemDisabled)
	} else if tone == "danger" {
		itemClass = clActionMenuItem.Merge(clActionMenuItemTone["danger"])
	}

	attrs := []g.Node{
		h.Class(itemClass.Compile()), h.Role("menuitem"), g.Attr("tabindex", "-1"),
		g.Attr("data-action", "click->action-menu#select keydown.enter->action-menu#select keydown.space->action-menu#select"),
	}
	if disabled {
		attrs = append(attrs, g.Attr("aria-disabled", "true"))
	}
	if !disabled {
		attrs = append(attrs, htmxAttrs(contracts.HTMXProps{
			Get: item.HxGet, Post: item.HxPost, Delete: item.HxDelete,
			Target: item.HxTarget, Swap: item.HxSwap, Confirm: item.HxConfirm,
		})...)
		if item.Action != "" {
			attrs = append(attrs, g.Attr("data-item-action", item.Action))
		}
	}
	trustedAttrs := make(map[string]string, len(traceAttrs)+len(item.Attrs))
	for name, value := range traceAttrs {
		trustedAttrs[name] = value
	}
	for name, value := range item.Attrs {
		trustedAttrs[name] = value
	}
	attrs = append(attrs, attrPairs(trustedAttrs)...)

	content := make([]g.Node, 0, 2)
	if item.Icon != "" {
		content = append(content, h.Span(
			h.Class(clActionMenuItemIcon.Compile()),
			Icon(atoms.IconProps{Name: item.Icon, Size: "sm", Tone: "neutral"}),
		))
	}
	content = append(content, h.Span(h.Class(clActionMenuItemLabel.Compile()), g.Text(item.Label)))
	if item.Href != "" && !disabled {
		attrs = append(attrs, h.Href(item.Href))
		return h.A(append(attrs, content...)...)
	}
	attrs = append(attrs, h.Type("button"))
	if disabled {
		attrs = append(attrs, h.Disabled())
	}
	return h.Button(append(attrs, content...)...)
}

// DrawerSlots is the trusted Go composition seam for rich drawer regions.
// Portable clients use the string fields on DrawerProps; server-rendered Go
// applications use these slots without reimplementing the overlay shell.
type DrawerSlots struct {
	Header []g.Node
	Body   []g.Node
	Footer []g.Node
}

// Drawer renders a portable drawer contract.
func Drawer(p molecules.DrawerProps) g.Node {
	return DrawerWithSlots(p, DrawerSlots{})
}

// DrawerWithSlots renders the canonical accessible modal edge panel. The
// drawer controller owns open/close state, focus containment, focus return,
// and document scroll locking.
func DrawerWithSlots(p molecules.DrawerProps, slots DrawerSlots) g.Node {
	position := strings.ToLower(strings.TrimSpace(p.Position))
	if _, ok := clDrawerPosition[position]; !ok {
		position = "right"
	}
	size := strings.ToLower(strings.TrimSpace(p.Size))
	if _, ok := clDrawerWidth[size]; !ok {
		size = "medium"
	}
	closable := defaultTrue(p.Closable)
	closeOnOverlay := defaultTrue(p.CloseOnOverlay)
	closeOnEscape := defaultTrue(p.CloseOnEscape)
	showOverlay := defaultTrue(p.ShowOverlay)
	open := p.Open && !p.Hidden
	state := "closed"
	if open {
		state = "open"
	}

	panelClass := clDrawerPanel.Merge(clDrawerPosition[position])
	if position == "bottom" {
		panelClass = panelClass.Merge(clDrawerBottomSize[size])
	} else {
		panelClass = panelClass.Merge(clDrawerWidth[size])
	}

	titleID := ""
	if p.ID != "" && p.Title != "" && len(slots.Header) == 0 {
		titleID = p.ID + "-title"
	}
	accessibleName := strings.TrimSpace(p.AriaLabel)
	if accessibleName == "" {
		accessibleName = strings.TrimSpace(p.Title)
	}
	if accessibleName == "" {
		accessibleName = "Drawer"
	}
	closeLabel := strings.TrimSpace(p.CloseLabel)
	if closeLabel == "" {
		closeLabel = "Close"
	}

	headerContent := slots.Header
	if len(headerContent) == 0 && (p.Title != "" || p.Description != "") {
		text := make([]g.Node, 0, 2)
		if p.Title != "" {
			title := []g.Node{
				h.Class(clDrawerTitle.Compile()),
				g.Attr("data-drawer-title", ""),
				g.Text(p.Title),
			}
			if titleID != "" {
				title = append([]g.Node{h.ID(titleID)}, title...)
			}
			text = append(text, h.H2(title...))
		}
		if p.Description != "" {
			text = append(text, h.P(
				h.Class(clDrawerDescription.Compile()),
				g.Attr("data-drawer-description", ""),
				g.Text(p.Description),
			))
		}
		headerContent = []g.Node{h.Div(h.Class(clDrawerTitleBlock.Compile()), g.Group(text))}
	}
	header := make([]g.Node, 0, len(headerContent)+1)
	header = append(header, headerContent...)
	if closable {
		closeButton := []g.Node{
			h.Type("button"), h.Class(clDrawerClose.Compile()),
			g.Attr("data-action", "click->drawer#close"),
			g.Attr("data-drawer-close", ""),
			g.Attr("aria-label", closeLabel),
		}
		if p.Disabled {
			closeButton = append(closeButton, h.Disabled())
		}
		closeButton = append(closeButton, Icon(atoms.IconProps{Name: "x", Size: "md"}))
		header = append(header, h.Button(closeButton...))
	}

	body := slots.Body
	if len(body) == 0 && p.Body != "" {
		body = []g.Node{g.Text(p.Body)}
	}
	footer := slots.Footer
	if len(footer) == 0 && p.Footer != "" {
		footer = []g.Node{g.Text(p.Footer)}
	}

	panel := []g.Node{
		h.Class(panelClass.Compile()),
		g.Attr("data-drawer-panel", ""),
		g.Attr("data-state", state),
		g.Attr("aria-hidden", boolText(!open)),
		g.Attr("tabindex", "-1"),
	}
	if !open {
		panel = append(panel, g.Attr("hidden"))
	}
	if len(header) > 0 {
		panel = append(panel, h.Div(h.Class(clDrawerHeader.Compile()), g.Group(header)))
	}
	panel = append(panel, h.Div(h.Class(clDrawerBody.Compile()), g.Attr("data-drawer-body", ""), g.Group(body)))
	if len(footer) > 0 {
		panel = append(panel, h.Div(
			h.Class(clDrawerFooter.Compile()),
			g.Attr("data-drawer-footer", ""),
			g.Group(footer),
		))
	}

	rootProps := p.ComponentProps
	rootProps.Class = ""
	rootProps.Disabled = false
	rootProps.Hidden = false
	root := baseAttrs(rootProps)
	root = append(root,
		classes(clDrawerRoot.Compile(), p.Class),
		g.Attr("data-component", "drawer"),
		g.Attr("data-controller", "drawer"),
		g.Attr("data-drawer-position-value", position),
		g.Attr("data-drawer-open-value", boolText(open)),
		g.Attr("data-drawer-close-on-escape-value", boolText(closeOnEscape)),
		g.Attr("data-state", state),
		h.Role("dialog"), g.Attr("aria-modal", "true"),
		g.Attr("aria-hidden", boolText(!open)),
	)
	if p.ID != "" {
		root = append(root, g.Attr("data-drawer-id-value", p.ID))
	}
	if titleID != "" && p.AriaLabel == "" {
		root = append(root, g.Attr("aria-labelledby", titleID))
	} else {
		root = append(root, g.Attr("aria-label", accessibleName))
	}
	if p.OpenOnSwap {
		root = append(root, g.Attr("data-action", "htmx:afterSwap->drawer#openDrawer"))
	}
	if p.Disabled {
		root = append(root, g.Attr("aria-disabled", "true"))
	}
	if !open {
		root = append(root, g.Attr("hidden"))
	}
	if showOverlay {
		overlay := []g.Node{
			h.Class(clDrawerOverlay.Compile()),
			g.Attr("data-drawer-backdrop", ""),
			g.Attr("data-state", state),
			g.Attr("aria-hidden", "true"),
		}
		if closeOnOverlay && !p.Disabled {
			overlay = append(overlay, g.Attr("data-action", "click->drawer#close"))
		}
		if !open {
			overlay = append(overlay, g.Attr("hidden"))
		}
		root = append(root, h.Div(overlay...))
	}
	root = append(root, h.Div(panel...))
	return h.Div(root...)
}

func defaultTrue(value *bool) bool {
	return value == nil || *value
}

// ModalSlots is the trusted Go composition seam for rich dialog regions.
// Portable clients use ModalProps body/footer strings; server-rendered Go
// applications use these slots without recreating modal chrome or behavior.
type ModalSlots struct {
	Header []g.Node
	Body   []g.Node
	Footer []g.Node
}

// Modal renders a portable modal contract.
func Modal(p molecules.ModalProps) g.Node {
	return ModalWithSlots(p, ModalSlots{})
}

// ModalWithSlots renders the canonical modal root. Deferred roots deliberately
// contain no panel: HTMX swaps a ModalPanelWithSlots response into the root and
// the shared htmx-modal controller opens it after the swap.
func ModalWithSlots(p molecules.ModalProps, slots ModalSlots) g.Node {
	size := modalSize(p.Size)
	open := p.Open && !p.Hidden
	state := "closed"
	if open {
		state = "open"
	}
	centered := defaultTrue(p.Centered)
	closeOnEscape := defaultTrue(p.CloseOnEscape)
	clearOnClose := p.Deferred
	if p.ClearOnClose != nil {
		clearOnClose = *p.ClearOnClose
	}

	rootClass := clModalRoot.Merge(clModalCentered)
	if !centered {
		rootClass = clModalRoot.Merge(clModalBottomSheet)
	}
	rootProps := p.ComponentProps
	rootProps.Class = ""
	rootProps.Disabled = false
	rootProps.Hidden = false
	root := baseAttrs(rootProps)
	root = append(root,
		classes(rootClass.Compile(), p.Class),
		g.Attr("data-component", "modal"),
		g.Attr("data-controller", "htmx-modal"),
		g.Attr("data-htmx-modal-open-value", boolText(open)),
		g.Attr("data-htmx-modal-close-on-escape-value", boolText(closeOnEscape)),
		g.Attr("data-htmx-modal-clear-on-close-value", boolText(clearOnClose)),
		g.Attr("data-state", state),
		g.Attr("aria-hidden", boolText(!open)),
		g.Attr("tabindex", "-1"),
	)
	if p.Disabled {
		root = append(root, g.Attr("aria-disabled", "true"))
	}
	if !open {
		root = append(root, g.Attr("hidden"), h.Style("display:none"))
	}

	actions := make([]string, 0, 2)
	if p.OpenOnSwap {
		actions = append(actions, "htmx:afterSwap->htmx-modal#show")
	}
	if p.Deferred && defaultTrue(p.CloseOnOverlay) && !p.Disabled {
		actions = append(actions, "click->htmx-modal#backdropClick")
	}
	if len(actions) > 0 {
		root = append(root, g.Attr("data-action", strings.Join(actions, " ")))
	}
	if p.Deferred {
		return h.Div(root...)
	}

	root = append(root, h.Role("dialog"), g.Attr("aria-modal", "true"))
	titleID := modalTitleID(p, slots)
	accessibleName := modalAccessibleName(p)
	if titleID != "" && strings.TrimSpace(p.AriaLabel) == "" {
		root = append(root, g.Attr("aria-labelledby", titleID))
	} else {
		root = append(root, g.Attr("aria-label", accessibleName))
	}
	if defaultTrue(p.ShowOverlay) {
		overlay := []g.Node{
			h.Class(clModalOverlay.Compile()),
			g.Attr("data-modal-backdrop", ""),
			g.Attr("aria-hidden", "true"),
		}
		if defaultTrue(p.CloseOnOverlay) && !p.Disabled {
			overlay = append(overlay, g.Attr("data-action", "click->htmx-modal#close"))
		}
		root = append(root, h.Div(overlay...))
	}
	root = append(root, modalPanel(p, slots, false, size))
	return h.Div(root...)
}

// ModalPanelWithSlots renders the panel fragment returned by an HTMX endpoint.
// It is a complete dialog because its deferred parent is only a swap target.
func ModalPanelWithSlots(p molecules.ModalProps, slots ModalSlots) g.Node {
	return modalPanel(p, slots, true, modalSize(p.Size))
}

// ModalPanel renders a server-loaded modal panel with one rich body node.
func ModalPanel(p molecules.ModalProps, body g.Node) g.Node {
	return ModalPanelWithSlots(p, ModalSlots{Body: []g.Node{body}})
}

func modalPanel(p molecules.ModalProps, slots ModalSlots, standalone bool, size string) g.Node {
	closable := defaultTrue(p.Closable)
	showClose := defaultTrue(p.ShowClose)
	titleID := modalTitleID(p, slots)
	closeLabel := strings.TrimSpace(p.CloseLabel)
	if closeLabel == "" {
		closeLabel = "Close"
	}

	headerContent := slots.Header
	if len(headerContent) == 0 && (p.Title != "" || p.Description != "") {
		text := make([]g.Node, 0, 2)
		if p.Title != "" {
			title := []g.Node{h.Class(clModalTitle.Compile()), g.Text(p.Title)}
			if titleID != "" {
				title = append([]g.Node{h.ID(titleID)}, title...)
			}
			text = append(text, h.H2(title...))
		}
		if p.Description != "" {
			text = append(text, h.P(h.Class(clModalDescription.Compile()), g.Text(p.Description)))
		}
		headerContent = []g.Node{h.Div(h.Class(clModalTitleBlock.Compile()), g.Group(text))}
	}
	header := append([]g.Node(nil), headerContent...)
	if closable && showClose {
		header = append(header, ModalCloseButton(closeLabel, ""))
	}

	body := slots.Body
	if len(body) == 0 && p.Body != "" {
		body = []g.Node{g.Text(p.Body)}
	}
	footer := slots.Footer
	if len(footer) == 0 && p.Footer != "" {
		footer = []g.Node{g.Text(p.Footer)}
	}

	panel := []g.Node{
		h.Class(clModalPanel.Merge(clModalPanelSize[size]).Compile()),
		g.Attr("data-modal-panel", ""),
		g.Attr("data-action", "click->htmx-modal#stopPropagation"),
		g.Attr("tabindex", "-1"),
	}
	if standalone {
		panel = append(panel, h.Role("dialog"), g.Attr("aria-modal", "true"))
		if titleID != "" && strings.TrimSpace(p.AriaLabel) == "" {
			panel = append(panel, g.Attr("aria-labelledby", titleID))
		} else {
			panel = append(panel, g.Attr("aria-label", modalAccessibleName(p)))
		}
	}
	if len(header) > 0 {
		panel = append(panel, h.Div(h.Class(clModalHeader.Compile()), g.Group(header)))
	}
	panel = append(panel, h.Div(
		h.Class(clModalBody.Compile()),
		g.Attr("data-modal-body", ""),
		g.Group(body),
	))
	if len(footer) > 0 {
		panel = append(panel, h.Div(
			h.Class(clModalFooter.Compile()),
			g.Attr("data-modal-footer", ""),
			g.Group(footer),
		))
	}
	return h.Div(panel...)
}

func modalSize(value string) string {
	size := strings.ToLower(strings.TrimSpace(value))
	if size == "fullscreen" {
		size = "full"
	}
	if _, ok := clModalPanelSize[size]; !ok {
		return "medium"
	}
	return size
}

func modalTitleID(p molecules.ModalProps, slots ModalSlots) string {
	if p.ID == "" || p.Title == "" || len(slots.Header) > 0 {
		return ""
	}
	return p.ID + "-title"
}

func modalAccessibleName(p molecules.ModalProps) string {
	if label := strings.TrimSpace(p.AriaLabel); label != "" {
		return label
	}
	if title := strings.TrimSpace(p.Title); title != "" {
		return title
	}
	return "Dialog"
}

// ModalCloseButton creates the canonical controller-backed icon close action.
func ModalCloseButton(label, class string) g.Node {
	if strings.TrimSpace(label) == "" {
		label = "Close"
	}
	return h.Button(
		h.Type("button"), classes(clModalClose.Compile(), class),
		g.Attr("data-action", "click->htmx-modal#close"),
		g.Attr("data-modal-close", ""), g.Attr("aria-label", label),
		Icon(atoms.IconProps{Name: "x", Size: "md"}),
	)
}

// ModalCancelButton creates a text dismissal action for modal footers.
func ModalCancelButton(label, class string) g.Node {
	if strings.TrimSpace(label) == "" {
		label = "Cancel"
	}
	return h.Button(
		h.Type("button"), classes(clModalCancel.Compile(), class),
		g.Attr("data-action", "click->htmx-modal#close"),
		g.Attr("data-modal-cancel", ""), g.Text(label),
	)
}

// FileUploadSlots is the trusted Go composition seam for branded drop-zone
// content. Portable clients use the localized copy fields on FileUploadProps;
// product renderers can replace only the prompt without rebuilding the input,
// validation, progress, list, or controller contract.
type FileUploadSlots struct {
	Content []g.Node
}

// FileUpload renders the canonical progressively enhanced upload field.
// Without JavaScript it remains a native multipart file input. When UploadURL
// is set, the shared controller uploads one file directly and stores the
// returned public URL in the named hidden input.
func FileUpload(p molecules.FileUploadProps) g.Node {
	return FileUploadWithSlots(p, FileUploadSlots{})
}

// FileUploadWithSlots renders FileUpload with optional custom drop-zone copy.
func FileUploadWithSlots(p molecules.FileUploadProps, slots FileUploadSlots) g.Node {
	name := strings.TrimSpace(p.Name)
	if name == "" {
		name = "file"
	}
	inputID := name
	if id := strings.TrimSpace(p.ID); id != "" {
		inputID = id + "-input"
	}
	dropZone := defaultTrue(p.DropZone)
	showList := defaultTrue(p.ShowList)
	remote := strings.TrimSpace(p.UploadURL) != ""
	multiple := p.Multiple && !remote

	accessibleLabel := strings.TrimSpace(p.AriaLabel)
	if accessibleLabel == "" {
		accessibleLabel = strings.TrimSpace(p.Label)
	}
	if accessibleLabel == "" {
		accessibleLabel = name
	}
	promptLabel := fileUploadDefault(p.PromptLabel, "Click to upload")
	dropLabel := fileUploadDefault(p.DropLabel, "or drag and drop")
	chooseLabel := fileUploadTemplate(
		fileUploadDefault(p.ChooseLabel, "Choose {label}"),
		map[string]string{"label": accessibleLabel},
	)
	dropZoneLabel := fileUploadTemplate(
		"Upload {label}",
		map[string]string{"label": accessibleLabel},
	)
	loadingLabel := fileUploadDefault(p.LoadingLabel, "Uploading...")
	removeLabel := fileUploadDefault(p.RemoveLabel, "Remove file")
	uploadedLabel := fileUploadDefault(p.UploadedLabel, "Uploaded file")
	missingError := fileUploadDefault(p.MissingError, "Choose a file")
	tooLargeError := fileUploadDefault(p.TooLargeError, "{name} exceeds the maximum file size ({maxSize})")
	typeError := fileUploadDefault(p.TypeError, "{name} is not an accepted file type")

	describedBy := ""
	if strings.TrimSpace(p.HelpText) != "" {
		describedBy = inputID + "-help"
	}
	inputAttrs := []g.Node{
		h.ID(inputID), h.Type("file"),
		g.Attr("data-file-upload-input", "true"),
		g.Attr("data-action", "change->file-upload#handleChange"),
		g.Attr("aria-label", chooseLabel),
	}
	if !remote {
		inputAttrs = append(inputAttrs, h.Name(name))
	}
	if accept := strings.TrimSpace(p.Accept); accept != "" {
		inputAttrs = append(inputAttrs, g.Attr("accept", accept))
	}
	if multiple {
		inputAttrs = append(inputAttrs, g.Attr("multiple", ""))
	}
	if p.Required && !remote {
		inputAttrs = append(inputAttrs, h.Required())
	}
	if p.Disabled {
		inputAttrs = append(inputAttrs, h.Disabled())
	}
	if describedBy != "" {
		inputAttrs = append(inputAttrs, g.Attr("aria-describedby", describedBy))
	}

	inputClass := clFileUploadInputHidden.Compile()
	if !dropZone {
		inputClass = clFileUploadInput.Compile()
	}
	inputAttrs = append(inputAttrs, h.Class(inputClass))

	controls := make([]g.Node, 0, 4)
	if remote {
		valueAttrs := []g.Node{
			h.Type("hidden"), h.Name(name), h.Value(p.Value),
			g.Attr("data-file-upload-value-input", "true"),
		}
		if p.Disabled {
			valueAttrs = append(valueAttrs, h.Disabled())
		}
		controls = append(controls, h.Input(valueAttrs...))
	}
	fileInput := h.Input(inputAttrs...)
	if dropZone {
		zoneClass := clFileUploadDropZone
		if p.Disabled {
			zoneClass = zoneClass.Merge(clFileUploadDropZoneDisabled)
		}
		promptContent := slots.Content
		if len(promptContent) == 0 {
			promptContent = []g.Node{
				h.Span(
					h.Class(clFileUploadIcon.Compile()),
					Icon(atoms.IconProps{Name: "upload", Size: "xl", Tone: "neutral"}),
				),
				h.P(
					h.Class(clFileUploadPromptText.Compile()),
					h.Span(h.Class(clFileUploadPromptAction.Compile()), g.Text(promptLabel)),
					g.Text(" "+dropLabel),
				),
			}
			if accept := strings.TrimSpace(p.Accept); accept != "" {
				promptContent = append(promptContent, h.P(h.Class(clFileUploadHint.Compile()), g.Text(accept)))
			}
			if p.MaxSize > 0 {
				maxSizeLabel := fileUploadTemplate(
					fileUploadDefault(p.MaxSizeLabel, "Max size: {size}"),
					map[string]string{"size": fileUploadFormatBytes(p.MaxSize)},
				)
				promptContent = append(promptContent, h.P(h.Class(clFileUploadHint.Compile()), g.Text(maxSizeLabel)))
			}
		}
		controls = append(controls,
			fileInput,
			h.Label(
				h.For(inputID), h.Class(zoneClass.Compile()),
				g.Attr("aria-label", dropZoneLabel),
				g.Attr("data-file-upload-dropzone", "true"),
				g.Attr("data-action", "dragover->file-upload#dragOver dragleave->file-upload#dragLeave drop->file-upload#drop"),
				h.Div(
					h.Class(clFileUploadDropZoneInner.Compile()),
					h.Div(
						h.Class(clFileUploadLoading.Compile()),
						g.Attr("data-file-upload-loading", "true"), g.Attr("hidden", ""),
						h.Span(
							h.Class(clFileUploadLoadingIcon.Compile()),
							Spinner(atoms.SpinnerProps{Label: loadingLabel, Size: "lg"}),
						),
						h.Span(h.Class(clFileUploadPromptText.Compile()), g.Text(loadingLabel)),
					),
					h.Div(
						h.Class(clFileUploadPrompt.Compile()),
						g.Attr("data-file-upload-prompt", "true"),
						g.Group(promptContent),
					),
				),
			),
		)
	} else {
		controls = append(controls, fileInput)
	}
	if showList {
		controls = append(controls, h.Div(
			h.Class(clFileUploadList.Compile()),
			g.Attr("data-file-upload-list", "true"),
			g.Attr("data-action", "click->file-upload#clickList"),
			g.Attr("data-part", "file-list"), g.Attr("hidden", ""),
			h.P(
				h.Class(clFileUploadError.Compile()),
				g.Attr("data-file-upload-error", "true"), g.Attr("hidden", ""),
			),
		))
	}

	rootProps := p.ComponentProps
	rootProps.Disabled = false
	rootProps.Class = ""
	root := baseAttrs(rootProps)
	root = append(root,
		classes(clFileUploadRoot.Compile(), p.Class),
		g.Attr("data-component", "fileupload"),
		g.Attr("data-controller", "file-upload"),
		g.Attr("data-state", "idle"),
		g.Attr("data-file-upload-multiple-value", strconv.FormatBool(multiple)),
		g.Attr("data-file-upload-disabled-value", strconv.FormatBool(p.Disabled)),
		g.Attr("data-file-upload-preview-value", strconv.FormatBool(p.Preview)),
		g.Attr("data-file-upload-max-size-value", strconv.FormatInt(p.MaxSize, 10)),
		g.Attr("data-file-upload-remove-label-value", removeLabel),
		g.Attr("data-file-upload-uploaded-label-value", uploadedLabel),
		g.Attr("data-file-upload-missing-error-value", missingError),
		g.Attr("data-file-upload-too-large-error-value", tooLargeError),
		g.Attr("data-file-upload-type-error-value", typeError),
	)
	if p.Disabled {
		root = append(root, g.Attr("aria-disabled", "true"))
	}
	if accept := strings.TrimSpace(p.Accept); accept != "" {
		root = append(root, g.Attr("data-file-upload-accept-value", accept))
	}
	if remote {
		root = append(root,
			g.Attr("data-file-upload-upload-url-value", strings.TrimSpace(p.UploadURL)),
			g.Attr("data-file-upload-upload-category-value", strings.TrimSpace(p.UploadCategory)),
		)
		if value := strings.TrimSpace(p.Value); value != "" {
			root = append(root, g.Attr("data-file-upload-current-url-value", value))
		}
		if currentName := strings.TrimSpace(p.CurrentName); currentName != "" {
			root = append(root, g.Attr("data-file-upload-current-name-value", currentName))
		}
	}
	if hasHTMXEnhancement(p.HTMXProps) {
		enhancement := p.HTMXProps
		if enhancement.Trigger == "" {
			enhancement.Trigger = "change from:[data-file-upload-input]"
		}
		if enhancement.Include == "" {
			enhancement.Include = "find [data-file-upload-input]"
		}
		root = append(root, htmxAttrs(enhancement)...)
		root = append(root, g.Attr("hx-encoding", "multipart/form-data"))
	}
	if label := strings.TrimSpace(p.Label); label != "" {
		root = append(root, Label(atoms.LabelProps{Text: label, For: inputID, Required: p.Required}))
	}
	root = append(root, g.Group(controls))
	if helpText := strings.TrimSpace(p.HelpText); helpText != "" {
		root = append(root, h.P(h.ID(inputID+"-help"), h.Class(clHelp.Compile()), g.Text(helpText)))
	}
	return h.Div(root...)
}

func fileUploadDefault(value, fallback string) string {
	if value = strings.TrimSpace(value); value != "" {
		return value
	}
	return fallback
}

func fileUploadTemplate(template string, values map[string]string) string {
	replacements := make([]string, 0, len(values)*2)
	for key, value := range values {
		replacements = append(replacements, "{"+key+"}", value)
	}
	return strings.NewReplacer(replacements...).Replace(template)
}

func fileUploadFormatBytes(bytes int64) string {
	const (
		kilobyte = 1024
		megabyte = 1024 * kilobyte
		gigabyte = 1024 * megabyte
	)
	switch {
	case bytes >= gigabyte:
		return fmt.Sprintf("%.1f GB", float64(bytes)/float64(gigabyte))
	case bytes >= megabyte:
		return fmt.Sprintf("%.1f MB", float64(bytes)/float64(megabyte))
	case bytes >= kilobyte:
		return fmt.Sprintf("%.1f KB", float64(bytes)/float64(kilobyte))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

// ModalForm closes its owning modal after a successful HTMX request.
func ModalForm(attrs ...g.Node) g.Node {
	nodes := make([]g.Node, 0, len(attrs)+1)
	nodes = append(nodes, attrs...)
	nodes = append(nodes, g.Attr("data-action", "htmx:afterRequest->htmx-modal#closeOnSuccess"))
	return h.Form(nodes...)
}

// Autocomplete renders the canonical progressively enhanced APG combobox.
// The visible query and submitted value are deliberately separate: selecting
// a human-readable option updates the hidden field consumed by forms. Static
// options and server-backed HTMX results share one controller DOM contract.
func Autocomplete(p molecules.AutocompleteProps) g.Node {
	minChars := p.MinChars
	if minChars <= 0 {
		minChars = 2
	}
	debounce := p.Debounce
	if debounce <= 0 {
		debounce = 300
	}
	name := strings.TrimSpace(p.Name)
	inputID := "autocomplete-input"
	hiddenID := "autocomplete-value"
	listboxID := "autocomplete-listbox"
	indicatorID := "autocomplete-indicator"
	if name != "" {
		inputID = name + "-input"
		hiddenID = name
		listboxID = name + "-listbox"
		indicatorID = name + "-indicator"
	}
	queryName := strings.TrimSpace(p.QueryName)
	if queryName == "" {
		queryName = "q"
	}

	invalid := p.Invalid || p.Error != ""
	inputClass := clInput.Merge(clInputTone["neutral"]).Merge(clInputSize["md"]).Merge(clInputPadStart)
	if invalid {
		inputClass = clInput.Merge(clInputError).Merge(clInputSize["md"]).Merge(clInputPadStart)
	}
	serverBacked := strings.TrimSpace(p.SearchURL) != "" || strings.TrimSpace(p.HTMXProps.Get) != ""
	if serverBacked {
		inputClass = inputClass.Merge(clInputPadEnd)
	}

	describedBy := make([]string, 0, 2)
	if p.Error != "" {
		describedBy = append(describedBy, inputID+"-error")
	}
	if p.HelpText != "" {
		describedBy = append(describedBy, inputID+"-help")
	}
	inputAttrs := []g.Node{
		h.ID(inputID), h.Type("text"), h.Name(queryName), h.Value(p.DisplayValue),
		h.Class(inputClass.Compile()), h.AutoComplete("off"),
		g.Attr("role", "combobox"), g.Attr("aria-autocomplete", "list"),
		g.Attr("aria-controls", listboxID), g.Attr("aria-expanded", "false"),
		g.Attr("data-autocomplete-input", "true"),
		g.Attr("data-action", "input->autocomplete#handleInput focus->autocomplete#focusInput keydown->autocomplete#handleKeydown"),
	}
	if p.Placeholder != "" {
		inputAttrs = append(inputAttrs, h.Placeholder(p.Placeholder))
	}
	if p.Required {
		inputAttrs = append(inputAttrs, h.Required())
	}
	if p.Disabled {
		inputAttrs = append(inputAttrs, h.Disabled())
	}
	if invalid {
		inputAttrs = append(inputAttrs, g.Attr("aria-invalid", "true"))
	}
	if len(describedBy) > 0 {
		inputAttrs = append(inputAttrs, g.Attr("aria-describedby", strings.Join(describedBy, " ")))
	}
	if p.Label == "" && name != "" {
		inputAttrs = append(inputAttrs, g.Attr("aria-label", name))
	}
	if serverBacked {
		enhancement := p.HTMXProps
		if enhancement.Get == "" {
			enhancement.Get = p.SearchURL
		}
		if enhancement.Trigger == "" {
			enhancement.Trigger = fmt.Sprintf("keyup changed delay:%dms", debounce)
		}
		if enhancement.Target == "" {
			enhancement.Target = "#" + listboxID
		}
		if enhancement.Swap == "" {
			enhancement.Swap = "innerHTML"
		}
		if enhancement.Indicator == "" {
			enhancement.Indicator = "#" + indicatorID
		}
		inputAttrs = append(inputAttrs, htmxAttrs(enhancement)...)
	}

	hiddenAttrs := []g.Node{
		h.ID(hiddenID), h.Type("hidden"), h.Name(name), h.Value(p.Value),
		g.Attr("data-autocomplete-hidden-input", "true"),
	}
	if p.Disabled {
		hiddenAttrs = append(hiddenAttrs, h.Disabled())
	}

	options := make([]g.Node, 0, len(p.Options))
	for index, option := range p.Options {
		label := strings.TrimSpace(option.Label)
		if label == "" {
			label = option.Value
		}
		value := strings.TrimSpace(option.Value)
		if value == "" {
			value = label
		}
		attrs := []g.Node{
			h.ID(fmt.Sprintf("%s-option-%d", listboxID, index)), h.Role("option"),
			g.Attr("tabindex", "-1"), g.Attr("aria-selected", "false"),
			g.Attr("data-autocomplete-option", "true"),
			g.Attr("data-autocomplete-label", label),
			g.Attr("data-autocomplete-value", value),
			h.Class(clAutocompleteOption.Compile()), g.Text(label),
		}
		if option.Disabled {
			attrs = append(attrs, g.Attr("aria-disabled", "true"))
		}
		options = append(options, h.Li(attrs...))
	}
	listboxLabel := strings.TrimSpace(p.Label)
	if listboxLabel == "" {
		listboxLabel = "Suggestions"
	}
	panel := h.Div(
		h.Class(clAutocompletePanel.Compile()), h.Style("max-height:15rem"), g.Attr("hidden"),
		g.Attr("aria-hidden", "true"), g.Attr("data-autocomplete-panel", "true"),
		h.Ul(
			h.ID(listboxID), h.Role("listbox"), h.Class(clAutocompleteList.Compile()),
			g.Attr("aria-label", listboxLabel+" suggestions"),
			g.Attr("data-autocomplete-listbox", "true"),
			g.Attr("data-action", "htmx:afterSwap->autocomplete#resultsLoaded click->autocomplete#selectFromEvent mousemove->autocomplete#highlightFromEvent"),
			g.Group(options),
		),
	)
	control := []g.Node{
		h.Class(clAutocompleteControl.Compile()),
		h.Input(hiddenAttrs...),
		h.Span(h.Class(clInputIconStart.Compile()), Icon(atoms.IconProps{Name: "search", Size: "md", Tone: "neutral"})),
		h.Input(inputAttrs...),
	}
	if serverBacked {
		control = append(control, h.Span(
			h.ID(indicatorID),
			h.Class("htmx-indicator "+clAutocompleteIndicator.Compile()),
			g.Attr("aria-hidden", "true"),
			h.Span(h.Class(clAutocompleteSpinner.Compile())),
		))
	}
	control = append(control, panel)

	rootProps := p.ComponentProps
	rootProps.Disabled = false
	rootProps.Class = ""
	root := baseAttrs(rootProps)
	root = append(root,
		classes(clFieldWrap.Compile(), p.Class),
		g.Attr("data-component", "autocomplete"),
		g.Attr("data-controller", "autocomplete"),
		g.Attr("data-autocomplete-min-chars-value", strconv.Itoa(minChars)),
		g.Attr("data-autocomplete-static-value", strconv.FormatBool(len(p.Options) > 0)),
		g.Attr("data-autocomplete-initial-query-value", p.DisplayValue),
		g.Attr("data-autocomplete-initial-value-value", p.Value),
	)
	if p.Label != "" {
		root = append(root, Label(atoms.LabelProps{Text: p.Label, For: inputID, Required: p.Required}))
	}
	root = append(root, h.Div(control...))
	if p.Error != "" {
		root = append(root, h.P(h.ID(inputID+"-error"), h.Class(clFieldErr.Compile()), h.Role("alert"), g.Text(p.Error)))
	}
	if p.HelpText != "" {
		root = append(root, h.P(h.ID(inputID+"-help"), h.Class(clHelp.Compile()), g.Text(p.HelpText)))
	}
	return h.Div(root...)
}

// Dropdown renders the canonical progressively enhanced select-only listbox.
// The submitted hidden input is present without JavaScript; the shared
// dropdown controller adds disclosure, filtering, keyboard navigation, and
// single- or multi-selection behavior.
func Dropdown(p molecules.DropdownProps) g.Node {
	placeholder := strings.TrimSpace(p.Placeholder)
	if placeholder == "" {
		placeholder = "Select an option"
	}
	searchLabel := strings.TrimSpace(p.SearchLabel)
	if searchLabel == "" {
		searchLabel = "Search options"
	}
	clearLabel := strings.TrimSpace(p.ClearLabel)
	if clearLabel == "" {
		clearLabel = "Clear selection"
	}
	openLabel := strings.TrimSpace(p.OpenLabel)
	if openLabel == "" {
		openLabel = "Open options"
	}
	size := strings.ToLower(strings.TrimSpace(p.Size))
	switch size {
	case "small":
		size = "sm"
	case "medium", "":
		size = "md"
	case "large":
		size = "lg"
	}
	if _, ok := clDropdownTriggerSize[size]; !ok {
		size = "md"
	}

	selected := dropdownSelectedValues(p)
	encodedSelected, _ := json.Marshal(selected)
	triggerText := dropdownTriggerText(placeholder, p.Multiple, selected, p.Options)
	accessibleName := strings.TrimSpace(p.AriaLabel)
	if accessibleName == "" {
		accessibleName = strings.TrimSpace(p.Label)
	}
	if accessibleName == "" {
		accessibleName = strings.TrimSpace(p.Name)
	}
	if accessibleName == "" {
		accessibleName = placeholder
	}

	panelID := ""
	if id := strings.TrimSpace(p.ID); id != "" {
		panelID = id + "-panel"
	}
	trigger := []g.Node{
		h.Type("button"),
		h.Class(clDropdownButton.Merge(clDropdownButtonSize[size]).Compile()),
		g.Attr("data-dropdown-button", "true"),
		g.Attr("data-action", "click->dropdown#toggle"),
		g.Attr("aria-haspopup", "listbox"),
		g.Attr("aria-expanded", "false"),
		g.Attr("aria-label", accessibleName),
		g.Attr("data-state", "closed"),
	}
	if panelID != "" {
		trigger = append(trigger, g.Attr("aria-controls", panelID))
	}
	if p.Disabled {
		trigger = append(trigger, h.Disabled())
	}
	trigger = append(trigger, h.Span(
		h.Class(clDropdownTriggerLabel.Compile()),
		g.Attr("data-dropdown-trigger-label", "true"),
		g.Text(triggerText),
	))

	actions := make([]g.Node, 0, 2)
	if p.Clearable {
		clear := []g.Node{
			h.Type("button"), h.Class(clDropdownIconButton.Compile()),
			g.Attr("aria-label", clearLabel), g.Attr("tabindex", "-1"),
			g.Attr("data-dropdown-clear-button", "true"),
			g.Attr("data-action", "click->dropdown#clear"),
		}
		if len(selected) == 0 {
			clear = append(clear, g.Attr("hidden"))
		}
		if p.Disabled {
			clear = append(clear, h.Disabled())
		}
		clear = append(clear, Icon(atoms.IconProps{Name: "x-mark", Size: "sm", Tone: "neutral"}))
		actions = append(actions, h.Button(clear...))
	}
	chevron := []g.Node{
		h.Type("button"), h.Class(clDropdownChevron.Compile()),
		g.Attr("aria-label", openLabel), g.Attr("tabindex", "-1"),
		g.Attr("data-dropdown-chevron-button", "true"),
		g.Attr("data-action", "click->dropdown#toggle"),
	}
	if p.Disabled {
		chevron = append(chevron, h.Disabled())
	}
	chevron = append(chevron, Icon(atoms.IconProps{Name: "chevron-down", Size: "sm", Tone: "neutral"}))
	actions = append(actions, h.Button(chevron...))

	panel := []g.Node{
		h.Class(clDropdownPanel.Compile()),
		g.Attr("data-dropdown-panel", "true"),
		g.Attr("data-state", "closed"),
		g.Attr("hidden"),
		h.Role("listbox"),
		g.Attr("aria-hidden", "true"),
		g.Attr("aria-multiselectable", strconv.FormatBool(p.Multiple)),
	}
	if panelID != "" {
		panel = append(panel, h.ID(panelID))
	}
	if p.Searchable {
		panel = append(panel, h.Div(
			h.Class(clDropdownSearchWrap.Compile()),
			h.Input(
				h.Type("search"), h.Class(clDropdownSearch.Compile()),
				h.Placeholder(searchLabel), g.Attr("aria-label", searchLabel),
				g.Attr("autocomplete", "off"),
				g.Attr("data-dropdown-search-input", "true"),
				g.Attr("data-action", "input->dropdown#filterOptions"),
			),
		))
	}
	panel = append(panel, dropdownOptionsList(p.Options, selected))

	rootProps := p.ComponentProps
	rootProps.Disabled = false
	rootProps.Class = ""
	root := baseAttrs(rootProps)
	root = append(root,
		classes(clDropdownRoot.Compile(), p.Class),
		g.Attr("data-component", "dropdown"),
		g.Attr("data-controller", "dropdown"),
		g.Attr("data-dropdown-contract", "1"),
		g.Attr("data-dropdown-kind-value", "select"),
		g.Attr("data-dropdown-multiple-value", strconv.FormatBool(p.Multiple)),
		g.Attr("data-dropdown-searchable-value", strconv.FormatBool(p.Searchable)),
		g.Attr("data-dropdown-placeholder-value", placeholder),
		g.Attr("data-dropdown-selected-value", string(encodedSelected)),
		g.Attr("data-state", "closed"),
	)
	if p.Name != "" {
		root = append(root, g.Attr("data-dropdown-input-name-value", p.Name))
	}
	if p.Disabled {
		root = append(root, g.Attr("aria-disabled", "true"))
	}
	enhancement := p.HTMXProps
	if enhancement.Trigger == "" && (enhancement.Get != "" || enhancement.Post != "" || enhancement.Put != "" || enhancement.Patch != "" || enhancement.Delete != "") {
		enhancement.Trigger = "dropdown:changed"
	}
	root = append(root, htmxAttrs(enhancement)...)
	if p.Label != "" {
		root = append(root, h.Div(h.Class(clLabel.Compile()), g.Text(p.Label)))
	}
	root = append(root, dropdownHiddenInputs(p.Name, p.Multiple, p.Disabled, selected)...)
	root = append(root,
		h.Div(
			h.Class(clDropdownTrigger.Merge(clDropdownTriggerSize[size]).Compile()),
			h.Button(trigger...),
			h.Div(h.Class(clDropdownTriggerActions.Compile()), g.Group(actions)),
		),
		h.Div(panel...),
	)
	return h.Div(root...)
}

func dropdownSelectedValues(p molecules.DropdownProps) []string {
	candidates := p.Selected
	if len(candidates) == 0 && strings.TrimSpace(p.Value) != "" {
		candidates = []string{p.Value}
	}
	selected := make([]string, 0, len(candidates))
	seen := make(map[string]struct{}, len(candidates))
	for _, candidate := range candidates {
		value := strings.TrimSpace(candidate)
		if value == "" {
			continue
		}
		if _, duplicate := seen[value]; duplicate {
			continue
		}
		seen[value] = struct{}{}
		selected = append(selected, value)
		if !p.Multiple {
			break
		}
	}
	return selected
}

func dropdownTriggerText(placeholder string, multiple bool, selected []string, options []molecules.Option) string {
	if len(selected) == 0 {
		return placeholder
	}
	labels := make(map[string]string, len(options))
	for _, option := range options {
		value := strings.TrimSpace(option.Value)
		if value == "" {
			continue
		}
		label := strings.TrimSpace(option.Label)
		if label == "" {
			label = value
		}
		labels[value] = label
	}
	visible := make([]string, 0, len(selected))
	for _, value := range selected {
		label := labels[value]
		if label == "" {
			label = value
		}
		visible = append(visible, label)
	}
	if multiple {
		return strings.Join(visible, ", ")
	}
	return visible[0]
}

func dropdownHiddenInputs(name string, multiple, disabled bool, selected []string) []g.Node {
	if strings.TrimSpace(name) == "" {
		return nil
	}
	if multiple {
		inputs := make([]g.Node, 0, len(selected))
		for _, value := range selected {
			attrs := []g.Node{
				h.Type("hidden"), h.Name(name), h.Value(value),
				g.Attr("data-dropdown-hidden-value", "true"),
			}
			if disabled {
				attrs = append(attrs, h.Disabled())
			}
			inputs = append(inputs, h.Input(attrs...))
		}
		return []g.Node{h.Div(g.Attr("data-dropdown-hidden-inputs", "true"), g.Group(inputs))}
	}
	value := ""
	if len(selected) > 0 {
		value = selected[0]
	}
	attrs := []g.Node{
		h.Type("hidden"), h.Name(name), h.Value(value),
		g.Attr("data-dropdown-hidden-input", "true"),
	}
	if disabled {
		attrs = append(attrs, h.Disabled())
	}
	return []g.Node{h.Input(attrs...)}
}

type dropdownOptionGroup struct {
	name    string
	options []molecules.Option
}

func dropdownOptionsList(options []molecules.Option, selected []string) g.Node {
	selectedSet := make(map[string]struct{}, len(selected))
	for _, value := range selected {
		selectedSet[value] = struct{}{}
	}
	groups := make([]dropdownOptionGroup, 0)
	indexes := make(map[string]int)
	for _, option := range options {
		group := strings.TrimSpace(option.Group)
		index, exists := indexes[group]
		if !exists {
			index = len(groups)
			indexes[group] = index
			groups = append(groups, dropdownOptionGroup{name: group})
		}
		groups[index].options = append(groups[index].options, option)
	}
	nodes := make([]g.Node, 0, len(groups))
	for _, group := range groups {
		options := make([]g.Node, 0, len(group.options)+1)
		groupAttrs := []g.Node{h.Role("presentation")}
		if group.name != "" {
			groupAttrs = []g.Node{h.Role("group"), g.Attr("aria-label", group.name)}
			options = append(options, h.Div(
				h.Class(clDropdownGroupLabel.Compile()), h.Role("presentation"), g.Text(group.name),
			))
		}
		for _, option := range group.options {
			_, isSelected := selectedSet[strings.TrimSpace(option.Value)]
			options = append(options, dropdownOption(option, isSelected))
		}
		nodes = append(nodes, h.Div(append(groupAttrs, g.Group(options))...))
	}
	return h.Div(h.Class(clDropdownOptions.Compile()), h.Role("presentation"), g.Group(nodes))
}

func dropdownOption(option molecules.Option, selected bool) g.Node {
	value := strings.TrimSpace(option.Value)
	label := strings.TrimSpace(option.Label)
	if label == "" {
		label = value
	}
	optionClass := clDropdownOption
	if selected {
		optionClass = optionClass.Merge(clDropdownOptionSelected)
	}
	if option.Disabled {
		optionClass = optionClass.Merge(clDropdownOptionDisabled)
	}
	attrs := []g.Node{
		h.Type("button"), h.Role("option"), g.Attr("tabindex", "-1"),
		h.Class(optionClass.Compile()),
		g.Attr("data-dropdown-option", value),
		g.Attr("data-dropdown-label", label),
		g.Attr("data-dropdown-search", strings.ToLower(label)),
		g.Attr("data-action", "click->dropdown#choose keydown.enter->dropdown#choose keydown.space->dropdown#choose"),
		g.Attr("aria-selected", strconv.FormatBool(selected)),
	}
	if option.Disabled {
		attrs = append(attrs, h.Disabled(), g.Attr("aria-disabled", "true"))
	}
	mark := []g.Node{
		h.Class(clDropdownOptionMark.Compile()),
		g.Attr("data-dropdown-checkmark", "true"),
	}
	if !selected {
		mark = append(mark, g.Attr("hidden"))
	}
	mark = append(mark, Icon(atoms.IconProps{Name: "check", Size: "sm", Tone: "brand"}))
	attrs = append(attrs, h.Span(mark...))
	spacer := []g.Node{
		h.Class(clDropdownOptionSpacer.Compile()),
		g.Attr("data-dropdown-spacer", "true"),
	}
	if selected {
		spacer = append(spacer, g.Attr("hidden"))
	}
	attrs = append(attrs, h.Span(spacer...))
	if iconName := strings.TrimSpace(option.Icon); iconName != "" {
		attrs = append(attrs, h.Span(
			h.Class(clDropdownOptionIcon.Compile()),
			g.Attr("data-dropdown-option-icon", iconName),
			Icon(atoms.IconProps{Name: iconName, Size: "sm", Tone: "neutral"}),
		))
	}
	attrs = append(attrs, h.Span(h.Class(clDropdownOptionLabel.Compile()), g.Text(label)))
	return h.Button(attrs...)
}

// DatePicker renders the canonical native date field. Browser-native date
// selection remains the progressive enhancement; the shared renderer owns the
// accessible label, validation, constraints, calendar adornment, and HTMX
// request contract.
func DatePicker(p molecules.DatePickerProps) g.Node {
	placeholder := strings.TrimSpace(p.Placeholder)
	if placeholder == "" {
		placeholder = "Select a date"
	}
	enhancement := p.HTMXProps
	if enhancement.Trigger == "" && (enhancement.Get != "" || enhancement.Post != "") {
		enhancement.Trigger = "change"
	}
	componentProps := p.ComponentProps
	if format := strings.TrimSpace(p.Format); format != "" {
		attrs := make(map[string]string, len(componentProps.Attrs)+1)
		for name, value := range componentProps.Attrs {
			attrs[name] = value
		}
		attrs["data-date-format"] = format
		componentProps.Attrs = attrs
	}
	return inputFieldWithSlots(
		"datepicker",
		atoms.InputProps{
			ComponentProps: componentProps,
			HTMXProps:      enhancement,
			Name:           p.Name,
			Type:           "date",
			Value:          p.Value,
			Placeholder:    placeholder,
			Label:          p.Label,
			HelpText:       p.HelpText,
			Error:          p.Error,
			Invalid:        p.Invalid,
			Required:       p.Required,
			Min:            p.Min,
			Max:            p.Max,
		},
		[]g.Node{Icon(atoms.IconProps{Name: "calendar", Size: "md", Tone: "neutral"})},
		nil,
	)
}

// TableSlots is the trusted Go composition seam for rich web cells and
// server-driven sorting. Portable table data remains in TableProps; callers
// only opt into these callbacks when a cell needs real markup.
type TableSlots struct {
	Cell             func(molecules.TableRow, molecules.TableColumn) g.Node
	CellAttrs        func(molecules.TableRow, molecules.TableColumn) []g.Node
	RowAttrs         func(molecules.TableRow) []g.Node
	SortURL          func(molecules.TableColumn) string
	SortState        func(molecules.TableColumn) string
	SortButtonAttrs  func(molecules.TableColumn) []g.Node
	SelectAllLabel   string
	SelectRowLabel   func(molecules.TableRow) string
	SelectRowChecked func(molecules.TableRow) bool
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
			if slots.SortButtonAttrs != nil {
				button = append(button, slots.SortButtonAttrs(c)...)
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
		if slots.RowAttrs != nil {
			cells = append(cells, slots.RowAttrs(r)...)
		}
		if p.Selectable {
			label := "Select row"
			if slots.SelectRowLabel != nil {
				label = fallbackText(slots.SelectRowLabel(r), label)
			}
			input := []g.Node{h.Class(clCheckbox.Compile()), h.Type("checkbox"),
				g.Attr("data-pk-select", r.ID), g.Attr("aria-label", label)}
			if slots.SelectRowChecked != nil && slots.SelectRowChecked(r) {
				input = append(input, h.Checked())
			}
			cells = append(cells, h.Td(h.Class(tdClass.Compile()), h.Input(input...)))
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
			if slots.CellAttrs != nil {
				td = append(td, slots.CellAttrs(r, c)...)
			}
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
	wrap = append(wrap,
		classes(clTableWrap.Compile(), p.Class),
		g.Attr("data-component", "table"),
		h.Table(
			h.Class(clTable.Compile()),
			h.THead(head...),
			h.TBody(bodyRows...),
		),
	)
	return h.Div(wrap...)
}

// CardSlots is the trusted Go composition seam for the three structural card
// regions. Portable title/description/media data remains in CardProps.
type CardSlots struct {
	Header  []g.Node
	Content []g.Node
	Footer  []g.Node
}

// Card renders free-form card children. Call CardWithSlots when header,
// content, and footer need the canonical section treatment.
func Card(p molecules.CardProps, children ...g.Node) g.Node {
	return cardNode(p, CardSlots{}, children)
}

// CardWithSlots renders canonical header/content/footer regions without a
// downstream wrapper or style implementation.
func CardWithSlots(p molecules.CardProps, slots CardSlots) g.Node {
	return cardNode(p, slots, nil)
}

func cardNode(p molecules.CardProps, slots CardSlots, children []g.Node) g.Node {
	sectioned := len(slots.Header)+len(slots.Content)+len(slots.Footer) > 0
	rootClass := cardRootClasses(p, sectioned)
	nodes := baseAttrs(p.ComponentProps, htmxAttrs(p.HTMXProps)...)
	nodes = append(nodes,
		classes(rootClass, p.Class),
		g.Attr("data-component", "card"),
	)

	body := make([]g.Node, 0, 5)
	if p.Image != "" && normalizedCardImagePosition(p.ImagePosition) == "top" {
		body = append(body, cardImage(p, false))
	}
	if sectioned {
		padding := cardPaddingClasses(p.Padding, true)
		header := slots.Header
		if len(header) == 0 && (p.Title != "" || p.Description != "") {
			header = cardTextHeader(p)
		}
		if len(header) > 0 {
			body = append(body, h.Div(
				h.Class(strings.TrimSpace(padding+" "+clCardHeader.Compile())),
				g.Group(header),
			))
		}
		if len(slots.Content) > 0 {
			body = append(body, h.Div(h.Class(padding), g.Group(slots.Content)))
		}
		if len(slots.Footer) > 0 {
			body = append(body, h.Div(
				h.Class(strings.TrimSpace(padding+" "+clCardFooter.Compile())),
				g.Group(slots.Footer),
			))
		}
	} else {
		body = append(body, cardTextHeader(p)...)
		body = append(body, children...)
	}
	if p.Image != "" && normalizedCardImagePosition(p.ImagePosition) == "bottom" {
		body = append(body, cardImage(p, false))
	}

	position := normalizedCardImagePosition(p.ImagePosition)
	if p.Image != "" && (position == "left" || position == "right") {
		content := h.Div(h.Class(clCardVertical.Compile()), g.Group(body))
		image := cardImage(p, true)
		if position == "left" {
			body = []g.Node{h.Div(h.Class(clCardHorizontal.Compile()), image, content)}
		} else {
			body = []g.Node{h.Div(h.Class(clCardHorizontal.Compile()), content, image)}
		}
	}

	nodes = append(nodes, body...)
	if p.Clickable && p.Href != "" {
		return h.A(append(nodes, h.Href(p.Href))...)
	}
	return h.Article(nodes...)
}

func cardTextHeader(p molecules.CardProps) []g.Node {
	var nodes []g.Node
	if p.Title != "" {
		nodes = append(nodes, h.P(h.Class(clCardTitle.Compile()), g.Text(p.Title)))
	}
	if p.Description != "" {
		nodes = append(nodes, h.P(h.Class(clCardDesc.Compile()), g.Text(p.Description)))
	}
	return nodes
}

func cardImage(p molecules.CardProps, horizontal bool) g.Node {
	className := clCardImageVertical.Compile()
	if horizontal {
		className = clCardImageHorizontal.Compile()
	}
	return h.Img(h.Src(p.Image), h.Alt(p.ImageAlt), h.Class(className))
}

func normalizedCardImagePosition(position string) string {
	switch strings.ToLower(strings.TrimSpace(position)) {
	case "bottom", "left", "right":
		return strings.ToLower(strings.TrimSpace(position))
	default:
		return "top"
	}
}

func cardPaddingClasses(padding string, sectioned bool) string {
	switch strings.ToLower(strings.TrimSpace(padding)) {
	case "none":
		return clCardPadNone.Compile()
	case "small":
		return clCardPadSmall.Compile()
	case "large":
		return clCardPadLarge.Compile()
	case "medium":
		return clCardPadMedium.Compile()
	default:
		if sectioned {
			return clCardPadMedium.Compile()
		}
		return clCardPadDefault.Compile()
	}
}

func cardRootClasses(p molecules.CardProps, sectioned bool) string {
	cl := clCardFrame
	if sectioned {
		cl = cl.Merge(clCardSectioned)
	} else {
		switch strings.ToLower(strings.TrimSpace(p.Padding)) {
		case "none":
			cl = cl.Merge(clCardPadNone)
		case "small":
			cl = cl.Merge(clCardPadSmall)
		case "medium":
			cl = cl.Merge(clCardPadMedium)
		case "large":
			cl = cl.Merge(clCardPadLarge)
		default:
			cl = cl.Merge(clCardPadDefault)
		}
	}

	variant := strings.ToLower(strings.TrimSpace(p.Variant))
	if variant != "elevated" && variant != "plain" {
		cl = cl.Merge(clCardBorder)
	}
	shadow := strings.ToLower(strings.TrimSpace(p.Shadow))
	if shadow == "" {
		switch variant {
		case "outlined", "plain":
			shadow = "none"
		case "elevated":
			shadow = "medium"
		default:
			shadow = "small"
		}
	}
	switch shadow {
	case "medium":
		cl = cl.Merge(clCardShadowMedium)
	case "large":
		cl = cl.Merge(clCardShadowLarge)
	case "none":
	default:
		cl = cl.Merge(clCardShadowSmall)
	}
	if p.Hoverable {
		cl = cl.Merge(clCardHoverable)
	}
	if p.Clickable {
		cl = cl.Merge(clCardClickable)
	}
	return cl.Compile()
}

// AccordionSection is one rich-content disclosure projected into an
// Accordion. Portable callers can use AccordionProps.Items; trusted Go and
// delivery-slot callers use Content to preserve arbitrary child composition.
type AccordionSection struct {
	ID          string
	Title       string
	Subtitle    string
	Icon        string
	Content     []g.Node
	DefaultOpen bool
	Open        bool
	Disabled    bool
}

// AccordionSlots carries the ordered sections projected into an Accordion.
type AccordionSlots struct {
	Sections []AccordionSection
}

// Accordion renders portable string-backed sections from AccordionProps.
func Accordion(p molecules.AccordionProps) g.Node {
	sections := make([]AccordionSection, 0, len(p.Items))
	for _, item := range p.Items {
		id := item.ID
		if strings.TrimSpace(id) == "" {
			id = item.Key
		}
		var content []g.Node
		if item.Content != "" {
			content = []g.Node{g.Text(item.Content)}
		}
		sections = append(sections, AccordionSection{
			ID:          id,
			Title:       item.Title,
			Subtitle:    item.Subtitle,
			Icon:        item.Icon,
			Content:     content,
			DefaultOpen: item.DefaultOpen,
			Open:        item.Open,
			Disabled:    item.Disabled,
		})
	}
	return AccordionWithSlots(p, AccordionSlots{Sections: sections})
}

// AccordionWithSections is the concise API for rich panel content.
func AccordionWithSections(p molecules.AccordionProps, sections ...AccordionSection) g.Node {
	return AccordionWithSlots(p, AccordionSlots{Sections: sections})
}

// AccordionWithSlots renders the canonical accessible disclosure structure
// and the controller contract shared by server-rendered application surfaces.
func AccordionWithSlots(p molecules.AccordionProps, slots AccordionSlots) g.Node {
	sections := normalizedAccordionSections(slots.Sections)
	if len(sections) == 0 {
		return nil
	}

	defaultOpen := accordionDefaultOpen(p.DefaultOpen, sections)
	openSet := make(map[string]struct{}, len(defaultOpen))
	for _, id := range defaultOpen {
		openSet[id] = struct{}{}
	}
	encodedOpen, _ := json.Marshal(defaultOpen)

	rootClass := clAccordionRoot
	if !p.Flush {
		if p.Bordered == nil || *p.Bordered {
			rootClass = rootClass.Merge(clAccordionBordered)
		} else {
			rootClass = rootClass.Merge(clAccordionUnbordered)
		}
	}
	rootID := strings.TrimSpace(p.ID)
	if rootID == "" {
		rootID = "accordion"
	}

	children := make([]g.Node, 0, len(sections)+6)
	rootProps := p.ComponentProps
	rootProps.Disabled = false
	children = append(children,
		baseAttrs(rootProps)...,
	)
	children = append(children,
		classes(rootClass.Compile(), p.Class),
		g.Attr("data-component", "accordion"),
		g.Attr("data-controller", "accordion"),
		g.Attr("data-accordion-multiple-value", boolText(p.Multiple)),
		g.Attr("data-accordion-open-items-value", string(encodedOpen)),
	)
	for index, section := range sections {
		_, open := openSet[section.ID]
		children = append(children, accordionSectionNode(rootID, section, index, open, p.Disabled))
	}
	return h.Div(children...)
}

func normalizedAccordionSections(sections []AccordionSection) []AccordionSection {
	normalized := make([]AccordionSection, len(sections))
	used := make(map[string]int, len(sections))
	for index, section := range sections {
		baseID := strings.TrimSpace(section.ID)
		if baseID == "" {
			baseID = fmt.Sprintf("item-%d", index+1)
		}
		used[baseID]++
		section.ID = baseID
		if used[baseID] > 1 {
			section.ID = fmt.Sprintf("%s-%d", baseID, used[baseID])
		}
		normalized[index] = section
	}
	return normalized
}

func accordionDefaultOpen(explicit []string, sections []AccordionSection) []string {
	seen := make(map[string]bool, len(explicit)+len(sections))
	open := make([]string, 0, len(explicit)+len(sections))
	appendID := func(id string) {
		id = strings.TrimSpace(id)
		if id == "" || seen[id] {
			return
		}
		seen[id] = true
		open = append(open, id)
	}
	for _, id := range explicit {
		appendID(id)
	}
	for _, section := range sections {
		if section.DefaultOpen || section.Open {
			appendID(section.ID)
		}
	}
	return open
}

func accordionSectionNode(
	rootID string,
	section AccordionSection,
	index int,
	open bool,
	rootDisabled bool,
) g.Node {
	headerID := rootID + "-header-" + section.ID
	panelID := rootID + "-panel-" + section.ID
	disabled := rootDisabled || section.Disabled
	state := "closed"
	if open {
		state = "open"
	}

	itemAttrs := []g.Node{
		g.Attr("data-accordion-item", section.ID),
		g.Attr("data-state", state),
	}
	if index > 0 {
		itemAttrs = append(itemAttrs, h.Class(clAccordionSeparator.Compile()))
	}

	titleBlock := []g.Node{
		h.Span(h.Class(clAccordionTitle.Compile()), g.Text(section.Title)),
	}
	if section.Subtitle != "" {
		titleBlock = append(titleBlock,
			h.Span(h.Class(clAccordionSubtitle.Compile()), g.Text(section.Subtitle)),
		)
	}
	lead := []g.Node{h.Class(clAccordionLead.Compile())}
	if section.Icon != "" {
		lead = append(lead,
			h.Span(h.Class(clAccordionItemIcon.Compile()), icon(section.Icon)),
		)
	}
	lead = append(lead, h.Span(h.Class(clAccordionTitleBlock.Compile()), g.Group(titleBlock)))

	chevronClass := clAccordionChevron
	if open {
		chevronClass = chevronClass.Merge(clAccordionChevronOpen)
	}
	triggerClass := clAccordionTrigger
	if disabled {
		triggerClass = triggerClass.Merge(clAccordionTriggerDisabled)
	}
	trigger := []g.Node{
		h.Type("button"),
		h.ID(headerID),
		h.Class(triggerClass.Compile()),
		g.Attr("aria-controls", panelID),
		g.Attr("aria-expanded", boolText(open)),
		g.Attr("aria-disabled", boolText(disabled)),
		g.Attr("data-accordion-id", section.ID),
		g.Attr("data-action", "click->accordion#toggle keydown.enter->accordion#toggle keydown.space->accordion#toggle"),
	}
	if disabled {
		trigger = append(trigger, h.Disabled())
	}
	trigger = append(trigger,
		h.Span(lead...),
		h.Span(
			h.Class(chevronClass.Compile()),
			g.Attr("data-accordion-icon", section.ID),
			icon("chevron-down"),
		),
	)

	panel := []g.Node{
		h.ID(panelID),
		h.Role("region"),
		g.Attr("aria-labelledby", headerID),
		g.Attr("data-accordion-panel", section.ID),
		g.Attr("aria-hidden", boolText(!open)),
	}
	if !open {
		panel = append(panel, g.Attr("hidden"))
	}
	panel = append(panel, h.Div(h.Class(clAccordionPanel.Compile()), g.Group(section.Content)))

	itemAttrs = append(itemAttrs,
		h.H3(h.Button(trigger...)),
		h.Div(panel...),
	)
	return h.Div(itemAttrs...)
}

// Stepper renders a progress-navigation landmark from the portable Stepper
// contract. Statuses omitted (or left pending) are derived from CurrentStep;
// explicit active, completed, and error states remain authoritative.
func Stepper(p molecules.StepperProps) g.Node {
	steps := normalizedStepperSteps(p.Steps, p.CurrentStep)
	if len(steps) == 0 {
		return nil
	}

	orientation := "horizontal"
	if p.Orientation == "vertical" {
		orientation = "vertical"
	}
	label := fallbackText(p.NavigationLabel, "Progress")

	rootProps := p.ComponentProps
	rootProps.Disabled = false
	attrs := append(baseAttrs(rootProps),
		g.Attr("data-component", "stepper"),
		g.Attr("data-orientation", orientation),
		g.Attr("aria-label", label),
		g.Attr("aria-disabled", boolText(p.Disabled)),
	)
	if strings.TrimSpace(p.Class) != "" {
		attrs = append(attrs, h.Class(p.Class))
	}

	listClass := clStepperListHorizontal
	if orientation == "vertical" {
		listClass = clStepperListVertical
	}
	items := make([]g.Node, 0, len(steps))
	for index, step := range steps {
		items = append(items, stepperStepNode(p, step, index, len(steps), orientation))
	}

	attrs = append(attrs, h.Ol(h.Class(listClass.Compile()), g.Group(items)))
	return h.Nav(attrs...)
}

func normalizedStepperSteps(steps []molecules.StepItem, currentStep int) []molecules.StepItem {
	if len(steps) == 0 {
		return nil
	}
	if currentStep < 0 {
		currentStep = 0
	}
	if currentStep >= len(steps) {
		currentStep = len(steps) - 1
	}
	normalized := make([]molecules.StepItem, len(steps))
	for index, step := range steps {
		status := strings.ToLower(strings.TrimSpace(step.Status))
		switch status {
		case "active", "completed", "error":
		default:
			switch {
			case index < currentStep:
				status = "completed"
			case index == currentStep:
				status = "active"
			default:
				status = "pending"
			}
		}
		step.Status = status
		normalized[index] = step
	}
	return normalized
}

func stepperStepNode(
	p molecules.StepperProps,
	step molecules.StepItem,
	index int,
	total int,
	orientation string,
) g.Node {
	isLast := index == total-1
	itemClass := clStepperItemHorizontal
	if isLast {
		itemClass = clStepperItemLast
	}
	if orientation == "vertical" {
		itemClass = clStepperVerticalItem
	}
	itemAttrs := append(stepperDataAttrs(step, index), h.Class(itemClass.Compile()))
	if step.Status == "active" {
		itemAttrs = append(itemAttrs, g.Attr("aria-current", "step"))
	}

	if orientation == "vertical" {
		icons := []g.Node{
			h.Class(clStepperVerticalIcons.Compile()),
			stepperIndicator(p, step, index),
		}
		if !isLast {
			icons = append(icons, stepperConnector(step.Status, true))
		}
		text := []g.Node{
			h.Class(clStepperVerticalText.Compile()),
			stepperLabel(p.Compact, step, false),
		}
		if step.Description != "" {
			text = append(text, h.P(
				h.Class(clStepperDescription.Merge(clStepperDescriptionState[step.Status]).Compile()),
				g.Attr("data-stepper-label-description", ""),
				g.Text(step.Description),
			))
		}
		row := stepperRow(p, step, index,
			clStepperRowVertical,
			h.Div(icons...),
			h.Div(text...),
		)
		return h.Li(append(itemAttrs, row)...)
	}

	row := stepperRow(p, step, index,
		clStepperRowHorizontal,
		stepperIndicator(p, step, index),
		stepperLabel(p.Compact, step, true),
	)
	children := []g.Node{row}
	if !isLast {
		children = append(children, stepperConnector(step.Status, false))
	}
	return h.Li(append(itemAttrs, children...)...)
}

func stepperRow(
	p molecules.StepperProps,
	step molecules.StepItem,
	index int,
	class tw.ClassList,
	children ...g.Node,
) g.Node {
	row := append([]g.Node{h.Class(class.Compile())}, children...)
	if strings.TrimSpace(p.StepAction) == "" {
		return h.Div(row...)
	}
	button := []g.Node{
		h.Type("button"),
		h.Class(class.Compile()),
		g.Attr("data-action", strings.TrimSpace(p.StepAction)),
		g.Attr("aria-label", fmt.Sprintf("Go to step %d: %s", index+1, step.Label)),
	}
	button = append(button, stepperDataAttrs(step, index)...)
	if p.Disabled {
		button = append(button, h.Disabled(), g.Attr("aria-disabled", "true"))
	}
	button = append(button, children...)
	return h.Button(button...)
}

func stepperIndicator(p molecules.StepperProps, step molecules.StepItem, index int) g.Node {
	indicatorClass := clStepperIndicator.Merge(clStepperIndicatorState[step.Status])
	if p.Compact {
		indicatorClass = indicatorClass.Merge(clStepperIndicatorCompact)
	} else {
		indicatorClass = indicatorClass.Merge(clStepperIndicatorRegular)
	}
	clickable := p.Clickable && strings.TrimSpace(p.StepAction) == ""
	if clickable {
		indicatorClass = indicatorClass.Merge(clStepperClickable)
	}
	if p.Disabled {
		indicatorClass = indicatorClass.Merge(clStepperDisabled)
	}

	content := stepperIndicatorContent(step, index)
	attrs := []g.Node{
		h.Class(indicatorClass.Compile()),
		g.Attr("data-stepper-indicator", ""),
	}
	if clickable {
		attrs = append(attrs,
			h.Type("button"),
			g.Attr("data-step", itoa(index)),
			g.Attr("aria-label", fmt.Sprintf("Go to step %d: %s", index+1, step.Label)),
		)
		attrs = append(attrs, stepperDataAttrs(step, index)...)
		if p.Disabled {
			attrs = append(attrs, h.Disabled(), g.Attr("aria-disabled", "true"))
		}
		attrs = append(attrs, content)
		return h.Button(attrs...)
	}
	attrs = append(attrs, g.Attr("aria-hidden", "true"), content)
	return h.Span(attrs...)
}

func stepperIndicatorContent(step molecules.StepItem, index int) g.Node {
	switch step.Status {
	case "completed":
		return Icon(atoms.IconProps{Name: "check", Size: "sm"})
	case "error":
		return Icon(atoms.IconProps{Name: "x-mark", Size: "sm"})
	default:
		if validIconName(step.Icon) {
			return Icon(atoms.IconProps{Name: step.Icon, Size: "sm"})
		}
		return h.Span(h.Class(clStepperGlyph.Compile()), g.Text(itoa(index+1)))
	}
}

func stepperLabel(compact bool, step molecules.StepItem, inlineDescription bool) g.Node {
	labelClass := clStepperLabel
	if compact {
		labelClass = clStepperLabelCompact
	}
	labelClass = labelClass.Merge(clStepperLabelState[step.Status])
	attrs := []g.Node{
		h.Class(labelClass.Compile()),
		g.Attr("data-stepper-label", ""),
	}
	if inlineDescription && step.Description != "" {
		attrs = append(attrs,
			h.Span(
				h.Class(clStepperLabelBlock.Compile()),
				g.Attr("data-stepper-label-text", ""),
				g.Text(step.Label),
			),
			h.Span(
				h.Class(clStepperDescription.Merge(clStepperDescriptionState[step.Status]).Compile()),
				g.Attr("data-stepper-label-description", ""),
				g.Text(step.Description),
			),
		)
		return h.Div(attrs...)
	}
	attrs = append(attrs, g.Text(step.Label))
	return h.Span(attrs...)
}

func stepperConnector(status string, vertical bool) g.Node {
	connectorStatus := "pending"
	if status == "completed" {
		connectorStatus = "completed"
	}
	class := clStepperConnectorHorizontal
	if vertical {
		class = clStepperConnectorVertical
	}
	class = class.Merge(clStepperConnectorState[connectorStatus])
	return h.Div(
		h.Class(class.Compile()),
		g.Attr("data-stepper-connector", ""),
		g.Attr("data-connector-status", connectorStatus),
		g.Attr("aria-hidden", "true"),
	)
}

func stepperDataAttrs(step molecules.StepItem, index int) []g.Node {
	attrs := []g.Node{
		g.Attr("data-step-status", step.Status),
		g.Attr("data-step-index", itoa(index)),
	}
	if key := strings.TrimSpace(step.Key); key != "" {
		attrs = append(attrs, g.Attr("data-step-key", key))
	}
	return attrs
}

// SidebarSlots carries trusted rich composition around the portable navigation
// model. Items is a compatibility seam for already-rendered navigation nodes;
// new portable callers should prefer SidebarProps.Items or Sections.
type SidebarSlots struct {
	Brand  []g.Node
	Items  []g.Node
	Footer []g.Node
}

// Sidebar renders the portable sidebar model.
func Sidebar(p molecules.SidebarProps) g.Node {
	return SidebarWithSlots(p, SidebarSlots{})
}

// SidebarWithSlots renders the canonical persistent navigation surface while
// preserving rich brand and footer composition for trusted Go callers.
func SidebarWithSlots(p molecules.SidebarProps, slots SidebarSlots) g.Node {
	flavor := "admin"
	if p.Flavor == "content" {
		flavor = "content"
	}
	collapsible := p.Collapsible || p.Collapsed

	rootClass := clSidebarRootAdmin
	columnClass := clSidebarColumnAdmin
	brandClass := clSidebarBrandAdmin
	navWrapClass := clSidebarNavWrapAdmin
	navClass := clSidebarNavAdmin
	footerClass := clSidebarFooterAdmin
	if flavor == "content" {
		rootClass = clSidebarRootContent
		columnClass = clSidebarColumnContent
		brandClass = clSidebarBrandContent
		navWrapClass = clSidebarNavWrapContent
		navClass = clSidebarNavContent
	} else if p.Collapsed {
		rootClass = rootClass.Merge(clSidebarWidthCollapsed)
	} else {
		rootClass = rootClass.Merge(clSidebarWidthExpanded)
	}
	if flavor == "content" {
		footerClass = clSidebarFooterContent
	}
	if p.Disabled {
		rootClass = rootClass.Merge(clSidebarDisabled)
	}

	sidebarID := strings.TrimSpace(p.ID)
	if sidebarID == "" {
		sidebarID = flavor + "-sidebar"
	}
	if collapsible && !strings.HasSuffix(sidebarID, "-collapsible") {
		sidebarID += "-collapsible"
	}
	label := strings.TrimSpace(p.NavigationLabel)
	if label == "" {
		if flavor == "content" {
			label = "Content navigation"
		} else {
			label = "Admin navigation"
		}
	}

	rootProps := p.ComponentProps
	rootProps.ID = ""
	rootProps.Disabled = false
	attrs := append(baseAttrs(rootProps),
		h.ID(sidebarID),
		classes(rootClass.Compile(), p.Class),
		g.Attr("data-component", "sidebar"),
		g.Attr("data-sidebar-flavor", flavor),
		g.Attr("data-sidebar-collapsible", boolText(collapsible)),
		g.Attr("data-sidebar-collapsed", boolText(p.Collapsed)),
		g.Attr("data-state", sidebarState(p.Collapsed)),
		g.Attr("aria-disabled", boolText(p.Disabled)),
	)
	if collapsible {
		attrs = append(attrs, g.Attr("aria-expanded", boolText(!p.Collapsed)))
	}

	column := []g.Node{h.Class(columnClass.Compile())}
	if brand := sidebarBrand(p, slots.Brand, flavor, brandClass); brand != nil {
		column = append(column, brand)
	}
	column = append(column, h.Div(
		h.Class(navWrapClass.Compile()),
		h.Nav(
			h.Class(navClass.Compile()),
			g.Attr("aria-label", label),
			sidebarNavigation(p, slots.Items, flavor),
		),
	))
	if len(slots.Footer) > 0 {
		column = append(column, h.Footer(
			h.Class(footerClass.Compile()),
			g.Attr("data-sidebar-footer", ""),
			g.Group(slots.Footer),
		))
	}

	attrs = append(attrs, h.Div(
		h.Class(clSidebarInner.Compile()),
		h.Div(column...),
	))
	return h.Aside(attrs...)
}

func sidebarState(collapsed bool) string {
	if collapsed {
		return "collapsed"
	}
	return "expanded"
}

func sidebarBrand(
	p molecules.SidebarProps,
	brand []g.Node,
	flavor string,
	class tw.ClassList,
) g.Node {
	if len(brand) > 0 {
		return h.Header(
			h.Class(class.Compile()),
			g.Attr("data-sidebar-brand", ""),
			g.Group(brand),
		)
	}
	label := strings.TrimSpace(p.BrandLabel)
	if label == "" && flavor == "admin" {
		label = "Admin"
	}
	if label == "" {
		return nil
	}
	href := strings.TrimSpace(p.BrandHref)
	if href == "" {
		href = "/admin"
	}
	return h.Header(
		h.Class(class.Compile()),
		g.Attr("data-sidebar-brand", ""),
		h.A(
			h.Href(href),
			h.Class(clSidebarBrandLink.Compile()),
			h.Span(h.Class(clSidebarBrandText.Compile()), g.Text(label)),
		),
	)
}

func sidebarNavigation(p molecules.SidebarProps, richItems []g.Node, flavor string) g.Node {
	if len(richItems) > 0 {
		return g.Group(richItems)
	}
	if len(p.Sections) > 0 {
		sections := make([]g.Node, 0, len(p.Sections))
		for _, section := range p.Sections {
			sections = append(sections, sidebarSectionNode(p, section, flavor))
		}
		return g.Group(sections)
	}
	items := make([]g.Node, 0, len(p.Items))
	for _, item := range p.Items {
		items = append(items, sidebarItemNode(p, item, flavor, 0))
	}
	return h.Ul(h.Class(clSidebarSectionList.Compile()), g.Group(items))
}

func sidebarSectionNode(
	p molecules.SidebarProps,
	section molecules.SidebarSection,
	flavor string,
) g.Node {
	sectionID := sidebarResolvedID(section.ID, section.Label)
	attrs := []g.Node{
		h.Class(clSidebarSection.Compile()),
		g.Attr("data-sidebar-section", sectionID),
		g.Attr("data-sidebar-search-section", "true"),
	}
	if tone := strings.TrimSpace(section.Tone); tone != "" {
		attrs = append(attrs, g.Attr("data-sidebar-tone", tone))
	}
	if searchText := strings.TrimSpace(section.SearchText); searchText != "" {
		attrs = append(attrs, g.Attr("data-sidebar-search-text", searchText))
	}
	if section.Label != "" || section.Glyph != "" {
		header := []g.Node{
			h.Class(clSidebarSectionHeader.Compile()),
			g.Attr("data-sidebar-section-header", ""),
		}
		if section.Glyph != "" {
			header = append(header, h.Span(
				h.Class(clSidebarSectionGlyph.Compile()),
				g.Attr("data-sidebar-section-glyph", "true"),
				g.Attr("aria-hidden", "true"),
				g.Text(section.Glyph),
			))
		}
		if section.Label != "" {
			header = append(header, h.Span(
				g.Attr("data-sidebar-section-label", "true"),
				g.Text(section.Label),
			))
		}
		attrs = append(attrs, h.Div(header...))
	}
	items := make([]g.Node, 0, len(section.Items))
	for _, item := range section.Items {
		items = append(items, sidebarItemNode(p, item, flavor, 0))
	}
	attrs = append(attrs, h.Ul(h.Class(clSidebarSectionList.Compile()), g.Group(items)))
	return h.Section(attrs...)
}

func sidebarItemNode(
	p molecules.SidebarProps,
	item molecules.SidebarItem,
	flavor string,
	depth int,
) g.Node {
	active := sidebarItemActive(item, p.Current)
	disabled := p.Disabled || item.Disabled
	itemID := sidebarResolvedID(item.ID, item.Label)
	attrs := []g.Node{
		g.Attr("data-sidebar-item", itemID),
		g.Attr("data-sidebar-depth", itoa(depth)),
		g.Attr("data-active", boolText(active)),
		g.Attr("data-has-badge", boolText(item.Badge != "")),
	}
	if active {
		attrs = append(attrs, g.Attr("data-state", "active"))
	} else {
		attrs = append(attrs, g.Attr("data-state", "idle"))
	}
	if disabled {
		attrs = append(attrs, g.Attr("data-disabled", "true"))
	}
	if searchText := strings.TrimSpace(item.SearchText); searchText != "" {
		attrs = append(attrs,
			g.Attr("data-sidebar-search-item", "true"),
			g.Attr("data-sidebar-search-text", searchText),
		)
	}
	attrs = append(attrs, attrPairs(item.Attrs)...)

	linkClass := clSidebarLinkAdmin
	activeClass := clSidebarLinkActiveAdmin
	idleClass := clSidebarLinkIdleAdmin
	prefixClass := clSidebarPrefixAdmin
	labelClass := clSidebarLabelVisible
	if flavor == "content" {
		linkClass = clSidebarLinkContent
		activeClass = clSidebarLinkActiveContent
		idleClass = clSidebarLinkIdleContent
		prefixClass = clSidebarPrefixContent
		labelClass = clSidebarLabelContent
	}
	if p.Collapsed && flavor == "admin" {
		linkClass = linkClass.Merge(clSidebarLinkPadCollapsed)
		labelClass = clSidebarLabelHidden
	} else {
		linkClass = linkClass.Merge(clSidebarLinkPadExpanded)
	}
	if active {
		linkClass = linkClass.Merge(activeClass)
	} else {
		linkClass = linkClass.Merge(idleClass)
	}
	if disabled {
		linkClass = linkClass.Merge(clSidebarItemDisabled)
	}

	content := []g.Node{h.Class(linkClass.Compile()), g.Attr("data-nav", itemID)}
	if active {
		content = append(content, g.Attr("aria-current", "page"))
	}
	if disabled {
		content = append(content, g.Attr("aria-disabled", "true"))
	}
	if p.Collapsed && flavor == "admin" {
		content = append(content, g.Attr("aria-label", item.Label))
	}
	if validIconName(item.Icon) {
		content = append(content, h.Span(
			g.Attr("data-sidebar-item-icon", ""),
			icon(item.Icon),
		))
	}
	if item.Prefix != "" && !(p.Collapsed && flavor == "admin") {
		content = append(content, h.Span(
			h.Class(prefixClass.Compile()),
			g.Attr("data-sidebar-item-prefix", "true"),
			g.Text(item.Prefix),
		))
	}
	content = append(content, h.Span(
		h.Class(labelClass.Compile()),
		g.Attr("data-sidebar-item-label", "true"),
		g.Text(item.Label),
	))
	if item.Badge != "" {
		content = append(content, Badge(atoms.BadgeProps{
			Label: item.Badge, Variant: item.BadgeVariant, Size: "sm",
		}))
	}
	if len(item.Children) > 0 {
		content = append(content, h.Span(
			g.Attr("data-sidebar-item-chevron", ""),
			icon("chevron-down"),
		))
	}

	var itemLink g.Node
	if strings.TrimSpace(item.Href) != "" && !disabled {
		itemLink = h.A(append([]g.Node{h.Href(item.Href)}, content...)...)
	} else {
		itemLink = h.Span(content...)
	}
	if len(item.Children) == 0 {
		attrs = append(attrs, itemLink)
		return h.Li(attrs...)
	}

	children := make([]g.Node, 0, len(item.Children))
	for _, child := range item.Children {
		children = append(children, sidebarItemNode(p, child, flavor, depth+1))
	}
	attrs = append(attrs, h.Div(
		h.Class(clSidebarNestedGroup.Compile()),
		itemLink,
		h.Ul(
			h.Class(clSidebarNestedIndent.Compile()),
			g.Attr("data-sidebar-submenu", itemID),
			g.Group(children),
		),
	))
	return h.Li(attrs...)
}

func sidebarItemActive(item molecules.SidebarItem, current string) bool {
	if current == "" {
		if item.Active {
			return true
		}
	} else if current == item.Href ||
		(item.Href != "" && item.Href != "/admin" && strings.HasPrefix(current, item.Href)) {
		return true
	}
	for _, child := range item.Children {
		if sidebarItemActive(child, current) {
			return true
		}
	}
	return false
}

func sidebarResolvedID(explicit string, fallback string) string {
	if id := strings.TrimSpace(explicit); id != "" {
		return id
	}
	var normalized strings.Builder
	lastDash := false
	for _, char := range strings.ToLower(strings.TrimSpace(fallback)) {
		if (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') {
			normalized.WriteRune(char)
			lastDash = false
			continue
		}
		if !lastDash && normalized.Len() > 0 {
			normalized.WriteByte('-')
			lastDash = true
		}
	}
	id := strings.Trim(normalized.String(), "-")
	if id == "" {
		return "item"
	}
	return id
}

// Breadcrumb renders molecules.BreadcrumbProps as an aria-labelled trail; the
// current page is text, not a link, and carries aria-current.
func Breadcrumb(p molecules.BreadcrumbProps) g.Node {
	if len(p.Items) == 0 {
		return nil
	}
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
	if p.Class != "" {
		nav = append(nav, h.Class(p.Class))
	}
	nav = append(nav, htmxAttrs(p.HTMXProps)...)
	nav = append(nav,
		g.Attr("data-component", "breadcrumb"),
		g.Attr("aria-label", "Breadcrumb"),
		h.Ol(items...),
	)
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
	if mode != "previous-next" && mode != "load-more" {
		return g.Text("")
	}
	previousURL, nextURL, cursorIntent, validCursorIntent := paginationCursorURLs(p, mode)
	if cursorIntent && !validCursorIntent {
		return g.Text("")
	}
	if cursorIntent && previousURL == "" && nextURL == "" {
		return g.Text("")
	}

	if previousURL != "" || nextURL != "" {
		items := baseAttrs(p.ComponentProps)
		items = append(
			items,
			g.Attr("aria-label", navigationLabel),
			g.Attr("data-component", "cursor-pagination"),
			g.Attr("data-cursor-mode", mode),
			classes(clPagination.Compile(), p.Class),
		)
		if mode == "previous-next" && previousURL != "" {
			items = append(items, paginationCursorLink(
				p,
				previousURL,
				previousLabel,
				"backward",
				"prev",
			))
		}
		if nextURL != "" {
			label := nextLabel
			if mode == "load-more" {
				label = loadMoreLabel
			}
			items = append(items, paginationCursorLink(
				p,
				nextURL,
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

const maxPaginationCursorBytes = 16 << 10

func paginationCursorURLs(
	p molecules.PaginationProps,
	mode string,
) (previousURL string, nextURL string, intent bool, valid bool) {
	intent = p.PreviousURL != "" || p.NextURL != "" ||
		p.PreviousCursor != "" || p.NextCursor != "" ||
		p.BeforeParameter != "" || p.AfterParameter != ""
	if !intent {
		return "", "", false, true
	}

	previousURL, nextURL = p.PreviousURL, p.NextURL
	if previousURL != "" && !safePaginationURL(previousURL) {
		return "", "", true, false
	}
	if nextURL != "" && !safePaginationURL(nextURL) {
		return "", "", true, false
	}

	if p.PreviousCursor == "" && p.NextCursor == "" {
		return previousURL, nextURL, true, previousURL != "" || nextURL != ""
	}
	beforeParameter := p.BeforeParameter
	if beforeParameter == "" {
		beforeParameter = "before"
	}
	afterParameter := p.AfterParameter
	if afterParameter == "" {
		afterParameter = "after"
	}
	if !validPaginationQueryParameter(beforeParameter) ||
		!validPaginationQueryParameter(afterParameter) ||
		beforeParameter == afterParameter {
		return "", "", true, false
	}
	if !safePaginationURL(p.BaseURL) {
		return "", "", true, false
	}
	if mode == "previous-next" && previousURL == "" && p.PreviousCursor != "" {
		var ok bool
		previousURL, ok = paginationCursorURL(p.BaseURL, beforeParameter, afterParameter, p.PreviousCursor)
		if !ok {
			return "", "", true, false
		}
	}
	if nextURL == "" && p.NextCursor != "" {
		var ok bool
		nextURL, ok = paginationCursorURL(p.BaseURL, afterParameter, beforeParameter, p.NextCursor)
		if !ok {
			return "", "", true, false
		}
	}
	return previousURL, nextURL, true, true
}

func paginationCursorURL(baseURL string, parameter string, opposite string, cursor string) (string, bool) {
	if cursor == "" || strings.TrimSpace(cursor) == "" || len(cursor) > maxPaginationCursorBytes {
		return "", false
	}
	parsed, err := url.Parse(baseURL)
	if err != nil {
		return "", false
	}
	query, err := url.ParseQuery(parsed.RawQuery)
	if err != nil {
		return "", false
	}
	query.Del(parameter)
	query.Del(opposite)
	query.Set(parameter, cursor)
	parsed.RawQuery = query.Encode()
	return parsed.String(), true
}

func safePaginationURL(raw string) bool {
	if raw == "" || strings.TrimSpace(raw) != raw {
		return false
	}
	parsed, err := url.Parse(raw)
	if err != nil {
		return false
	}
	switch strings.ToLower(parsed.Scheme) {
	case "":
		return true
	case "http", "https":
		return parsed.Host != ""
	default:
		return false
	}
}

func validPaginationQueryParameter(name string) bool {
	if name == "" || len(name) > 64 {
		return false
	}
	for index, char := range name {
		if (char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(index > 0 && char >= '0' && char <= '9') ||
			(index > 0 && (char == '_' || char == '-' || char == '.')) {
			continue
		}
		return false
	}
	return true
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
		children = append(children, ButtonWithSlots(
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
			ButtonSlots{IconStart: []g.Node{Icon(atoms.IconProps{Name: "x-mark", Size: "sm", Tone: "neutral"})}},
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

// TabSlot is one trusted Go tab-panel composition. Portable navigation-only
// tabs remain available through TabsProps.Items; rich panel bodies use this
// slot rather than a second component implementation.
type TabSlot struct {
	ID       string
	Label    string
	Icon     string
	Badge    string
	Disabled bool
	HxGet    string
	Content  []g.Node
}

// TabsSlots carries the ordered tab panels projected into TabsWithSlots.
type TabsSlots struct {
	Tabs []TabSlot
}

type tabsStyle struct {
	root, list, button, active, inactive string
}

// Tabs renders portable navigation tabs. Use TabsWithPanels when each tab owns
// a panel body on the current page.
func Tabs(p molecules.TabsProps) g.Node {
	style := resolveTabsStyle(p)
	active := activeItemKey(p.Items, p.ActiveTab)
	if p.Disabled {
		active = ""
	}
	items := []g.Node{
		h.Class(style.list),
		h.Role("tablist"),
		g.Attr("aria-orientation", tabsOrientation(p.Orientation)),
	}
	for _, it := range p.Items {
		disabled := it.Disabled || p.Disabled
		isActive := it.Key == active && !disabled
		stateClass := style.inactive
		if isActive {
			stateClass = style.active
		}
		className := strings.TrimSpace(style.button + " " + stateClass)
		if disabled {
			className += " " + clTabsDisabled.Compile()
		}
		node := []g.Node{
			h.Class(className),
			h.Role("tab"),
			g.Attr("aria-selected", boolText(isActive)),
			g.Attr("aria-disabled", boolText(disabled)),
		}
		if it.Icon != "" {
			node = append(node, h.Span(h.Class(clTabsIcon.Compile()), icon(it.Icon)))
		}
		node = append(node, g.Text(it.Label))
		if it.Badge != "" {
			node = append(node, h.Span(h.Class(clTabsBadge.Compile()), g.Text(it.Badge)))
		}
		if it.URL != "" && !disabled {
			items = append(items, h.A(append(node, h.Href(it.URL))...))
			continue
		}
		button := append(node, h.Type("button"))
		if disabled {
			button = append(button, h.Disabled())
		}
		items = append(items, h.Button(button...))
	}
	rootProps := p.ComponentProps
	rootProps.Disabled = false
	nav := baseAttrs(rootProps)
	nav = append(nav, classes(style.root, p.Class), g.Attr("data-component", "tabs"))
	nav = append(nav, items...)
	return h.Nav(nav...)
}

// TabsWithPanels is the concise application API for controller-backed tabs.
func TabsWithPanels(p molecules.TabsProps, tabs ...TabSlot) g.Node {
	return TabsWithSlots(p, TabsSlots{Tabs: tabs})
}

// TabsWithSlots renders tabs and their panels from one canonical contract.
func TabsWithSlots(p molecules.TabsProps, slots TabsSlots) g.Node {
	if len(slots.Tabs) == 0 {
		return nil
	}

	style := resolveTabsStyle(p)
	active := activeSlotID(slots.Tabs, p.ActiveTab)
	if p.Disabled {
		active = ""
	}
	rootID := p.ID
	if rootID == "" {
		rootID = "tabs"
		if active != "" {
			rootID += "-" + active
		}
	}

	tabList := []g.Node{
		h.Class(style.list),
		h.Role("tablist"),
		g.Attr("aria-orientation", tabsOrientation(p.Orientation)),
	}
	panels := []g.Node{h.Class(clTabsPanels.Compile())}
	for _, tab := range slots.Tabs {
		disabled := tab.Disabled || p.Disabled
		isActive := tab.ID == active && !disabled
		panelID := rootID + "-panel-" + tab.ID
		tabID := rootID + "-tab-" + tab.ID
		stateClass := style.inactive
		if isActive {
			stateClass = style.active
		}
		buttonClass := strings.TrimSpace(style.button + " " + stateClass)
		if disabled {
			buttonClass += " " + clTabsDisabled.Compile()
		}
		button := []g.Node{
			h.Class(buttonClass), h.Type("button"), h.Role("tab"), h.ID(tabID),
			g.Attr("aria-controls", panelID),
			g.Attr("aria-selected", boolText(isActive)),
			g.Attr("aria-disabled", boolText(disabled)),
			g.Attr("tabindex", activeTabIndex(isActive)),
			g.Attr("data-tabs-tab", tab.ID),
			g.Attr("data-tabs-active-classes", style.active),
			g.Attr("data-tabs-inactive-classes", style.inactive),
			g.Attr("data-action", "click->tabs#activate"),
		}
		if disabled {
			button = append(button, h.Disabled())
		}
		if tab.Icon != "" {
			button = append(button, h.Span(h.Class(clTabsIcon.Compile()), icon(tab.Icon)))
		}
		button = append(button, h.Span(g.Text(tab.Label)))
		if tab.Badge != "" {
			button = append(button, h.Span(h.Class(clTabsBadge.Compile()), g.Text(tab.Badge)))
		}
		tabList = append(tabList, h.Button(button...))

		panel := []g.Node{
			h.Class(clTabsPanel.Compile()), h.ID(panelID), h.Role("tabpanel"),
			g.Attr("aria-labelledby", tabID),
			g.Attr("aria-hidden", boolText(!isActive)),
			g.Attr("data-tabs-panel", tab.ID),
			g.Attr("data-state", tabState(isActive)),
		}
		if !isActive {
			panel = append(panel, g.Attr("hidden"))
		}
		hxGet := tab.HxGet
		if hxGet == "" {
			hxGet = p.HxGet
		}
		if hxGet != "" {
			panel = append(panel,
				g.Attr("data-tabs-lazy", "true"),
				g.Attr("hx-get", hxGet),
				g.Attr("hx-trigger", "tabs:activate from:this once"),
				g.Attr("hx-swap", "innerHTML"),
				h.Div(h.Class(clTabsLazy.Compile()),
					h.Span(h.Class(clTabsLazyLabel.Compile()), g.Text(fallbackText(p.LoadingLabel, "Loading...")))),
			)
		} else {
			panel = append(panel, tab.Content...)
		}
		panels = append(panels, h.Div(panel...))
	}

	rootProps := p.ComponentProps
	rootProps.Disabled = false
	root := baseAttrs(rootProps)
	root = append(root,
		classes(style.root, p.Class),
		g.Attr("data-component", "tabs"),
		g.Attr("data-controller", "tabs"),
		g.Attr("data-tabs-contract", "1"),
		g.Attr("data-tabs-active-tab-value", active),
		h.Div(tabList...),
		h.Div(panels...),
	)
	return h.Div(root...)
}

func resolveTabsStyle(p molecules.TabsProps) tabsStyle {
	orientation := tabsOrientation(p.Orientation)
	root := clTabsRoot.Merge(clTabsRootHorizontal)
	list := clTabsListBase.Merge(clTabsListHorizontal)
	button := clTabsButtonBase
	if orientation == "vertical" {
		root = clTabsRoot.Merge(clTabsRootVertical)
		list = clTabsListBase.Merge(clTabsListVertical)
	}
	if p.Variant == "pills" {
		button = button.Merge(clTabsButtonPills)
		return tabsStyle{root.Compile(), list.Compile(), button.Compile(), clTabsPillsActive.Compile(), clTabsPillsIdle.Compile()}
	}
	if orientation == "vertical" {
		list = list.Merge(clTabsListUnderlineVertical)
		button = button.Merge(clTabsButtonUnderlineVertical)
	} else {
		list = list.Merge(clTabsListUnderlineHorizontal)
		button = button.Merge(clTabsButtonUnderlineHorizontal)
	}
	return tabsStyle{root.Compile(), list.Compile(), button.Compile(), clTabsUnderlineActive.Compile(), clTabsUnderlineIdle.Compile()}
}

func tabsOrientation(value string) string {
	if value == "vertical" {
		return value
	}
	return "horizontal"
}

func activeItemKey(items []molecules.TabItem, requested string) string {
	for _, item := range items {
		if item.Key == requested && !item.Disabled {
			return item.Key
		}
	}
	for _, item := range items {
		if !item.Disabled {
			return item.Key
		}
	}
	return ""
}

func activeSlotID(tabs []TabSlot, requested string) string {
	for _, tab := range tabs {
		if tab.ID == requested && !tab.Disabled {
			return tab.ID
		}
	}
	for _, tab := range tabs {
		if !tab.Disabled {
			return tab.ID
		}
	}
	return ""
}

func boolText(value bool) string {
	if value {
		return "true"
	}
	return "false"
}

func activeTabIndex(active bool) string {
	if active {
		return "0"
	}
	return "-1"
}

func tabState(active bool) string {
	if active {
		return "active"
	}
	return "inactive"
}
