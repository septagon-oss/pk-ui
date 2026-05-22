package molecules

import "github.com/septagon-oss/platformkit-ui/contracts"

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
	Label    string `json:"label"`
	Value    string `json:"value"`
	Disabled bool   `json:"disabled,omitempty"`
}
