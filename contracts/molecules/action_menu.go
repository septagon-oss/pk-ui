package molecules

// Implements: REQ-011.
// Per: ADR-0031.
// Discipline: C-14.
import "github.com/septagon-oss/pk-ui/contracts"

// ActionMenuProps defines properties for a context/action menu.
type ActionMenuProps struct {
	contracts.ComponentProps

	TriggerLabel string              `json:"triggerLabel,omitempty"`
	TriggerIcon  string              `json:"triggerIcon,omitempty"` // default "ellipsis-vertical"
	Align        string              `json:"align,omitempty"`       // start, end
	Width        string              `json:"width,omitempty"`       // Tailwind width class
	Sections     []ActionMenuSection `json:"sections,omitempty"`
}

// ActionMenuSection groups items under an optional heading.
type ActionMenuSection struct {
	Label string           `json:"label,omitempty"`
	Items []ActionMenuItem `json:"items"`
}

// ActionMenuItem represents a single menu entry.
type ActionMenuItem struct {
	Label     string `json:"label"`
	Icon      string `json:"icon,omitempty"`
	Href      string `json:"href,omitempty"`
	Danger    bool   `json:"danger,omitempty"`
	Disabled  bool   `json:"disabled,omitempty"`
	HxGet     string `json:"hxGet,omitempty"`
	HxPost    string `json:"hxPost,omitempty"`
	HxDelete  string `json:"hxDelete,omitempty"`
	HxTarget  string `json:"hxTarget,omitempty"`
	HxConfirm string `json:"hxConfirm,omitempty"`
}

// ToMap converts ActionMenuProps to map[string]any for unified Component construction.
func (p ActionMenuProps) ToMap() map[string]any { return propsToMap(p) }
