package atoms

import "github.com/septagon-oss/platformkit-ui/contracts"

// ButtonProps defines the platform-agnostic properties for a Button component.
type ButtonProps struct {
	contracts.ComponentProps
	contracts.HTMXProps

	Text         string `json:"text"`
	Variant      string `json:"variant,omitempty"` // primary, secondary, success, warning, error, info, outline, ghost, link
	Size         string `json:"size,omitempty"`    // xs, small, medium, large, xl, 2xl
	Type         string `json:"type,omitempty"`    // button, submit, reset
	Loading      bool   `json:"loading,omitempty"`
	FullWidth    bool   `json:"fullWidth,omitempty"`
	Icon         string `json:"icon,omitempty"`
	IconPosition string `json:"iconPosition,omitempty"` // left, right
	IconStyle    string `json:"iconStyle,omitempty"`    // outline, solid
}

// ToMap converts ButtonProps to map[string]any for unified Component construction.
func (p ButtonProps) ToMap() map[string]any { return propsToMap(p) }
