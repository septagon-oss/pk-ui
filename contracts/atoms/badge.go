package atoms

// Implements: REQ-011.
// Per: ADR-0031.
// Discipline: C-14.
import "github.com/septagon-oss/pk-ui/contracts"

// BadgeProps defines the platform-agnostic properties for a Badge component.
type BadgeProps struct {
	contracts.ComponentProps

	Label       string `json:"label"`
	Variant     string `json:"variant,omitempty"` // primary, secondary, outline
	Tone        string `json:"tone,omitempty"`    // neutral, brand, success, warning, danger, info
	Size        string `json:"size,omitempty"`    // xs, sm, md, lg, xl, 2xl
	Dot         bool   `json:"dot,omitempty"`     // show status dot before the label
	Count       int    `json:"count,omitempty"`   // positive count, visually capped to 99+
	Removable   bool   `json:"removable,omitempty"`
	RemoveLabel string `json:"removeLabel,omitempty"` // localized remove-button label
	Live        bool   `json:"live,omitempty"`        // polite status announcement
}

// ToMap converts BadgeProps to map[string]any for unified Component construction.
func (p BadgeProps) ToMap() map[string]any { return propsToMap(p) }
