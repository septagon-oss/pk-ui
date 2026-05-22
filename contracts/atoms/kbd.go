package atoms

import "github.com/septagon-oss/platformkit-ui/contracts"

// KbdProps defines properties for a keyboard shortcut display chip.
type KbdProps struct {
	contracts.ComponentProps

	Keys []string `json:"keys"`           // Key labels (e.g. ["⌘", "K"])
	Size string   `json:"size,omitempty"` // xs, small, medium, large
}

// ToMap converts KbdProps to map[string]any for unified Component construction.
func (p KbdProps) ToMap() map[string]any { return propsToMap(p) }
