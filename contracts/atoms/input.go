package atoms

// Implements: REQ-011.
// Per: ADR-0031.
// Discipline: C-14.
import "github.com/septagon-oss/pk-ui/contracts"

// InputProps defines the platform-agnostic properties for an Input component.
type InputProps struct {
	contracts.ComponentProps
	contracts.HTMXProps

	Name         string `json:"name"`
	Type         string `json:"type,omitempty"` // text, email, password, number, tel, url, search, date, time
	Value        string `json:"value,omitempty"`
	Placeholder  string `json:"placeholder,omitempty"`
	Label        string `json:"label,omitempty"`
	HelpText     string `json:"helpText,omitempty"`
	Error        string `json:"error,omitempty"`
	Invalid      bool   `json:"invalid,omitempty"`
	Required     bool   `json:"required,omitempty"`
	ReadOnly     bool   `json:"readOnly,omitempty"`
	AutoFocus    bool   `json:"autoFocus,omitempty"`
	Min          string `json:"min,omitempty"`
	Max          string `json:"max,omitempty"`
	Step         string `json:"step,omitempty"`
	MinLength    int    `json:"minLength,omitempty"`
	MaxLength    int    `json:"maxLength,omitempty"`
	Pattern      string `json:"pattern,omitempty"`
	Autocomplete string `json:"autocomplete,omitempty"`
	Size         string `json:"size,omitempty"` // sm, md, lg
	Tone         string `json:"tone,omitempty"` // neutral, success, warning, danger
	FullWidth    bool   `json:"fullWidth,omitempty"`
}

// ToMap converts InputProps to map[string]any for unified Component construction.
func (p InputProps) ToMap() map[string]any { return propsToMap(p) }
