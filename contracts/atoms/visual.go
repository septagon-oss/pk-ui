package atoms

// Implements: REQ-011.
// Per: ADR-0031.
// Discipline: C-14.
import "github.com/septagon-oss/pk-ui/contracts"

// IconProps defines a provider-neutral system glyph.
type IconProps struct {
	contracts.ComponentProps

	Name      string `json:"name"`
	Size      string `json:"size,omitempty"`   // xs, sm, md, lg, xl, 2xl
	Tone      string `json:"tone,omitempty"`   // neutral, brand, success, warning, danger, info
	Weight    string `json:"weight,omitempty"` // outline; extension providers may add governed weights
	AriaLabel string `json:"ariaLabel,omitempty"`
}

// ToMap converts IconProps to map[string]any for unified Component construction.
func (p IconProps) ToMap() map[string]any { return propsToMap(p) }

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
	Size           string `json:"size,omitempty"`           // xs, sm, md, lg, xl, 2xl (canonical control scale; legacy words small/medium/large still accepted on render)
	Shape          string `json:"shape,omitempty"`          // circle, rounded, square, pill
	Tone           string `json:"tone,omitempty"`           // neutral, brand, success, warning, danger, info
	Status         string `json:"status,omitempty"`         // online, offline, busy, away
	StatusPosition string `json:"statusPosition,omitempty"` // top-right, bottom-right, top-left, bottom-left
	StatusLabel    string `json:"statusLabel,omitempty"`    // localized accessible presence label
}

// ToMap converts AvatarProps to map[string]any for unified Component construction.
func (p AvatarProps) ToMap() map[string]any { return propsToMap(p) }

// LinkProps defines properties for a hyperlink.
type LinkProps struct {
	contracts.ComponentProps
	contracts.HTMXProps

	Label    string `json:"label"`
	Href     string `json:"href"`
	External bool   `json:"external,omitempty"` // opens in new tab
	Variant  string `json:"variant,omitempty"`  // primary, secondary, text, underline
	Target   string `json:"target,omitempty"`
	Rel      string `json:"rel,omitempty"`
}

// ToMap converts LinkProps to map[string]any for unified Component construction.
func (p LinkProps) ToMap() map[string]any { return propsToMap(p) }

// TagProps defines properties for an interactive tag/chip.
type TagProps struct {
	contracts.ComponentProps

	Label       string `json:"label"`
	Tone        string `json:"tone,omitempty"` // neutral, brand, success, warning, danger, info
	Removable   bool   `json:"removable,omitempty"`
	Selected    bool   `json:"selected,omitempty"`
	OnRemoveURL string `json:"onRemoveURL,omitempty"`
}

// ToMap converts TagProps to map[string]any for unified Component construction.
func (p TagProps) ToMap() map[string]any { return propsToMap(p) }
