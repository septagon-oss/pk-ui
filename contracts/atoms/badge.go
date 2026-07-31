package atoms

// Implements: REQ-011.
// Per: ADR-0031.
// Discipline: C-14.
import "github.com/septagon-oss/pk-ui/contracts"

// BadgeProps defines the platform-agnostic properties for a Badge component.
type BadgeProps struct {
	contracts.ComponentProps

	Label   string `json:"label"`
	Variant string `json:"variant,omitempty"` // primary, secondary, outline
	Tone    string `json:"tone,omitempty"`    // neutral, brand, success, warning, danger, info
	Size    string `json:"size,omitempty"`    // xs, sm, md, lg, xl, 2xl
	Dot     bool   `json:"dot,omitempty"`     // show status dot before the label
}

// ToMap converts BadgeProps to map[string]any for unified Component construction.
func (p BadgeProps) ToMap() map[string]any { return propsToMap(p) }
