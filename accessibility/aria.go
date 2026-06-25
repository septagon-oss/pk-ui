package accessibility

import (
	"fmt"
	"strings"
	"sync/atomic"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

// ARIARole represents standard ARIA roles
type ARIARole string

const (
	// Widget roles
	ARIARoleButton      ARIARole = "button"
	ARIARoleLink        ARIARole = "link"
	ARIARoleCheckbox    ARIARole = "checkbox"
	ARIARoleCombobox    ARIARole = "combobox"
	ARIARoleListbox     ARIARole = "listbox"
	ARIARoleOption      ARIARole = "option"
	ARIARoleProgressbar ARIARole = "progressbar"
	ARIARoleRadio       ARIARole = "radio"
	ARIARoleScrollbar   ARIARole = "scrollbar"
	ARIARoleSlider      ARIARole = "slider"
	ARIARoleSpinbutton  ARIARole = "spinbutton"
	ARIARoleSwitch      ARIARole = "switch"
	ARIARoleTab         ARIARole = "tab"
	ARIARoleTabpanel    ARIARole = "tabpanel"
	ARIARoleTextbox     ARIARole = "textbox"

	// Composite roles
	ARIARoleComboboxListbox ARIARole = "combobox listbox"
	ARIARoleGrid            ARIARole = "grid"
	ARIARoleTablist         ARIARole = "tablist"
	ARIARoleTree            ARIARole = "tree"
	ARIARoleTreegrid        ARIARole = "treegrid"

	// Document structure roles
	ARIARoleApplication  ARIARole = "application"
	ARIARoleDocument     ARIARole = "document"
	ARIARoleGroup        ARIARole = "group"
	ARIARoleHeading      ARIARole = "heading"
	ARIARoleImg          ARIARole = "img"
	ARIARoleList         ARIARole = "list"
	ARIARoleListitem     ARIARole = "listitem"
	ARIARoleMath         ARIARole = "math"
	ARIARoleNote         ARIARole = "note"
	ARIARolePresentation ARIARole = "presentation"
	ARIARoleRegion       ARIARole = "region"
	ARIARoleSeparator    ARIARole = "separator"
	ARIARoleToolbar      ARIARole = "toolbar"

	// Landmark roles
	ARIARoleBanner        ARIARole = "banner"
	ARIARoleComplementary ARIARole = "complementary"
	ARIARoleContentinfo   ARIARole = "contentinfo"
	ARIARoleForm          ARIARole = "form"
	ARIARoleMain          ARIARole = "main"
	ARIARoleNavigation    ARIARole = "navigation"
	ARIARoleSearch        ARIARole = "search"

	// Live region roles
	ARIARoleAlert   ARIARole = "alert"
	ARIARoleLog     ARIARole = "log"
	ARIARoleMarquee ARIARole = "marquee"
	ARIARoleStatus  ARIARole = "status"
	ARIARoleTimer   ARIARole = "timer"

	// Window roles
	ARIARoleAlertdialog ARIARole = "alertdialog"
	ARIARoleDialog      ARIARole = "dialog"
)

// ARIAState represents ARIA state values
type ARIAState string

const (
	ARIAStateTrue  ARIAState = "true"
	ARIAStateFalse ARIAState = "false"
	ARIAStateMixed ARIAState = "mixed"
)

// ARIAProperties contains ARIA properties for components
type ARIAProperties struct {
	// Accessible name and description
	Label       string `json:"label,omitempty"`
	LabelledBy  string `json:"labelledby,omitempty"`
	DescribedBy string `json:"describedby,omitempty"`

	// Widget attributes
	Checked  *ARIAState `json:"checked,omitempty"`
	Disabled *ARIAState `json:"disabled,omitempty"`
	Expanded *ARIAState `json:"expanded,omitempty"`
	HasPopup string     `json:"haspopup,omitempty"`
	Hidden   *ARIAState `json:"hidden,omitempty"`
	Invalid  *ARIAState `json:"invalid,omitempty"`
	Pressed  *ARIAState `json:"pressed,omitempty"`
	ReadOnly *ARIAState `json:"readonly,omitempty"`
	Required *ARIAState `json:"required,omitempty"`
	Selected *ARIAState `json:"selected,omitempty"`

	// Live region attributes
	Live     string     `json:"live,omitempty"`
	Atomic   *ARIAState `json:"atomic,omitempty"`
	Busy     *ARIAState `json:"busy,omitempty"`
	Relevant string     `json:"relevant,omitempty"`

	// Relationship attributes
	Controls         string `json:"controls,omitempty"`
	Owns             string `json:"owns,omitempty"`
	ActiveDescendant string `json:"activedescendant,omitempty"`

	// Drag and drop attributes
	Dropeffect string     `json:"dropeffect,omitempty"`
	Grabbed    *ARIAState `json:"grabbed,omitempty"`

	// Value attributes
	ValueMin  *float64 `json:"valuemin,omitempty"`
	ValueMax  *float64 `json:"valuemax,omitempty"`
	ValueNow  *float64 `json:"valuenow,omitempty"`
	ValueText string   `json:"valuetext,omitempty"`

	// Position attributes
	Level    *int `json:"level,omitempty"`
	PosInSet *int `json:"posinset,omitempty"`
	SetSize  *int `json:"setsize,omitempty"`

	// Input attributes
	Autocomplete    string     `json:"autocomplete,omitempty"`
	Multiline       *ARIAState `json:"multiline,omitempty"`
	Multiselectable *ARIAState `json:"multiselectable,omitempty"`
	Orientation     string     `json:"orientation,omitempty"`
	Placeholder     string     `json:"placeholder,omitempty"`

	// Global attributes
	Role     ARIARole `json:"role,omitempty"`
	TabIndex *int     `json:"tabindex,omitempty"`
}

// NewARIAProperties creates a new ARIA properties instance
func NewARIAProperties() *ARIAProperties {
	return &ARIAProperties{}
}

// SetLabel sets the aria-label attribute
func (ap *ARIAProperties) SetLabel(label string) *ARIAProperties {
	ap.Label = label
	return ap
}

// SetLabelledBy sets the aria-labelledby attribute
func (ap *ARIAProperties) SetLabelledBy(id string) *ARIAProperties {
	ap.LabelledBy = id
	return ap
}

// SetDescribedBy sets the aria-describedby attribute
func (ap *ARIAProperties) SetDescribedBy(id string) *ARIAProperties {
	ap.DescribedBy = id
	return ap
}

// SetRole sets the role attribute
func (ap *ARIAProperties) SetRole(role ARIARole) *ARIAProperties {
	ap.Role = role
	return ap
}

// SetChecked sets the aria-checked attribute
func (ap *ARIAProperties) SetChecked(checked bool) *ARIAProperties {
	if checked {
		ap.Checked = &[]ARIAState{ARIAStateTrue}[0]
	} else {
		ap.Checked = &[]ARIAState{ARIAStateFalse}[0]
	}
	return ap
}

// SetCheckedMixed sets the aria-checked attribute to mixed
func (ap *ARIAProperties) SetCheckedMixed() *ARIAProperties {
	ap.Checked = &[]ARIAState{ARIAStateMixed}[0]
	return ap
}

// SetDisabled sets the aria-disabled attribute
func (ap *ARIAProperties) SetDisabled(disabled bool) *ARIAProperties {
	if disabled {
		ap.Disabled = &[]ARIAState{ARIAStateTrue}[0]
	} else {
		ap.Disabled = &[]ARIAState{ARIAStateFalse}[0]
	}
	return ap
}

// SetExpanded sets the aria-expanded attribute
func (ap *ARIAProperties) SetExpanded(expanded bool) *ARIAProperties {
	if expanded {
		ap.Expanded = &[]ARIAState{ARIAStateTrue}[0]
	} else {
		ap.Expanded = &[]ARIAState{ARIAStateFalse}[0]
	}
	return ap
}

// SetHasPopup sets the aria-haspopup attribute
func (ap *ARIAProperties) SetHasPopup(popup string) *ARIAProperties {
	ap.HasPopup = popup
	return ap
}

// SetHidden sets the aria-hidden attribute
func (ap *ARIAProperties) SetHidden(hidden bool) *ARIAProperties {
	if hidden {
		ap.Hidden = &[]ARIAState{ARIAStateTrue}[0]
	} else {
		ap.Hidden = &[]ARIAState{ARIAStateFalse}[0]
	}
	return ap
}

// SetInvalid sets the aria-invalid attribute
func (ap *ARIAProperties) SetInvalid(invalid bool) *ARIAProperties {
	if invalid {
		ap.Invalid = &[]ARIAState{ARIAStateTrue}[0]
	} else {
		ap.Invalid = &[]ARIAState{ARIAStateFalse}[0]
	}
	return ap
}

// SetPressed sets the aria-pressed attribute
func (ap *ARIAProperties) SetPressed(pressed bool) *ARIAProperties {
	if pressed {
		ap.Pressed = &[]ARIAState{ARIAStateTrue}[0]
	} else {
		ap.Pressed = &[]ARIAState{ARIAStateFalse}[0]
	}
	return ap
}

// SetReadOnly sets the aria-readonly attribute
func (ap *ARIAProperties) SetReadOnly(readonly bool) *ARIAProperties {
	if readonly {
		ap.ReadOnly = &[]ARIAState{ARIAStateTrue}[0]
	} else {
		ap.ReadOnly = &[]ARIAState{ARIAStateFalse}[0]
	}
	return ap
}

// SetRequired sets the aria-required attribute
func (ap *ARIAProperties) SetRequired(required bool) *ARIAProperties {
	if required {
		ap.Required = &[]ARIAState{ARIAStateTrue}[0]
	} else {
		ap.Required = &[]ARIAState{ARIAStateFalse}[0]
	}
	return ap
}

// SetSelected sets the aria-selected attribute
func (ap *ARIAProperties) SetSelected(selected bool) *ARIAProperties {
	if selected {
		ap.Selected = &[]ARIAState{ARIAStateTrue}[0]
	} else {
		ap.Selected = &[]ARIAState{ARIAStateFalse}[0]
	}
	return ap
}

// SetLive sets the aria-live attribute
func (ap *ARIAProperties) SetLive(live string) *ARIAProperties {
	ap.Live = live
	return ap
}

// SetControls sets the aria-controls attribute
func (ap *ARIAProperties) SetControls(controls string) *ARIAProperties {
	ap.Controls = controls
	return ap
}

// SetActiveDescendant sets the aria-activedescendant attribute
func (ap *ARIAProperties) SetActiveDescendant(id string) *ARIAProperties {
	ap.ActiveDescendant = id
	return ap
}

// SetValueRange sets the aria-valuemin, aria-valuemax, and aria-valuenow attributes
func (ap *ARIAProperties) SetValueRange(min, max, now float64) *ARIAProperties {
	ap.ValueMin = &min
	ap.ValueMax = &max
	ap.ValueNow = &now
	return ap
}

// SetValueText sets the aria-valuetext attribute
func (ap *ARIAProperties) SetValueText(text string) *ARIAProperties {
	ap.ValueText = text
	return ap
}

// SetLevel sets the aria-level attribute
func (ap *ARIAProperties) SetLevel(level int) *ARIAProperties {
	ap.Level = &level
	return ap
}

// SetPosition sets the aria-posinset and aria-setsize attributes
func (ap *ARIAProperties) SetPosition(posInSet, setSize int) *ARIAProperties {
	ap.PosInSet = &posInSet
	ap.SetSize = &setSize
	return ap
}

// SetTabIndex sets the tabindex attribute
func (ap *ARIAProperties) SetTabIndex(tabIndex int) *ARIAProperties {
	ap.TabIndex = &tabIndex
	return ap
}

// ToAttributes converts ARIA properties to gomponents attributes
func (ap *ARIAProperties) ToAttributes() []g.Node {
	var attrs []g.Node

	if ap.Role != "" {
		attrs = append(attrs, g.Attr("role", string(ap.Role)))
	}

	if ap.Label != "" {
		attrs = append(attrs, g.Attr("aria-label", ap.Label))
	}

	if ap.LabelledBy != "" {
		attrs = append(attrs, g.Attr("aria-labelledby", ap.LabelledBy))
	}

	if ap.DescribedBy != "" {
		attrs = append(attrs, g.Attr("aria-describedby", ap.DescribedBy))
	}

	if ap.Checked != nil {
		attrs = append(attrs, g.Attr("aria-checked", string(*ap.Checked)))
	}

	if ap.Disabled != nil {
		attrs = append(attrs, g.Attr("aria-disabled", string(*ap.Disabled)))
	}

	if ap.Expanded != nil {
		attrs = append(attrs, g.Attr("aria-expanded", string(*ap.Expanded)))
	}

	if ap.HasPopup != "" {
		attrs = append(attrs, g.Attr("aria-haspopup", ap.HasPopup))
	}

	if ap.Hidden != nil {
		attrs = append(attrs, g.Attr("aria-hidden", string(*ap.Hidden)))
	}

	if ap.Invalid != nil {
		attrs = append(attrs, g.Attr("aria-invalid", string(*ap.Invalid)))
	}

	if ap.Pressed != nil {
		attrs = append(attrs, g.Attr("aria-pressed", string(*ap.Pressed)))
	}

	if ap.ReadOnly != nil {
		attrs = append(attrs, g.Attr("aria-readonly", string(*ap.ReadOnly)))
	}

	if ap.Required != nil {
		attrs = append(attrs, g.Attr("aria-required", string(*ap.Required)))
	}

	if ap.Selected != nil {
		attrs = append(attrs, g.Attr("aria-selected", string(*ap.Selected)))
	}

	if ap.Live != "" {
		attrs = append(attrs, g.Attr("aria-live", ap.Live))
	}

	if ap.Atomic != nil {
		attrs = append(attrs, g.Attr("aria-atomic", string(*ap.Atomic)))
	}

	if ap.Busy != nil {
		attrs = append(attrs, g.Attr("aria-busy", string(*ap.Busy)))
	}

	if ap.Relevant != "" {
		attrs = append(attrs, g.Attr("aria-relevant", ap.Relevant))
	}

	if ap.Controls != "" {
		attrs = append(attrs, g.Attr("aria-controls", ap.Controls))
	}

	if ap.Owns != "" {
		attrs = append(attrs, g.Attr("aria-owns", ap.Owns))
	}

	if ap.ActiveDescendant != "" {
		attrs = append(attrs, g.Attr("aria-activedescendant", ap.ActiveDescendant))
	}

	if ap.ValueMin != nil {
		attrs = append(attrs, g.Attr("aria-valuemin", fmt.Sprintf("%.2f", *ap.ValueMin)))
	}

	if ap.ValueMax != nil {
		attrs = append(attrs, g.Attr("aria-valuemax", fmt.Sprintf("%.2f", *ap.ValueMax)))
	}

	if ap.ValueNow != nil {
		attrs = append(attrs, g.Attr("aria-valuenow", fmt.Sprintf("%.2f", *ap.ValueNow)))
	}

	if ap.ValueText != "" {
		attrs = append(attrs, g.Attr("aria-valuetext", ap.ValueText))
	}

	if ap.Level != nil {
		attrs = append(attrs, g.Attr("aria-level", fmt.Sprintf("%d", *ap.Level)))
	}

	if ap.PosInSet != nil {
		attrs = append(attrs, g.Attr("aria-posinset", fmt.Sprintf("%d", *ap.PosInSet)))
	}

	if ap.SetSize != nil {
		attrs = append(attrs, g.Attr("aria-setsize", fmt.Sprintf("%d", *ap.SetSize)))
	}

	if ap.Autocomplete != "" {
		attrs = append(attrs, g.Attr("aria-autocomplete", ap.Autocomplete))
	}

	if ap.Multiline != nil {
		attrs = append(attrs, g.Attr("aria-multiline", string(*ap.Multiline)))
	}

	if ap.Multiselectable != nil {
		attrs = append(attrs, g.Attr("aria-multiselectable", string(*ap.Multiselectable)))
	}

	if ap.Orientation != "" {
		attrs = append(attrs, g.Attr("aria-orientation", ap.Orientation))
	}

	if ap.Placeholder != "" {
		attrs = append(attrs, g.Attr("aria-placeholder", ap.Placeholder))
	}

	if ap.TabIndex != nil {
		attrs = append(attrs, g.Attr("tabindex", fmt.Sprintf("%d", *ap.TabIndex)))
	}

	return attrs
}

// AccessibilityManager manages accessibility features for components
type AccessibilityManager struct {
	announcements []string
	focusStack    []string
}

// NewAccessibilityManager creates a new accessibility manager
func NewAccessibilityManager() *AccessibilityManager {
	return &AccessibilityManager{
		announcements: make([]string, 0),
		focusStack:    make([]string, 0),
	}
}

// Announce adds an announcement for screen readers
func (am *AccessibilityManager) Announce(message string) {
	am.announcements = append(am.announcements, message)
}

// GetAnnouncements returns all pending announcements
func (am *AccessibilityManager) GetAnnouncements() []string {
	announcements := make([]string, len(am.announcements))
	copy(announcements, am.announcements)
	am.announcements = am.announcements[:0] // Clear announcements
	return announcements
}

// PushFocus adds a component to the focus stack
func (am *AccessibilityManager) PushFocus(componentID string) {
	am.focusStack = append(am.focusStack, componentID)
}

// PopFocus removes the last component from the focus stack
func (am *AccessibilityManager) PopFocus() string {
	if len(am.focusStack) == 0 {
		return ""
	}

	last := am.focusStack[len(am.focusStack)-1]
	am.focusStack = am.focusStack[:len(am.focusStack)-1]
	return last
}

// GetFocusTarget returns the current focus target
func (am *AccessibilityManager) GetFocusTarget() string {
	if len(am.focusStack) == 0 {
		return ""
	}
	return am.focusStack[len(am.focusStack)-1]
}

// AccessibilityValidator validates accessibility properties
type AccessibilityValidator struct{}

// ValidateARIA validates ARIA properties for common issues
func (av *AccessibilityValidator) ValidateARIA(props *ARIAProperties) []string {
	var warnings []string

	// Check for missing accessible name
	if props.Role != "" && props.Label == "" && props.LabelledBy == "" {
		switch props.Role {
		case ARIARoleButton, ARIARoleLink, ARIARoleTextbox:
			warnings = append(warnings, "Interactive element missing accessible name (aria-label or aria-labelledby)")
		}
	}

	// Check for invalid role combinations
	if props.Role == ARIARoleButton && props.Checked != nil {
		warnings = append(warnings, "Button role should not have aria-checked. Consider using role='switch' instead")
	}

	// Check for missing required properties
	if props.Role == ARIARoleProgressbar {
		if props.ValueMin == nil || props.ValueMax == nil {
			warnings = append(warnings, "Progress bar missing aria-valuemin or aria-valuemax")
		}
	}

	if props.Role == ARIARoleSlider {
		if props.ValueMin == nil || props.ValueMax == nil || props.ValueNow == nil {
			warnings = append(warnings, "Slider missing aria-valuemin, aria-valuemax, or aria-valuenow")
		}
	}

	// Check for conflicting states
	if props.Hidden != nil && *props.Hidden == ARIAStateTrue && props.TabIndex != nil && *props.TabIndex >= 0 {
		warnings = append(warnings, "Hidden element should not be focusable (tabindex >= 0)")
	}

	return warnings
}

// LiveRegion creates a live region for announcements
type LiveRegion struct {
	ID       string
	Level    string // "polite", "assertive", or "off"
	Atomic   bool
	Relevant string
}

// NewLiveRegion creates a new live region
func NewLiveRegion(id, level string) *LiveRegion {
	return &LiveRegion{
		ID:       id,
		Level:    level,
		Atomic:   false,
		Relevant: "additions text",
	}
}

// Render renders the live region as a gomponents Node
func (lr *LiveRegion) Render() g.Node {
	attrs := []g.Node{
		h.ID(lr.ID),
		g.Attr("aria-live", lr.Level),
		g.Attr("aria-relevant", lr.Relevant),
		h.Class("sr-only"), // Screen reader only
	}

	if lr.Atomic {
		attrs = append(attrs, g.Attr("aria-atomic", "true"))
	}

	return h.Div(g.Group(attrs))
}

// SkipLink creates a skip navigation link
type SkipLink struct {
	Text   string
	Target string
	Class  string
}

// NewSkipLink creates a new skip link
func NewSkipLink(text, target string) *SkipLink {
	return &SkipLink{
		Text:   text,
		Target: target,
		Class:  "sr-only focus:not-sr-only focus:absolute focus:left-4 focus:top-4 focus:z-50 focus:rounded-md focus:border focus:border-border-primary focus:bg-surface-primary focus:px-3 focus:py-2 focus:text-fg-primary focus:shadow-lg focus:outline-none focus:ring-2 focus:ring-ring-brand",
	}
}

// Render renders the skip link as a gomponents Node
func (sl *SkipLink) Render() g.Node {
	return h.A(
		h.Href("#"+sl.Target),
		h.Class(sl.Class),
		g.Text(sl.Text),
	)
}

// FocusTrap manages focus trapping for modals and dialogs
type FocusTrap struct {
	ContainerID        string
	AutoFocus          bool
	RestoreFocus       bool
	FocusableSelectors []string
}

// NewFocusTrap creates a new focus trap
func NewFocusTrap(containerID string) *FocusTrap {
	return &FocusTrap{
		ContainerID:  containerID,
		AutoFocus:    true,
		RestoreFocus: true,
		FocusableSelectors: []string{
			"button:not([disabled])",
			"[href]",
			"input:not([disabled])",
			"select:not([disabled])",
			"textarea:not([disabled])",
			"[tabindex]:not([tabindex='-1'])",
		},
	}
}

// GetJavaScript returns JavaScript code for the focus trap
func (ft *FocusTrap) GetJavaScript() string {
	return fmt.Sprintf(`
(function() {
	const container = document.getElementById('%s');
	if (!container) return;

	const focusableElements = container.querySelectorAll('%s');
	const firstFocusable = focusableElements[0];
	const lastFocusable = focusableElements[focusableElements.length - 1];

	function trapFocus(e) {
		if (e.key === 'Tab') {
			if (e.shiftKey) {
				if (document.activeElement === firstFocusable) {
					e.preventDefault();
					lastFocusable.focus();
				}
			} else {
				if (document.activeElement === lastFocusable) {
					e.preventDefault();
					firstFocusable.focus();
				}
			}
		}
		if (e.key === 'Escape') {
			// Handle escape key
			container.dispatchEvent(new CustomEvent('focus-trap-escape'));
		}
	}

	container.addEventListener('keydown', trapFocus);
	
	%s

	// Return cleanup function
	return function cleanup() {
		container.removeEventListener('keydown', trapFocus);
	};
})();
`, ft.ContainerID,
		strings.Join(ft.FocusableSelectors, ","),
		func() string {
			if ft.AutoFocus {
				return "if (firstFocusable) firstFocusable.focus();"
			}
			return ""
		}())
}

// ScreenReaderText creates screen reader only text
func ScreenReaderText(text string) g.Node {
	return h.Span(
		h.Class("sr-only"),
		g.Text(text),
	)
}

// VisuallyHidden creates visually hidden content that's still accessible
func VisuallyHidden(content g.Node) g.Node {
	return h.Span(
		h.Class("sr-only"),
		content,
	)
}

// ColorContrastValidator validates color contrast ratios
type ColorContrastValidator struct{}

// ContrastRatio represents a color contrast ratio
type ContrastRatio float64

const (
	ContrastRatioAA      ContrastRatio = 4.5 // WCAG AA standard
	ContrastRatioAAA     ContrastRatio = 7.0 // WCAG AAA standard
	ContrastRatioAALarge ContrastRatio = 3.0 // WCAG AA for large text
)

// Common accessibility utility functions

// idCounter provides monotonically increasing IDs for accessibility elements.
var idCounter atomic.Uint64

// GenerateID generates a unique ID for accessibility purposes.
// Each call returns a distinct value, safe for concurrent use.
func GenerateID(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, idCounter.Add(1))
}

// AddScreenReaderSupport adds screen reader support to existing elements
func AddScreenReaderSupport(element g.Node, label, description string) g.Node {
	attrs := []g.Node{element}

	if label != "" {
		attrs = append(attrs, g.Attr("aria-label", label))
	}

	if description != "" {
		descID := GenerateID("desc")
		attrs = append(attrs,
			g.Attr("aria-describedby", descID),
			h.Span(h.ID(descID), h.Class("sr-only"), g.Text(description)),
		)
	}

	return g.Group(attrs)
}

// Global accessibility manager instance
var GlobalAccessibilityManager = NewAccessibilityManager()

// Common ARIA property builders
func ARIAButton(label string) *ARIAProperties {
	return NewARIAProperties().SetRole(ARIARoleButton).SetLabel(label)
}

func ARIACheckbox(label string, checked bool) *ARIAProperties {
	return NewARIAProperties().SetRole(ARIARoleCheckbox).SetLabel(label).SetChecked(checked)
}

func ARIATextbox(label string, required bool) *ARIAProperties {
	props := NewARIAProperties().SetRole(ARIARoleTextbox).SetLabel(label)
	if required {
		props.SetRequired(true)
	}
	return props
}

func ARIACombobox(label string, expanded bool) *ARIAProperties {
	return NewARIAProperties().SetRole(ARIARoleCombobox).SetLabel(label).SetExpanded(expanded)
}

func ARIADialog(label string) *ARIAProperties {
	return NewARIAProperties().SetRole(ARIARoleDialog).SetLabel(label)
}

func ARIAAlert(message string) *ARIAProperties {
	return NewARIAProperties().SetRole(ARIARoleAlert).SetLabel(message).SetLive("assertive")
}

func ARIAStatus(message string) *ARIAProperties {
	return NewARIAProperties().SetRole(ARIARoleStatus).SetLabel(message).SetLive("polite")
}
