package molecules

// Implements: REQ-011.
// Per: ADR-0031.
// Discipline: C-14.
import "github.com/septagon-oss/pk-ui/contracts"

// FormProps defines platform-agnostic properties for a Form component.
type FormProps struct {
	contracts.ComponentProps
	contracts.HTMXProps

	Action    string      `json:"action,omitempty"`
	Method    string      `json:"method,omitempty"` // GET, POST
	Fields    []FormField `json:"fields,omitempty"`
	SubmitURL string      `json:"submitURL,omitempty"`
}

// ToMap converts FormProps to map[string]any for unified Component construction.
func (p FormProps) ToMap() map[string]any { return propsToMap(p) }

// CheckboxGroupProps defines a native multiple-choice fieldset. Required is a
// group-level semantic and must also be enforced by the server because HTML
// has no native "at least one checkbox" constraint.
type CheckboxGroupProps struct {
	contracts.ComponentProps
	contracts.HTMXProps

	Name        string   `json:"name"`
	Label       string   `json:"label,omitempty"`
	Description string   `json:"description,omitempty"`
	AriaLabel   string   `json:"ariaLabel,omitempty"`
	Options     []Option `json:"options"`
	Selected    []string `json:"selected,omitempty"`
	Required    bool     `json:"required,omitempty"`
	Orientation string   `json:"orientation,omitempty"` // vertical, horizontal
	Error       string   `json:"error,omitempty"`
}

// ToMap converts CheckboxGroupProps to map[string]any for unified Component construction.
func (p CheckboxGroupProps) ToMap() map[string]any { return propsToMap(p) }

// RadioGroupProps defines a native single-choice fieldset. Browsers own its
// arrow-key, Tab, Space, required, and submission behavior.
type RadioGroupProps struct {
	contracts.ComponentProps
	contracts.HTMXProps

	Name        string   `json:"name"`
	Label       string   `json:"label,omitempty"`
	Description string   `json:"description,omitempty"`
	AriaLabel   string   `json:"ariaLabel,omitempty"`
	Options     []Option `json:"options"`
	Value       string   `json:"value,omitempty"`
	Required    bool     `json:"required,omitempty"`
	Orientation string   `json:"orientation,omitempty"` // vertical, horizontal
	Error       string   `json:"error,omitempty"`
}

// ToMap converts RadioGroupProps to map[string]any for unified Component construction.
func (p RadioGroupProps) ToMap() map[string]any { return propsToMap(p) }

// FormField represents a field within a form.
type FormField struct {
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Label       string   `json:"label,omitempty"`
	Placeholder string   `json:"placeholder,omitempty"`
	Value       string   `json:"value,omitempty"`
	Required    bool     `json:"required,omitempty"`
	HelpText    string   `json:"helpText,omitempty"`
	Options     []Option `json:"options,omitempty"` // for select/radio
}

// Option represents a selectable option.
type Option struct {
	Label       string `json:"label"`
	Value       string `json:"value"`
	Description string `json:"description,omitempty"`
	Icon        string `json:"icon,omitempty"`
	Group       string `json:"group,omitempty"`
	Disabled    bool   `json:"disabled,omitempty"`
}
