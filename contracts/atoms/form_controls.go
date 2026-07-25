package atoms

// Implements: REQ-011.
// Per: ADR-0031.
// Discipline: C-14.
import "github.com/septagon-oss/pk-ui/contracts"

// CheckboxProps defines properties for a checkbox input.
type CheckboxProps struct {
	contracts.ComponentProps

	Name          string `json:"name"`
	Label         string `json:"label"`
	Checked       bool   `json:"checked,omitempty"`
	Indeterminate bool   `json:"indeterminate,omitempty"`
	Value         string `json:"value,omitempty"`
	Required      bool   `json:"required,omitempty"`
}

// ToMap converts CheckboxProps to map[string]any for unified Component construction.
func (p CheckboxProps) ToMap() map[string]any { return propsToMap(p) }

// RadioProps defines properties for a radio button.
type RadioProps struct {
	contracts.ComponentProps

	Name     string `json:"name"`
	Label    string `json:"label"`
	Value    string `json:"value"`
	Checked  bool   `json:"checked,omitempty"`
	Required bool   `json:"required,omitempty"`
}

// ToMap converts RadioProps to map[string]any for unified Component construction.
func (p RadioProps) ToMap() map[string]any { return propsToMap(p) }

// ToggleProps defines properties for a toggle/switch.
type ToggleProps struct {
	contracts.ComponentProps

	Name    string `json:"name"`
	Label   string `json:"label"`
	Checked bool   `json:"checked,omitempty"`
	Size    string `json:"size,omitempty"` // small, medium, large
}

// ToMap converts ToggleProps to map[string]any for unified Component construction.
func (p ToggleProps) ToMap() map[string]any { return propsToMap(p) }

// SliderProps defines properties for a range slider.
type SliderProps struct {
	contracts.ComponentProps

	Name      string  `json:"name"`
	Min       float64 `json:"min"`
	Max       float64 `json:"max"`
	Step      float64 `json:"step,omitempty"`
	Value     float64 `json:"value,omitempty"`
	ShowValue bool    `json:"showValue,omitempty"`
	Label     string  `json:"label,omitempty"`
}

// ToMap converts SliderProps to map[string]any for unified Component construction.
func (p SliderProps) ToMap() map[string]any { return propsToMap(p) }
