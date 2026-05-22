package atoms

import "github.com/septagon-oss/platformkit-ui/contracts"

// BadgeProps defines the platform-agnostic properties for a Badge component.
type BadgeProps struct {
	contracts.ComponentProps

	Text    string `json:"text"`
	Variant string `json:"variant,omitempty"` // primary, secondary, success, warning, error, info
	Size    string `json:"size,omitempty"`    // small, medium, large
	Dot     bool   `json:"dot,omitempty"`     // show status dot instead of text
	Icon    string `json:"icon,omitempty"`
}

// ToMap converts BadgeProps to map[string]any for unified Component construction.
func (p BadgeProps) ToMap() map[string]any { return propsToMap(p) }
