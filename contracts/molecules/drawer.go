package molecules

import "github.com/septagon-oss/platformkit-ui/contracts"

// DrawerProps defines properties for a slide-out panel overlay.
type DrawerProps struct {
	contracts.ComponentProps

	Title          string `json:"title,omitempty"`
	Description    string `json:"description,omitempty"`
	Position       string `json:"position,omitempty"` // left, right, bottom
	Size           string `json:"size,omitempty"`     // small, medium, large, xl, full
	Closable       bool   `json:"closable,omitempty"`
	CloseOnOverlay bool   `json:"closeOnOverlay,omitempty"`
	CloseOnEscape  bool   `json:"closeOnEscape,omitempty"`
	ShowOverlay    bool   `json:"showOverlay,omitempty"`
}

// ToMap converts DrawerProps to map[string]any for unified Component construction.
func (p DrawerProps) ToMap() map[string]any { return propsToMap(p) }
