package atoms

// Implements: REQ-011.
// Per: ADR-0031.
// Discipline: C-14.
import "github.com/septagon-oss/pk-ui/contracts"

// AlertProps defines properties for a persistent inline status message.
type AlertProps struct {
	contracts.ComponentProps

	Message     string `json:"message"`
	Title       string `json:"title,omitempty"`
	Tone        string `json:"tone,omitempty"` // info, success, warning, danger
	Dismissible bool   `json:"dismissible,omitempty"`
	Bordered    bool   `json:"bordered,omitempty"`
	Compact     bool   `json:"compact,omitempty"`
}

// ToMap converts AlertProps to map[string]any for unified Component construction.
func (p AlertProps) ToMap() map[string]any { return propsToMap(p) }
