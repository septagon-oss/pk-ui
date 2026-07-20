package atoms

// Implements: REQ-011.
// Per: ADR-0031.
// Discipline: C-14.
import "github.com/septagon-oss/platformkit-ui/contracts"

// DividerProps defines properties for a divider/separator.
type DividerProps struct {
	contracts.ComponentProps

	Orientation string `json:"orientation,omitempty"` // horizontal, vertical
	Text        string `json:"text,omitempty"`        // optional label (e.g., "OR")
}

// ToMap converts DividerProps to map[string]any for unified Component construction.
func (p DividerProps) ToMap() map[string]any { return propsToMap(p) }

// AvatarProps defines properties for an avatar.
type AvatarProps struct {
	contracts.ComponentProps

	Src            string `json:"src,omitempty"`
	Alt            string `json:"alt,omitempty"`
	Name           string `json:"name,omitempty"`
	Initials       string `json:"initials,omitempty"`
	FallbackIcon   string `json:"fallbackIcon,omitempty"`
	Size           string `json:"size,omitempty"`           // xs, small, medium, large, xl, 2xl
	Shape          string `json:"shape,omitempty"`          // circle, rounded, square, pill
	Status         string `json:"status,omitempty"`         // online, offline, busy, away
	StatusPosition string `json:"statusPosition,omitempty"` // top-right, bottom-right, top-left, bottom-left
}

// ToMap converts AvatarProps to map[string]any for unified Component construction.
func (p AvatarProps) ToMap() map[string]any { return propsToMap(p) }

// LinkProps defines properties for a hyperlink.
type LinkProps struct {
	contracts.ComponentProps
	contracts.HTMXProps

	Text              string `json:"text"`
	Href              string `json:"href"`
	External          bool   `json:"external,omitempty"`          // opens in new tab
	TrailingAdornment string `json:"trailingAdornment,omitempty"` // optional trailing icon or token
	Variant           string `json:"variant,omitempty"`           // primary, secondary, text, underline
	Target            string `json:"target,omitempty"`
	Rel               string `json:"rel,omitempty"`
}

// ToMap converts LinkProps to map[string]any for unified Component construction.
func (p LinkProps) ToMap() map[string]any { return propsToMap(p) }

// TagProps defines properties for an interactive tag/chip.
type TagProps struct {
	contracts.ComponentProps

	Text        string `json:"text"`
	Icon        string `json:"icon,omitempty"`
	Removable   bool   `json:"removable,omitempty"`
	Selected    bool   `json:"selected,omitempty"`
	Variant     string `json:"variant,omitempty"` // default, primary, success, warning, error
	OnRemoveURL string `json:"onRemoveURL,omitempty"`
}

// ToMap converts TagProps to map[string]any for unified Component construction.
func (p TagProps) ToMap() map[string]any { return propsToMap(p) }
