package atoms

// Implements: REQ-011.
// Per: ADR-0031.
// Discipline: C-14.
import "github.com/septagon-oss/pk-ui/contracts"

// SelectProps defines platform-agnostic properties for a native single-value
// Select component. It mirrors InputProps where the concepts overlap so form
// builders can treat text-like and choice-like fields uniformly.
type SelectProps struct {
	contracts.ComponentProps
	contracts.HTMXProps

	Name        string         `json:"name"`
	Label       string         `json:"label,omitempty"`
	Value       string         `json:"value,omitempty"`
	Values      []string       `json:"values,omitempty"`
	Placeholder string         `json:"placeholder,omitempty"` // rendered as a disabled-free empty option
	Options     []SelectOption `json:"options"`
	Required    bool           `json:"required,omitempty"`
	Multiple    bool           `json:"multiple,omitempty"`
	VisibleRows int            `json:"visibleRows,omitempty"`
	FullWidth   bool           `json:"fullWidth,omitempty"`
	HelpText    string         `json:"helpText,omitempty"`
	Error       string         `json:"error,omitempty"`
}

// SelectOption is one choice in a Select.
type SelectOption struct {
	Label       string `json:"label"`
	Value       string `json:"value"`
	Description string `json:"description,omitempty"`
	Group       string `json:"group,omitempty"`
	Disabled    bool   `json:"disabled,omitempty"`
}

// ToMap converts SelectProps to map[string]any for unified Component construction.
func (p SelectProps) ToMap() map[string]any { return propsToMap(p) }
