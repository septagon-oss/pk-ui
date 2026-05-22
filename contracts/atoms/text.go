package atoms

import "github.com/septagon-oss/platformkit-ui/contracts"

// TextProps defines the platform-agnostic properties for a Text component.
type TextProps struct {
	contracts.ComponentProps

	Content  string `json:"content"`
	Size     string `json:"size,omitempty"`     // xs, small, medium, large, xl
	Weight   string `json:"weight,omitempty"`   // normal, medium, semibold, bold
	Color    string `json:"color,omitempty"`    // primary, secondary, tertiary, success, warning, error, muted
	Truncate bool   `json:"truncate,omitempty"` // truncate with ellipsis
	Lines    int    `json:"lines,omitempty"`    // line clamp
}

// ToMap converts TextProps to map[string]any for unified Component construction.
func (p TextProps) ToMap() map[string]any { return propsToMap(p) }

// HeadingProps defines properties for heading elements (H1-H6).
type HeadingProps struct {
	contracts.ComponentProps

	Text     string `json:"text"`
	Level    int    `json:"level"`            // 1-6
	Anchor   string `json:"anchor,omitempty"` // optional anchor ID
	Truncate bool   `json:"truncate,omitempty"`
}

// ToMap converts HeadingProps to map[string]any for unified Component construction.
func (p HeadingProps) ToMap() map[string]any { return propsToMap(p) }

// LabelProps defines properties for form labels.
type LabelProps struct {
	contracts.ComponentProps

	Text     string `json:"text"`
	For      string `json:"for,omitempty"` // associated input ID
	Required bool   `json:"required,omitempty"`
}

// ToMap converts LabelProps to map[string]any for unified Component construction.
func (p LabelProps) ToMap() map[string]any { return propsToMap(p) }
