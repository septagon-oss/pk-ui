package molecules

// Implements: REQ-011.
// Per: ADR-0031.
// Discipline: C-14.
import "github.com/septagon-oss/pk-ui/contracts"

// DrawerProps defines properties for a slide-out panel overlay.
type DrawerProps struct {
	contracts.ComponentProps

	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Body        string `json:"body,omitempty"`
	Footer      string `json:"footer,omitempty"`
	AriaLabel   string `json:"ariaLabel,omitempty"`
	CloseLabel  string `json:"closeLabel,omitempty"`
	Position    string `json:"position,omitempty"` // left, right, bottom
	Size        string `json:"size,omitempty"`     // small, medium, large, xl, full

	// Pointer booleans preserve the intended default-true behavior while still
	// allowing portable clients to explicitly disable an affordance.
	Closable       *bool `json:"closable,omitempty"`
	CloseOnOverlay *bool `json:"closeOnOverlay,omitempty"`
	CloseOnEscape  *bool `json:"closeOnEscape,omitempty"`
	ShowOverlay    *bool `json:"showOverlay,omitempty"`
	Open           bool  `json:"open,omitempty"`
	OpenOnSwap     bool  `json:"openOnSwap,omitempty"`
}

// ToMap converts DrawerProps to map[string]any for unified Component construction.
func (p DrawerProps) ToMap() map[string]any { return propsToMap(p) }
