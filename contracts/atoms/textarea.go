package atoms

// Implements: REQ-011.
// Per: ADR-0031.
// Discipline: C-14.
import "github.com/septagon-oss/platformkit-ui/contracts"

// TextareaProps defines properties for a multi-line text input.
type TextareaProps struct {
	contracts.ComponentProps
	contracts.HTMXProps

	Name         string `json:"name"`
	Placeholder  string `json:"placeholder,omitempty"`
	Value        string `json:"value,omitempty"`
	Label        string `json:"label,omitempty"`
	HelperText   string `json:"helperText,omitempty"`
	ErrorMessage string `json:"errorMessage,omitempty"`
	Required     bool   `json:"required,omitempty"`
	ReadOnly     bool   `json:"readOnly,omitempty"`
	Rows         int    `json:"rows,omitempty"`
	MinRows      int    `json:"minRows,omitempty"`
	MaxRows      int    `json:"maxRows,omitempty"`
	MaxLength    int    `json:"maxLength,omitempty"`
	ShowCount    bool   `json:"showCount,omitempty"`
	AutoResize   bool   `json:"autoResize,omitempty"`
	FullWidth    bool   `json:"fullWidth,omitempty"`
}

// ToMap converts TextareaProps to map[string]any for unified Component construction.
func (p TextareaProps) ToMap() map[string]any { return propsToMap(p) }
