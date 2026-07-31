package atoms

// Implements: REQ-011.
// Per: ADR-0031.
// Discipline: C-14.
import "github.com/septagon-oss/pk-ui/contracts"

// EmptyStateProps defines properties for an empty data state placeholder.
type EmptyStateProps struct {
	contracts.ComponentProps

	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Compact     bool   `json:"compact,omitempty"`
	Bordered    bool   `json:"bordered,omitempty"`
}

// ToMap converts EmptyStateProps to map[string]any for unified Component construction.
func (p EmptyStateProps) ToMap() map[string]any { return propsToMap(p) }
