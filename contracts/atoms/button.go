package atoms

// Implements: REQ-011.
// Per: ADR-0031.
// Discipline: C-14.
import "github.com/septagon-oss/pk-ui/contracts"

// ButtonProps defines the platform-agnostic properties for a Button component.
type ButtonProps struct {
	contracts.ComponentProps
	contracts.HTMXProps

	Label     string `json:"label"`
	Variant   string `json:"variant,omitempty"` // primary, secondary, outline, ghost, link
	Tone      string `json:"tone,omitempty"`    // neutral, brand, success, warning, danger, info
	Size      string `json:"size,omitempty"`    // xs, sm, md, lg, xl, 2xl
	Type      string `json:"type,omitempty"`    // button, submit, reset
	Loading   bool   `json:"loading,omitempty"`
	FullWidth bool   `json:"fullWidth,omitempty"`
	IconOnly  bool   `json:"iconOnly,omitempty"`
	AriaLabel string `json:"ariaLabel,omitempty"`
}

// ToMap converts ButtonProps to map[string]any for unified Component construction.
func (p ButtonProps) ToMap() map[string]any { return propsToMap(p) }
